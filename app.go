package demo

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gourses/demo/api"
	"github.com/gourses/demo/storage"
)

func Run() error {
	const (
		addr    = "127.0.0.1:8080"
		connStr = "postgres://user:passwd@localhost:5432/db"
	)

	s, err := storage.New(connStr)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	a := api.New(addr, s)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt)
		<-exit
		a.Close()
	}()

	return a.Run()
}
