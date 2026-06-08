.PHONY: build build-all test clean install lint fmt

BINARY_NAME=smartmodelpicker
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS=-ldflags "-X github.com/user/smartmodelpicker-cli/cmd.version=${VERSION} \
	-X github.com/user/smartmodelpicker-cli/cmd.buildTime=${BUILD_TIME} \
	-X github.com/user/smartmodelpicker-cli/cmd.gitCommit=${GIT_COMMIT}"

# 默认构建当前平台
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} .

# 构建所有平台
build-all:
	@echo "🚀 构建全平台二进制文件..."
	mkdir -p dist
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 .
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 .
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 .
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 .
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe .
	
	@echo "✅ 构建完成！"
	@ls -lh dist/

# 运行测试
test:
	go test -v ./...

# 代码格式化
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run ./...

# 清理构建产物
clean:
	rm -f ${BINARY_NAME}
	rm -rf dist/

# 本地安装
install: build
	cp ${BINARY_NAME} ${GOPATH}/bin/ 2>/dev/null || cp ${BINARY_NAME} ~/go/bin/ 2>/dev/null || echo "请手动将 ${BINARY_NAME} 添加到 PATH"

# 开发运行
run: build
	./${BINARY_NAME}

# 显示版本
version: build
	./${BINARY_NAME} version
