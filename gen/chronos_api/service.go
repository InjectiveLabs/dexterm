// Code generated by goa v3.1.1, DO NOT EDIT.
//
// ChronosAPI service
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package chronosapi

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// ChronosAPI implements historical data API for e.g. TradingView.
type Service interface {
	// Get a list of all instruments.
	SymbolInfo(context.Context, *SymbolInfoPayload) (res *TradingViewSymbolInfoResponse, err error)
	// Request for history bars. Each property of the response object is treated as
	// a table column.
	History(context.Context, *HistoryPayload) (res *HistoryResponse, err error)
	// Get history of past fill events, filtered by trade pair name
	FillsHistory(context.Context, *FillsHistoryPayload) (res []*FillEvent, err error)
	// Gets market summary for the latest interval (hour, day, month)
	MarketSummary(context.Context, *MarketSummaryPayload) (res *MarketSummaryResponse, err error)
	// Request for futures asset prices history bars. Each property of the response
	// object is treated as a table column.
	FuturesHistory(context.Context, *FuturesHistoryPayload) (res *FuturesHistoryResponse, err error)
	// Get history of past fill events, filtered by trade pair name
	FuturesFillsHistory(context.Context, *FuturesFillsHistoryPayload) (res []*FuturesFillEvent, err error)
	// Gets futures market summary for the latest interval (hour, day, month)
	FuturesMarketSummary(context.Context, *FuturesMarketSummaryPayload) (res *FuturesMarketSummaryResponse, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "ChronosAPI"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [7]string{"symbolInfo", "history", "fillsHistory", "marketSummary", "futuresHistory", "futuresFillsHistory", "futuresMarketSummary"}

// SymbolInfoPayload is the payload type of the ChronosAPI service symbolInfo
// method.
type SymbolInfoPayload struct {
	// ID of a symbol group. It is only required if you use groups of symbols to
	// restrict access to instrument's data.
	Group *string
}

// TradingViewSymbolInfoResponse is the result type of the ChronosAPI service
// symbolInfo method.
type TradingViewSymbolInfoResponse struct {
	// Status of the response.
	S string
	// Error message.
	Errmsg *string
	// This is the name of the symbol - a string that the users will see. It should
	// contain uppercase letters, numbers, a dot or an underscore. Also, it will be
	// used for data requests if you are not using tickers.
	Symbol []string
	// Description of a symbol. Will be displayed in the chart legend for this
	// symbol.
	Description []string
	// Symbol currency, also named as counter currency. If a symbol is a currency
	// pair, then the currency field has to contain the second currency of this
	// pair. For example, USD is a currency for EURUSD ticker. Fiat currency must
	// meet the ISO 4217 standard. The default value is null.
	Currency []string
	// Short name of exchange where this symbol is listed.
	ExchangeListed []string
	// Short name of exchange where this symbol is traded.
	ExchangeTraded []string
	// Minimal integer price change.
	Minmovement []int
	// Indicates how many decimal points the price has. For example, if the price
	// has 2 decimal points (ex., 300.01), then pricescale is 100. If it has 3
	// decimals, then pricescale is 1000 etc. If the price doesn't have decimals,
	// set pricescale to 1
	Pricescale []int
	// Timezone of the exchange for this symbol. We expect to get the name of the
	// time zone in olsondb format.
	Timezone []string
	// Symbol type (forex/stock, crypto etc.).
	Type []string
	// Bitcoin and other cryptocurrencies: the session string should be 24x7
	SessionRegular []string
	// For currency pairs only. This field contains the first currency of the pair.
	// For example, base currency for EURUSD ticker is EUR. Fiat currency must meet
	// the ISO 4217 standard.
	BaseCurrency []string
	// This is a number for complex price formatting cases.
	Minmov2 []int
	// Boolean showing whether this symbol wants to have complex price formatting
	// (see minmov2) or not. The default value is false.
	Fractional []bool
	// Root of the features. It's required for futures symbol types only. Provide a
	// null value for other symbol types. The default value is null.
	Root []string
	// Short description of the futures root that will be displayed in the symbol
	// search. It's required for futures only. Provide a null value for other
	// symbol types. The default value is null.
	RootDescription []string
	// Boolean value showing whether the symbol includes intraday (minutes)
	// historical data.
	HasIntraday []bool
	// Boolean showing whether the symbol includes volume data or not. The default
	// value is false.
	HasNoVolume []bool
	// Boolean value showing whether the symbol is CFD. The base instrument type is
	// set using the type field.
	IsCfd []bool
	// This is a unique identifier for this particular symbol in your symbology. If
	// you specify this property then its value will be used for all data requests
	// for this symbol.
	Ticker []string
	// The boolean value showing whether data feed has its own daily resolution
	// bars or not.
	HasDaily []bool
	// This is an array containing intraday resolutions (in minutes) that the data
	// feed may provide
	IntradayMultipliers []string
	// The boolean value showing whether data feed has its own weekly and monthly
	// resolution bars or not.
	HasWeeklyAndMonthly []bool
	// The currency value of a single whole unit price change in the instrument's
	// currency. If the value is not provided it is assumed to be 1.
	Pointvalue []int
	// Expiration of the futures in the following format: YYYYMMDD. Required for
	// futures type symbols only.
	Expiration []int
	// The principle of building bars. The default value is trade.
	BarSource []string
	// The principle of bar alignment. The default value is none.
	BarTransform []string
	// Is used to create the zero-volume bars in the absence of any trades
	BarFillgaps []bool
}

// HistoryPayload is the payload type of the ChronosAPI service history method.
type HistoryPayload struct {
	// Symbol name or ticker.
	Symbol string
	// Symbol resolution. Possible resolutions are daily (D or 1D, 2D ... ), weekly
	// (1W, 2W ...), monthly (1M, 2M...) and an intra-day resolution – minutes(1, 2
	// ...).
	Resolution string
	// Unix timestamp (UTC) of the leftmost required bar, including from
	From *int
	// Unix timestamp (UTC) of the rightmost required bar, including to. It can be
	// in the future. In this case, the rightmost required bar is the latest
	// available bar.
	To int
	// Number of bars (higher priority than from) starting with to. If countback is
	// set, from should be ignored.
	Countback *int
}

// HistoryResponse is the result type of the ChronosAPI service history method.
type HistoryResponse struct {
	// Status of the response.
	S string
	// Error message.
	Errmsg *string
	// Unix time of the next bar if there is no data in the requested period
	// (optional).
	Nb *int
	// Bar time, Unix timestamp (UTC). Daily bars should only have the date part,
	// time should be 0.
	T []int
	// Open price.
	O []float64
	// High price.
	H []float64
	// Low price.
	L []float64
	// Close price.
	C []float64
	// Volume.
	V []float64
}

// FillsHistoryPayload is the payload type of the ChronosAPI service
// fillsHistory method.
type FillsHistoryPayload struct {
	// Account address to get related fill events
	Account *string
	// Trade pair name
	TradePair string
}

// MarketSummaryPayload is the payload type of the ChronosAPI service
// marketSummary method.
type MarketSummaryPayload struct {
	// Trade pair name
	TradePair string
	// Specify the resolution
	Resolution string
}

// MarketSummaryResponse is the result type of the ChronosAPI service
// marketSummary method.
type MarketSummaryResponse struct {
	// Open price.
	Open float64
	// High price.
	High float64
	// Low price.
	Low float64
	// Volume.
	Volume float64
	// Current price based on latest fill event.
	Price float64
	// Change percent from the previous period on the same resolution.
	Change float64
}

// FuturesHistoryPayload is the payload type of the ChronosAPI service
// futuresHistory method.
type FuturesHistoryPayload struct {
	// ID of the derivative market
	MarketID string
	// Symbol resolution. Possible resolutions are daily (D or 1D, 2D ... ), weekly
	// (1W, 2W ...), monthly (1M, 2M...) and an intra-day resolution – minutes(1, 2
	// ...).
	Resolution string
	// Unix timestamp (UTC) of the leftmost required bar, including from
	From *int
	// Unix timestamp (UTC) of the rightmost required bar, including to. It can be
	// in the future. In this case, the rightmost required bar is the latest
	// available bar.
	To int
	// Number of bars (higher priority than from) starting with to. If countback is
	// set, from should be ignored.
	Countback *int
}

// FuturesHistoryResponse is the result type of the ChronosAPI service
// futuresHistory method.
type FuturesHistoryResponse struct {
	// Status of the response.
	S string
	// Error message.
	Errmsg *string
	// Unix time of the next bar if there is no data in the requested period
	// (optional).
	Nb *int
	// Bar time, Unix timestamp (UTC). Daily bars should only have the date part,
	// time should be 0.
	T []int
	// Open price.
	O []float64
	// High price.
	H []float64
	// Low price.
	L []float64
	// Close price.
	C []float64
	// Volume.
	V []float64
}

// FuturesFillsHistoryPayload is the payload type of the ChronosAPI service
// futuresFillsHistory method.
type FuturesFillsHistoryPayload struct {
	// Account address to get related fill events
	Account *string
	// Market ID of the futures pair
	MarketID string
}

// FuturesMarketSummaryPayload is the payload type of the ChronosAPI service
// futuresMarketSummary method.
type FuturesMarketSummaryPayload struct {
	// Market ID of the futures pair
	MarketID string
	// Specify the resolution
	Resolution string
}

// FuturesMarketSummaryResponse is the result type of the ChronosAPI service
// futuresMarketSummary method.
type FuturesMarketSummaryResponse struct {
	// Open price.
	Open float64
	// High price.
	High float64
	// Low price.
	Low float64
	// Volume.
	Volume float64
	// Current price based on latest fill event.
	Price float64
	// Change percent from the previous period on the same resolution.
	Change float64
}

type FillEvent struct {
	// Account's side in the trade
	Side string
	// UNIX timestamp of the fill event
	Ts int64
	// Filled amount in quote currency
	Size float64
	// Filled amount in base currency
	Filled float64
	// Price in quote currency
	Price float64
	// Transaction hash related to this fill
	TxHash *string
}

type FuturesFillEvent struct {
	// Account's side in the trade
	Side string
	// UNIX timestamp of the fill event
	Ts int64
	// Filled amount in quote currency
	Size float64
	// Filled amount in base currency
	Filled float64
	// Price in quote currency
	Price float64
	// Transaction hash related to this fill
	TxHash *string
}

// MakeBadRequest builds a goa.ServiceError from an error.
func MakeBadRequest(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "bad_request",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "not_found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeInternal builds a goa.ServiceError from an error.
func MakeInternal(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "internal",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
