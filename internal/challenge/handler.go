/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package challenge

import (
	"log"
	"net/http"

	"github.com/CTFxd/ctfxd-server/internal/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
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
		log.Printf("challenge: error(%v)\n", err)
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
		log.Printf("challenge: error(%v)\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "challenge not found"})
		return
	}

	challenge.Flag = ""
	c.JSON(http.StatusOK, challenge)
}

func (h *Handler) GetSolves(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	challenge, err := h.service.GetChallenge(ctx, id)
	if err != nil {
		log.Printf("challenge: error(%v)\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "challenge not found"})
		return
	}

	c.JSON(http.StatusOK, challenge.Solves)
}

func (h *Handler) GetFlag(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	challenge, err := h.service.GetChallenge(ctx, id)
	if err != nil {
		log.Printf("challenge: error(%v)\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "challenge not found"})
		return
	}

	c.JSON(http.StatusOK, challenge.Flag)
}

func (h *Handler) Post(c *gin.Context) {
	var req Challenge
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("challenge: error(%v)\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// set author from JWT context
	if req.Author == "" {
		req.Author = auth.GetUserEmail(c)
	}

	err := h.service.CreateChallenge(c.Request.Context(), &req)
	if err != nil {
		log.Printf("challenge: error(%v)\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create challeng"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var update bson.M

	if err := c.ShouldBindJSON(&update); err != nil {
		log.Printf("challenge: error(%v)\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if err := h.service.UpdateChallenge(c.Request.Context(), id, update); err != nil {
		log.Printf("challenge: error(%v)\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteChallenge(c.Request.Context(), id); err != nil {
		log.Printf("challenge: error(%v)\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete challenge"})
	}

	c.Status(http.StatusNoContent)
}
