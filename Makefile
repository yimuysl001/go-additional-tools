BINARY="myapp"
VERSION=1.0.0
Vendor="YourCompany"
ProgramName="MyAwesomeApp"
BUILD_DATE=$(shell date +%FT%T%z)
# 构建参数
LDFLAGS=-ldflags="-X main.Version=${VERSION} -X main.Vendor=${Vendor} -X main.ProgramName=${ProgramName}  -X main.BuildDate=${BUILD_DATE} -w -s"



# 默认目标
.PHONY: all
all: build

# 构建项目
.PHONY: build
build:
    go build -o ${BINARY} ${LDFLAGS} main.go

# 运行测试
.PHONY: test
test:
    go test -v ./...

# 清理构建文件
.PHONY: clean
clean:
    if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

# 安装依赖
.PHONY: deps
deps:
    go mod download

# 代码格式化
.PHONY: fmt
fmt:
    go fmt ./...