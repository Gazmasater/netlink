BINARY := main

all: $(BINARY)
	@true

$(BINARY): cmd/main.go
	@go build -o $(BINARY) cmd/main.go

clean:
	@rm -f $(BINARY)

run: all
	@sudo ./$(BINARY)

