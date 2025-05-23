BIN_DIR := bin
CLI := $(BIN_DIR)/squirrel
SERVER := $(BIN_DIR)/server

.PHONY: all cli server clean

all: cli server

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

cli: $(CLI)

$(CLI): main.go | $(BIN_DIR)
	go build -o $(CLI)

server: $(SERVER)

$(SERVER): cmd/server/main.go cmd/server/withpq.go | $(BIN_DIR)
	go build -tags=integration -o $(SERVER) ./cmd/server

clean:
	rm -rf $(BIN_DIR)
