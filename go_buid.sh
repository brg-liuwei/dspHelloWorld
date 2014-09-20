cur_path=`pwd`
export GOPATH=$cur_path:$GOPATH
#export CGO_CFLAGS="-I${cur_path}/lib/leveldb/include"
#export CGO_LDFLAGS="-L${cur_path}/lib/leveldb"

go build -o bin/dspServer src/server.go
go build -o bin/cli src/cli.go
