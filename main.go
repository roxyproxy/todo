package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
	"github.com/joho/godotenv"
	"net"
	"net/http"
	"os"
	"todo/logger"
	"todo/server/grpcsrv"
	"todo/server/httpsrv"
	"todo/service"
	"todo/storage/postgres"

	conf "todo/config"
)

func main() {
	log := logger.New(os.Stderr)
	if err := godotenv.Load(); err != nil {
		log.Warning("No .env file found")
	}
	config := conf.New()

	dbpool, err := pgxpool.Connect(context.Background(), config.DBUrl)
	if err != nil {
		log.Errorf("Unable to connection to database: %v\n", err)
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
		// server := grpcsrv.NewGrpcServer(inmemory.NewInMemoryStorage(), config, log)
		server := grpcsrv.NewGrpcServer(service, config, log)
		err = server.Serve(lis)
		if err != nil {
			log.Errorf("grpc failed to serve: %v", err)
		}
	}()

	// server := httpsrv.NewHTTPServer(inmemory.NewInMemoryStorage(), config, log)
	server := httpsrv.NewHTTPServer(service, config, log)
	log.Infof("http server listening at %v", config.HTTPPort)

	err = http.ListenAndServe(config.HTTPPort, server.Serve)
	if err != nil {
		log.Errorf("http failed to serve: %w", err)
	}
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
