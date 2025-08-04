CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/server service/main.go
chmod 777 build/server

commandName=gdp
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/gdp command/main.go
chmod 777 build/gdp
