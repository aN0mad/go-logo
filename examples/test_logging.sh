#!/bin/bash

# Create logs directory if it doesn't exist
mkdir -p ./logs

echo "===== Running Text Logger Example ====="
go run ./examples/text_logger.go

echo -e "\n===== Running JSON Logger Example ====="
go run ./examples/json_logger.go

echo -e "\n===== Running Pretty JSON Logger Example ====="
go run ./examples/pretty_json_logger.go

echo -e "\n===== Running Channel Logger Example ====="
go run ./examples/channel_logger.go

echo -e "\n===== Running File Logger Example ====="
go run ./examples/file_logger.go

echo -e "\n===== Running Context Logger Example ====="
go run ./examples/context_logger.go

echo -e "\n===== Running Comprehensive Logger Example ====="
go run ./examples/all_features_logger.go

echo -e "\n===== All examples completed! ====="
echo "Log files can be found in ./logs/"
ls -la ./logs/