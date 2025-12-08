include .env
export $(shell sed 's/=.*//' .env)

#======================================================
# HTTP Server
#======================================================
PROJECT=go-draft
APP_VERSION=v0.0.1
APP_NAME=go-draft

PROJECT_ROOT_PATH=../../../..
INSTALL_PATH=/usr/local/bin

BUILD_PLATFORM=linux
BIN_PATH=bin/$(APP_NAME)/linux

APP_BINARY=$(BIN_PATH)/$(APP_NAME)
APP_INSTALL=$(INSTALL_PATH)/$(APP_NAME)

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
	CGO_ENABLED=0 GOOS=$(BUILD_PLATFORM) GOARCH=amd64 go build -ldflags="-s -w" -o $(APP_BINARY) -v src/cmd/$(APP_NAME)/main.go && \
	chmod +x $(APP_BINARY);

## install-app: Установить приложение
install-app:
	sudo ln -sf $(shell pwd)/$(PROJECT_ROOT_PATH)/$(APP_BINARY) $(APP_INSTALL);

## uninstall-app: Удалить приложение
uninstall-app:
	sudo rm -f $(APP_INSTALL);
