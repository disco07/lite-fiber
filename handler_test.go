package lite

import (
	"bytes"
	"github.com/go-lite/lite/errors"
	"github.com/go-lite/lite/mime"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type HandlerTestSuite struct {
	suite.Suite
}

func (suite *HandlerTestSuite) SetupTest() {
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

type requestApplicationJSON struct {
	ID uint64 `lite:"path=id"`
}

type responserequestApplicationJSON struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_ApplicationJSON_Requests() {
	app := New()
	Get(app, "/foo/:id", func(c *ContextWithRequest[requestApplicationJSON]) (responserequestApplicationJSON, error) {
		req, err := c.Requests()
		if err != nil {
			return responserequestApplicationJSON{}, err
		}

		return responserequestApplicationJSON{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo/123", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

type requestApplicationXML struct {
	ID uint64 `lite:"path=id"`
}

type responseApplicationXML struct {
	ID      uint64 `xml:"id"`
	Message string `xml:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_ApplicationXML_Requests() {
	app := New()
	Get(app, "/foo/:id", func(c *ContextWithRequest[requestApplicationXML]) (responseApplicationXML, error) {
		req, err := c.Requests()
		if err != nil {
			return responseApplicationXML{}, err
		}

		c.SetContentType(mime.ApplicationXML)

		return responseApplicationXML{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo/123", nil)
	req.Header.Set("Content-Type", "application/xml")
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), "<responseApplicationXML><id>123</id><message>Hello World</message></responseApplicationXML>", utils.UnsafeString(body))
}

type requestApplicationJSONError struct {
	ID uint64 `lite:"path=id"`
}

type responseApplicationJSONError struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_ApplicationJSON_Requests_Error() {
	app := New()
	Get(app, "/foo/:id", func(c *ContextWithRequest[requestApplicationJSONError]) (responseApplicationJSONError, error) {
		return responseApplicationJSONError{}, assert.AnError
	})

	req := httptest.NewRequest("GET", "/foo/123", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type requestPath struct {
	ID uint64 `lite:"path=id"`
}

type responsePath struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Path_Error() {

	app := New()
	Get(app, "/foo/:id", func(c *ContextWithRequest[requestPath]) (responsePath, error) {
		req, err := c.Requests()
		if err != nil {
			return responsePath{}, err
		}

		if req.ID == 0 {
			return responsePath{}, errors.NewBadRequestError("ID is required")
		}

		return responsePath{}, nil
	})

	req := httptest.NewRequest("GET", "/foo/abc", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 400")
}

type requestQuery struct {
	ID uint64 `lite:"query=id"`
}

type responseQuery struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Query() {
	app := New()
	Get(app, "/foo", func(c *ContextWithRequest[requestQuery]) (responseQuery, error) {
		req, err := c.Requests()
		if err != nil {
			return responseQuery{}, err
		}

		return responseQuery{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo?id=123", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

func (suite *HandlerTestSuite) TestContextWithRequest_Query_Error() {
	app := New()
	Get(app, "/foo", func(c *ContextWithRequest[requestPath]) (responsePath, error) {
		return responsePath{}, assert.AnError
	})

	req := httptest.NewRequest("GET", "/foo?id=abc", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type requestHeader struct {
	ID uint64 `lite:"header=id"`
}

type responseHeader struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Header() {

	app := New()
	Get(app, "/foo", func(c *ContextWithRequest[requestHeader]) (responseHeader, error) {
		req, err := c.Requests()
		if err != nil {
			return responseHeader{}, err
		}

		return responseHeader{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo", nil)
	req.Header.Set("id", "123")
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

type requestHeaderErr struct {
	ID uint64 `lite:"header=id"`
}

type responseHeaderErr struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Header_Error() {
	app := New()
	Get(app, "/foo", func(c *ContextWithRequest[requestHeaderErr]) (responseHeaderErr, error) {
		return responseHeaderErr{}, assert.AnError
	})

	req := httptest.NewRequest("GET", "/foo", nil)
	req.Header.Set("id", "abc")
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type reqBody struct {
	ID float64 `json:"id" xml:"id"`
}

type requestBody struct {
	Body reqBody `lite:"req=body"`
}

type responseBody struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body() {

	app := New()
	Post(app, "/foo", func(c *ContextWithRequest[requestBody]) (responseBody, error) {
		req, err := c.Requests()
		if err != nil {
			return responseBody{}, err
		}

		return responseBody{
			ID:      req.Body.ID,
			Message: "Hello World",
		}, nil
	})

	bodyJSON := `{"id":123}`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 200")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

type requestBodyError struct {
	Body reqBody `lite:"req=body"`
}

type responseBodyError struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_Error() {
	app := New()
	Post(app, "/foo", func(c *ContextWithRequest[requestBodyError]) (responseBodyError, error) {
		return responseBodyError{}, assert.AnError
	})

	bodyJSON := `{"id":"abc"}`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type requestBodyXML struct {
	Body reqBody `lite:"req=body"`
}

type responseBodyXML struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_ApplicationXML() {
	app := New()
	Post(app, "/foo", func(c *ContextWithRequest[requestBodyXML]) (responseBodyXML, error) {
		req, err := c.Requests()
		if err != nil {
			return responseBodyXML{}, err
		}

		return responseBodyXML{
			ID:      req.Body.ID,
			Message: "Hello World",
		}, nil
	})

	bdyXML := `<request><id>123</id></request>`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(bdyXML))
	req.Header.Set("Content-Type", "application/xml")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 200")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

type requestBodyXMLError struct {
	Body reqBody `lite:"req=body"`
}

type responseBodyXMLError struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_ApplicationXML_Error() {
	app := New()
	Post(app, "/foo", func(c *ContextWithRequest[requestBodyXMLError]) (responseBodyXMLError, error) {
		return responseBodyXMLError{}, assert.AnError
	})

	bodyXML := `<request><id>abc</id></request>`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(bodyXML))
	req.Header.Set("Content-Type", "application/xml")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type requestBodyXMLInvalid struct {
	Body reqBody `lite:"req=body"`
}

type responseBodyXMLInvalid struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_ApplicationXML_Invalid() {
	app := New()
	Post(app, "/foo", func(c *ContextWithRequest[requestBodyXMLInvalid]) (responseBodyXMLInvalid, error) {
		return responseBodyXMLInvalid{}, assert.AnError
	})

	requestBody := `<request><id>abc</id>`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/xml")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type requestBodyPut struct {
	Body reqBody `lite:"req=body"`
}

type responseBodyPut struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_Put() {
	app := New()
	Put(app, "/foo", func(c *ContextWithRequest[requestBodyPut]) (responseBodyPut, error) {
		req, err := c.Requests()
		if err != nil {
			return responseBodyPut{}, err
		}

		return responseBodyPut{
			ID:      req.Body.ID,
			Message: "Hello World",
		}, nil
	})

	bodyJSON := `{"id":123}`
	req := httptest.NewRequest("PUT", "/foo", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

type requestDelete struct {
	ID uint64 `lite:"path=id"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Delete() {
	app := New()
	Delete(app, "/foo/:id", func(c *ContextWithRequest[requestDelete]) (ret struct{}, err error) {
		req, err := c.Requests()
		if err != nil {
			return
		}

		if req.ID == 0 {
			err = errors.NewBadRequestError("ID is required")

			return
		}

		return
	})

	req := httptest.NewRequest("DELETE", "/foo/123", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 204, resp.StatusCode, "Expected status code 204")
}

type requestPatch struct {
	ID uint64 `lite:"path=id"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Patch() {
	app := New()
	Patch(app, "/foo/:id", func(c *ContextWithRequest[requestPatch]) (ret struct{}, err error) {
		req, err := c.Requests()
		if err != nil {
			return
		}

		if req.ID == 0 {
			err = errors.NewBadRequestError("ID is required")

			return
		}

		return
	})

	req := httptest.NewRequest("PATCH", "/foo/123", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

type requestPatchError struct {
	ID uint64 `lite:"path=id"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_PatchError() {
	app := New()
	Patch(app, "/foo/:id", func(c *ContextWithRequest[requestPatchError]) (ret struct{}, err error) {
		req, err := c.Requests()
		if err != nil {
			return
		}

		if req.ID == 0 {
			err = errors.NewBadRequestError("ID is required")

			return
		}

		return
	})

	req := httptest.NewRequest("PATCH", "/foo/0", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 400, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Head() {
	app := New()
	Head(app, "/foo", func(c *ContextNoRequest) (ret struct{}, err error) {
		c.SetContentType(mime.ApplicationJSON)
		c.Status(200)

		return
	})

	req := httptest.NewRequest("HEAD", "/foo", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Options() {
	app := New()
	Options(app, "/foo/", func(c *ContextNoRequest) (ret struct{}, err error) {

		return
	})

	req := httptest.NewRequest("OPTIONS", "/foo", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Trace() {
	app := New()
	Trace(app, "/foo/", func(c *ContextNoRequest) (ret struct{}, err error) {

		return
	})

	req := httptest.NewRequest("TRACE", "/foo", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Connect() {
	app := New()
	Connect(app, "/foo/", func(c *ContextNoRequest) (ret struct{}, err error) {

		return
	})

	req := httptest.NewRequest("CONNECT", "/foo", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

type fakeContext[r any] struct {
	*ContextNoRequest
}

func (c fakeContext[r]) Requests() (r, error) {
	var ret r
	return ret, nil
}

func (c fakeContext[r]) Status(status int) Context[r] {
	return nil
}

func (c fakeContext[r]) SetContentType(extension string, charset ...string) Context[r] {
	return nil
}

func (suite *HandlerTestSuite) TestCustomContext() {
	assert.Panicsf(suite.T(), func() {
		c := newLiteContext[request, fakeContext[request]](ContextNoRequest{})
		assert.Nil(suite.T(), c)
	}, "unknown type")
}

type requestRoute struct {
	ID uint64 `lite:"path=id"`
}

type responseRoute struct {
	ID uint64 `json:"id"`
}

func (suite *HandlerTestSuite) TestRegisterRoute() {
	app := New()

	assert.Panicsf(suite.T(), func() {
		_ = registerRoute[requestRoute, responseRoute](
			app,
			Route[requestRoute, responseRoute]{
				path:        "/foo/:id",
				method:      "GET",
				contentType: "application/json",
				statusCode:  200,
			},
			nil,
			func(ctx *fiber.Ctx) error {
				return nil
			},
		)
	}, "unknown parameter type")

}

func (suite *HandlerTestSuite) TestContextWithRequest_FullBody() {
	app := New()

	Post(app, "/test/:id/:is_admin", func(c *ContextWithRequest[testRequest]) (testResponse, error) {
		req, err := c.Requests()
		if err != nil {
			return testResponse{}, err
		}

		err = c.SaveFile(req.Body.File, "./logo/lite.png")
		if err != nil {
			return testResponse{}, err
		}

		return testResponse{
			ID:        req.Params.ID,
			FirstName: req.Body.Metadata.FirstName,
			LastName:  req.Body.Metadata.LastName,
		}, nil
	})

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Ajouter le fichier
	fileWriter, err := w.CreateFormFile("file", "lite.png")
	if err != nil {
		suite.T().Fatalf("Failed to create form file: %s", err)
	}

	file, err := os.Open("./logo/lite.png")
	if err != nil {
		suite.T().Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		suite.T().Fatalf("Failed to copy file: %s", err)
	}

	// Ajouter le JSON
	metadataWriter, err := w.CreateFormField("metadata")
	if err != nil {
		suite.T().Fatalf("Failed to create form field: %s", err)
	}
	data := `{"first_name":"John","last_name":"Doe"}`
	_, err = metadataWriter.Write([]byte(data))
	if err != nil {
		suite.T().Fatalf("Failed to write metadata: %s", err)
	}

	nameWriter, err := w.CreateFormField("name")
	if err != nil {
		suite.T().Fatalf("Failed to create form field: %s", err)
	}
	_, err = nameWriter.Write([]byte("test"))
	if err != nil {
		suite.T().Fatalf("Failed to write name: %s", err)
	}

	w.Close()

	req := httptest.NewRequest("POST", "/test/123/true", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 201")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"first_name":"John","last_name":"Doe", "name":""}`, utils.UnsafeString(body))

	spec, err := app.SaveOpenAPISpec()
	assert.NoError(suite.T(), err)

	expected := `components:
    parameters:
        cookie:
            in: cookie
            name: cookie
            schema:
                $ref: '#/components/schemas/cookie'
        filter:
            in: query
            name: filter
            schema:
                $ref: '#/components/schemas/filter'
        id:
            in: path
            name: id
            required: true
            schema:
                $ref: '#/components/schemas/id'
        is_admin:
            in: path
            name: is_admin
            required: true
            schema:
                $ref: '#/components/schemas/is_admin'
    schemas:
        bodyRequest:
            properties:
                file:
                    format: byte
                    type: string
                metadata:
                    properties:
                        first_name:
                            type: string
                        last_name:
                            type: string
                    type: object
                name:
                    type: string
            required:
                - name
                - file
            type: object
        cookie:
            properties:
                Domain:
                    type: string
                Expires:
                    format: date-time
                    type: string
                HttpOnly:
                    type: boolean
                MaxAge:
                    type: integer
                Name:
                    type: string
                Path:
                    type: string
                Raw:
                    type: string
                RawExpires:
                    type: string
                SameSite:
                    type: integer
                Secure:
                    type: boolean
                Unparsed:
                    items:
                        type: string
                    type: array
                Value:
                    type: string
            type: object
        filter:
            type: string
        httpGenericError:
            properties:
                id:
                    type: string
                message:
                    type: string
                status:
                    type: integer
            type: object
        id:
            maximum: 1.8446744073709552e+19
            minimum: 0
            type: integer
        is_admin:
            type: string
        testResponse:
            properties:
                first_name:
                    type: string
                id:
                    maximum: 1.8446744073709552e+19
                    minimum: 0
                    type: integer
                last_name:
                    type: string
                name:
                    type: string
            required:
                - id
                - name
                - first_name
                - last_name
            type: object
info:
    description: OpenAPI
    title: OpenAPI
    version: 0.0.1
openapi: 3.0.3
paths:
    /test/{id}/{is_admin}:
        post:
            operationId: POST/test/:id/:is_admin
            parameters:
                - $ref: '#/components/parameters/id'
                - $ref: '#/components/parameters/is_admin'
                - $ref: '#/components/parameters/filter'
                - $ref: '#/components/parameters/cookie'
            requestBody:
                content:
                    multipart/form-data:
                        schema:
                            $ref: '#/components/schemas/bodyRequest'
            responses:
                "201":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/testResponse'
                    description: OK
                "400":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Bad Request
                "401":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Unauthorized
                "404":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Not Found
                "409":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Conflict
                "500":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Internal Server Error`

	assert.YAMLEqf(suite.T(), expected, string(spec), "openapi generated spec")
}
