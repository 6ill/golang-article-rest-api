package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/6ill/go-article-rest-api/internal/helper"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func DatabaseInit(v *viper.Viper) *sql.DB {
	dsn := fmt.Sprintf("%s?sslmode=disable", v.GetString("DB_DSN"))
	fmt.Printf("\ndsn: %s", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		helper.Logger(helper.LoggerLevelPanic, "Unable to connect database", err)
	}

	db.SetMaxOpenConns(v.GetInt("MAX_OPEN_CONNS"))
	db.SetMaxIdleConns(v.GetInt("MAX_IDLE_CONNS"))
	duration, err := time.ParseDuration(v.GetString("MAX_IDLE_TIME"))
	if err != nil {
		duration = 15 * time.Minute
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		helper.Logger(helper.LoggerLevelError, "Could not ping the database", err)
	}

	helper.Logger(helper.LoggerLevelInfo, "MySQL is connected", nil)

	return db
}
