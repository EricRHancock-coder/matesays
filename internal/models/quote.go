// Package models contains: Quote is the quote model mapped to the quotes table
package models

import (
	"time"
)

type Quote struct {
	ID          int64     `json:"id"`
	Quote       string    `json:"quote"`
	CreatedTime time.Time `json:"created_time"`
}
