package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/motru4/legendary-currency-wallet/services/exchanger/internal/models"
	"github.com/motru4/legendary-currency-wallet/services/exchanger/internal/repositories"
)

type ExchangeRatesPostgresRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

// compile-time проверка, что структура реализует интерфейс
var _ repositories.ExchangeRatesRepository = (*ExchangeRatesPostgresRepository)(nil)

func NewExchangeRatesPostgresRepository(db *sqlx.DB, log *slog.Logger) *ExchangeRatesPostgresRepository {
	return &ExchangeRatesPostgresRepository{
		db:  db,
		log: log.With("component", "exchange_rates_repo"),
	}
}

func (r *ExchangeRatesPostgresRepository) GetRate(
	ctx context.Context,
	base, target string,
) (*models.ExchangeRate, error) {
	const query = `
		SELECT base_currency, target_currency, rate, updated_at
		FROM exchange_rates
		WHERE base_currency = $1 AND target_currency = $2
	`

	var m models.ExchangeRate

	err := r.db.GetContext(ctx, &m, query, base, target)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Info("rate not found", "base", base, "target", target)
			return nil, repositories.ErrRateNotFound
		}

		r.log.Error("failed to get rate",
			"base", base,
			"target", target,
			"err", err,
		)
		return nil, err
	}

	r.log.Debug("rate loaded from db",
		"base", m.BaseCurrency,
		"target", m.TargetCurrency,
		"rate", m.Rate,
		"updated_at", m.UpdatedAt,
	)

	return &m, nil
}

func (r *ExchangeRatesPostgresRepository) GetRatesForBase(
	ctx context.Context,
	base string,
) ([]models.ExchangeRate, error) {
	const query = `
		SELECT base_currency, target_currency, rate, updated_at
		FROM exchange_rates
		WHERE base_currency = $1
		ORDER BY target_currency
	`

	var rates []models.ExchangeRate

	if err := r.db.SelectContext(ctx, &rates, query, base); err != nil {
		r.log.Error("failed to query rates for base", "base", base, "err", err)
		return nil, err
	}

	if len(rates) == 0 {
		r.log.Info("no rates for base", "base", base)
	}

	r.log.Debug("rates loaded for base", "base", base, "count", len(rates))

	return rates, nil
}
