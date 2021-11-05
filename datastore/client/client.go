// Copyright 2021 PGHQ. All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package client provides resources for data persistence and retrieval.
package client

import (
	"context"
	"time"
)

// Client represents a client for operating on a database.
type Client interface {
	Connect() error
	Filter() Filter
	Query() Query
	Add() Add
	Update() Update
	Remove() Remove
	Transaction(ctx context.Context) (Transaction, error)
}

// Transaction represents a database transaction.
type Transaction interface {
	Execute(statement Encoder) (int, error)
	Commit() error
	Rollback() error
}

// Add represents a command to add items to the collection
type Add interface {
	Encoder
	To(collection string) Add
	Item(value map[string]interface{}) Add
	Query(query Query) Add
	Execute(ctx context.Context) (int, error)
}

// Update represents a command to update items in the collection
type Update interface {
	Encoder
	In(collection string) Update
	Item(value map[string]interface{}) Update
	Filter(filter Filter) Update
	Execute(ctx context.Context) (int, error)
}

// Remove represents a command to remove items from the collection
type Remove interface {
	Encoder
	From(collection string) Remove
	Filter(filter Filter) Remove
	Order(by string) Remove
	After(key string, value *time.Time) Remove
	Execute(ctx context.Context) (int, error)
}

// Encoder represents a statement encoder
type Encoder interface {
	Statement() (string, []interface{}, error)
}

// Query represents a query builder
type Query interface {
	Encoder
	Secondary() Query
	From(collection string) Query
	And(collection string, args ...interface{}) Query
	Filter(filter Filter) Query
	Order(by string) Query
	First(first int) Query
	After(key string, value *time.Time) Query
	Return(key string, args ...interface{}) Query
	Execute(ctx context.Context, dst interface{}) error
}

// Filter represents criteria for querying the collection
type Filter interface {
	IsNil() bool
	Eq(key string, value interface{}) Filter
	Lt(key string, value interface{}) Filter
	Gt(key string, value interface{}) Filter
	NotEq(key string, value interface{}) Filter
	BeginsWith(key string, prefix string) Filter
	EndsWith(key string, suffix string) Filter
	Contains(key string, value interface{}) Filter
	NotContains(key string, value interface{}) Filter
	Or(another Filter) Filter
	And(another Filter) Filter
}
