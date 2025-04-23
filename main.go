/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package main

import (
  "context"
  "errors"
  "fmt"
  "log"
  "net/http"
  "os"
  "os/signal"
  "regexp"
  "strconv"
  "sync"
  "time"

  "github.com/CTFxd/ctfxd-server/api/handler"
  "github.com/CTFxd/ctfxd-server/internal/auth"
  "github.com/CTFxd/ctfxd-server/internal/challenge"
  "github.com/CTFxd/ctfxd-server/internal/user"
  "github.com/CTFxd/ctfxd-server/pkg/db"
  "github.com/gin-gonic/gin"
  "github.com/joho/godotenv"
)

const (
  DEFAULT_DB_NAME        string = "ctdxd"
  DEFAULT_SERV_PORT             = "8080"
  DEFAULT_ROUTINE_PERIOD        = "30s"
)

type ServerConfig struct {
  mongodbUri     string
  secretPhrase   []byte
  dbName         string
  superuserEmail string
  superuserPass  string
  host           string
  port           string
  routinePeriod  time.Duration
  trustedProxies []string
}

func main() {
  serverConfigs, err := loadServerConfigs()
  if err != nil {
    log.Fatalln(err)
  }

  auth.JwtKey = serverConfigs.secretPhrase

  mongoClient := db.NewMongodbInit(serverConfigs.mongodbUri, serverConfigs.dbName)
  defer mongoClient.Close()

  userRepo := user.NewRepository(mongoClient.Database)
  userService := user.NewService(userRepo)
  userHandler := user.NewHandler(userService)

  status := createSuperUser(userService, serverConfigs.superuserEmail, serverConfigs.superuserPass)
  if status == false {
    log.Fatalf("failed to create superuser(id:%s password: %s)\n", serverConfigs.superuserEmail, serverConfigs.superuserPass)
  }

  challengeRepo := challenge.NewRepository(mongoClient.Database)
  challengeService := challenge.NewService(challengeRepo)
  challengeHandler := challenge.NewHandler(challengeService)

  router := gin.Default()
  router.SetTrustedProxies(nil)

  apiV1 := router.Group("/api/v1")

  // ping route
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "Welcome to CTFxd Server!",
    })
  })

  apiV1.GET("/reference", apiReferenceGen())

  handler.SetupUserRoutes(apiV1, userHandler)
  handler.SetupChallengeRoutes(apiV1, challengeHandler)

  srv := &http.Server{
    Addr:    fmt.Sprintf("%s:%s", serverConfigs.host, serverConfigs.port),
    Handler: router,
  }

  log.Println("DEBUG DEBUG", srv.Addr)

  errChan := make(chan error, 1)
  quit := make(chan os.Signal, 1)

  signal.Notify(quit, os.Interrupt)

  var wg sync.WaitGroup
  wg.Add(2)

  go func() {
    defer wg.Done()
    log.Printf("Starting the server at ':%s'", serverConfigs.port)
    if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
      log.Printf("Server: listen error(%v)\n", err)
      errChan <- err
    }
  }()

  cleanerCtx, cleanerCancel := context.WithCancel(context.Background())
  defer cleanerCancel()

  go func() {
    defer wg.Done()
    cleanOrphanFileUploadsRoutine(challengeService, cleanerCtx, serverConfigs.routinePeriod)
  }()

  select {
  case err := <-errChan:
    log.Printf("Server failed to start: %v", err)
    cleanerCancel()
    wg.Wait()
    os.Exit(1)
  case sig := <-quit:
    log.Printf("Received shutdown signal(%v)\n", sig)
    cleanerCancel()
    log.Println("Shutting down the server...")
  }

  serverClosingCtx, serverClosingCancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer serverClosingCancel()

  if err := srv.Shutdown(serverClosingCtx); err != nil {
    log.Fatal("Server forced to shutdown:", err)
  }

  wg.Wait()
}

func loadServerConfigs() (*ServerConfig, error) {
  re := regexp.MustCompile(`^(\d+)([hms]{1})$`)

  err := godotenv.Load()
  if err != nil {
    return nil, errors.New("error loading .env file")
  }

  serverConfig := new(ServerConfig)

  // check for MONGODB_URI
  mongodbUri, ok := os.LookupEnv("MONGODB_URI")
  if !ok || mongodbUri == "" {
    return nil, errors.New("error \"MONGODB_URI\" not found")
  }
  serverConfig.mongodbUri = mongodbUri

  // check for SECRET_PHRASE
  secretPhrase, ok := os.LookupEnv("SECRET_PHRASE")
  if !ok || secretPhrase == "" {
    return nil, errors.New("error \"SECRET_PHRASE\" not found")
  }
  serverConfig.secretPhrase = []byte(secretPhrase)

  // check for DB_NAME (Database Name)
  dbName, ok := os.LookupEnv("DB_NAME")
  if !ok || dbName == "" {
    dbName = DEFAULT_DB_NAME
    log.Printf("warning: DB_NAME not found! using default(%s)\n", dbName)
  }
  serverConfig.dbName = dbName

  // check for SUPERUSER_EMAIL (Superuser Email)
  superUserId, ok := os.LookupEnv("SUPERUSER_EMAIL")
  if !ok || superUserId == "" {
    return nil, errors.New("error: SUPERUSER_EMAIL not found!")
  }
  serverConfig.superuserEmail = superUserId

  // check for SUPERUSER_PASS (Superuser Password)
  superUserPass, ok := os.LookupEnv("SUPERUSER_PASS")
  if !ok || superUserPass == "" {
    return nil, errors.New("error: SUPERUSER_PASS not found!")
  }
  serverConfig.superuserPass = superUserPass

  // check for HOST name
  host, ok := os.LookupEnv("HOST")
  if !ok {
    log.Println("warning: PORT not found!")
  }
  serverConfig.host = host

  // check for PORT number
  port, ok := os.LookupEnv("PORT")
  if !ok || port == "" {
    port = DEFAULT_SERV_PORT
    log.Printf("warning: PORT not found! using default(%s)\n", port)
  }
  serverConfig.port = port

  // check for ROUTINE_PERIOD (the routine interval period)
  timePeriod, ok := os.LookupEnv("ROUTINE_PERIOD")
  if !ok || timePeriod == "" {
    timePeriod = DEFAULT_ROUTINE_PERIOD
    log.Printf("warning: ROUTINE_PERIOD not found! using default(%s)\n", timePeriod)
  }

  timePeriodMatch := re.FindStringSubmatch(timePeriod)
  if timePeriodMatch == nil {
    return nil, errors.New("error: invalid ROUTINE_PERIOD format!")
  }

  tick, err := strconv.ParseUint(timePeriodMatch[1], 10, 64)
  if err != nil {
    return nil, errors.New("error: invalid time period format")
  }

  serverConfig.routinePeriod = time.Duration(tick)

  switch timePeriodMatch[2] {
  case "h":
    serverConfig.routinePeriod *= time.Hour
  case "m":
    serverConfig.routinePeriod *= time.Minute
  case "s":
    serverConfig.routinePeriod *= time.Second
  }

  fmt.Printf("time period: %v", serverConfig.routinePeriod)

  return serverConfig, nil
}

func createSuperUser(userService *user.Service, email, password string) bool {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  err := userService.Register(ctx, email, password, true)
  if err == user.ErrUserExists {
    return true
  }

  if err != nil {
    return false
  }

  return true
}

func cleanOrphanFileUploadsRoutine(service *challenge.Service, ctx context.Context, period time.Duration) {
  ticker := time.NewTicker(period)
  defer ticker.Stop()

  for {
    select {
    case <-ctx.Done():
      log.Println("Stopping file cleaner...")
      return
    case <-ticker.C:
      serv_ctx, cancel := context.WithTimeout(context.Background(), period)
      defer cancel()

      log.Println("Running file cleaner...")
      fileNames, err := service.CleanOrphanFileUploads(serv_ctx)
      if err != nil {
        log.Printf("Error: file clearner: %v\n", err)
      } else if len(fileNames) > 0 {
        log.Printf("Cleaned files successfully: %v\n", fileNames)
      } else {
        log.Println("No files to be cleaned")
      }
    }
  }
}
