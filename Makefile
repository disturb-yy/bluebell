.PHONY: all build run gotool clean help  # 声明要执行的命令名
BINARY="bluebell"  # 编译生成的可执行文件的名称

all: gotool build  # 可选的要生成 targets（all） 需要的文件或者是目标(gotool build)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/${BINARY}  # 编译成linux平台的代码

run:
	@go run ./main.go conf/config.yaml  # @代表不显示执行的命令

gotool:
	go fmt ./  # 格式化代码
	go vet ./  # 语法检查

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi  # 如果有该文件，进行删除操作

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"