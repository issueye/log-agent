export CGO_ENABLED=0
export GOARCH=amd64
export GOOS=linux

go mod vendor

TAG="Beta"
BRANCH=$(git symbolic-ref --short -q HEAD)
COMMIT=$(git rev-parse --verify HEAD)
NOW=$(date '+%FT%T%z')

VERSION="v0.1.1-${TAG}"
APPNAME="logAgent-${VERSION}"
DESCRIPTION="日志采集服务"

go build -o bin/${APPNAME} -ldflags "-X demo/build.AppName=Demo \
-X github.com/issueye/log-agent/internal/initialize.Branch=${BRANCH} \
-X github.com/issueye/log-agent/internal/initialize.Commit=${COMMIT} \
-X github.com/issueye/log-agent/internal/initialize.Date=${NOW} \
-X github.com/issueye/log-agent/internal/initialize.AppName=${DESCRIPTION} \
-X github.com/issueye/log-agent/internal/initialize.Version=${VERSION}" main.go