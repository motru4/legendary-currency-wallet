package repositories

import (
	"context"
	"errors"

	"github.com/motru4/legendary-currency-wallet/services/exchanger/internal/models"
)

var (
	ErrRateNotFound = errors.New("exchange rate not found")
)

type ExchangeRatesRepository interface {
	// GetRate возвращает курс для конкретной пары валют.
	// base и target — строковые коды валют: "USD", "EUR", "RUB".
	GetRate(ctx context.Context, base, target string) (*models.ExchangeRate, error)

	// GetRatesForBase возвращает все курсы для указанной базовой валюты.
	// Например, для "USD" вернёт курсы USD->EUR, USD->RUB (если они есть).
	GetRatesForBase(ctx context.Context, base string) ([]models.ExchangeRate, error)
}
