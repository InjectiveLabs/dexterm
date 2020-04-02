// Code generated by goa v3.1.1, DO NOT EDIT.
//
// HTTP request path constructors for the RelayerAPI service.
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package client

import (
	"fmt"
)

// AssetPairsRelayerAPIPath returns the URL path to the RelayerAPI service assetPairs HTTP endpoint.
func AssetPairsRelayerAPIPath() string {
	return "/api/sra/v3/asset_pairs"
}

// OrdersRelayerAPIPath returns the URL path to the RelayerAPI service orders HTTP endpoint.
func OrdersRelayerAPIPath() string {
	return "/api/sra/v3/orders"
}

// OrderByHashRelayerAPIPath returns the URL path to the RelayerAPI service orderByHash HTTP endpoint.
func OrderByHashRelayerAPIPath(orderHash string) string {
	return fmt.Sprintf("/api/sra/v3/order/%v", orderHash)
}

// OrderbookRelayerAPIPath returns the URL path to the RelayerAPI service orderbook HTTP endpoint.
func OrderbookRelayerAPIPath() string {
	return "/api/sra/v3/orderbook"
}

// OrderConfigRelayerAPIPath returns the URL path to the RelayerAPI service orderConfig HTTP endpoint.
func OrderConfigRelayerAPIPath() string {
	return "/api/sra/v3/order_config"
}

// FeeRecipientsRelayerAPIPath returns the URL path to the RelayerAPI service feeRecipients HTTP endpoint.
func FeeRecipientsRelayerAPIPath() string {
	return "/api/sra/v3/fee_recipients"
}

// PostOrderRelayerAPIPath returns the URL path to the RelayerAPI service postOrder HTTP endpoint.
func PostOrderRelayerAPIPath() string {
	return "/api/sra/v3/order"
}