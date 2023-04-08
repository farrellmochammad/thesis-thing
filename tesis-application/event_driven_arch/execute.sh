#!/bin/bash

function cleanup {
  echo "Cleaning up child processes..."
  kill -- -$$
}

# Trap the interrupt signal
trap cleanup INT

# run ci-connector service
cd ci-connector-eda
go run main.go -port :8083 -rethink localhost:28015 &
go run main.go -port :8085 -rethink localhost:28017 &
cd ..

# run ci-connector-eda-event-router service
cd ci-connector-eda-event-router
go run main.go &
cd ..

# run bi-fast-eda-event-router service
cd bi-fast-eda-event-router
go run main.go &
cd ..

# run bi-fast-eda service
cd bi-fast-eda
go run main.go &
cd ..

# run bi-fast-hub-eda service
cd bi-fast-hub-eda
go run main.go &
cd ..

# run bi-fast-hub-eda service
cd prm-eda
go run main.go &
cd ..


# Wait for child processes to complete
wait