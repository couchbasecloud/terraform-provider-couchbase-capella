package api

type GetCertificateResponse struct {
	Certificate string `json:"certificate"`
}

// GetCertificatesResponse defines the model for a GetCertificatesResponse.
type GetCertificatesResponse struct {
	Data GetCertificateResponse `json:"data"`
}
