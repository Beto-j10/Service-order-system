package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Conn export
var Conn *gorm.DB

// Init connect to databe
func Init() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432"
	var err error
	Conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	initData()
}
