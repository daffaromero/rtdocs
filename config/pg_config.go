package config

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"rtdocs/utils"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

var (
	host               = utils.GetEnv("DB_HOST")
	port               = utils.GetEnv("DB_PORT")
	username           = utils.GetEnv("DB_USERNAME")
	password           = utils.GetEnv("DB_PASSWORD")
	dbName             = utils.GetEnv("DB_NAME")
	minConns           = utils.GetEnv("DB_MIN_CONNS")
	maxConns           = utils.GetEnv("DB_MAX_CONNS")
	TimeOutDuration, _ = strconv.Atoi(utils.GetEnv("DB_CONNECTION_TIMEOUT"))
)

func NewPostgresDatabase() *pgxpool.Pool {
	logger := utils.NewLogger()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, dbName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("Failed to parse config: " + dsn)
		return nil
	}

	minConnsInt, err := strconv.Atoi(minConns)
	if err != nil {
		logger.Error("DB_MIN_CONNS expected to be integer minimum connections " + minConns)
	}
	maxConnsInt, err := strconv.Atoi(maxConns)
	if err != nil {
		logger.Error("DB_MAX_CONNS expected to be integer maximum connections" + maxConns)
	}

	poolConfig.MinConns = int32(minConnsInt)
	poolConfig.MaxConns = int32(maxConnsInt)

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Error("Failed to connect " + dsn)
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping the database to verify connection
	if err := pool.Ping(c); err != nil {
		logger.Error("Database connection failed during ping: ", err)
		pool.Close()
		return nil
	}

	logger.Infow("Connected to database: ", "dsn", dsn)

	return pool
}
