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

func SetupChallengeRoutes(router *gin.Engine, challengeHandler *challenge.Handler) {
	apiV1 := router.Group("/api/v1")

	// protected routes (requires login)
	protected := apiV1.Group("")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/challenges", challengeHandler.List)
		protected.GET("/challenge/:id", challengeHandler.Get)
	}
}
