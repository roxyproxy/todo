package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"todo/logger"
	"todo/metrics"
	"todo/server/grpcsrv"
	"todo/server/httpsrv"
	"todo/service"
	"todo/storage/postgres"

	conf "todo/config"
)

func init() {
	prometheus.Register(metrics.TotalRequests)
	//prometheus.Register(metrics.ResponseStatus)
	//prometheus.Register(metrics.HttpDuration)
}

func main() {
	log := logger.New(os.Stderr)
	if err := godotenv.Load(); err != nil {
		log.Warning("No .env file found")
	}
	config := conf.New()

	dbpool, err := pgxpool.Connect(context.Background(), config.DBUrl)
	if err != nil {
		log.Errorf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	migrateDatabase(context.Background(), dbpool, log)

	service := service.NewService(postgres.NewPostgresStorage(dbpool), config)

	go func() {
		lis, err := net.Listen("tcp", config.GrpcPort)
		if err != nil {
			log.Errorf("failed to listen: %v", err)
		}
		log.Infof("grpc server listening at %v", lis.Addr())
		server := grpcsrv.NewGrpcServer(service, config, log)
		err = server.Serve(lis)
		if err != nil {
			log.Errorf("grpc failed to serve: %v", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-done
		log.Infof("system call:%+v", oscall)
		cancel()
	}()

	server := httpsrv.NewHTTPServer(service, config, log)
	httpServer := &http.Server{
		Addr:    config.HTTPPort,
		Handler: server.Serve,
	}

	go func() {
		err = httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("http failed to serve: %v", err)
		}
	}()
	log.Infof("http server listening at %v", config.HTTPPort)

	<-ctx.Done()

	if err = httpServer.Shutdown(context.Background()); err != nil {
		log.Infof("http server Shutdown Failed:%+s", err)
	}

	log.Infof("http server exited properly")

}

func migrateDatabase(ctx context.Context, dbpool *pgxpool.Pool, log logger.Logger) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), "schema_version")
	if err != nil {
		log.Errorf("Unable to create a migrator: %v", err)
	}

	err = migrator.LoadMigrations("./storage/postgres/migrations")
	if err != nil {
		log.Errorf("Unable to load migrations: %v", err)
	}

	err = migrator.Migrate(ctx)

	if err != nil {
		log.Errorf("Unable to migrate: %v", err)
	}

	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		log.Errorf("Unable to get current schema version: %v", err)
	}

	log.Infof("Migration done. Current schema version: %v", ver)
}

// mockgen -destination=server/mocks/mockstore.go -package=mockstore todo/storage Storage
// docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
// psql -U postgres
