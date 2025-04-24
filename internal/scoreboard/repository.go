/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package scoreboard

import (
  "context"

  "github.com/CTFxd/ctfxd-server/internal/submission"
  "go.mongodb.org/mongo-driver/v2/bson"
  "go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
  submisRepo *submission.Repository
}

func NewRepository(submisRepo *submission.Repository) *Repository {
  repo := new(Repository)
  repo.submisRepo = submisRepo

  return repo
}

func (r *Repository) GetScoreboard(ctx context.Context) ([]Score, error) {
  pipeline := mongo.Pipeline{
    bson.D{{Key: "$lookup", Value: bson.M{
      "from":         "challenges",
      "localField":   "challenge_id",
      "foreignField": "_id",
      "as":           "challenge",
    }}},
    bson.D{{Key: "$unwind", Value: "$challenge"}},

    // group by user_id (to get the aggregated scores)
    bson.D{{Key: "$group", Value: bson.M{
      "_id":        "$user_id",
      "score":      bson.M{"$sum": "$challenge.points"},
      "last_solve": bson.M{"$max": "$timestamp"},
    }}},

    // join users to get the email
    bson.D{{Key: "$lookup", Value: bson.M{
      "from":         "users",
      "localField":   "_id", // _id -> user_id from group (Should Work ?)
      "foreignField": "_id",
      "as":           "user",
    }}},
    bson.D{{Key: "$unwind", Value: "$user"}},

    // projection
    bson.D{{Key: "$project", Value: bson.M{
      "_id":        0,
      "user_id":    "$_id",
      "score":      1,
      "last_solve": 1,
      "email":      "$user.email",
    }}},

    // sort the final result
    bson.D{{Key: "$sort", Value: bson.D{
      {Key: "score", Value: -1},
      {Key: "last_solve", Value: 1},
    }}},
  }

  raw, err := r.submisRepo.AggregateSubmission(ctx, pipeline)
  if err != nil {
    return nil, err
  }

  var scores []Score
  for _, doc := range raw {
    var s Score
    bsonBytes, _ := bson.Marshal(doc)
    err = bson.Unmarshal(bsonBytes, &s)
    if err == nil {
      scores = append(scores, s)
    }
  }

  return scores, nil
}
