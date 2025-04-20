/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

//go:generate go run github.com/go-swagger/go-swagger/cmd/swagger@latest generate spec --scan-models -o ./docs/swagger.yml

// Package api CTFxd APIs.
//
// These CTFxd APIs do things better.
//
// swagger: '2.0'
//
// Schemes: http
// BasePath: /api/v1
// Version: 0.0.1
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
  "log"
  "net/http"

  "github.com/MarceloPetrucio/go-scalar-api-reference"
  "github.com/gin-gonic/gin"
)

func apiReferenceGen() gin.HandlerFunc {
  return func(c *gin.Context) {
    htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
      SpecURL: "./docs/swagger.yml",
      CustomOptions: scalar.CustomOptions{
        PageTitle: "CTFxd API",
      },
      DarkMode: true,
    })

    if err != nil {
      log.Printf("doc: error: %v\n", err)
      c.AbortWithError(http.StatusInternalServerError, err)
      return
    }

    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
  }
}
