#!/bin/bash

# Run swag init to generate swagger docs
# echo "Generating Swagger documentation"
# swag init
# if [ $? -ne 0 ]; then
#     echo "Failed to generate Swagger documentation"
#     exit 1
# fi

# Run Go application
echo "Starting the Go application..."
go run .
if [ $? -ne 0 ]; then
    echo "Failed to start the Go application"
    exit 1
fi