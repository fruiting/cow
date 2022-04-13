package internal

//go:generate mockgen -destination=mock/tarantool_client.go -package=mock github.com/tarantool/go-tarantool Connector
//go:generate mockgen -destination=mock/http_server.go -package=mock net/http ResponseWriter
