/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package auth

import (
  "errors"
  "log"
  "net/http"
  "strings"

  "github.com/gin-gonic/gin"
)

const (
  ContextUserID = "user_id"
  ContextEmail  = "email"
  ContextRole   = "role"
)

func AuthMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    token, err := extractJWT(c)
    if err != nil {
      log.Printf("JWT validation failed: %v\n", err)
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
      return
    }

    claims, err := validateJWT(token)
    if err != nil {
      log.Printf("JWT validation failed: %v\n", err)
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
      return
    }

    setUserContext(c, claims)

    c.Next()
  }
}

func AdminMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    role := GetUserRole(c)
    if role == "" || role != "admin" {
      c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin privilege not found"})
      return
    }

    c.Next()
  }
}

func GetUserID(c *gin.Context) string {
  val, exists := c.Get(ContextUserID)
  if !exists {
    return ""
  }

  return val.(string)
}

func GetUserEmail(c *gin.Context) string {
  val, exists := c.Get(ContextEmail)
  if !exists {
    return ""
  }

  return val.(string)
}

func GetUserRole(c *gin.Context) string {
  val, exists := c.Get(ContextRole)
  if !exists {
    return ""
  }

  return val.(string)
}

func extractJWT(c *gin.Context) (string, error) {
  authHeader := c.GetHeader("Authorization")
  if authHeader == "" {
    return "", errors.New("Missing auth token")
  }

  headerParts := strings.Split(authHeader, " ")
  if len(headerParts) != 2 || headerParts[0] != "Bearer" {
    return "", errors.New("Invalid auth header")
  }

  return headerParts[1], nil
}

func validateJWT(token string) (*Claims, error) {
  claims, err := ParseJWT(token)
  if err != nil {
    return nil, errors.New("Invalid or expired token")
  }

  return claims, nil
}

func setUserContext(c *gin.Context, claims *Claims) {
  c.Set(ContextUserID, claims.UserID)
  c.Set(ContextEmail, claims.Email)
  c.Set(ContextRole, claims.Role)
}
