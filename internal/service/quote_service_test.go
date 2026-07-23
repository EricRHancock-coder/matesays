package service

import (
		"context"
		"errors"
		"testing"
		"time"

		"matesays/internal/models"
		"matesays/internal/repository"
)

// fakeRepo is a hand-rolled mock, no testify needed.
type fakeRepo struct {
		quotes []models.Quote
		nextID int64
}

func (f *fakeRepo) GetAllQuotes(ctx context.Context) ([]models.Quote, error) {
		return f.quotes, nil
}

func (f *fakeRepo) GetByID(ctx context.Context, id int64) (*models.Quote, error) {
		for i := range f.quotes {
					if f.quotes[i].ID == id {
								return &f.quotes[i], nil
					}
		}
		return nil, repository.ErrNotFound
}

func (f *fakeRepo) AddQuote(ctx context.Context, text string) (*models.Quote, error) {
		f.nextID++
		q := models.Quote{ID: f.nextID, Quote: text, CreatedTime: time.Now()}
		f.quotes = append(f.quotes, q)
		qCopy := q // return a copy so caller doesn't mutate slice
		return &qCopy, nil
}

func (f *fakeRepo) GetRandom(ctx context.Context) (*models.Quote, error) {
		if len(f.quotes) == 0 {
					return nil, repository.ErrNotFound
		}
		return &f.quotes[0], nil
}

func (f *fakeRepo) Delete(ctx context.Context, id int64) error {
		for i, q := range f.quotes {
					if q.ID == id {
								f.quotes = append(f.quotes[:i], f.quotes[i+1:]...)
								return nil
					}
		}
		return repository.ErrNotFound
}

func (f *fakeRepo) DeleteAll(ctx context.Context) (int64, error) {
		n := int64(len(f.quotes))
		f.quotes = nil
		return n, nil
}

func newFakeRepo() *fakeRepo {
		return &fakeRepo{
					quotes: []models.Quote{
								{ID: 1, Quote: "hello", CreatedTime: time.Now()},
								{ID: 2, Quote: "world", CreatedTime: time.Now()},
					},
					nextID: 2,
		}
}

// TestQuoteService_GetQuote_Found tests retrieving an existing quote.
func TestQuoteService_GetQuote_Found(t *testing.T) {
		ctx := context.Background()
		svc := NewQuoteService(newFakeRepo())

		quote, err := svc.GetQuote(ctx, 1)
		if err != nil {
					t.Fatalf("unexpected error: %v", err)
		}
		if quote.Quote != "hello" {
					t.Errorf("got %q, want hello", quote.Quote)
		}
		if quote.ID != 1 {
					t.Errorf("got ID %d, want 1", quote.ID)
		}
}

// TestQuoteService_GetQuote_NotFound tests retrieving a missing quote.
func TestQuoteService_GetQuote_NotFound(t *testing.T) {
		ctx := context.Background()
		svc := NewQuoteService(newFakeRepo())

		_, err := svc.GetQuote(ctx, 99)
		if !errors.Is(err, repository.ErrNotFound) {
					t.Fatalf("expected ErrNotFound, got %v", err)
		}
}

// TestQuoteService_GetAllQuotes tests listing all quotes.
func TestQuoteService_GetAllQuotes(t *testing.T) {
		ctx := context.Background()
		svc := NewQuoteService(newFakeRepo())

		quotes, err := svc.GetAllQuotes(ctx)
		if err != nil {
					t.Fatalf("unexpected error: %v", err)
		}
		if len(quotes) != 2 {
					t.Errorf("expected 2 quotes, got %d", len(quotes))
		}
}

// TestQuoteService_AddQuote tests adding a new quote.
func TestQuoteService_AddQuote(t *testing.T) {
		ctx := context.Background()
		repo := newFakeRepo()
		svc := NewQuoteService(repo)

		quote, err := svc.AddQuote(ctx, "something mate said")
		if err != nil {
					t.Fatalf("unexpected error: %v", err)
		}
		if quote.Quote != "something mate said" {
					t.Errorf("got %q, want something mate said", quote.Quote)
		}
		if quote.ID != 3 {
					t.Errorf("got ID %d, expected 3", quote.ID)
		}

		// Verify it was added
		quotes, _ := repo.GetAllQuotes(ctx)
		if len(quotes) != 3 {
					t.Errorf("expected 3 quotes in repo, got %d", len(quotes))
		}
}

// TestQuoteService_GetRandomQuote tests getting a random quote.
func TestQuoteService_GetRandomQuote(t *testing.T) {
		ctx := context.Background()
		svc := NewQuoteService(newFakeRepo())

		quote, err := svc.GetRandomQuote(ctx)
		if err != nil {
					t.Fatalf("unexpected error: %v", err)
		}
		if quote == nil {
					t.Fatal("expected quote, got nil")
		}
		if quote.ID != 1 {
					t.Errorf("got ID %d, expected 1 (first in fake)", quote.ID)
		}
}

// TestQuoteService_GetRandomQuote_Empty tests random on empty repo.
func TestQuoteService_GetRandomQuote_Empty(t *testing.T) {
		ctx := context.Background()
		svc := NewQuoteService(&fakeRepo{})

		_, err := svc.GetRandomQuote(ctx)
		if !errors.Is(err, repository.ErrNotFound) {
					t.Fatalf("expected ErrNotFound, got %v", err)
		}
}

// TestQuoteService_DeleteQuote tests deleting an existing quote.
func TestQuoteService_DeleteQuote(t *testing.T) {
		ctx := context.Background()
		repo := newFakeRepo()
		svc := NewQuoteService(repo)

		if err := svc.DeleteQuote(ctx, 1); err != nil {
					t.Fatalf("unexpected error: %v", err)
		}

		// Verify it's gone
		_, err := svc.GetQuote(ctx, 1)
		if !errors.Is(err, repository.ErrNotFound) {
					t.Fatalf("expected ErrNotFound after delete, got %v", err)
		}
}

// TestQuoteService_DeleteQuote_NotFound tests deleting a missing quote.
func TestQuoteService_DeleteQuote_NotFound(t *testing.T) {
		ctx := context.Background()
		svc := NewQuoteService(newFakeRepo())

		err := svc.DeleteQuote(ctx, 99)
		if !errors.Is(err, repository.ErrNotFound) {
					t.Fatalf("expected ErrNotFound, got %v", err)
		}
}

// TestQuoteService_DeleteAllQuotes tests deleting all quotes.
func TestQuoteService_DeleteAllQuotes(t *testing.T) {
		ctx := context.Background()
		repo := newFakeRepo()
		svc := NewQuoteService(repo)

		count, err := svc.DeleteAllQuotes(ctx)
		if err != nil {
					t.Fatalf("unexpected error: %v", err)
		}
		if count != 2 {
					t.Errorf("got %d deleted, want 2", count)
		}

		quotes, _ := repo.GetAllQuotes(ctx)
		if len(quotes) != 0 {
					t.Errorf("expected 0 quotes, got %d", len(quotes))
		}
}

// TestQuoteService_DeleteAllQuotes_Empty tests delete-all on empty repo.
func TestQuoteService_DeleteAllQuotes_Empty(t *testing.T) {
		ctx := context.Background()
		svc := NewQuoteService(&fakeRepo{})

		count, err := svc.DeleteAllQuotes(ctx)
		if err != nil {
					t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
					t.Errorf("got %d deleted, want 0", count)
		}
}
