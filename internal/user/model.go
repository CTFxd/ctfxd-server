/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package user

import (
  "go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
  ID       bson.ObjectID `bson:"_id,omitempty" json:"id"`
  Email    string        `bson:"email" json:"email"`
  Password string        `bson:"password,omitempty" json:"-"`
  Role     string        `bson:"role" json:"-"`
}
