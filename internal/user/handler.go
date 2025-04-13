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

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func NewHandler(service *Service) *Handler {
	handler := new(Handler)
	handler.service = service
	return handler
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			time.Sleep(200 * time.Millisecond)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		} else {
			log.Printf("login error (internal): %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}

	token, err := auth.GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		log.Printf("error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate auth token"})
		return
	}

	// TODO: c.SetCookie("token", token, 24*60*60, "/", "", true, true)
	c.SetCookie("token", token, 24*60*60, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) GetMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":    auth.GetUserID(c),
		"email": auth.GetUserEmail(c),
	})
}
