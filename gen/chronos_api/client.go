// Code generated by goa v3.1.1, DO NOT EDIT.
//
// ChronosAPI client
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package chronosapi

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "ChronosAPI" service client.
type Client struct {
	SymbolInfoEndpoint           goa.Endpoint
	HistoryEndpoint              goa.Endpoint
	FillsHistoryEndpoint         goa.Endpoint
	MarketSummaryEndpoint        goa.Endpoint
	FuturesHistoryEndpoint       goa.Endpoint
	FuturesFillsHistoryEndpoint  goa.Endpoint
	FuturesMarketSummaryEndpoint goa.Endpoint
}

// NewClient initializes a "ChronosAPI" service client given the endpoints.
func NewClient(symbolInfo, history, fillsHistory, marketSummary, futuresHistory, futuresFillsHistory, futuresMarketSummary goa.Endpoint) *Client {
	return &Client{
		SymbolInfoEndpoint:           symbolInfo,
		HistoryEndpoint:              history,
		FillsHistoryEndpoint:         fillsHistory,
		MarketSummaryEndpoint:        marketSummary,
		FuturesHistoryEndpoint:       futuresHistory,
		FuturesFillsHistoryEndpoint:  futuresFillsHistory,
		FuturesMarketSummaryEndpoint: futuresMarketSummary,
	}
}

// SymbolInfo calls the "symbolInfo" endpoint of the "ChronosAPI" service.
func (c *Client) SymbolInfo(ctx context.Context, p *SymbolInfoPayload) (res *TradingViewSymbolInfoResponse, err error) {
	var ires interface{}
	ires, err = c.SymbolInfoEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*TradingViewSymbolInfoResponse), nil
}

// History calls the "history" endpoint of the "ChronosAPI" service.
func (c *Client) History(ctx context.Context, p *HistoryPayload) (res *HistoryResponse, err error) {
	var ires interface{}
	ires, err = c.HistoryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*HistoryResponse), nil
}

// FillsHistory calls the "fillsHistory" endpoint of the "ChronosAPI" service.
func (c *Client) FillsHistory(ctx context.Context, p *FillsHistoryPayload) (res []*FillEvent, err error) {
	var ires interface{}
	ires, err = c.FillsHistoryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.([]*FillEvent), nil
}

// MarketSummary calls the "marketSummary" endpoint of the "ChronosAPI" service.
func (c *Client) MarketSummary(ctx context.Context, p *MarketSummaryPayload) (res *MarketSummaryResponse, err error) {
	var ires interface{}
	ires, err = c.MarketSummaryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*MarketSummaryResponse), nil
}

// FuturesHistory calls the "futuresHistory" endpoint of the "ChronosAPI"
// service.
func (c *Client) FuturesHistory(ctx context.Context, p *FuturesHistoryPayload) (res *FuturesHistoryResponse, err error) {
	var ires interface{}
	ires, err = c.FuturesHistoryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*FuturesHistoryResponse), nil
}

// FuturesFillsHistory calls the "futuresFillsHistory" endpoint of the
// "ChronosAPI" service.
func (c *Client) FuturesFillsHistory(ctx context.Context, p *FuturesFillsHistoryPayload) (res []*FuturesFillEvent, err error) {
	var ires interface{}
	ires, err = c.FuturesFillsHistoryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.([]*FuturesFillEvent), nil
}

// FuturesMarketSummary calls the "futuresMarketSummary" endpoint of the
// "ChronosAPI" service.
func (c *Client) FuturesMarketSummary(ctx context.Context, p *FuturesMarketSummaryPayload) (res *FuturesMarketSummaryResponse, err error) {
	var ires interface{}
	ires, err = c.FuturesMarketSummaryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*FuturesMarketSummaryResponse), nil
}