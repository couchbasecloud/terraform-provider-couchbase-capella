package acceptance_tests

func clusterResourceName() string {
	return providerName + "_cluster"
}

func allowlistResourceName() string {
	return providerName + "_allowlist"
}

func appServiceResourceName() string {
	return providerName + "_app_service"
}

func organizationResourceName() string {
	return providerName + "_organization"
}
