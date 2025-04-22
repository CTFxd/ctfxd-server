/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package challenge

import (
  "context"
  "errors"
  "mime/multipart"
  "os"
  "path/filepath"

  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/v2/bson"
)

var (
  ErrFileNotFound = errors.New("file not found")
)

type Service struct {
  repo        *Repository
  fileService *FileService
}

func NewService(repo *Repository) *Service {
  serv := new(Service)
  fileserv := NewFileService()

  serv.repo = repo
  serv.fileService = fileserv

  return serv
}

func (s *Service) ListChallenges(ctx context.Context) ([]Challenge, error) {
  return s.repo.GetAll(ctx)
}

func (s *Service) GetChallenge(ctx context.Context, id string) (*Challenge, error) {
  return s.repo.GetByID(ctx, id)
}

func (s *Service) CreateChallenge(ctx context.Context, c *Challenge) error {
  return s.repo.Create(ctx, c)
}

func (s *Service) UpdateChallenge(ctx context.Context, id string, update *UpdateChallengeRequest) error {
  updateDoc := bson.M{}

  if update.Title != nil {
    updateDoc["title"] = *update.Title
  }
  if update.Category != nil {
    updateDoc["category"] = *update.Category
  }
  if update.Description != nil {
    updateDoc["description"] = *update.Description
  }
  if update.Points != nil {
    updateDoc["points"] = *update.Points
  }
  if update.State != nil {
    updateDoc["state"] = *update.State
  }
  if update.Type != nil {
    updateDoc["type"] = *update.Type
  }
  if update.Solves != nil {
    updateDoc["solves"] = *update.Solves
  }
  if update.Flag != nil {
    updateDoc["flag"] = *update.Flag
  }
  if update.Author != nil {
    updateDoc["author"] = *update.Author
  }

  return s.repo.Update(ctx, id, bson.M{"$set": updateDoc})
}

func (s *Service) DeleteChallenge(ctx context.Context, id string) error {
  return s.repo.Delete(ctx, id)
}

func (s *Service) CreateChallengeWithFiles(ctx context.Context, c *Challenge, form *multipart.Form, gc *gin.Context) error {
  uploadedFiles, err := s.fileService.processUploads(form.File["files"], gc)
  if err != nil {
    return err
  }

  c.Files = uploadedFiles

  if err := s.repo.Create(ctx, c); err != nil {
    // s.fileService.cleanupFiles(uploadedFiles)
    return err
  }

  return nil
}

func (s *Service) UpdateChallengeFile(ctx context.Context, id string, fileUUID string, form *multipart.Form, gc *gin.Context) error {
  challenge, err := s.repo.GetByID(ctx, id)
  if err != nil {
    return err
  }

  fileIndx := -1
  for k, f := range challenge.Files {
    if f.UUID != fileUUID {
      continue
    }

    fileIndx = k
    break
  }

  newFiles, err := s.fileService.processUploads(form.File["files"], gc)
  if err != nil || len(newFiles) != 1 {
    return ErrFileNotFound
  }

  // oldFile := challenge.Files[fileIndx]
  challenge.Files[fileIndx] = newFiles[0]

  update := bson.M{"$set": bson.M{"files": challenge.Files}}
  if err := s.repo.Update(ctx, id, update); err != nil {
    // s.fileService.cleanupFiles(newFiles)
    return err
  }

  // Not need as the cleanup handler with take care of it
  // s.fileService.cleanupFiles([]FileMeta{oldFile})

  return nil
}

func (s *Service) AddChallengeFile(ctx context.Context, id string, form *multipart.Form, gc *gin.Context) ([]FileMeta, error) {
  newFiles, err := s.fileService.processUploads(form.File["files"], gc)
  if err != nil {
    return nil, err
  }

  var uploadedFiles []FileMeta
  for _, newFile := range newFiles {
    update := bson.M{"$push": bson.M{"files": newFile}}
    if err := s.repo.Update(ctx, id, update); err != nil {
      // s.fileService.cleanupFiles([]FileMeta{newFile})
      return uploadedFiles, err
    }

    uploadedFiles = append(uploadedFiles, newFile)
  }

  return uploadedFiles, nil
}

func (s *Service) DeleteChallengeFile(ctx context.Context, id string, fileUUID string, gc *gin.Context) error {
  challenge, err := s.repo.GetByID(ctx, id)
  if err != nil {
    return err
  }

  var deletedFile FileMeta
  var remainingFiles []FileMeta

  fileIndx := -1
  for i, f := range challenge.Files {
    if f.UUID == fileUUID {
      deletedFile = f
      fileIndx = i
      break
    }
  }

  if deletedFile.UUID == "" {
    return ErrFileNotFound
  }

  remainingFiles = challenge.Files[:fileIndx]
  if fileIndx < len(challenge.Files)-1 {
    remainingFiles = append(remainingFiles, challenge.Files[fileIndx+1:]...)
  }

  update := bson.M{"$set": bson.M{"files": remainingFiles}}
  if err := s.repo.Update(ctx, id, update); err != nil {
    return err
  }

  // Not need as the cleanup handler with take care of it
  // s.fileService.cleanupFiles([]FileMeta{deletedFile})

  return nil
}

func (s *Service) CleanOrphanFileUploads(ctx context.Context) ([]string, error) {
  challenges, err := s.ListChallenges(ctx)
  if err != nil {
    return nil, err
  }

  hasParentFiles := make(map[string]int)
  for _, challenge := range challenges {
    for _, f := range challenge.Files {
      hasParentFiles[f.UUID] = 1
    }
  }

  files, err := os.ReadDir(s.fileService.uploadDir)
  if err != nil {
    return nil, err
  }

  var removedFiles []string

  for _, file := range files {
    if file.IsDir() {
      continue
    }
    if _, ok := hasParentFiles[file.Name()]; !ok {
      os.Remove(filepath.Join(s.fileService.uploadDir, file.Name()))
      removedFiles = append(removedFiles, file.Name())
    }
  }

  return removedFiles, nil
}
