# go-logo
A golang logging library that allows logging to the console, to a channel or both.

## Makefile
There is a makefile in the repository that makes building and cleaning the folder structure easier

```bash
root@f6a960422947:/workspace# make help
 help: print this help message
 tidy: format code and tidy modfile
 build: build the unix version
 buildwin: build the windows version
 all: build all applications for unix and windows
 clean: clean the repository
```

## Build
The easiest way is to compile with the included Docker container or just on Linux in general. Golang cross-compilation makes this easy.
### Linux
```bash
make build
```

### Windows
```bash
make buildwin
```

## Usage
```bash
/workspace/bin ./go-logo.elf
Initialized logger
[14:11:09] time=2025-07-29T14:11:09.752Z level=DEBUG msg=This is a debug message source=main.go:21
[14:11:09] time=2025-07-29T14:11:09.752Z level=INFO msg=This is an info message source=main.go:22
[14:11:09] time=2025-07-29T14:11:09.752Z level=WARN msg=This is a warning message source=main.go:23
[14:11:09] time=2025-07-29T14:11:09.753Z level=ERROR msg=This is an error message source=main.go:24
[14:11:09] time=2025-07-29T14:11:09.753Z level=FATAL msg=This is a fatal message - It will exit the program source=main.go:25
exit status 1
```

### Colored output
![Colorful logs](assets/colored_logs.png)

## Example
Example source code is given in the [examples](examples/) directory.

### Logging to console
```golang
package main

import (
	"fmt"
	"log/slog"
	logger "logo/logo"
)

func main() {

    // Init the logger
	logger.Init(
		logger.AddSource(),
		logger.SetLevel(slog.LevelDebug),
		logger.AddFileOutput("logs/app.log", 10, 3, 30, true),
	)

    // Return logger for use
	log := logger.L()
	fmt.Println("Initialized logger")

    // Log messages
	log.Trace("This is a trace message - It will not show up")
	log.Debug("This is a debug message")
	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")
	log.Fatal("This is a fatal message - It will exit the program")
}
```

# TODO
- Cleanup .bak files