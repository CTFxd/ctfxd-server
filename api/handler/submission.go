/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package handler

import (
  "github.com/CTFxd/ctfxd-server/internal/auth"
  "github.com/CTFxd/ctfxd-server/internal/submission"
  "github.com/gin-gonic/gin"
)

func SetupSubmissionRoutes(apiGrp *gin.RouterGroup, submissionHandler *submission.Handler) {
  // protected group (all routes require auth)
  protected := apiGrp.Group("")
  protected.Use(auth.AuthMiddleware())
  {
    protected.POST("/submit", submissionHandler.Submit)
  }
}
