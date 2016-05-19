package main

import (
	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/lszanto/links/config"
	"github.com/lszanto/links/controllers"
	"github.com/lszanto/links/middleware"
)

func main() {
	// set error holder
	var err error

	// CONFIG

	// open a config file
	configFile, _ := os.Open("app/config.json")

	// decode config
	decoder := json.NewDecoder(configFile)

	// create config object
	config := config.Config{}

	// decode
	err = decoder.Decode(&config)

	// GORM DATABASE

	// create db connection
	db, err := gorm.Open(config.DatabaseEngine, config.DatabaseString)

	if err != nil {
		panic("failed to connect to database")
	}

	// controllers
	lc := controllers.NewLinkController(db, config)
	uc := controllers.NewUserController(db, config)

	// ROUTER

	// setup router
	router := gin.Default()

	// add routes
	router.POST("/login", uc.Login)
	router.POST("/link", middleware.JWTVerify(config.SecretKey), lc.Post)
	router.GET("/link/:id", lc.Get)

	// SET STATIC DIR, START SERVER

	// setup static folder
	router.Static("/assets", "./assets")

	// run server
	router.Run()
}
