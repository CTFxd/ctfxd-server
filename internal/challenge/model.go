/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package challenge

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Challenge struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string        `bson:"title" json:"title"`
	Category    string        `bson:"category" json:"category"`
	Description string        `bson:"description" json:"description"`
	Points      int           `bson:"points" json:"points"`
	Flag        string        `bson:"flag" json:"-"`
}
