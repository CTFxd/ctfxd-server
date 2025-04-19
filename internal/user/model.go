/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package user

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// swagger:model User
type User struct {
	// User ID in hex format
	// example: 5f7d8f9e0c1d2e3f4a5b6c7d
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	// User's email address
	// example: user@example.com
	Email string `bson:"email" json:"email"`

	// Password (not exposed in API responses)
	Password string `bson:"password,omitempty" json:"-"`

	// User role (user/admin)
	Role string `bson:"role" json:"-"`
}
