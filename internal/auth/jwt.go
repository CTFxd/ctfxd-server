/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package auth

import (
  "errors"
  "fmt"
  "time"

  "github.com/golang-jwt/jwt/v5"
)

var JwtKey []byte

type Claims struct {
  UserID string `json:"user_id"`
  Email  string `json:"email"`
  Role   string `json:"role"`
  jwt.RegisteredClaims
}

func GenerateJWT(userID, email, role string) (string, error) {
  claims := &Claims{
    UserID: userID,
    Email:  email,
    Role:   role,
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
      IssuedAt:  jwt.NewNumericDate(time.Now()),
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(JwtKey)
}

func ParseJWT(tokenStr string) (*Claims, error) {
  claims := new(Claims)

  parser := jwt.NewParser(jwt.WithLeeway(10 * time.Second))
  token, err := parser.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return JwtKey, nil
  })

  if err != nil || !token.Valid {
    return nil, errors.New("invalid JWT")
  }

  if claims.ExpiresAt == nil {
    return nil, errors.New("token must have expiration time")
  }

  return claims, nil
}
