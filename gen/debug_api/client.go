// Code generated by goa v3.1.1, DO NOT EDIT.
//
// DebugAPI client
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package debugapi

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "DebugAPI" service client.
type Client struct {
	VersionEndpoint goa.Endpoint
}

// NewClient initializes a "DebugAPI" service client given the endpoints.
func NewClient(version goa.Endpoint) *Client {
	return &Client{
		VersionEndpoint: version,
	}
}

// Version calls the "version" endpoint of the "DebugAPI" service.
func (c *Client) Version(ctx context.Context) (res *VersionResult, err error) {
	var ires interface{}
	ires, err = c.VersionEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(*VersionResult), nil
}