package api

type CreateAllowedCIDRRequest struct {
	Cidr string `json:"cidr"`

	Comment string `json:"comment,omitempty"`

	ExpiresAt string `json:"expiresAt"`
}

type AppServiceAllowedCIDRResponse struct {
	Id string `json:"id"`

	Cidr string `json:"cidr"`

	Comment string `json:"comment,omitempty"`

	ExpiresAt string `json:"expiresAt"`

	Status string `json:"status"`

	Type string `json:"type"`

	Audit CouchbaseAuditData `json:"audit"`
}
type ListAppServiceAllowedCIDRResponse struct {
	Data []AppServiceAllowedCIDRResponse `json:"data"`
}
