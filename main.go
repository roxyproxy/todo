package main

import (
	"github.com/joho/godotenv"
	"net"
	"net/http"
	"os"
	"todo/logger"
	"todo/server/grpcsrv"
	"todo/server/httpsrv"
	"todo/storage/inmemory"

	conf "todo/config"
)

const (
	grpcPort = ":5005"
	httpPort = ":5006"
)

func main() {
	log := logger.New(os.Stderr)
	if err := godotenv.Load(); err != nil {
		log.Warning("No .env file found")
	}

	go func() {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Errorf("failed to listen: %v", err)
		}
		log.Infof("grpc server listening at %v", lis.Addr())
		server := grpcsrv.NewGrpcServer(inmemory.NewInMemoryStorage(), conf.New(), log)
		err = server.Serve(lis)
		if err != nil {
			log.Errorf("grpc failed to serve: %v", err)
		}
	}()

	server := httpsrv.NewHttpServer(inmemory.NewInMemoryStorage(), conf.New(), log)
	log.Infof("http server listening at %v", httpPort)

	err := http.ListenAndServe(httpPort, server.Serve)
	if err != nil {
		log.Error("http failed to serve: %v", err)
	}

}
