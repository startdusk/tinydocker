CMD=go
APP_NAME=tinydocker

.PHONY: all
all: clean build 

.PHONY: build
build:
	$(CMD) build

.PHONY: clean
clean: 
	rm -rf $(APP_NAME)
