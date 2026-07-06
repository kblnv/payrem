PAYR_CORE_SRC = ./cmd/payr/payr.go
PAYR_PLUGINS_SRC = ./plugins

PAYR_CORE_OUTPUT ?= ./build/payr
PAYR_PLUGINS_OUTPUT ?= ./build/plugins

.PHONY: core plugins run clean fmt

core:
	go build -o $(PAYR_CORE_OUTPUT) $(PAYR_CORE_SRC)

plugins:
	@for dir in $(PAYR_PLUGINS_SRC)/*/ ; do \
		name=$$(basename $$dir); \
		mkdir -p $(PAYR_PLUGINS_OUTPUT)/$$name; \
		go build -buildmode=plugin -o $(PAYR_PLUGINS_OUTPUT)/$$name/$$name.so $$dir; \
	done

clean:
	rm ./build/*

fmt:
	go fmt payr/...
