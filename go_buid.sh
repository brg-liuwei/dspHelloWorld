cur_path=`pwd`
export GOPATH=$cur_path:$GOPATH

go build -o bin/dspServer src/server.go
go build -o bin/cli src/cli.go
