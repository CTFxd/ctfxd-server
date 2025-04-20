/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package main

import (
  "context"
  "log"
  "net/http"
  "os"

  "github.com/CTFxd/ctfxd-server/api/handler"
  "github.com/CTFxd/ctfxd-server/internal/challenge"
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

  superUserId, ok := os.LookupEnv("SUPERUSER_EMAIL")
  if !ok || superUserId == "" {
    log.Fatalln("error: SUPERUSER_EMAIL not found!")
  }

  superUserPass, ok := os.LookupEnv("SUPERUSER_PASS")
  if !ok || superUserPass == "" {
    log.Fatalln("error: SUPERUSER_PASS not found!")
  }

  userRepo := user.NewRepository(mongoClient.Database)
  userService := user.NewService(userRepo)
  userHandler := user.NewHandler(userService)

  status := createSuperUser(userService, superUserId, superUserPass)
  if status == false {
    log.Fatalf("failed to create superuser(id:%s password: %s)\n", superUserId, superUserPass)
  }

  challengeRepo := challenge.NewRepository(mongoClient.Database)
  challengeService := challenge.NewService(challengeRepo)
  challengeHandler := challenge.NewHandler(challengeService)

  router := gin.Default()
  apiV1 := router.Group("/api/v1")

  // ping route
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

  apiV1.GET("/reference", apiReferenceGen())

  handler.SetupUserRoutes(apiV1, userHandler)
  handler.SetupChallengeRoutes(apiV1, challengeHandler)

  router.Run()
}

func createSuperUser(userService *user.Service, email, password string) bool {
  err := userService.Register(context.TODO(), email, password, true)
  if err == user.ErrUserExists {
    return true
  }

  if err != nil {
    return false
  }

  return true
}
