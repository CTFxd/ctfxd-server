/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package handler

import (
  "github.com/CTFxd/ctfxd-server/internal/auth"
  "github.com/CTFxd/ctfxd-server/internal/challenge"
  "github.com/gin-gonic/gin"
)

func SetupChallengeRoutes(apiGrp *gin.RouterGroup, challengeHandler *challenge.Handler) {

  public := apiGrp.Group("/challenge")
  {
    public.GET("", challengeHandler.GetChallenges)
    public.GET("/:id", challengeHandler.GetChallenge)
  }

  // protected routes (requires login)
  protected := apiGrp.Group("/challenge")
  protected.Use(auth.AuthMiddleware())
  {
    protected.GET("/:id/solves", challengeHandler.GetSolves)

    admin := protected.Group("")
    admin.Use(auth.AdminMiddleware())
    {
      admin.POST("", challengeHandler.CreateChallenge)
      admin.PATCH("/:id", challengeHandler.UpdateChallenge)
      admin.DELETE("/:id", challengeHandler.DeleteChallenge)
      admin.GET("/:id/flag", challengeHandler.GetFlag)

      admin.POST("/:id/file", challengeHandler.AddChallengeFile)
      admin.PUT("/:id/file/:uuid", challengeHandler.UpdateChallengeFile)
      admin.DELETE("/:id/file/:uuid", challengeHandler.DeleteChallengeFile)
    }
  }
}
