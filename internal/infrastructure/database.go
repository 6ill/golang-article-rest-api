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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		helper.Logger(helper.LoggerLevelError, "Could not ping the database", err)
	}

	helper.Logger(helper.LoggerLevelInfo, "MySQL is connected", nil)

	return db
}
