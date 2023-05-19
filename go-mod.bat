rmdir /s/q vendor

set GOPROXY=https://goproxy.cn


:: 更新依赖
go mod tidy
:: 更新本地依赖
go mod vendor