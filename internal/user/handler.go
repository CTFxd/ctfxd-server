/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package user

import (
  "errors"
  "log"
  "net/http"
  "time"

  "github.com/gin-gonic/gin"

  "github.com/CTFxd/ctfxd-server/internal/auth"
)

type Handler struct {
  service *Service
}

// swagger:model RegisterRequest
type RegisterRequest struct {
  // Email of the user
  // required: true
  // example: user@example.com
  Email string `json:"email" binding:"required,email"`

  // Password must be at least 8 characters long
  // required: true
  // example: password123
  Password string `json:"password" binding:"required,min=8"`
}

// swagger:model LoginRequest
type LoginRequest struct {
  // Email of the user
  // required: true
  // example: user@example.com
  Email string `json:"email" binding:"required,email"`

  // Password of the user
  // required: true
  // example: password123
  Password string `json:"password" binding:"required,min=8"`
}

func NewHandler(service *Service) *Handler {
  handler := new(Handler)
  handler.service = service
  return handler
}

// swagger:operation POST /register users registerUser
// ---
// tags: [users]
// description: Register a new user account
//
// parameters:
//   - name: body
//     in: body
//     description: User registration information
//     required: true
//     schema: {$ref: "#/definitions/RegisterRequest"}
//
// responses:
//
//  201:
//    description: User created successfully
//  400:
//    description: Bad request
//    schema:
//      type: object
//      properties:
//        error:
//          type: string
//          example: "Invalid email format"
//  409:
//    description: Conflict
//    schema:
//      type: object
//      properties:
//        error:
//          type: string
//          example: "user already exists"
//  500:
//    description: Internal server error
//    schema:
//      type: object
//      properties:
//        error:
//          type: string
//          example: "failed to create user"
func (h *Handler) RegisterUser(c *gin.Context) {
  h.register(c, false)
}

// swagger:operation POST /admin/register admin registerAdmin
// ---
// tags: [admin]
// description: Register a new admin account (requires admin privileges)
// security:
// - bearerAuth: []
// parameters:
//   - name: body
//     in: body
//     required: true
//     schema: {$ref: "#/definitions/RegisterRequest"}
//
// responses:
//
//  201:
//    description: Admin user created successfully
//  400:
//    description: Bad request
//    schema:
//      type: object
//      properties:
//        error:
//          type: string
//          example: "Invalid password format"
//  409:
//    description: Conflict
//    schema:
//      type: object
//      properties:
//        error:
//          type: string
//          example: "user already exists"
//  500:
//    description: Internal server error
//    schema:
//      type: object
//      properties:
//        error:
//          type: string
//          example: "failed to create admin user"
func (h *Handler) RegisterAdmin(c *gin.Context) {
  h.register(c, true)
}

func (h *Handler) register(c *gin.Context, isAdmin bool) {
  var req RegisterRequest

  if err := c.ShouldBindJSON(&req); err != nil {
    log.Printf("register:(Invalid JSON Binding) error(%v)\n", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  err := h.service.Register(c.Request.Context(), req.Email, req.Password, isAdmin)
  if err != nil {
    if errors.Is(err, ErrUserExists) {
      log.Printf("register: error: %v\n", err)
      c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
      return
    }

    log.Printf("register: error(internal): %v\n", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
    return
  }

  c.Status(http.StatusCreated)
}

// swagger:operation POST /login users login
// ---
// tags: [users]
// description: Authenticate a user
// parameters:
//   - name: body
//     in: body
//     required: true
//     schema: {$ref: "#/definitions/LoginRequest"}
//
// responses:
//
//  200:
//    description: Successfully authenticated
//    schema:
//      type: object
//      properties:
//        token:
//          type: string
//          description: JWT access token
//          example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
//        user:
//          type: object
//          description: Authenticated user details
//          properties:
//            id:
//              type: string
//              description: User ID in hex format
//              example: 5f7d8f9e0c1d2e3f4a5b6c7d
//            email:
//              type: string
//              format: email
//              example: user@example.com
//            role:
//              type: string
//              enum: [user, admin]
//              example: user
//  400:
//    description: Invalid request format
//    schema:
//      type: object
//      properties:
//        error:
//          type: string
//          example: invalid email format
//  401:
//    description: Authentication failed
//    schema:
//      type: object
//      properties:
//        error:
//          type: string
//          example: invalid credentials
//  500:
//    description: Internal server error
//    schema:
//      type: object
//      properties:
//        error:
//          type: string
//          example: failed to generate auth token
func (h *Handler) Login(c *gin.Context) {
  var req LoginRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    log.Printf("login:(Invalid JSON Binding) error(%v)\n", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  user, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
  if err != nil {
    if errors.Is(err, ErrInvalidCredentials) {
      time.Sleep(200 * time.Millisecond)
      log.Printf("login: error: %v\n", err)
      c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
    } else {
      log.Printf("login: error(internal): %v\n", err)
      c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    }
    return
  }

  token, err := auth.GenerateJWT(user.ID.Hex(), user.Email, user.Role)
  if err != nil {
    log.Printf("login: error JWT: %v\n", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate auth token"})
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "token": token,
    "user": gin.H{
      "id":    user.ID.Hex(),
      "email": user.Email,
      "role":  user.Role,
    },
  })
}

// swagger:operation GET /me users getMe
// ---
// tags: [users]
// description: Get current authenticated user details
// security:
// - bearerAuth: []
// responses:
//
//  200:
//    description: Successfully retrieved user details
//    schema:
//      type: object
//      properties:
//        id:
//          type: string
//          description: User ID in hex format
//          example: 5f7d8f9e0c1d2e3f4a5b6c7d
//        email:
//          type: string
//          format: email
//          example: user@example.com
//        role:
//          type: string
//          enum: [user, admin]
//          example: user
//  401:
//    description: Unauthorized - missing or invalid token
//    schema:
//      type: object
//      properties:
//        error:
//          type: string
//          example: unauthorized
func (h *Handler) GetMe(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
    "id":    auth.GetUserID(c),
    "email": auth.GetUserEmail(c),
    "role":  auth.GetUserRole(c),
  })
}
