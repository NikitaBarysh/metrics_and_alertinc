package service

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

var durationSleep = map[int]time.Duration{
	1: time.Second,
	2: 3 * time.Second,
	3: 5 * time.Second,
}

func Retry(fn func() error, attempt int) {
	err := fn()

	if attempt > 3 {
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		attempt++
		time.Sleep(durationSleep[attempt])
		Retry(fn, attempt)
	}

	if errors.Is(err, pgx.ErrTxClosed) {
		attempt++
		time.Sleep(durationSleep[attempt])
		Retry(fn, attempt)
	}

	if errors.Is(err, pgx.ErrTxCommitRollback) {
		attempt++
		time.Sleep(durationSleep[attempt])
		Retry(fn, attempt)
	}

	if errors.Is(err, sql.ErrConnDone) {
		attempt++
		time.Sleep(durationSleep[attempt])
		Retry(fn, attempt)
	}

	if errors.Is(err, sql.ErrNoRows) {
		attempt++
		time.Sleep(durationSleep[attempt])
		Retry(fn, attempt)
	}

	if errors.Is(err, sql.ErrTxDone) {
		attempt++
		time.Sleep(durationSleep[attempt])
		Retry(fn, attempt)
	}
}
