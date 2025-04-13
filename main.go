/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/CTFxd/ctfxd-server/api/handler"
	"github.com/CTFxd/ctfxd-server/internal/user"
	"github.com/CTFxd/ctfxd-server/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok || dbName == "" {
		dbName = "ctdxd"
		log.Printf("error: DB_NAME not found! using default(%s)\n", dbName)
	}

	mongoClient := db.NewMongodbInit(os.Getenv("MONGODB_URI"), dbName)
	defer mongoClient.Close()

	userRepo := user.NewRepository(mongoClient.Database)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	router := gin.Default()

	// ping route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	handler.SetupUserRoutes(router, userHandler)

	router.Run()
}
