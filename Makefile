OUT_DIR := output

.PHONY: build run build-linux build-mac build-windows build-all

update-dependencies:
	go mod tidy

build: clean update-dependencies
	go build -o $(OUT_DIR)/cilikube cmd/server/main.go

run: build
	./$(OUT_DIR)/cilikube

clean:
	rm -rf $(OUT_DIR)

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUT_DIR)/cilikube-linux-amd64 cmd/server/main.go

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o $(OUT_DIR)/cilikube-mac-amd64 cmd/server/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(OUT_DIR)/cilikube-windows-amd64.exe cmd/server/main.go

build-all:
	make build-linux
	make build-mac
	make build-windows