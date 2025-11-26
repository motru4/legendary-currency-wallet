package models

import "time"

type ExchangeRate struct {
	BaseCurrency   string    `db:"base_currency"`
	TargetCurrency string    `db:"target_currency"`
	Rate           float64   `db:"rate"`
	UpdatedAt      time.Time `db:"updated_at"`
}
