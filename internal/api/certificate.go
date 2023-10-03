package api

// GetCertificateResponse defines model for GetCertificateResponse.
type GetCertificateResponse struct {
	// Certificate is the certificate of the capella cluster
	Certificate string `json:"certificate"`
}
