package main

import (
	"github.com/ehickey08/GinRESTAPI/database"
	"github.com/ehickey08/GinRESTAPI/recommendations"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const basePath = "/api"

func main() {
	database.SetupDatabase()
	database.SetupMongoDatabase()
	router := gin.Default()
	recommendations.SetupRoutes(router, basePath)
	log.Fatal(router.Run())
}
