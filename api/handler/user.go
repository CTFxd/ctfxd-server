/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package handler

import (
	"github.com/CTFxd/ctfxd-server/internal/auth"
	"github.com/CTFxd/ctfxd-server/internal/user"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userHandler *user.Handler) {
	apiV1 := router.Group("/api/v1")

	// public routes (no auth) Probably ???
	public := apiV1.Group("")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}

	// protected group (all routes require auth)
	protected := apiV1.Group("")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/me", userHandler.GetMe)
	}
}
