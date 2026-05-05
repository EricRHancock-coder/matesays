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
	"flag"
	"fmt"
	"os"
)

type Quote struct {
	_Index int    `json:"index"`
	_Quote string `json:"quote"`
}

var quotes []Quote

func main() {
	var sFlag string
	// qFlag := flag.String("q", "", "Prints quote associated with the number passed in.")
	flag.StringVar(&sFlag, "s", "", "Stores the quoted string and prints the number associated with it.")
	flag.StringVar(&sFlag, "store", "", "Stores the quoted string and prints the number associated with it.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-s string] [-q number]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "This program stores and retrieves quotes that my roommate says.\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()

	if sFlag == "" {
		fmt.Fprintf(os.Stderr, "Error, cannot pass in empty quote.\n")
		os.Exit(1)
	}

	addQuote(sFlag)
	fmt.Printf("Printing all quotes...\n")
	printAllQuotes()
}

func createNewQuote(index int, quote string) *Quote {
	newQuote := new(Quote)
	newQuote._Index = index
	newQuote._Quote = quote

	return newQuote
}

// storeQuote adds the quote to the quotes array
func addQuote(quote string) int {
	// loop over slice and get the max index and add one to it
	newIndex := 0
	for _, tmpQuote := range quotes {
		if tmpQuote._Index >= newIndex {
			newIndex = tmpQuote._Index + 1
		}
	}

	newQuote := createNewQuote(newIndex, quote)
	quotes = append(quotes, *newQuote)

	return newIndex
}

func printAllQuotes() {
	for _, tmpQuote := range quotes {
		fmt.Printf("%d: %s\n", tmpQuote._Index, tmpQuote._Quote)
	}
}

func saveQuotes(quotes []string) {
}
