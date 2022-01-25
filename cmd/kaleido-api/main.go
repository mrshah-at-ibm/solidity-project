package main

import (
	"fmt"

	"github.com/mrshah-at-ibm/kaleido-project/pkg/executer"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/routes"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		panic(errors.Wrap(err, "Error configuring logger"))
	}

	defer logger.Sync()

	e, err := executer.NewExecuter(logger)
	if err != nil {
		panic(errors.Wrap(err, "Unable to setup executer"))
	}

	err = e.DeployContract()
	if err != nil {
		panic(errors.Wrap(err, "Unable to setup contract"))
	}

	r, err := routes.New(logger, e)
	r.SetupRoutes()
	if err != nil {
		panic(errors.Wrap(err, "Error setting up app"))
	}
	fmt.Println("Starting Server")
	err = r.StartServer()
	if err != nil {
		panic(errors.Wrap(err, "Error starting server"))
	}

	fmt.Println("Server Running")
}
