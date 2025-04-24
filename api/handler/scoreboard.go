/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package handler

import (
  "github.com/CTFxd/ctfxd-server/internal/scoreboard"
  "github.com/gin-gonic/gin"
)

func SetupScoreboardRoutes(apiGrp *gin.RouterGroup, scoreboardHandler *scoreboard.Handler) {
  public := apiGrp.Group("")
  {
    public.GET("/scoreboard", scoreboardHandler.Get)
  }
}
