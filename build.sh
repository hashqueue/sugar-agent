set -e
echo 'Start build astools binary executable files...'
version=$(go version)
if [[ "$version" =~ "go version" ]]; then
  echo "$version"
else
  echo "Error: please install Go first."
  exit 4
fi
GOOS=linux GOARCH=amd64 go build -o sugar-agent_amd64 cmd/main.go
GOOS=linux GOARCH=arm64 go build -o sugar-agent_arm64 cmd/main.go
echo "build done."
ls -larth ./sugar-agent*
