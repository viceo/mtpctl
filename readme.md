## Build

Here is an example on how to build this project
```bash
export CGO_ENABLED=0 
export GOOS="linux" 
export GOARCH="amd64"
go build -trimpath -o .build/tplibcmd
```

