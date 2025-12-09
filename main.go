package main

import (
	"log/slog"
	"os"

	"sunny_land/src/engine/core"
)

func main() {
	minLevel := slog.LevelDebug
	options := &slog.HandlerOptions{
		Level: minLevel,
	}
	handler := slog.NewTextHandler(os.Stdout, options)
	slog.SetDefault(slog.New(handler))

	g := core.NewGameApp()
	g.Run()
}
