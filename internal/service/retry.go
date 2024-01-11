// Package service - содержит внутреннею логику приложения
package service

import (
	"database/sql"
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
)

var durationSleep = map[int]time.Duration{
	1: time.Second,
	2: 3 * time.Second,
	3: 5 * time.Second,
}

// Retry - попытка сделать повторные запросы при ошибке
func Retry(fn func() error, attempt int) {
	err := fn()
	var netErr net.Error
	var urlErr url.Error

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

	if errors.As(err, &netErr) && netErr.Timeout() {
		attempt++
		time.Sleep(durationSleep[attempt])
		Retry(fn, attempt)
	}

	if errors.Is(err, &urlErr) {
		attempt++
		time.Sleep(durationSleep[attempt])
		Retry(fn, attempt)
	}
}
