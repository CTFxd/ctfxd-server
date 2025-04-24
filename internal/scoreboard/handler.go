/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package scoreboard

import (
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
)

type Handler struct {
  service *Service
}

func NewHandler(service *Service) *Handler {
  handler := new(Handler)
  handler.service = service
  return handler
}

func (h *Handler) Get(c *gin.Context) {
  scores, err := h.service.GetScoreboard(c.Request.Context())
  if err != nil {
    log.Printf("scoreboard: error(%v)\n", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch scoreboard"})
    return
  }

  c.JSON(http.StatusOK, scores)
}
