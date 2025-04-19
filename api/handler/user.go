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

func SetupUserRoutes(apiGrp *gin.RouterGroup, userHandler *user.Handler) {

	// public routes (no auth) Probably ???
	public := apiGrp.Group("")
	{
		public.POST("/register", userHandler.RegisterUser)
		public.POST("/login", userHandler.Login)
	}

	// protected group (all routes require auth)
	protected := apiGrp.Group("")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/me", userHandler.GetMe)
	}

	admin := protected.Group("")
	admin.Use(auth.AdminMiddleware())
	{
		admin.POST("/admin/register", userHandler.RegisterAdmin)
	}
}
