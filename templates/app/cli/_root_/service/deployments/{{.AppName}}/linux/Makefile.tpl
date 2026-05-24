include ../../configs/_make_/config/project.mk
include ../../configs/_make_/config/go-build.mk
include ../../configs/_make_/lib/go-build/cache-ram-drive.mk
include ../../configs/_make_/lib/go-build/progress-bar.mk

#======================================================
# {{.AppName}}
#======================================================
APP_NAME={{.AppName}}
APP_LOG_NAME={{.AppName}}

PROJECT_ROOT_PATH=../../..
INSTALL_PATH=/usr/local/bin

BUILD_PLATFORM=linux
BIN_PATH=bin/$(APP_NAME)/linux

APP_BINARY=$(BIN_PATH)/$(APP_NAME)
APP_INSTALL=$(INSTALL_PATH)/$(APP_NAME)

define GoBuildApp
	cd $(PROJECT_ROOT_PATH)/ && \
	go mod tidy && \
	$(GO_CACHE_ENV) \
	$(GO_TMPDIR_ENV) \
	GO111MODULE=on CGO_ENABLED=1 \
	CC=/usr/bin/gcc CXX=/usr/bin/g++ \
	GOOS=$(BUILD_PLATFORM) GOARCH=amd64 \
	go build $(1) -x -ldflags="-s -w" -a -installsuffix nocgo \
	-o "$(APP_BINARY)" \
	./src/cmd/$(APP_NAME) 2>&1
endef

#=========================================================================
#
#=========================================================================
.PHONY: help build

help: Makefile
	@echo "Выберите опцию сборки "$(BINARY_NAME)":"
	@sed -n 's/^##//p' $< | column -s ':' |  sed -e 's/^/ /'

## all: Скомпилировать приложение и установить
all: uninstall-app build-app install-app

## build-app: Скомпилировать приложение
build-app:
	cd $(PROJECT_ROOT_PATH)/ && \
	go mod tidy && \
	CGO_ENABLED=0 GOOS=$(BUILD_PLATFORM) GOARCH=amd64 go build -ldflags="-s -w" -o $(APP_BINARY) -v ./src/cmd/$(APP_NAME) && \
	chmod +x $(APP_BINARY) && \
	ls -lh $(APP_BINARY) && \
	$(APP_BINARY) --version;

## build-app-progress: Скомпилировать (sudo apt install gawk)
build-app-progress:
	$(call GoBuildProgress,GoBuildApp,$(APP_BINARY_NAME)) && \
	chmod +x $(APP_BINARY) && \
	ls -lh $(APP_BINARY) && \
	$(APP_BINARY) --version;

## install-app: Установить приложение
install-app:
	sudo mkdir -p /var/log/$(APP_LOG_NAME) && \
	sudo chmod 777 /var/log/$(APP_LOG_NAME) && \
	sudo ln -sf $(shell pwd)/$(PROJECT_ROOT_PATH)/$(APP_BINARY) $(APP_INSTALL);

## uninstall-app: Удалить приложение
uninstall-app:
	sudo rm -f $(APP_INSTALL);