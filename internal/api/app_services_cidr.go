package api

type CreateAllowedCIDRRequest struct {
	Cidr string `json:"cidr"`

	Comment string `json:"comment,omitempty"`

	ExpiresAt string `json:"expiresAt"`
}
