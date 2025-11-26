package grpc

import (
	"errors"

	"github.com/motru4/legendary-currency-wallet/proto/exchange"
)

func currencyToString(c exchange.Currency) (string, error) {
	switch c {
	case exchange.Currency_USD:
		return "USD", nil
	case exchange.Currency_EUR:
		return "EUR", nil
	case exchange.Currency_RUB:
		return "RUB", nil
	case exchange.Currency_CURRENCY_UNSPECIFIED:
		return "", errors.New("currency unspecified")
	default:
		return "", errors.New("unknown currency")
	}
}

func stringToCurrency(s string) (exchange.Currency, error) {
	switch s {
	case "USD":
		return exchange.Currency_USD, nil
	case "EUR":
		return exchange.Currency_EUR, nil
	case "RUB":
		return exchange.Currency_RUB, nil
	default:
		return exchange.Currency_CURRENCY_UNSPECIFIED, errors.New("unsupported currency: " + s)
	}
}
