set -e
echo 'Start build astools binary executable files...'
version=$(go version)
if [[ "$version" =~ "go version" ]]; then
  echo "$version"
else
  echo "Error: please install Go first."
  exit 4
fi
# Go语言中的 CGO_ENABLED 环境变量是用于控制 Go 是否开启 CGO（C语言调用Go的接口）的开关。当 CGO_ENABLED 设置为 0 时，
# Go将禁用对因特尔指令集、C库、系统工具链的依赖，也就是禁用了 CGO。这时，Go只能使用纯Go代码，不能调用C语言库等外部资源。
# 当CGO_ENABLED=1，进行编译时会将文件中引用libc的库（比如常用的net包），以动态链接的方式生成目标文件。
# 当CGO_ENABLED=0，进行编译时则会把在目标文件中未定义的符号（外部函数）一起链接到可执行文件中。
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o sugar-agent_amd64 cmd/main.go
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o sugar-agent_arm64 cmd/main.go
echo "build done."
ls -larth ./sugar-agent*
