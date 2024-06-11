package main

import (
	"errors"
	"log"
	"os"

	"github.com/disco07/lite"
	"github.com/disco07/lite/examples/parameters"
	"github.com/disco07/lite/examples/returns"
)

// Define example handler
func getHandler(c *lite.ContextWithRequest[parameters.GetReq]) (returns.GetResponse, error) {
	request, err := c.Requests()
	if err != nil {
		return returns.GetResponse{}, err
	}

	if request.Params == "test" {
		return returns.GetResponse{}, errors.New("test is not valid name")
	}

	return returns.GetResponse{
		Message: "Hello World!, " + request.Params,
	}, nil
}

func postHandler(c *lite.ContextWithRequest[parameters.CreateReq]) (returns.CreateResponse, error) {
	body, err := c.Requests()
	if err != nil {
		return returns.CreateResponse{}, err
	}

	if body.Body.FirstName == "" {
		return returns.CreateResponse{}, errors.New("first_name are required")
	}

	return returns.CreateResponse{
		ID:        body.Params.ID,
		FirstName: body.Body.FirstName,
		LastName:  body.Body.LastName,
	}, nil
}

func getArrayHandler(_ *lite.ContextWithRequest[parameters.GetArrayReq]) (returns.GetArrayReturnsResponse, error) {
	res := make([]returns.Ret, 0)

	value := "value"
	res = append(res, returns.Ret{
		Message: "Hello World!",
		Embed: returns.Embed{
			Key:        "key",
			ValueEmbed: &value,
		},
	},
		returns.Ret{
			Message: "Hello World 2!",
			Embed: returns.Embed{
				Key: "key2",
			},
		},
	)

	return res, nil
}

func main() {
	app := lite.NewApp()

	lite.Get(app, "/example/:name", getHandler)

	lite.Post(app, "/example/:id", postHandler).
		OperationID("createExample").
		Description("Create example").
		AddTags("example")

	lite.Get(app, "/example", getArrayHandler)

	app.AddServer("http://localhost:6000", "example server")

	yamlBytes, err := app.SaveOpenAPISpec()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("./examples/api/openapi.yaml")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	_, err = f.Write(yamlBytes)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.Listen(":6000"))
}
