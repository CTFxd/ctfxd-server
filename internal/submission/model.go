/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package submission

import (
  "time"

  "go.mongodb.org/mongo-driver/v2/bson"
)

type Submission struct {
  ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
  UserID      bson.ObjectID `bson:"user_id,omitempty" json:"user_id"`
  Email       string        `bson:"email" json:"email"`
  ChallengeID bson.ObjectID `bson:"challenge_id" json:"challenge_id"`
  Timestamp   time.Time     `bson:"timestamp" json:"timestamp"`
}
