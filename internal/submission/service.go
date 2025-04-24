/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package submission

import (
  "context"
  "errors"
  "sync"
  "time"

  "github.com/CTFxd/ctfxd-server/internal/challenge"
  "go.mongodb.org/mongo-driver/v2/bson"
)

var (
  ErrAlreadySolved = errors.New("already solved")
  ErrIncorrectFlag = errors.New("incorrect flag")
)

var LastSuccessSubmission = struct {
  SubmissionTime time.Time
  Mtx            sync.RWMutex
}{}

type Service struct {
  repo          *Repository
  challengeServ *challenge.Service
}

func NewService(repo *Repository, challengeServ *challenge.Service) *Service {
  serv := new(Service)

  serv.repo = repo
  serv.challengeServ = challengeServ

  return serv
}

func (s *Service) Submit(ctx context.Context, userID, email, challengeID, submittedFlag string) error {
  challenge, err := s.challengeServ.GetChallenge(ctx, challengeID)
  if err != nil {
    return err
  }

  if challenge.Flag != submittedFlag {
    return ErrIncorrectFlag
  }

  solved, err := s.repo.HasSolved(ctx, email, challengeID)
  if err != nil {
    return err
  }
  if solved {
    return ErrAlreadySolved
  }

  userObjID, err := bson.ObjectIDFromHex(userID)
  if err != nil {
    return err
  }

  chalObjID, err := bson.ObjectIDFromHex(challengeID)
  if err != nil {
    return err
  }

  sub := &Submission{
    UserID:      userObjID,
    Email:       email,
    ChallengeID: chalObjID,
    Timestamp:   time.Now().UTC(),
  }

  err = s.repo.Create(ctx, sub)
  if err != nil {
    return err
  }

  LastSuccessSubmission.Mtx.Lock()
  defer LastSuccessSubmission.Mtx.Unlock()
  LastSuccessSubmission.SubmissionTime = sub.Timestamp

  return nil
}
