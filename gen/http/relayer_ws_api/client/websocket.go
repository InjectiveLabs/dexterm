// Code generated by goa v3.1.1, DO NOT EDIT.
//
// RelayerWsAPI WebSocket client streaming
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package client

import (
	"io"

	relayerwsapi "github.com/InjectiveLabs/injective-core/api/gen/relayer_ws_api"
	"github.com/gorilla/websocket"
	goahttp "goa.design/goa/v3/http"
)

// ConnConfigurer holds the websocket connection configurer functions for the
// streaming endpoints in "RelayerWsAPI" service.
type ConnConfigurer struct {
	OrdersStreamingFn goahttp.ConnConfigureFunc
}

// OrdersStreamingClientStream implements the
// relayerwsapi.OrdersStreamingClientStream interface.
type OrdersStreamingClientStream struct {
	// conn is the underlying websocket connection.
	conn *websocket.Conn
}

// NewConnConfigurer initializes the websocket connection configurer function
// with fn for all the streaming endpoints in "RelayerWsAPI" service.
func NewConnConfigurer(fn goahttp.ConnConfigureFunc) *ConnConfigurer {
	return &ConnConfigurer{
		OrdersStreamingFn: fn,
	}
}

// Recv reads instances of "relayerwsapi.OrdersStreamingResult" from the
// "ordersStreaming" endpoint websocket connection.
func (s *OrdersStreamingClientStream) Recv() (*relayerwsapi.OrdersStreamingResult, error) {
	var (
		rv   *relayerwsapi.OrdersStreamingResult
		body OrdersStreamingResponseBody
		err  error
	)
	err = s.conn.ReadJSON(&body)
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		s.conn.Close()
		return rv, io.EOF
	}
	if err != nil {
		return rv, err
	}
	err = ValidateOrdersStreamingResponseBody(&body)
	if err != nil {
		return rv, err
	}
	res := NewOrdersStreamingResultOK(&body)
	return res, nil
}
