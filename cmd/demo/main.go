package main

import (
	"log"

	"github.com/gourses/demo"
	"golang.org/x/exp/slog"
)

func main() {
	slog.Info("Starting App")
	if err := demo.Run(); err != nil {
		log.Fatal(err)
	}
	slog.Info("App exited successfully")
}
