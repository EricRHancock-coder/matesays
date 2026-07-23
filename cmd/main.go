package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"matesays/internal/db"
	"matesays/internal/repository"
	"matesays/internal/service"
)

const sqlliteDBPath = "./matesays.db"

type config struct {
	list      bool
	add       string
	get       int64
	random    bool
	delete    int64
	deleteAll bool
}

func main() {
	// setup cli handler flags
	var cfg config

	flag.BoolVar(&cfg.list, "list", false, "list all quotes")
	flag.StringVar(&cfg.add, "add", "", "store a quote")
	flag.Int64Var(&cfg.get, "get", 0, "get quote by ID")
	flag.BoolVar(&cfg.random, "random", false, "print random quote")
	flag.Int64Var(&cfg.delete, "delete", 0, "delete quote by ID")
	flag.BoolVar(&cfg.deleteAll, "delete-all", false, "delete all quotes")

	flag.Parse()

	// setup logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	// check parameters
	if !cfg.list && cfg.get == 0 && cfg.add == "" && !cfg.random && cfg.delete == 0 && !cfg.deleteAll {
		printUsage()
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run(cfg config) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	sqlDB, err := db.Open(ctx, sqlliteDBPath)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	repo := repository.NewSqlliteQuoteRepo(sqlDB)
	svc := service.NewQuoteService(repo)

	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// list all quotes
	if cfg.list {
		quotes, err := svc.GetAllQuotes(callCtx)
		if err != nil {
			return fmt.Errorf("failed to get all quotes: %w", err)
		}

		for _, value := range quotes {
			fmt.Printf("%d: %s - %s\n", value.ID, value.Quote, value.CreatedTime)
		}
	}

	// get quote by id
	if cfg.get != 0 {
		quote, err := svc.GetQuote(callCtx, cfg.get)
		if err != nil {
			return fmt.Errorf("failed to get quote: %w", err)
		}
		fmt.Printf("%d: %s - %s\n", quote.ID, quote.Quote, quote.CreatedTime)
	}

	// add a quote and print it
	if cfg.add != "" {
		quote, err := svc.AddQuote(callCtx, cfg.add)
		if err != nil {
			return fmt.Errorf("failed to add quote: %w", err)
		}
		fmt.Printf("%d: %s - %s\n", quote.ID, quote.Quote, quote.CreatedTime)
	}

	// get random quote
	if cfg.random {
		quote, err := svc.GetRandomQuote(callCtx)
		if err != nil {
			return fmt.Errorf("failed to get random quote: %w", err)
		}
		fmt.Printf("%d: %s - %s\n", quote.ID, quote.Quote, quote.CreatedTime)
	}

	// delete a single quote
	if cfg.delete != 0 {
		err := svc.DeleteQuote(callCtx, cfg.delete)
		if err != nil {
			return fmt.Errorf("failed to delete quote: %w", err)
		}
		fmt.Printf("deleted quote %d\n", cfg.delete)
	}

	// delete all quotes from database
	if cfg.deleteAll {
		deleteCount, err := svc.DeleteAllQuotes(callCtx)
		if err != nil {
			return fmt.Errorf("failed to delete all quotes: %w", err)
		}
		fmt.Printf("deleted all %d quotes\n", deleteCount)
	}

	return nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [--list] [--add <text>] ...\n", os.Args[0])
	flag.PrintDefaults()
}
