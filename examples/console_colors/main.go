package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aN0mad/go-logo/v2"
)

func main() {
	fmt.Println("=== Console Output with Auto Color Detection ===")
	log1 := logo.New(logo.WithConsole())
	log1.Debug("Debug message")
	log1.Info("Info message")
	log1.Warn("Warning message")
	log1.Error("Error message")

	fmt.Println("\n=== Output with Colors Always Enabled ===")
	buf := &bytes.Buffer{}
	handler := logo.NewConsoleHandlerWithWriter(buf, logo.ColorAlways, logo.LevelDebug)
	router := logo.NewRouter(handler)
	log2 := logo.New(logo.WithRouter(router))

	log2.Info("This has ANSI color codes")
	log2.Warn("Warning with colors")
	fmt.Print(buf.String())

	fmt.Println("\n=== Output with Colors Disabled ===")
	buf2 := &bytes.Buffer{}
	handler2 := logo.NewConsoleHandlerWithWriter(buf2, logo.ColorNever, logo.LevelDebug)
	router2 := logo.NewRouter(handler2)
	log3 := logo.New(logo.WithRouter(router2))

	log3.Info("Plain text, no colors")
	log3.Error("Error without ANSI codes")
	fmt.Print(buf2.String())

	fmt.Println("\n=== Output with Source Information ===")
	colorSupport := logo.NewColorSupport(os.Stdout, logo.ColorAuto)
	prettyFormatter := logo.NewPrettyFormatter(colorSupport, logo.WithSourceInfo(true))
	writerHandler := logo.NewWriterHandler(os.Stdout, prettyFormatter, logo.LevelDebug)
	router3 := logo.NewRouter(writerHandler)
	log4 := logo.New(logo.WithRouter(router3), logo.WithSource(true))

	log4.Info("This message includes source location")
	log4.Warn("Warning with file and line number")
}
