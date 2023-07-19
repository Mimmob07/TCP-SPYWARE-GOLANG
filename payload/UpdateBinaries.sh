rm windows/payload-amd64.exe
rm darwin/payload-amd64-darwin
rm linux/payload-amd64-linux
GOOS=windows GOARCH=amd64 go build -o windows/payload-amd64.exe main.go
GOOS=darwin GOARCH=amd64 go build -o darwin/payload-amd64-darwin main.go
GOOS=linux GOARCH=amd64 go build -o linux/payload-amd64-linux main.go