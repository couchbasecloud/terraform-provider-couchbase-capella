package appservice

// PublicAccess is the public access mode of an App Service.
type PublicAccess string

const (
	PublicAccessUnrestricted     PublicAccess = "unrestricted"
	PublicAccessRestricting      PublicAccess = "restricting"
	PublicAccessRestricted       PublicAccess = "restricted"
	PublicAccessUnrestricting    PublicAccess = "unrestricting"
	PublicAccessRestrictFailed   PublicAccess = "restrictFailed"
	PublicAccessUnrestrictFailed PublicAccess = "unrestrictFailed"
)

// Mode returns the mode that the last public access change requested: restricted or unrestricted.
func (p PublicAccess) Mode() PublicAccess {
	switch p {
	case PublicAccessRestricting, PublicAccessRestricted, PublicAccessRestrictFailed:
		return PublicAccessRestricted
	default:
		return PublicAccessUnrestricted
	}
}

// IsFailed reports whether the last public access change failed.
func (p PublicAccess) IsFailed() bool {
	return p == PublicAccessRestrictFailed || p == PublicAccessUnrestrictFailed
}
