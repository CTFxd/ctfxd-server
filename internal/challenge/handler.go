/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package challenge

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

func (h *Handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	challenges, err := h.service.ListChallenges(ctx)
	if err != nil {
		log.Printf("challenge: error(%v)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch challenges"})
		return
	}

	for i := range challenges {
		challenges[i].Flag = ""
	}

	c.JSON(http.StatusOK, challenges)
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	challenge, err := h.service.GetChallenge(ctx, id)
	if err != nil {
		log.Printf("challenge: error(%v)", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "challenge not found"})
		return
	}

	challenge.Flag = ""
	c.JSON(http.StatusOK, challenge)
}
