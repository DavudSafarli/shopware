package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"strings"

	"encoding/csv"
	"fmt"
	"io"
	"redirectware/internal"
	"redirectware/storage/postgres"

	_ "github.com/lib/pq"
)

type Entry struct {
	From string
	To   string
}

func main() {
	slog.Info("connecting to postgres", "dsn_set", os.Getenv("POSTGRES_URL") != "")
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	storage := postgres.New(db)

	entryChan := make(chan Entry)
	go func() {
		populate(os.Getenv("CSV_FILE"), entryChan)
		close(entryChan)
	}()

	for entry := range entryChan {
		rule, err := internal.NewFullMatchRule(entry.From, entry.To)
		if err != nil {
			slog.Error("failed to create full match rule", "from", entry.From, "to", entry.To, "err", err)
			break
		}
		if err := storage.AddFullMatchRule(context.Background(), rule); err != nil {
			slog.Error("failed to insert full match rule", "from", entry.From, "to", entry.To, "err", err)
			break
		}
	}

	storage.SetWelcomePageURL(context.Background(), "https://www.shopware.com/en/")
}

func populate(filepath string, entryChan chan<- Entry) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(fmt.Errorf("failed to open csv file: %w", err))
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read the header
	header, err := reader.Read()
	if err != nil {
		panic(fmt.Errorf("failed to read csv header: %w", err))
	}
	// Find columns for "from" and "to" (case-insensitive, just in case)
	fromIdx, toIdx := -1, -1
	for i, h := range header {
		lh := strings.ToLower(strings.TrimSpace(h))
		switch lh {
		case "from":
			fromIdx = i
		case "to":
			toIdx = i
		}
	}
	if fromIdx == -1 || toIdx == -1 {
		panic("CSV must have 'from' and 'to' columns")
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("error reading csv", "err", err)
			continue
		}

		from := strings.TrimSpace(record[fromIdx])
		to := strings.TrimSpace(record[toIdx])

		if from == "" || to == "" {
			continue // skip invalid rows
		}
		entryChan <- Entry{From: from, To: to}
	}

}
