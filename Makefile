APP_NAME := push
OUTPUT_DIR := build
CGO_ENABLED := 0

PLATFORMS := linux/amd64 linux/386 darwin/amd64 darwin/arm64 windows/amd64 windows/386 freebsd/amd64 linux/arm linux/arm64

all: $(PLATFORMS)

$(PLATFORMS):
	@mkdir -p $(OUTPUT_DIR)
	GOOS=$(firstword $(subst /, ,$@)) GOARCH=$(lastword $(subst /, ,$@)) CGO_ENABLED=$(CGO_ENABLED) go build -o $(OUTPUT_DIR)/$(APP_NAME)-$(firstword $(subst /, ,$@))-$(lastword $(subst /, ,$@)) .
	@echo "Built for $@"

clean:
	rm -rf $(OUTPUT_DIR)
	@echo "Cleaned up the build directory"

.PHONY: $(PLATFORMS) all clean
