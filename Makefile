WIRE := wire
WIRE_INPUTS := $(shell find . -type f -name "wire.go")
WIRE_OUTPUTS := $(WIRE_INPUTS:wire.go=wire_gen.go)


.PHONY: all
all: mod_tidy wire


# wire gen
%wire_gen.go: %/wire.go
	@echo "Running wire..."
	@$(WIRE) ./$(*D)
	@echo "Wire generation $(*D) completed successfully."

.PHONY: wire
wire: $(WIRE_OUTPUTS)
	@echo "Wire target executed."

# go mod tidy
.PHONY: mod_tidy
mod_tidy: ./go.sum
	@echo "go mode tidy target executed."

./go.sum: go.mod
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "go mod tidy completed successfully."