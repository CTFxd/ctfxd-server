/*
 * Copyright (c) 2025, Arka Mondal. All rights reserved.
 * Use of this source code is governed by a BSD-style license that
 * can be found in the LICENSE file.
 */

package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongodbInit(uri string, dbName string) *MongoClient {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	opts.SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("MongoDB server ping failed: %v", err)
	}

	log.Println("MongoDB sever: Connection successful")
	return &MongoClient{
		Client:   client,
		Database: client.Database(dbName),
	}
}

func (mc *MongoClient) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mc.Client.Disconnect(ctx); err != nil {
		log.Fatalf("MongoDB disconnect failed: %v", err)
	}

	log.Println("MongoDB sever: Disconnect successful")
}
