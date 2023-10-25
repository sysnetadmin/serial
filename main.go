package main

import (
	"couchbase/database"
	"couchbase/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

/////////////////////////////encoded token---------------------------------------------------------------------------
/////////////////////////////decoded token---------------------------------------------------------------------------

func main() {

	database.CouchbaseConnect()
	database.Migration()
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./static")

	app.Use(cors.New())
	routes.Routes(app)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
