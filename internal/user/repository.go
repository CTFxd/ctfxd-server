/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package user

import (
  "context"

  "go.mongodb.org/mongo-driver/v2/bson"
  "go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
  collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
  repo := new(Repository)
  repo.collection = db.Collection("users")

  return repo
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
  _, err := r.collection.InsertOne(ctx, user)

  return err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
  user := new(User)

  err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(user)
  if err != nil {
    return nil, err
  }

  return user, nil
}
