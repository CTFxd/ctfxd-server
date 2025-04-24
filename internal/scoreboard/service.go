/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package scoreboard

import (
  "context"
  "sync"
  "time"

  "github.com/CTFxd/ctfxd-server/internal/submission"
)

var scoreBoardCache = struct {
  ScoreBoard  []Score
  LastUpdated time.Time
  Mtx         sync.RWMutex
}{}

type scoreBoardWMeta struct {
}

type Service struct {
  repo *Repository
}

func NewService(repo *Repository) *Service {
  serv := new(Service)

  serv.repo = repo
  return serv
}

func (s *Service) GetScoreboard(ctx context.Context) ([]Score, error) {
  submission.LastSuccessSubmission.Mtx.RLock()
  isBefore := scoreBoardCache.LastUpdated.Equal(time.Time{}) || scoreBoardCache.LastUpdated.Before(submission.LastSuccessSubmission.SubmissionTime)
  submission.LastSuccessSubmission.Mtx.RUnlock()

  if !isBefore {
    scoreBoardCache.Mtx.RLock()
    defer scoreBoardCache.Mtx.RUnlock()

    return scoreBoardCache.ScoreBoard, nil
  }

  scores, err := s.repo.GetScoreboard(ctx)
  if err != nil {
    return nil, err
  }

  scoreBoardCache.Mtx.Lock()
  defer scoreBoardCache.Mtx.Unlock()
  scoreBoardCache.LastUpdated = time.Now().UTC()
  scoreBoardCache.ScoreBoard = scores

  return scores, nil
}
