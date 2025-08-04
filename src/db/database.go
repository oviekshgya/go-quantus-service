package db

import (
	"context"
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/viper"
	"go-quantus-service/src/config"
	"go-quantus-service/src/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var db = newPostgresDB()

func GetDB() *gorm.DB {
	return db.Session(&gorm.Session{})
}

func newPostgresDB() *gorm.DB {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
	)

	var err error
	connDB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatal("failed to connect to PostgreSQL:", err)
	}

	if err := connDB.AutoMigrate(entities.User{}); err != nil {
		log.Fatal("failed to auto migrate user:", err)
	}

	sqlDB, err := connDB.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan instance *sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(viper.GetInt("DB_MAX_IDLE_CONNECTIONS"))
	sqlDB.SetMaxOpenConns(viper.GetInt("DB_MAX_OPEN_CONNECTIONS"))
	sqlDB.SetConnMaxLifetime(time.Duration(viper.GetInt("DB_MAX_LIFE_TIME")) * time.Second)

	connDB.Session(&gorm.Session{
		AllowGlobalUpdate:    true,
		FullSaveAssociations: false,
	})

	nrTx := config.NRc.StartTransaction("DB Operation")
	nrCtx := newrelic.NewContext(context.Background(), nrTx)
	fmt.Println("PostgreSQL DB connected")
	return connDB.WithContext(nrCtx)
}
