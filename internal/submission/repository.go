/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package submission

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
  repo.collection = db.Collection("submissions")

  return repo
}

func (r *Repository) Create(ctx context.Context, s *Submission) error {
  _, err := r.collection.InsertOne(ctx, s)

  return err
}

func (r *Repository) HasSolved(ctx context.Context, email string, challengeID string) (bool, error) {
  objId, err := bson.ObjectIDFromHex(challengeID)
  if err != nil {
    return false, err
  }

  count, err := r.collection.CountDocuments(ctx, bson.M{"email": email, "challenge_id": objId})
  if err != nil {
    return false, err
  }

  return count > 0, nil
}

func (r *Repository) AggregateSubmission(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error) {
  cursor, err := r.collection.Aggregate(ctx, pipeline)
  if err != nil {
    return nil, err
  }

  var raw []bson.M
  if err := cursor.All(ctx, &raw); err != nil {
    return nil, err
  }

  return raw, nil
}
