/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package challenge

import (
  "encoding/json"
  "errors"
  "fmt"
  "log"
  "net/http"

  "github.com/CTFxd/ctfxd-server/internal/auth"
  "github.com/gin-gonic/gin"
)

type Handler struct {
  service *Service
}

type UpdateChallengeRequest struct {
  Title       *string `bson:"title" json:"title"`
  Category    *string `bson:"category" json:"category"`
  Description *string `bson:"description" json:"description"`
  Points      *int    `bson:"points" json:"points"`
  State       *string `bson:"state" json:"state"`
  Type        *string `bson:"type" json:"type"`
  Solves      *int    `bson:"solves" json:"solves"`
  Flag        *string `bson:"flag" json:"flag"`
  Author      *string `bson:"author" json:"author"`
}

func NewHandler(service *Service) *Handler {
  handler := new(Handler)
  handler.service = service
  return handler
}

func (h *Handler) GetChallenges(c *gin.Context) {
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

func (h *Handler) GetChallenge(c *gin.Context) {
  id := c.Param("id")

  challenge, err := h.service.GetChallenge(c.Request.Context(), id)
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

func (h *Handler) CreateChallenge(c *gin.Context) {
  if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
    log.Printf("challenge: error(%v)\n", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart-form"})
    return
  }

  var req Challenge
  if err := json.Unmarshal([]byte(c.PostForm("data")), &req); err != nil {
    log.Printf("challenge: error(%v)\n", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge data"})
    return
  }

  // set author from JWT context
  if req.Author == "" {
    req.Author = auth.GetUserEmail(c)
  }

  form, err := c.MultipartForm()
  if err != nil {
    log.Printf("challenge: error(%v)\n", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart-form"})
    return
  }

  if err := h.service.CreateChallengeWithFiles(c.Request.Context(), &req, form, c); err != nil {
    log.Printf("challenge: error(%v)\n", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create challenge"})
    return
  }

  c.Status(http.StatusCreated)
}

func (h *Handler) UpdateChallenge(c *gin.Context) {
  id := c.Param("id")
  var update UpdateChallengeRequest

  if err := c.ShouldBindJSON(&update); err != nil {
    log.Printf("challenge: error(%v)\n", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
    return
  }

  if err := h.service.UpdateChallenge(c.Request.Context(), id, &update); err != nil {
    log.Printf("challenge: error(%v)\n", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
    return
  }

  c.Status(http.StatusOK)
}

func (h *Handler) DeleteChallenge(c *gin.Context) {
  id := c.Param("id")

  if err := h.service.DeleteChallenge(c.Request.Context(), id); err != nil {
    log.Printf("challenge: error(%v)\n", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete challenge"})
    return
  }

  c.Status(http.StatusNoContent)
}

func (h *Handler) AddChallengeFile(c *gin.Context) {
  challengeID := c.Param("id")

  form, err := c.MultipartForm()
  if err != nil {
    log.Printf("challenge: error(%v)\n", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart-form"})
    return
  }

  uploadedFiles, err := h.service.AddChallengeFile(c.Request.Context(), challengeID, form, c)
  if err != nil && len(uploadedFiles) == 0 {
    log.Printf("challenge: error(%v)\n", err)
    if errors.Is(err, ErrNoFile) || errors.Is(err, ErrFileExcedLimit) {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    } else {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload challenge file"})
    }

    return
  }

  c.JSON(http.StatusCreated, uploadedFiles)
}

func (h *Handler) UpdateChallengeFile(c *gin.Context) {
  challengeID := c.Param("id")
  fileUUID := c.Param("uuid")

  form, err := c.MultipartForm()
  if err != nil {
    log.Printf("challenge: error(%v)\n", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart-form"})
    return
  }

  if err := h.service.UpdateChallengeFile(c.Request.Context(), challengeID, fileUUID, form, c); err != nil {
    log.Printf("challenge: error(%v)\n", err)
    if errors.Is(err, ErrFileNotFound) || errors.Is(err, ErrNoFile) || errors.Is(err, ErrFileExcedLimit) {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    } else {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update challenge file"})
    }
    return
  }

  c.Status(http.StatusOK)
}

func (h *Handler) DeleteChallengeFile(c *gin.Context) {
  challengeID := c.Param("id")
  fileUUID := c.Param("uuid")

  err := h.service.DeleteChallengeFile(c.Request.Context(), challengeID, fileUUID, c)
  if err != nil {
    log.Printf("challenge: error(%v)\n", err)
    if errors.Is(err, ErrFileNotFound) || errors.Is(err, ErrNoFile) || errors.Is(err, ErrFileExcedLimit) {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    } else {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete challenge file"})
    }
    return
  }

  c.Status(http.StatusNoContent)
}

func (h *Handler) DownloadChallengeFile(c *gin.Context) {
  challengeID := c.Param("id")
  fileUUID := c.Param("uuid")

  filePath, fileName, err := h.service.GetChallengeFile(c.Request.Context(), challengeID, fileUUID)
  if err != nil {
    log.Printf("challenge: error(%v)\n", err)
    if errors.Is(err, ErrFileNotFound) || errors.Is(err, ErrFileNotOnStorage) {
      c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
    } else {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to download file"})
    }
    return
  }

  c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
  c.Header("content-Type", "application/octet-stream")
  c.File(filePath)
}
