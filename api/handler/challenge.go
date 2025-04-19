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

	// protected routes (requires login)
	protected := apiGrp.Group("")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/challenges", challengeHandler.List)
		protected.GET("/challenge/:id", challengeHandler.Get)
		protected.GET("/challenge/:id/solves", challengeHandler.GetSolves)
	}

	admin := protected.Group("")
	admin.Use(auth.AdminMiddleware())
	{
		admin.POST("/challenge", challengeHandler.Post)
		admin.PATCH("/challenge/:id", challengeHandler.Update)
		admin.DELETE("/challenge/:id", challengeHandler.Delete)

		admin.GET("/challenge/:id/flag", challengeHandler.GetFlag)
	}
}
