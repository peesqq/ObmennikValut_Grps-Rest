package service

import (
	"context"
	"fmt"
	"github.com/peesqq/proto-exchange/proto"
)

type ExchangeService struct {
	proto.UnimplementedExchangeServiceServer
}

func (s *ExchangeService) GetExchangeRates(ctx context.Context, in *proto.Empty) (*proto.ExchangeRatesResponse, error) {
	rates := map[string]float32{
		"USD": 1.0,
		"EUR": 0.85,
		"RUB": 75.0,
	}
	return &proto.ExchangeRatesResponse{Rates: rates}, nil
}

func (s *ExchangeService) GetExchangeRateForCurrency(ctx context.Context, in *proto.CurrencyRequest) (*proto.ExchangeRateResponse, error) {
	rate := float32(0.85)
	return &proto.ExchangeRateResponse{
		FromCurrency: in.FromCurrency,
		ToCurrency:   in.ToCurrency,
		Rate:         rate,
	}, nil
}
func (s *ExchangeService) ConvertCurrency(ctx context.Context, req *proto.ConvertCurrencyRequest) (*proto.ConvertCurrencyResponse, error) {
	rates := map[string]float32{
		"USD_RUB": 102.0,
		"RUB_USD": 0.009,
		"USD_EUR": 0.97,
		"EUR_USD": 1.02,
	}

	key := fmt.Sprintf("%s_%s", req.FromCurrency, req.ToCurrency)
	rate, exists := rates[key]
	if !exists {
		return nil, fmt.Errorf("exchange rate not found for %s to %s", req.FromCurrency, req.ToCurrency)
	}

	convertedAmount := req.Amount * rate
	return &proto.ConvertCurrencyResponse{ConvertedAmount: convertedAmount}, nil
}
