package database

import (
	"couchbase/model"
	"fmt"
	"log"
	"os"

	"github.com/couchbase/gocb/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	Cluster *gocb.Cluster
	Err     error
	DB      *gorm.DB
)

func CouchbaseConnect() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("HOST")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	bucket := os.Getenv("BUCKET")

	Cluster, Err = gocb.Connect(fmt.Sprintf("couchbase://%s", host), gocb.ClusterOptions{
		Username: username,
		Password: password,
	})
	if Err != nil {
		fmt.Println("Error connecting to Couchbase:", Err)
		return
	}
	Cluster.Bucket(bucket)

}

func Migration() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbname := os.Getenv("SQLDBNAME")
	port := os.Getenv("SQLPORT")
	username := os.Getenv("SQLDBUSERNAME")
	password := os.Getenv("SQLDBPASSWORD")
	host := os.Getenv("SQLDBHOST")

	DSN := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", username, password, host, port, dbname)

	fmt.Println(DSN)

	DB, Err = gorm.Open(sqlserver.Open(DSN), &gorm.Config{})
	if Err != nil {
		fmt.Println("not connected in the database : ", Err)
	} else {
		fmt.Println("connected to the database")
	}
	DB.AutoMigrate(&model.Products{}, &model.Header{}, &model.Details{})
}
