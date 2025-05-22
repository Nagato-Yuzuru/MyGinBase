WIRE := wire
WIRE_INPUTS := $(shell find . -type f -name "wire.go")
WIRE_OUTPUTS := $(WIRE_INPUTS:wire.go=wire_gen.go)

STRINGER_INPUTS := $(shell grep -rl '//go:generate stringer' . --include='*.go')

# Calculate all expected stringer output files
# Example: if foo.go has '//go:generate stringer -type=MyType', output is mytype_string.go in foo.go's directory
STRINGER_OUTPUTS := $(foreach file,$(STRINGER_INPUTS), \
    $(addsuffix _string.go, \
        $(addprefix $(dir $(file)), \
            $(shell grep '//go:generate stringer' $(file) | sed -E 's/.*-type=([a-zA-Z0-9_]+).*/\1/' | tr '[:upper:]' '[:lower:]') \
        ) \
    ) \
)

.PHONY: all
all: mod_tidy wire stringer

# wire gen
# This rule is already incremental
%wire_gen.go: %/wire.go
	@echo "Running wire for $(dir $<)..."
	@$(WIRE) ./$(*D)
	@echo "Wire generation for $(dir $<) completed successfully."

.PHONY: wire
wire: $(WIRE_OUTPUTS)
	@echo "Wire targets are up to date."

# Define a template for stringer generation rules
# $(1) is the output file (e.g., mytype_string.go)
# $(2) is the input file (e.g., types.go)
define GENERATE_STRINGER_RULE
$(1): $(2)
	@echo "Running go generate for $(2) (to create $(1))"
	@go generate $(2)
endef

# Dynamically generate rules for each stringer output file
# For each input file that contains '//go:generate stringer'
$(foreach cur_input_file,$(STRINGER_INPUTS), \
    $(eval _generated_type_names_lower := $(shell grep '//go:generate stringer' $(cur_input_file) | sed -E 's/.*-type=([a-zA-Z0-9_]+).*/\1/' | tr '[:upper:]' '[:lower:]')) \
    $(foreach _type_name, $(_generated_type_names_lower), \
        $(eval _cur_output_file := $(addsuffix _string.go, $(addprefix $(dir $(cur_input_file)), $(_type_name)))) \
        $(eval $(call GENERATE_STRINGER_RULE,$(_cur_output_file),$(cur_input_file))) \
    ) \
)

.PHONY: stringer
stringer: $(STRINGER_OUTPUTS) # Depends on all stringer output files
	@echo "Stringer files are up to date."


# go mod tidy
.PHONY: mod_tidy
mod_tidy: ./go.sum
	@echo "go mod tidy target executed."

./go.sum: go.mod
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "go mod tidy completed successfully."