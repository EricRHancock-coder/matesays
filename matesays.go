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
	"log/slog"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type quote struct {
	_ID        int64
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
	var dsFlag, psFlag int64
	var dFlag, pFlag *bool
	// qFlag := flag.String("q", "", "Prints quote associated with the number passed in.")
	flag.StringVar(&sFlag, "s", "", "Stores the quoted string and prints the number associated with it.")
	pFlag = flag.Bool("pa", false, "Prints all stored quotes.")
	dFlag = flag.Bool("da", false, "Deletes all quotes stored.")
	flag.Int64Var(&dsFlag, "ds", 0, "Deletes a single quote stored.")
	flag.Int64Var(&psFlag, "ps", 0, "Prints a single quoted string and prints the number associated with it.")

	flag.Usage = func() {
		printUsage()
		os.Exit(1)
	}

	flag.Parse()

	if sFlag == "" && !*dFlag && !*pFlag && dsFlag == 0 && psFlag == 0 {
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

	db, dbErr = sql.Open("sqlite", "./matesays.db")
	if dbErr != nil {
		slog.Error("failed to open database", "error", dbErr)
		os.Exit(1)
	}
	defer db.Close()

	dbErr = nil
	dbErr = createTables()
	if dbErr != nil {
		slog.Error("error creating table", "error", dbErr)
		os.Exit(1)
	}

	// get all quotes from db
	dbErr = nil
	dbErr = getAllQuotes()
	if dbErr != nil {
		slog.Error("error getting quotes", "error", dbErr)
		os.Exit(1)
	}

	// print single quote
	if psFlag != 0 {
		dbErr = nil
		dbErr = printSingleQuote(psFlag)
		if dbErr != nil {
			slog.Error("Error retreiving quote", "error", dbErr)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if dsFlag != 0 {
		err := deleteQuote(dsFlag)
		if err != nil {
			slog.Error("error deleting quote", "error", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *dFlag {
		dbErr = nil
		dbErr = deleteAllQuotes()

		if dbErr != nil {
			slog.Error("error deleting all quotes", "error", dbErr)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if sFlag != "" {
		slog.Info("adding quote", "quote", sFlag)
		qt := createNewQuote(0, sFlag, time.Now())
		dbErr = insertQuote(qt)
		if dbErr != nil {
			slog.Error("error creating quote", "error", dbErr)
			os.Exit(1)
		}
		addQuote(qt)
		printQuote(qt)
		os.Exit(0)
	}

	slog.Info("Printing all quotes...\n")
	printAllQuotes()
}

func printQuote(qt *quote) {
	fmt.Printf("%d:%s %s\n", qt._ID, qt._quote, qt._timeStamp)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-s string] [-q number]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "This program stores and retrieves quotes that my roommate says.\n")
	flag.PrintDefaults()
}

func createNewQuote(index int64, newQt string, newTm time.Time) *quote {
	newQuote := new(quote)
	newQuote._quote = newQt
	newQuote._timeStamp = newTm
	newQuote._ID = index

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

func printSingleQuote(id int64) error {
	for _, tmpQuote := range quotes {
		if tmpQuote._ID == id {
			printQuote(&tmpQuote)
			return nil
		}
	}

	return fmt.Errorf("quote %d not found", id)
}

func printAllQuotes() {
	for _, tmpQuote := range quotes {
		printQuote(&tmpQuote)
	}
}

// **************** DATABASE STUFF *************

func createTables() error {
	_, dbErr = db.Exec("CREATE TABLE  if not exists quotes ( id integer primary key autoincrement, quote string, created_time datetime default current_timestamp )")
	if dbErr != nil {
		return dbErr
	}

	return nil
}

func insertQuote(iQuote *quote) error {
	query := `insert into quotes (quote, created_time) values (?,?)`
	iTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, query, iQuote._quote, iTime)
	if err != nil {
		return fmt.Errorf("failed to insert quote: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		slog.Info("Successfully inserted", "quote string", iQuote._quote, "id", id, "created time", iTime)
	}

	iQuote._ID = id
	iQuote._timeStamp = iTime

	return nil
}

func getAllQuotes() error {
	rows, err := db.Query("select id,quote,created_time from quotes")
	if err != nil {
		return err
	}
	defer rows.Close()

	var storedID int64
	var storedQuote string
	var storedTime time.Time
	for rows.Next() {
		if err := rows.Scan(&storedID, &storedQuote, &storedTime); err != nil {
			return err
		}

		newQuote := createNewQuote(storedID, storedQuote, storedTime)
		addQuote(newQuote)

	}

	return nil
}

func getQuote(id int64) (error, string, int64) {
	rows, err := db.Query("select id,quote,created_time from quotes where id = ?", id)
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

func deleteAllQuotes() error {
	query := "delete from quotes"

	result, err := db.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	slog.Info("Deleted all records!\n", "rows deleted", rowsAffected)

	return nil
}

func deleteQuote(id int64) error {
	rows, err := db.Query("delete from quotes where id = ?", id)
	if err != nil {
		return err
	}
	defer rows.Close()

	slog.Info("Deleted 1 record!\n")

	return nil
}
