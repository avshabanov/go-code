
.PHONY: all
all: vendor

vendor: dependencies

# Installs dependencies
.PHONY: dependencies
dependencies:
	@echo "Installing Glide and locked dependencies..."
	glide --version || go get -u -f github.com/Masterminds/glide
	glide install

# Cleans up produced artifacts
.PHONY: clean
clean:
	rm -rf .bin

# Cleans up dependencies and produced artifacts
.PHONY: purge
purge:
	make clean
	rm -rf vendor
