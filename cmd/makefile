BINARY := main

all: $(BINARY)
	@true

$(BINARY): main.go
	@go build -o $(BINARY) main.go

clean:
	@rm -f $(BINARY)

run: all
	@sudo ./$(BINARY)

