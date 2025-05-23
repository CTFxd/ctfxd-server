/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package challenge

import (
  "time"

  "go.mongodb.org/mongo-driver/v2/bson"
)

type Challenge struct {
  ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
  Title       string        `bson:"title" json:"title"`
  Category    string        `bson:"category" json:"category"`
  Description string        `bson:"description" json:"description"`
  Points      int           `bson:"points" json:"points"`
  State       string        `bson:"state" json:"state"`
  Type        string        `bson:"type" json:"type"`
  Solves      int           `bson:"solves" json:"solves"`
  Flag        string        `bson:"flag" json:"flag"`
  Author      string        `bson:"author,omitempty" json:"author,omitempty"`
  Files       []FileMeta    `bson:"files,omitempty" json:"files,omitempty"`
}

type FileMeta struct {
  UUID       string    `bson:"uuid" json:"uuid"`
  Name       string    `bson:"name" json:"name"`
  Size       int64     `bson:"size" json:"size"`
  UploadedAt time.Time `bson:"uploadedat" json:"uploadedat"`
}
