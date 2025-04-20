/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package challenge

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
  repo.collection = db.Collection("challenges")

  return repo
}

func (r *Repository) GetAll(ctx context.Context) ([]Challenge, error) {
  cursor, err := r.collection.Find(ctx, bson.M{})
  if err != nil {
    return nil, err
  }
  defer cursor.Close(ctx)

  var challenges []Challenge

  if err := cursor.All(ctx, &challenges); err != nil {
    return nil, err
  }

  return challenges, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Challenge, error) {
  objId, err := bson.ObjectIDFromHex(id)
  if err != nil {
    return nil, err
  }

  challenge := new(Challenge)
  err = r.collection.FindOne(ctx, bson.M{"_id": objId}).Decode(challenge)
  if err != nil {
    return nil, err
  }

  return challenge, nil
}

func (r *Repository) Create(ctx context.Context, c *Challenge) error {
  _, err := r.collection.InsertOne(ctx, c)

  return err
}

func (r *Repository) Update(ctx context.Context, id string, update bson.M) error {
  objId, err := bson.ObjectIDFromHex(id)
  if err != nil {
    return err
  }

  _, err = r.collection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
  return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
  objId, err := bson.ObjectIDFromHex(id)
  if err != nil {
    return err
  }

  _, err = r.collection.DeleteOne(ctx, bson.M{"_id": objId})
  return err
}
