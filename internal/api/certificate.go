package api

// GetCertificateResponse is the response received from the Capella V4 Public API when asked to get the certificate details for a Capella cluster.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// Project Viewer
// Database Data Reader/Writer
// Database Data Reader
//
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetCertificateResponse struct {
	// Certificate is the certificate of the capella cluster
	Certificate string `json:"certificate"`
}
