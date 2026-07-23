// Package service validates and passes data from and to repo
package service

import (
	"context"

	"matesays/internal/models"
	"matesays/internal/repository"
)

// QuoteService depends on the REPOSITORY INTERFACE, not a concrete type.
// This is dependency inversion in action.
type QuoteService struct {
	repo repository.QuoteRepository
}

func NewQuoteService(repo repository.QuoteRepository) *QuoteService {
	return &QuoteService{repo: repo}
}

func (s *QuoteService) GetAllQuotes(ctx context.Context) ([]models.Quote, error) {
	return s.repo.GetAllQuotes(ctx)
}

func (s *QuoteService) GetQuote(ctx context.Context, id int64) (*models.Quote, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *QuoteService) AddQuote(ctx context.Context, text string) (*models.Quote, error) {
	return s.repo.AddQuote(ctx, text)
}

func (s *QuoteService) GetRandomQuote(ctx context.Context) (*models.Quote, error) {
	return s.repo.GetRandom(ctx)
}

func (s *QuoteService) DeleteQuote(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *QuoteService) DeleteAllQuotes(ctx context.Context) (int64, error) {
	return s.repo.DeleteAll(ctx)
}
