package repository

import (
		"context"
		"errors"
		"testing"
		"time"

		"matesays/internal/db"
)

func TestQuotesRepo_CRUD(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		sqlDB, err := db.Open(ctx, ":memory:")
		if err != nil {
					t.Fatalf("opening database: %v", err)
		}
		defer sqlDB.Close()

		repo := NewSqlliteQuoteRepo(sqlDB)

		// Test 1: AddQuote
		quote, err := repo.AddQuote(ctx, "something mate said")
		if err != nil {
					t.Fatalf("add quote: %v", err)
		}
		if quote.Quote != "something mate said" {
					t.Errorf("quote text mismatch: got %q, want %q", quote.Quote, "something mate said")
		}
		if quote.ID <= 0 {
					t.Errorf("expected positive ID, got %d", quote.ID)
		}
		if quote.CreatedTime.IsZero() {
					t.Error("expected created time to be set")
		}

		// Test 2: GetByID - found
		got, err := repo.GetByID(ctx, quote.ID)
		if err != nil {
					t.Fatalf("get by id: %v", err)
		}
		if got.Quote != quote.Quote {
					t.Errorf("got quote %q, want %q", got.Quote, quote.Quote)
		}

		// Test 3: GetByID - not found
		_, err = repo.GetByID(ctx, 9999)
		if err == nil {
					t.Fatal("expected error for missing ID, got nil")
		}
		if !errors.Is(err, ErrNotFound) {
					t.Errorf("expected ErrNotFound, got %T: %v", err, err)
		}

		// Test 4: Add second quote
		second, err := repo.AddQuote(ctx, "another quote")
		if err != nil {
					t.Fatalf("add second quote: %v", err)
		}
		if second.ID <= quote.ID {
					t.Errorf("second ID %d should be greater than first ID %d", second.ID, quote.ID)
		}

		// Test 5: GetAllQuotes
		quotes, err := repo.GetAllQuotes(ctx)
		if err != nil {
					t.Fatalf("get all quotes: %v", err)
		}
		if len(quotes) != 2 {
					t.Errorf("expected 2 quotes, got %d", len(quotes))
		}

		// Test 6: GetRandom - should return one of the two quotes
		random, err := repo.GetRandom(ctx)
		if err != nil {
					t.Fatalf("get random quote: %v", err)
		}
		found := false
		for _, q := range quotes {
					if q.ID == random.ID {
								found = true
								break
					}
		}
		if !found {
					t.Errorf("random quote ID %d not found in quotes list", random.ID)
		}

		// Test 7: Delete - existing
		if err := repo.Delete(ctx, quote.ID); err != nil {
					t.Fatalf("delete existing: %v", err)
		}

		// Verify it's gone
		_, err = repo.GetByID(ctx, quote.ID)
		if !errors.Is(err, ErrNotFound) {
					t.Errorf("expected ErrNotFound after delete, got %v", err)
		}

		// Test 8: Delete - not found
		if err := repo.Delete(ctx, 9999); err == nil {
					t.Fatal("expected error deleting missing ID, got nil")
		} else if !errors.Is(err, ErrNotFound) {
					t.Errorf("expected ErrNotFound, got %T: %v", err, err)
		}

		// Test 9: DeleteAll
		beforeCount, _ := repo.GetAllQuotes(ctx)
		count, err := repo.DeleteAll(ctx)
		if err != nil {
					t.Fatalf("delete all: %v", err)
		}
		if int(count) != len(beforeCount) {
					t.Errorf("delete all returned %d, expected %d", count, len(beforeCount))
		}

		quotes, err = repo.GetAllQuotes(ctx)
		if err != nil {
					t.Fatalf("get all after delete all: %v", err)
		}
		if len(quotes) != 0 {
					t.Errorf("expected 0 quotes after delete all, got %d", len(quotes))
		}

		// Test 10: GetRandom - empty table
		_, err = repo.GetRandom(ctx)
		if !errors.Is(err, ErrNotFound) {
					t.Errorf("expected ErrNotFound on empty table, got %v", err)
		}
}

// TestQuotesRepo_AddQuote_EmptyString tests storing an empty string.
func TestQuotesRepo_AddQuote_EmptyString(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		sqlDB, err := db.Open(ctx, ":memory:")
		if err != nil {
					t.Fatalf("opening database: %v", err)
		}
		defer sqlDB.Close()

		repo := NewSqlliteQuoteRepo(sqlDB)

		quote, err := repo.AddQuote(ctx, "")
		if err != nil {
					t.Fatalf("add empty quote: %v", err)
		}
		if quote.Quote != "" {
					t.Errorf("expected empty string, got %q", quote.Quote)
		}
}
