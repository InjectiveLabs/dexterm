// Code generated by goa v3.1.1, DO NOT EDIT.
//
// DebugAPI HTTP client types
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-core/api/design -o ../

package client

import (
	debugapi "github.com/InjectiveLabs/dexterm/gen/debug_api"
	goa "goa.design/goa/v3/pkg"
)

// VersionResponseBody is the type of the "DebugAPI" service "version" endpoint
// HTTP response body.
type VersionResponseBody struct {
	// Relayerd code version.
	Version *string `form:"version,omitempty" json:"version,omitempty" xml:"version,omitempty"`
	// Additional meta data.
	MetaData map[string]string `form:"metaData,omitempty" json:"metaData,omitempty" xml:"metaData,omitempty"`
}

// NewVersionResultOK builds a "DebugAPI" service "version" endpoint result
// from a HTTP "OK" response.
func NewVersionResultOK(body *VersionResponseBody) *debugapi.VersionResult {
	v := &debugapi.VersionResult{
		Version: *body.Version,
	}
	if body.MetaData != nil {
		v.MetaData = make(map[string]string, len(body.MetaData))
		for key, val := range body.MetaData {
			tk := key
			tv := val
			v.MetaData[tk] = tv
		}
	}

	return v
}

// ValidateVersionResponseBody runs the validations defined on
// VersionResponseBody
func ValidateVersionResponseBody(body *VersionResponseBody) (err error) {
	if body.Version == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("version", "body"))
	}
	return
}
