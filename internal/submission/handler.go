/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package submission

import (
  "log"
  "net/http"

  "github.com/CTFxd/ctfxd-server/internal/auth"
  "github.com/gin-gonic/gin"
)

type SubmitRequest struct {
  ChallengeID string `json:"challenge_id" binding:"required"`
  Flag        string `json:"flag" binding:"required"`
}

type Handler struct {
  service *Service
}

func NewHandler(service *Service) *Handler {
  handler := new(Handler)
  handler.service = service
  return handler
}

func (h *Handler) Submit(c *gin.Context) {
  var req SubmitRequest

  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
    return
  }

  userID := auth.GetUserID(c)
  userEmail := auth.GetUserEmail(c)
  err := h.service.Submit(c.Request.Context(), userID, userEmail, req.ChallengeID, req.Flag)
  if err != nil {
    log.Printf("challenge: error(%v)\n", err)
  }

  switch err {
  case nil:
    c.JSON(http.StatusOK, gin.H{"status": "correct"})
  case ErrIncorrectFlag:
    c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect flag"})
  case ErrAlreadySolved:
    c.JSON(http.StatusConflict, gin.H{"error": "already solved"})
  default:
    c.JSON(http.StatusInternalServerError, gin.H{"error": "submission failed"})
  }
}
