package grpc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/motru4/legendary-currency-wallet/proto/exchange"
	"github.com/motru4/legendary-currency-wallet/services/exchanger/internal/repositories"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExchangeServiceServer struct {
	exchange.UnimplementedExchangeServiceServer

	repo repositories.ExchangeRatesRepository
	log  *slog.Logger
}

func NewExchangeServiceServer(repo repositories.ExchangeRatesRepository, log *slog.Logger) *ExchangeServiceServer {
	return &ExchangeServiceServer{
		repo: repo,
		log:  log.With("component", "exchange_service"),
	}
}

func (s *ExchangeServiceServer) GetExchangeRateForCurrency(
	ctx context.Context,
	req *exchange.CurrencyPairRequest,
) (*exchange.ExchangeRateResponse, error) {
	fromStr, err := currencyToString(req.FromCurrency)
	if err != nil {
		s.log.Warn("invalid from_currency in request",
			"from_currency", req.FromCurrency.String(),
			"err", err,
		)
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_currency: %v", err)
	}

	toStr, err := currencyToString(req.ToCurrency)
	if err != nil {
		s.log.Warn("invalid to_currency in request",
			"to_currency", req.ToCurrency.String(),
			"err", err,
		)
		return nil, status.Errorf(codes.InvalidArgument, "invalid to_currency: %v", err)
	}

	if fromStr == toStr {
		now := time.Now().Unix()
		s.log.Debug("same currency pair, returning 1.0 rate",
			"currency", fromStr,
		)

		return &exchange.ExchangeRateResponse{
			FromCurrency:  req.FromCurrency,
			ToCurrency:    req.ToCurrency,
			Rate:          1.0,
			UpdatedAtUnix: now,
		}, nil
	}

	rateModel, err := s.repo.GetRate(ctx, fromStr, toStr)
	if err != nil {
		if errors.Is(err, repositories.ErrRateNotFound) {
			s.log.Info("rate not found",
				"from", fromStr,
				"to", toStr,
			)
			return nil, status.Errorf(codes.NotFound, "rate not found for %s -> %s", fromStr, toStr)
		}

		s.log.Error("failed to get rate from repository",
			"from", fromStr,
			"to", toStr,
			"err", err,
		)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	resp := &exchange.ExchangeRateResponse{
		FromCurrency:  req.FromCurrency,
		ToCurrency:    req.ToCurrency,
		Rate:          rateModel.Rate,
		UpdatedAtUnix: time.Now().Unix(), //в данной реалицации данные в бд не обновляются, при добавлении этой функциональности использовать - rateModel.UpdatedAt.Unix(),
	}

	s.log.Debug("rate returned",
		"from", fromStr,
		"to", toStr,
		"rate", resp.Rate,
		"updated_at_unix", resp.UpdatedAtUnix,
	)

	return resp, nil
}

func (s *ExchangeServiceServer) GetExchangeRates(
	ctx context.Context,
	req *exchange.ExchangeRatesRequest,
) (*exchange.ExchangeRatesResponse, error) {
	baseStr, err := currencyToString(req.BaseCurrency)
	if err != nil {
		s.log.Warn("invalid base_currency in request",
			"base_currency", req.BaseCurrency.String(),
			"err", err,
		)
		return nil, status.Errorf(codes.InvalidArgument, "invalid base_currency: %v", err)
	}

	ratesModels, err := s.repo.GetRatesForBase(ctx, baseStr)
	if err != nil {
		s.log.Error("failed to get rates for base",
			"base", baseStr,
			"err", err,
		)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	ratesMap := make(map[string]float64, len(ratesModels)+1)

	for _, rm := range ratesModels {
		ratesMap[rm.TargetCurrency] = rm.Rate
	}

	ratesMap[baseStr] = 1.0

	resp := &exchange.ExchangeRatesResponse{
		BaseCurrency: req.BaseCurrency,
		Rates:        ratesMap,
	}

	s.log.Debug("rates returned for base",
		"base", baseStr,
		"count", len(ratesMap),
	)

	return resp, nil
}
