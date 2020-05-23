// Code generated by goa v3.1.1, DO NOT EDIT.
//
// RelayerAPI HTTP server
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package server

import (
	"context"
	"net/http"

	relayerapi "github.com/InjectiveLabs/injective-core/api/gen/relayer_api"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
)

// Server lists the RelayerAPI service endpoint HTTP handlers.
type Server struct {
	Mounts        []*MountPoint
	AssetPairs    http.Handler
	Orders        http.Handler
	OrderByHash   http.Handler
	Orderbook     http.Handler
	OrderConfig   http.Handler
	FeeRecipients http.Handler
	PostOrder     http.Handler
	CORS          http.Handler
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the design.
type ErrorNamer interface {
	ErrorName() string
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// New instantiates HTTP handlers for all the RelayerAPI service endpoints
// using the provided encoder and decoder. The handlers are mounted on the
// given mux using the HTTP verb and path defined in the design. errhandler is
// called whenever a response fails to be encoded. formatter is used to format
// errors returned by the service methods prior to encoding. Both errhandler
// and formatter are optional and can be nil.
func New(
	e *relayerapi.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"AssetPairs", "GET", "/api/sra/v3/asset_pairs"},
			{"Orders", "GET", "/api/sra/v3/orders"},
			{"OrderByHash", "GET", "/api/sra/v3/order/{orderHash}"},
			{"Orderbook", "GET", "/api/sra/v3/orderbook"},
			{"OrderConfig", "POST", "/api/sra/v3/order_config"},
			{"FeeRecipients", "GET", "/api/sra/v3/fee_recipients"},
			{"PostOrder", "POST", "/api/sra/v3/order"},
			{"CORS", "OPTIONS", "/api/sra/v3/asset_pairs"},
			{"CORS", "OPTIONS", "/api/sra/v3/orders"},
			{"CORS", "OPTIONS", "/api/sra/v3/order/{orderHash}"},
			{"CORS", "OPTIONS", "/api/sra/v3/orderbook"},
			{"CORS", "OPTIONS", "/api/sra/v3/order_config"},
			{"CORS", "OPTIONS", "/api/sra/v3/fee_recipients"},
			{"CORS", "OPTIONS", "/api/sra/v3/order"},
		},
		AssetPairs:    NewAssetPairsHandler(e.AssetPairs, mux, decoder, encoder, errhandler, formatter),
		Orders:        NewOrdersHandler(e.Orders, mux, decoder, encoder, errhandler, formatter),
		OrderByHash:   NewOrderByHashHandler(e.OrderByHash, mux, decoder, encoder, errhandler, formatter),
		Orderbook:     NewOrderbookHandler(e.Orderbook, mux, decoder, encoder, errhandler, formatter),
		OrderConfig:   NewOrderConfigHandler(e.OrderConfig, mux, decoder, encoder, errhandler, formatter),
		FeeRecipients: NewFeeRecipientsHandler(e.FeeRecipients, mux, decoder, encoder, errhandler, formatter),
		PostOrder:     NewPostOrderHandler(e.PostOrder, mux, decoder, encoder, errhandler, formatter),
		CORS:          NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "RelayerAPI" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.AssetPairs = m(s.AssetPairs)
	s.Orders = m(s.Orders)
	s.OrderByHash = m(s.OrderByHash)
	s.Orderbook = m(s.Orderbook)
	s.OrderConfig = m(s.OrderConfig)
	s.FeeRecipients = m(s.FeeRecipients)
	s.PostOrder = m(s.PostOrder)
	s.CORS = m(s.CORS)
}

// Mount configures the mux to serve the RelayerAPI endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountAssetPairsHandler(mux, h.AssetPairs)
	MountOrdersHandler(mux, h.Orders)
	MountOrderByHashHandler(mux, h.OrderByHash)
	MountOrderbookHandler(mux, h.Orderbook)
	MountOrderConfigHandler(mux, h.OrderConfig)
	MountFeeRecipientsHandler(mux, h.FeeRecipients)
	MountPostOrderHandler(mux, h.PostOrder)
	MountCORSHandler(mux, h.CORS)
}

// MountAssetPairsHandler configures the mux to serve the "RelayerAPI" service
// "assetPairs" endpoint.
func MountAssetPairsHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleRelayerAPIOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/sra/v3/asset_pairs", f)
}

// NewAssetPairsHandler creates a HTTP handler which loads the HTTP request and
// calls the "RelayerAPI" service "assetPairs" endpoint.
func NewAssetPairsHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeAssetPairsRequest(mux, decoder)
		encodeResponse = EncodeAssetPairsResponse(encoder)
		encodeError    = EncodeAssetPairsError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "assetPairs")
		ctx = context.WithValue(ctx, goa.ServiceKey, "RelayerAPI")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountOrdersHandler configures the mux to serve the "RelayerAPI" service
// "orders" endpoint.
func MountOrdersHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleRelayerAPIOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/sra/v3/orders", f)
}

// NewOrdersHandler creates a HTTP handler which loads the HTTP request and
// calls the "RelayerAPI" service "orders" endpoint.
func NewOrdersHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeOrdersRequest(mux, decoder)
		encodeResponse = EncodeOrdersResponse(encoder)
		encodeError    = EncodeOrdersError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "orders")
		ctx = context.WithValue(ctx, goa.ServiceKey, "RelayerAPI")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountOrderByHashHandler configures the mux to serve the "RelayerAPI" service
// "orderByHash" endpoint.
func MountOrderByHashHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleRelayerAPIOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/sra/v3/order/{orderHash}", f)
}

// NewOrderByHashHandler creates a HTTP handler which loads the HTTP request
// and calls the "RelayerAPI" service "orderByHash" endpoint.
func NewOrderByHashHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeOrderByHashRequest(mux, decoder)
		encodeResponse = EncodeOrderByHashResponse(encoder)
		encodeError    = EncodeOrderByHashError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "orderByHash")
		ctx = context.WithValue(ctx, goa.ServiceKey, "RelayerAPI")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountOrderbookHandler configures the mux to serve the "RelayerAPI" service
// "orderbook" endpoint.
func MountOrderbookHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleRelayerAPIOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/sra/v3/orderbook", f)
}

// NewOrderbookHandler creates a HTTP handler which loads the HTTP request and
// calls the "RelayerAPI" service "orderbook" endpoint.
func NewOrderbookHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeOrderbookRequest(mux, decoder)
		encodeResponse = EncodeOrderbookResponse(encoder)
		encodeError    = EncodeOrderbookError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "orderbook")
		ctx = context.WithValue(ctx, goa.ServiceKey, "RelayerAPI")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountOrderConfigHandler configures the mux to serve the "RelayerAPI" service
// "orderConfig" endpoint.
func MountOrderConfigHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleRelayerAPIOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/api/sra/v3/order_config", f)
}

// NewOrderConfigHandler creates a HTTP handler which loads the HTTP request
// and calls the "RelayerAPI" service "orderConfig" endpoint.
func NewOrderConfigHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeOrderConfigRequest(mux, decoder)
		encodeResponse = EncodeOrderConfigResponse(encoder)
		encodeError    = EncodeOrderConfigError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "orderConfig")
		ctx = context.WithValue(ctx, goa.ServiceKey, "RelayerAPI")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountFeeRecipientsHandler configures the mux to serve the "RelayerAPI"
// service "feeRecipients" endpoint.
func MountFeeRecipientsHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleRelayerAPIOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/sra/v3/fee_recipients", f)
}

// NewFeeRecipientsHandler creates a HTTP handler which loads the HTTP request
// and calls the "RelayerAPI" service "feeRecipients" endpoint.
func NewFeeRecipientsHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeFeeRecipientsRequest(mux, decoder)
		encodeResponse = EncodeFeeRecipientsResponse(encoder)
		encodeError    = EncodeFeeRecipientsError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "feeRecipients")
		ctx = context.WithValue(ctx, goa.ServiceKey, "RelayerAPI")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountPostOrderHandler configures the mux to serve the "RelayerAPI" service
// "postOrder" endpoint.
func MountPostOrderHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleRelayerAPIOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/api/sra/v3/order", f)
}

// NewPostOrderHandler creates a HTTP handler which loads the HTTP request and
// calls the "RelayerAPI" service "postOrder" endpoint.
func NewPostOrderHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodePostOrderRequest(mux, decoder)
		encodeResponse = EncodePostOrderResponse(encoder)
		encodeError    = EncodePostOrderError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "postOrder")
		ctx = context.WithValue(ctx, goa.ServiceKey, "RelayerAPI")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service RelayerAPI.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleRelayerAPIOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/api/sra/v3/asset_pairs", f)
	mux.Handle("OPTIONS", "/api/sra/v3/orders", f)
	mux.Handle("OPTIONS", "/api/sra/v3/order/{orderHash}", f)
	mux.Handle("OPTIONS", "/api/sra/v3/orderbook", f)
	mux.Handle("OPTIONS", "/api/sra/v3/order_config", f)
	mux.Handle("OPTIONS", "/api/sra/v3/fee_recipients", f)
	mux.Handle("OPTIONS", "/api/sra/v3/order", f)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// handleRelayerAPIOrigin applies the CORS response headers corresponding to
// the origin for the service RelayerAPI.
func handleRelayerAPIOrigin(h http.Handler) http.Handler {
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "*") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
