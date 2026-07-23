// Package repository contains all repositories for this application
// QuoteRepository contains quote data
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"matesays/internal/models"
)

var ErrNotFound = errors.New("not found")

type QuoteRepository interface {
	GetAllQuotes(ctx context.Context) ([]models.Quote, error)
	GetByID(ctx context.Context, id int64) (*models.Quote, error)
	AddQuote(ctx context.Context, text string) (*models.Quote, error)
	GetRandom(ctx context.Context) (*models.Quote, error)
	Delete(ctx context.Context, id int64) error
	DeleteAll(ctx context.Context) (int64, error)
}

type sqlliteQuoteRepo struct {
	db *sql.DB
}

func NewSqlliteQuoteRepo(db *sql.DB) QuoteRepository {
	return &sqlliteQuoteRepo{db: db}
}

func (r *sqlliteQuoteRepo) GetAllQuotes(ctx context.Context) ([]models.Quote, error) {
	var quotes []models.Quote

	rows, err := r.db.QueryContext(ctx, "select id,quote,created_time from quotes")
	if err != nil {
		return nil, fmt.Errorf("querying context: %w", err)
	}
	// Proper way to check for close() err and defer by using anonymous func
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			// todo: change log to xlog or something better
			slog.Error("closing rows", "error", cerr)
		}
	}()

	for rows.Next() {
		var q models.Quote
		if err := rows.Scan(&q.ID, &q.Quote, &q.CreatedTime); err != nil {
			return nil, fmt.Errorf("error while fetching quotes: %w", err)
		}

		quotes = append(quotes, q)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating quotes: %w", err)
	}

	return quotes, nil
}

func (r *sqlliteQuoteRepo) GetByID(ctx context.Context, id int64) (*models.Quote, error) {
	var q models.Quote

	err := r.db.QueryRowContext(ctx, "select id,quote,created_time from quotes where id = ? ", id).Scan(&q.ID, &q.Quote, &q.CreatedTime)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("query error: %w", err)
	}
	return &q, nil
}

func (r *sqlliteQuoteRepo) AddQuote(ctx context.Context, text string) (*models.Quote, error) {
	var quote models.Quote

	err := r.db.QueryRowContext(ctx, "insert into quotes (quote) values (?) returning id, quote, created_time", text).Scan(&quote.ID, &quote.Quote, &quote.CreatedTime)
	if err != nil {
		return nil, fmt.Errorf("failed to insert: %w", err)
	}

	return &quote, nil
}

func (r *sqlliteQuoteRepo) GetRandom(ctx context.Context) (*models.Quote, error) {
	var q models.Quote

	err := r.db.QueryRowContext(ctx, "select id,quote,created_time from quotes ORDER BY RANDOM() LIMIT 1").Scan(&q.ID, &q.Quote, &q.CreatedTime)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("query error: %w", err)
	}
	return &q, nil
}

func (r *sqlliteQuoteRepo) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM quotes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleting quote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *sqlliteQuoteRepo) DeleteAll(ctx context.Context) (int64, error) {
	result, err := r.db.ExecContext(ctx, "DELETE FROM quotes")
	if err != nil {
		return 0, fmt.Errorf("deleting quotes: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("checking rows affected: %w", err)
	}

	// perfectly fine to return 0 rows affected
	return rowsAffected, nil
}
