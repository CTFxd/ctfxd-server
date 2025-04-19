/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid eamil or password")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	serv := new(Service)
	serv.repo = repo

	return serv
}

func (s *Service) Register(ctx context.Context, email, password string, isAdmin bool) error {
	_, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return ErrUserExists
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Email:    email,
		Password: string(hashed),
		Role:     "user",
	}

	if isAdmin == true {
		user.Role = "admin"
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}

		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
