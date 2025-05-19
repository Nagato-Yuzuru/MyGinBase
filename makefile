WIRE := wire

./inject/wire_gen.go: ./inject/wire.go
	@echo "Generating wire files..."
	@$(WIRE) ./...
	@echo "Wire files generated successfully."
	@echo "Run 'go mod tidy' to clean up the go.mod file."

.PHONY: wire
wire: ./inject/wire_gen.go
	@echo "Wire target excuted."