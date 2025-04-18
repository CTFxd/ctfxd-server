/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package challenge

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	serv := new(Service)
	serv.repo = repo

	return serv
}

func (s *Service) ListChallenges(ctx context.Context) ([]Challenge, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetChallenge(ctx context.Context, id string) (*Challenge, error) {
	return s.repo.GetByID(ctx, id)
}
