package api

import (
	"net/url"
)

const (
	// AllowedCIDRAccessModeAdmin applies an allowed CIDR to the App Service Admin and Metrics APIs.
	AllowedCIDRAccessModeAdmin = "admin"

	// AllowedCIDRAccessModePublic applies an allowed CIDR to the App Service Public API.
	AllowedCIDRAccessModePublic = "public"
)

// CreateAllowedCIDRRequest is the request body for creating an allowed CIDR on an App Service.
type CreateAllowedCIDRRequest struct {
	// Cidr is the ip address range in CIDR notation that is allowed to access the App Service.
	Cidr string `json:"cidr"`

	// Comment is an optional comment for the allowed CIDR.
	Comment string `json:"comment,omitempty"`

	//ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt string `json:"expiresAt,omitempty"`

	// AccessMode is the App Service API that the allowed CIDR applies to. The API defaults to admin.
	AccessMode string `json:"accessMode,omitempty"`
}

// AppServiceAllowedCIDRResponse is the response structure for an individual allowed CIDR entry on an App Service.
type AppServiceAllowedCIDRResponse struct {
	// Id is the UUID generated when an allowed cidr is created.
	Id string `json:"id"`

	// Cidr is the ip address range in CIDR notation that is allowed to access the App Service.
	Cidr string `json:"cidr"`

	// Comment is an optional comment for the allowed CIDR.
	Comment string `json:"comment,omitempty"`

	//ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt string `json:"expiresAt"`

	// Status indicates the status of the allowed CIDR.
	Status string `json:"status"`

	// Type indicates the type of the allowed CIDR, e.g., "temporary" or "permanent".
	Type string `json:"type"`

	// Audit contains the audit information for the App service CIDR.
	Audit CouchbaseAuditData `json:"audit"`

	// AccessMode is the App Service API that the allowed CIDR applies to.
	AccessMode string `json:"accessMode"`
}

// AllowedCIDRsAccessModeQuery returns the query string that lists the allowed CIDRs with the given access modes. The
// API lists admin allowed CIDRs when the query has no access mode.
func AllowedCIDRsAccessModeQuery(accessModes []string) string {
	if len(accessModes) == 0 {
		return ""
	}
	return "?" + url.Values{"accessMode": accessModes}.Encode()
}

// ListAppServiceAllowedCIDRResponse is the response structure for listing allowed CIDRs on an App Service.
type ListAppServiceAllowedCIDRResponse struct {
	// Data contains the list of allowed CIDRs on the App Service.
	Data []AppServiceAllowedCIDRResponse `json:"data"`
}
