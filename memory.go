package memory

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/sony/gobreaker"
	"github.com/pkg/errors"
	_ "github.com/lib/pq"
)

const (
	MAX_HISTORY_DAYS = 30
	MAX_RETRIES      = 3
)

var cb *gobreaker.CircuitBreaker

func init() {
	st := gobreaker.Settings{
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("[%s] Breaker state transition: %v => %v\n", name, from, to)
		},
	}
	cb = gobreaker.NewCircuitBreaker(st)
}

type User struct {
	ID      string
	Context []string
	History []HistoryItem
}

type HistoryItem struct {
	ID        int
	Text      string
	Timestamp time.Time
	Sentiment string
	Intents   []string
}

func GetUser(id string) (*User, error) {
	user, err := cb.Execute(func() (interface{}, error) {
		log.Printf("Circuit Breaker: Closed, executing DB call\n")
		result, err := fetchUserWithRetry(id)
		if err != nil {
			log.Printf("Circuit Breaker: Execution failed: %v\n", err)
			return nil, err
		}
		log.Printf("Circuit Breaker: Execution successful\n")
		return result, nil
	})
	if err != nil {
		log.Printf("Circuit Breaker: Open, error: %v\n", err)
		return nil, err
	}
	return user.(*User), nil
}

func fetchUserWithRetry(id string) (*User, error) {
	var user *User

	err := retry.Do(
		func() error {
			var internalErr error
			user, internalErr = fetchUserFromDB(id)
			if internalErr != nil {
				if internalErr == sql.ErrNoRows {
					return retry.Unrecoverable(internalErr)
				}
				return internalErr
			}
			return nil
		},
		retry.Attempts(MAX_RETRIES),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("#%d: Error fetching user: %s\n", n, err)
		}),
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No user found with ID:", id)
			return nil, nil // Return nil user without error
		}
		return nil, err
	}

	return user, nil
}

func fetchUserFromDB(id string) (*User, error) {
	db, err := sql.Open("postgres", "connection_string")
	if err != nil {
		log.Println("Error opening database connection:", err)
		return nil, errors.Wrap(err, "error opening database connection")
	}
	defer db.Close()

	var context, intents string
	err = db.QueryRow("SELECT context FROM users WHERE id=$1", id).Scan(&context)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Println("Error fetching context:", err)
		return nil, errors.Wrap(err, "error fetching context")
	}

	contextSlice := strings.Split(context, ",")

	rows, err := db.Query("SELECT * FROM history WHERE user_id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &User{
				ID:      id,
				Context: contextSlice,
				History: []HistoryItem{},
			}, nil
		}
		log.Println("Error fetching history:", err)
		return nil, errors.Wrap(err, "error fetching history")
	}
	defer rows.Close()

	var history []HistoryItem
	for rows.Next() {
		var h HistoryItem
		err := rows.Scan(&h.ID, &h.Text, &h.Timestamp, &h.Sentiment, &intents)
		if err != nil {
			log.Println("Error scanning history row:", err)
		}
		h.Intents = strings.Split(intents, ",")
		history = append(history, h)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error during rows iteration:", err)
		return nil, errors.Wrap(err, "error during rows iteration")
	}

	return &User{
		ID:      id,
		Context: contextSlice,
		History: history,
	}, nil
}
