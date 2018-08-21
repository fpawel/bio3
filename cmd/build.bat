set GOARCH=amd64
go build -i -ldflags="-linkmode internal -H windowsgui -X 'main.buildtime=%TIME%' -X main.prod=1" -o build/amd64/prod.exe
go build -o debug.exe -ldflags="-linkmode internal -X 'main.buildtime=%TIME%' -X main.debug=1" -o build/amd64/debug.exe