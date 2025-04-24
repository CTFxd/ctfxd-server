/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package scoreboard

import (
  "time"

  "go.mongodb.org/mongo-driver/v2/bson"
)

type Score struct {
  UserID    bson.ObjectID `bson:"user_id" json:"user_id"`
  Email     string        `bson:"email" json:"email"`
  Score     int           `bson:"score" json:"score"`
  LastSolve time.Time     `bson:"last_solve" json:"last_solve"`
}
