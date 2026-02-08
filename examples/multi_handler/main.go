package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aN0mad/go-logo/v2"
)

func main() {
	textFormatter := logo.NewTextFormatter()
	consoleHandler := logo.NewWriterHandler(os.Stdout, textFormatter, logo.LevelInfo)
	debugBuffer := &bytes.Buffer{}
	debugHandler := logo.NewWriterHandler(debugBuffer, textFormatter, logo.LevelDebug)
	router := logo.NewRouter(consoleHandler, debugHandler)
	log := logo.New(logo.WithRouter(router), logo.WithName("multi-handler"))

	fmt.Println("=== Console Output (INFO and above) ===")
	log.Debug("This debug message only goes to the buffer")
	log.Info("This info message goes to both console and buffer")
	log.Warn("This warning goes to both outputs")
	log.Error("This error goes to both outputs")

	fmt.Println("\n=== Debug Buffer Content (DEBUG and above) ===")
	fmt.Print(debugBuffer.String())

	fmt.Println("\n=== Adding another handler dynamically ===")
	errorBuffer := &bytes.Buffer{}
	errorHandler := logo.NewWriterHandler(errorBuffer, textFormatter, logo.LevelError)
	router.AddHandler(errorHandler)

	log.Error("This error goes to all three handlers")
	log.Info("This info only goes to console and debug buffer")

	fmt.Println("\n=== Error Buffer Content (ERROR only) ===")
	fmt.Print(errorBuffer.String())

	log.Close()
}
