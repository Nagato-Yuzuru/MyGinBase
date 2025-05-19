WIRE := wire

.PHONY: all
all: mod_tidy wire


# wire gen
./inject/wire_gen.go: ./inject/wire.go
	@echo "Generating wire files..."
	@$(WIRE) ./...
	@echo "Wire files generated successfully."

.PHONY: wire
wire: ./inject/wire_gen.go
	@echo "Wire target executed."

# go mod tidy
.PHONY: mod_tidy
mod_tidy: ./go.sum
	@echo "go mode tidy target executed."

./go.sum: go.mod
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "go mod tidy completed successfully."