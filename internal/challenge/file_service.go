/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package challenge

import (
  "errors"
  "log"
  "mime/multipart"
  "os"
  "path/filepath"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/google/uuid"
)

const uploadDir = "uploads/challenge"

var (
  ErrNoFile         = errors.New("uploaded files not found")
  ErrFileExcedLimit = errors.New("file size exceeds limit")
)

type FileService struct {
  uploadDir string
}

func NewFileService() *FileService {
  fileService := new(FileService)

  fileService.uploadDir = uploadDir
  os.MkdirAll(uploadDir, 0755)

  return fileService
}

func (fs *FileService) processUploads(files []*multipart.FileHeader, c *gin.Context) ([]FileMeta, error) {

  if len(files) == 0 {
    return nil, ErrNoFile
  }

  maxSize := int64(10 << 20)
  for _, file := range files {
    if file.Size > maxSize {
      return nil, ErrFileExcedLimit
    }
  }

  var uploadedFiles []FileMeta
  var cleanups []string

  var err error
  defer func() {
    if err != nil {
      for _, path := range cleanups {
        os.Remove(path)
      }
      log.Printf("cleaning files: %d(%v)\n", len(cleanups), err)
    }
  }()

  for _, fileHeader := range files {
    uuid := uuid.NewString()
    filePath := filepath.Join(uploadDir, uuid)

    if err = c.SaveUploadedFile(fileHeader, filePath); err != nil {
      return nil, err
    }

    cleanups = append(cleanups, filePath)

    file, errFh := fileHeader.Open()
    if errFh != nil {
      err = errFh
      return nil, err
    }
    defer file.Close()

    uploadedFiles = append(uploadedFiles, FileMeta{
      UUID:       uuid,
      Name:       fileHeader.Filename,
      Size:       fileHeader.Size,
      UploadedAt: time.Now().UTC(),
    })
  }

  return uploadedFiles, nil
}

func (fs *FileService) cleanupFiles(files []FileMeta) {
  for _, file := range files {
    os.Remove(filepath.Join(uploadDir, file.UUID))
  }
}
