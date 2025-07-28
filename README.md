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
/workspace/bin ./logo
INFO lib: 📦 Starting the CLI tool...
DEBU lib: Initializing subsystems...
INFO lib: Processing data stream...
DEBU lib: Processing item 1
DEBU lib: Processing item 2
DEBU lib: Processing item 3
WARN lib: Potential issue with item 3
DEBU lib: Processing item 4
DEBU lib: Processing item 5
INFO lib: ✔️ Completed all tasks

INFO lib: 📦 Starting channel-based logger...
data: {"time":"2025-07-25T19:13:06.236221915Z","level":"INFO","message":"Heartbeat 1"}
data: {"time":"2025-07-25T19:13:07.237442147Z","level":"INFO","message":"Heartbeat 2"}
data: {"time":"2025-07-25T19:13:08.237713136Z","level":"INFO","message":"Heartbeat 3"}
data: {"time":"2025-07-25T19:13:09.238248676Z","level":"INFO","message":"Heartbeat 4"}
data: {"time":"2025-07-25T19:13:10.237674747Z","level":"INFO","message":"Heartbeat 5"}
data: {"time":"2025-07-25T19:13:10.237680332Z","level":"WARN","message":"Slow operation detected"}
data: {"time":"2025-07-25T19:13:11.237936559Z","level":"INFO","message":"Heartbeat 6"}
data: {"time":"2025-07-25T19:13:12.237335514Z","level":"INFO","message":"Heartbeat 7"}
data: {"time":"2025-07-25T19:13:13.237355836Z","level":"INFO","message":"Heartbeat 8"}
data: {"time":"2025-07-25T19:13:13.237360289Z","level":"ERROR","message":"Something went wrong"}
data: {"time":"2025-07-25T19:13:14.237777154Z","level":"INFO","message":"Heartbeat 9"}
data: {"time":"2025-07-25T19:13:15.237203518Z","level":"INFO","message":"Heartbeat 10"}
data: {"time":"2025-07-25T19:13:15.237209519Z","level":"INFO","message":"Shutting down logging loop"}
INFO lib: ✔️ Stopping channel-based logger...
```

## Example
### Logging to console
```golang
package main

import (
	logger "logo/logo"
    "time"
)

func main(){
    log := logger.New(nil, true) // Initialize logger with console output enabled

    log.Info("📦 Starting the CLI tool...")

    log.Debug("Initializing subsystems...")
    time.Sleep(500 * time.Millisecond)

    log.Info("Processing data stream...")
    for i := 1; i <= 5; i++ {
        log.Debug("Processing item %d", i)
        time.Sleep(250 * time.Millisecond)

        if i == 3 {
            log.Warn("Potential issue with item %d", i)
        }
    }

    log.Info("✔️ Completed all tasks")
}
```

### Logging to a channel
```golang
import (
    "encoding/json"
	"fmt"
	logger "logo/logo"
)
func main(){
    logChan := make(chan logger.LogMessage, 100) // Channel for log messages
	defer close(logChan)                         // Ensure channel is closed when done

	// Start a goroutine to listen for log messages
	// and print them to the console
	// This simulates a logging system that could be used in a larger application
	go func() {
		for {
			select {
			case msg := <-logChan:
				data, _ := json.Marshal(msg)
				fmt.Printf("data: %s\n", data)

			}
		}
	}()

    fmt.Println("Starting channel logging")
	appLog := logger.New(logChan, false) // Initialize logger with channel
	appLog.Info("Msg 1: %s", "hello")
    appLog.Info("Msg 2: %s", "world")
    fmt.Println("Stopping channel logging")
}
```

# TODO
- Update README with new example
- Cleanup .bak files