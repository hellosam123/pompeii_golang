EXE := pompeii
VERSION := 0.2.0
NAME := $(EXE)-$(VERSION)
OUT_DIR := bin

PLATFORMS := linux/amd64 darwin/arm64 windows/amd64

define build_tar
	$(eval BIN_SUFFIX := $(if $(filter windows,$(1)),.exe,))
	$(eval BINARY_NAME := $(NAME)-$(1)-$(2)$(BIN_SUFFIX))	
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o $(BINARY_NAME)
	tar -czvf $(BINARY_NAME).tar.gz $(BINARY_NAME)
	rm $(BINARY_NAME)
	mv $(BINARY_NAME).tar.gz $(OUT_DIR)/
	@echo "Successfully staged $(BINARY_NAME) into $(OUT_DIR)!"
endef

build-all:
	@$(foreach platform,$(PLATFORMS),\
		$(call build_tar,$(word 1,$(subst /, ,$(platform))),$(word 2,$(subst /, ,$(platform)))))

