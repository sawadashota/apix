package main

import (
	"fmt"
	"os"

	"github.com/sawadashota/apix"
	"github.com/sawadashota/apix/health"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	srv, err := newServer()
	if err != nil {
		return err
	}
	defer srv.Registry().DB().Close()

	return srv.Serve()
}

func newServer() (*apix.Server, error) {
	srv, err := apix.NewDefaultServer()
	if err != nil {
		return nil, err
	}

	// Add handlers blow
	srv.Router().Register(health.NewHandler(srv))

	return srv, nil
}
