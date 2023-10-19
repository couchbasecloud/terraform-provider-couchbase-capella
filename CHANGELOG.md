# 1.0.0 (Unreleased)

## October 18 2023
UPGRADE NOTES:
- Upgraded provider to use Golang 1.20+ and Terraform 1.5+
- Upgraded terraform-plugin-framework to 1.4.1 

BUG FIXES:
- Update cluster will no longer fail for PoC organization with developer pro plan.
- User removed from organization from Capella UI won't break the terraform provider.

ENHANCEMENTS:
- Detailed examples and README with sample outputs