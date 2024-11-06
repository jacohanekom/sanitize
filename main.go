package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"sanitize/controller"
	"sanitize/data"

	_ "sanitize/docs"
)

// @title           Sanitize Web Service
// @version         1.0
// @description     This is a microservice to sanitize strings based on the stored words, and to manage the stored words
// @termsOfService  http://swagger.io/terms/

// @contact.name   Jaco Hanekom
// @contact.email  jhanekom@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
//	@BasePath	/api/v1

var dbUsername = os.Getenv("dbUsername")
var dbPassword = os.Getenv("dbPassword")
var dbHost = os.Getenv("dbHost")
var dbPort = os.Getenv("dbPort")
var dbSchema = os.Getenv("dbSchema")

func main() {
	log.Println("Starting Service...")

	log.Println("Setting up database")
	db, err := data.Initialize(fmt.Sprintf("%s;%s://%s:%s@%s:%s?database=%s",
		"sqlserver", "sqlserver", dbUsername, dbPassword, dbHost, dbPort, dbSchema))
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	c := controller.NewController(&db)

	v1 := r.Group("/api/v1")
	{
		words := v1.Group("/words")
		{
			words.GET("", c.ListWords)
			words.PUT("", c.AddWords)
			words.POST("", c.UpdateWords)
			words.DELETE("", c.DeleteWords)
		}
		sanitized := v1.Group("/sanitize")
		{
			sanitized.POST("", c.Sanitize)
		}
	}

	if _, err := os.Stat("sql_sensitive_list.json"); err == nil {
		log.Printf("Moving sample data to imported")
		os.Rename("sql_sensitive_list.json", "sql_sensitive_list.imported")
	}

	log.Printf("Starting server on port %s", os.Getenv("PORT"))
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
