GO_BUILD_ENV:=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_BUILD_FLAGS:=-x -ldflags="-s -w" -trimpath -tags netgo
