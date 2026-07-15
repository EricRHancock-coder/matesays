// Package main stores and retrieves quotes that
// roomnate says.
// When nothing is passed in it prints a random
// quote prefixed by the quote number..
// When -q (--quote) is passed in followed by a
// number it retrieves the quote associated with
// that number. If the quote is not found it
// returns the error "Quote not found " followed
// by the nunber.
// When -s (--store) is passed in the string
// folowing -s is stored as a quote and the quote
// number is returned. If no string us passed in
// then stdin is used. If stdin is empty then
// error is returned "No quote to store!"
//
// Progroam will store quotes om disk in a flat
// file.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type quote struct {
	_ID        int
	_quote     string
	_timeStamp time.Time
}

var quotes []quote

var (
	db    *sql.DB
	dbErr error
)

func main() {
	var sFlag string
	// qFlag := flag.String("q", "", "Prints quote associated with the number passed in.")
	flag.StringVar(&sFlag, "s", "", "Stores the quoted string and prints the number associated with it.")
	flag.StringVar(&sFlag, "store", "", "Stores the quoted string and prints the number associated with it.")

	flag.Usage = func() {
		printUsage()
		os.Exit(1)
	}

	flag.Parse()

	if sFlag == "" {
		printUsage()
		os.Exit(1)
	}

	// create logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo, // Options LevelDebugm LevelInfo, LEvelWarnm LevelError
	}

	// create text handler or JSON handler
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	db, dbErr := sql.Open("sqlite", "./matesays.db")
	if dbErr != nil {
		slog.Error("failed to open database", "error", dbErr)
		os.Exit(1)
	}
	defer db.Close()

	dbErr = createTables()
	if dbErr != nil {
		slog.Error("error creating table", "error", dbErr)
	}

	addQuote(sFlag)
	addQuote("Dip turtle")
	addQuote("Another noise")
	log.Printf("Printing all quotes...\n")
	printAllQuotes()
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-s string] [-q number]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "This program stores and retrieves quotes that my roommate says.\n")
	flag.PrintDefaults()
}

func createNewQuote(index int, newQt string, newTm time.Time) *quote {
	newQuote := new(quote)
	newQuote._ID = index
	newQuote._quote = newQt
	newQuote._timeStamp = newTm

	return newQuote
}

// storeQuote adds the quote to the quotes array
func addQuote(newQt *quote) int {
	// loop over slice and get the max id and add one to it
	newID := 0

	newQuote := newQt
	quotes = append(quotes, *newQuote)

	return newID
}

func printAllQuotes() {
	for _, tmpQuote := range quotes {
		fmt.Printf("%d: %s\n", tmpQuote._ID, tmpQuote._quote)
	}
}

// **************** DATABASE STUFF *************

func createTables() error {
	_, dbErr = db.Exec("CREATE TABLE  if not exists quotes ( id integer primary key autoincrement, quote string, created_time datetime default current_timestamp );")
	if dbErr != nil {
		return dbErr
	}

	return nil
}

func insertQuote(quote string) (int64, error) {
	query := `insert into quotes (quote) values (?)`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, query, quote)
	if err != nil {
		return 0, fmt.Errorf("failed to insert quote: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		slog.Info("Successfully inserted {%s} id {%d}\n", quote, id)
	}

	return id, nil
}

func getAllQuotes() error {
	rows, err := db.Query("select id,quote,create_time from quotes")
	if err != nil {
		return err
	}
	defer rows.Close()

	var storedID int64
	var storedQuote string
	var storedTime time.Time
	for rows.Next() {
		if err := rows.Scan(&storedID, &storedQuotei, &storedTime); err != nil {
			return err
		}
	}

	return nil
}

func getQuote(id int64) (error, string, int64) {
	rows, err := db.Query("select id,quote,create_time from quotes where id = ?", id)
	if err != nil {
		return err, "", 0
	}
	defer rows.Close()

	var storedID int64
	var storedQuote string
	for rows.Next() {
		if err := rows.Scan(&storedID, &storedQuote); err != nil {
			return err, "", 0
		}
	}

	return nil, storedQuote, storedID
}
