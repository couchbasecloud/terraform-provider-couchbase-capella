# Changelog

## [v1.4.0](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.4.0) (2024-12-09)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.3.0...v1.4.0)

**Implemented enhancements:**
- \[AV-90715\] Add support for `zones` in Cluster [\#241](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/241) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-82735\] Add support for Flush Bucket Data [\#234](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/234) ([Lagher0](https://github.com/Lagher0))
- \[AV-76498\] Add support for GSI Index Management [\#233](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/233) ([l0n3star](https://github.com/l0n3star))
- \[AV-87139\] Update gorunner version [\#229](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/229)([a-atri](https://github.com/a-atri))
- \[AV-78889\] Add support for Azure VNET Peering [\#216](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/216) ([PaulomeeCb](https://github.com/PaulomeeCb))

**Fixed bugs:**
- \[AV-90385\] Add check for unsupported `storage` and `IOPS` values in case of `Azure Premium Disk` [\#240](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/240) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-87077\] Resolve false positives in acceptance tests by correcting the handling of computed values [\#227](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/202) ([l0n3star](https://github.com/l0n3star))
- \[AV-86105\] Initialize `autoexpansion` field with a null value [\#226](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/226) ([aniket-Kumar-c](https://github.com/aniket-Kumar-c))



## [v1.3.0](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.3.0) (2024-09-11)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.2.1...v1.3.0)

**Implemented enhancements:**
- \[AV-86845\] Deprecate the `configurationType` attribute in the cluster resource [\#222](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/222) ([aniket-Kumar-c](https://github.com/aniket-Kumar-c))
- \[AV-85326\] Display the cluster connection string for the cluster resource [\#219](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/219) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-82517\] Remove the configuration type from the cluster resource [#207](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/207) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-76499\] Support Capella System Events and Activity Logs [\#206](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/206) ([aniket-Kumar-c](https://github.com/aniket-Kumar-c))
- \[AV-78888\] Add Support for VPC peering for AWS and GCP [\#205](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/205) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-76500\] Support for Private Endpoints for AWS and Azure [\#202](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/202) ([l0n3star](https://github.com/l0n3star))
- \[AV-79411\] Add a Custom User-Agent HTTP header for Terraform Provider Client [\#200](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/200) ([l0n3star](https://github.com/l0n3star))

**Fixed bugs:**
- \[AV-85326\] Fix Private Endpoints examples [\#218](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/218) ([l0n3star](https://github.com/l0n3star))
- \[AV-84160\] Make cluster timezone optional [\#215](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/215) ([Lagher0](https://github.com/Lagher0))
- \[AV-83067\] Remove server version check for acceptance tests [\#210](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/210) ([l0n3star](https://github.com/l0n3star))
- \[AV-83066\] Resolve getting-started folder bugs [\#209](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/209) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-79969\] Use noop when resource is destroyed for audit logs [\#203](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/203) ([l0n3star](https://github.com/l0n3star))

**Documentation Enhancement:**
- \[AV-83061\] Update contributing.md with windows steps [#208](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/208) ([Talina06](https://github.com/Talina06))

**Merged pull requests:**
- Bump github.com/hashicorp/terraform-plugin-framework from 1.6.1 to 1.9.0 [#201](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/201) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-go from 0.22.1 to 0.23.0 [#193](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/193) ([dependabot[bot]](https://github.com/apps/dependabot))


## [v1.2.1](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.2.1) (2024-05-07)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.2.0...v1.2.1)

**Fixed bugs:**

- \[AV-77464\] fix export schema [\#191](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/191) ([l0n3star](https://github.com/l0n3star))

## [v1.2.0](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.2.0) (2024-04-25)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.1.0...v1.2.0)

**Implemented enhancements:**

- \[AV-73782\] Couchbase Server Audit Events Support [\#158](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/158) ([l0n3star](https://github.com/l0n3star))

## [v1.1.0](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.1.0) (2024-04-02)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.0.0...v1.1.0)

**Implemented enhancements:**

- \[AV-75849\] Cluster and App service On Demand On/Off [\#181](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/181) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-74366\] Bucket Collections [\#163](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/163) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-73526\] Add support for rate-limiting retries [\#157](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/157) ([ajsqr](https://github.com/ajsqr))
- \[AV-70846\] Import sample buckets [\#156](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/156) ([Lagher0](https://github.com/Lagher0))
- \[AV-73296\] Bucket Scopes [\#153](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/153) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-71279\] Support autoexpansion for Azure cluster\(s\) [\#143](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/143) ([l0n3star](https://github.com/l0n3star))

**Fixed bugs:**

- \[AV-70854\] Fixed optional fields during cluster creation | Azure & GCP [\#141](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/141) ([nidhi07kumar](https://github.com/nidhi07kumar))

**Merged pull requests:**

- Bump github.com/hashicorp/terraform-plugin-go from 0.21.0 to 0.22.1 [\#165](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/165) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-testing from 1.6.0 to 1.7.0 [\#162](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/162) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-framework from 1.5.0 to 1.6.1 [\#161](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/161) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/stretchr/testify from 1.8.4 to 1.9.0 [\#160](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/160) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-go from 0.20.0 to 0.21.0 [\#148](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/148) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/google/uuid from 1.5.0 to 1.6.0 [\#145](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/145) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-framework from 1.4.2 to 1.5.0 [\#142](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/142) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/couchbase/tools-common/functional from 1.1.1 to 1.2.0 [\#140](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/140) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/cloudflare/circl from 1.3.3 to 1.3.7 [\#139](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/139) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-go from 0.19.1 to 0.20.0 [\#138](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/138) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump golang.org/x/crypto from 0.16.0 to 0.17.0 [\#136](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/136) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/google/uuid from 1.4.0 to 1.5.0 [\#132](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/132) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.0.0](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.0.0) (2023-11-10)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/d9a774ce9a0731bd15a6ca9eb9b7ea4d7f4e1d33...v1.0.0)

**Merged pull requests:**

- \[AV-67237\] Code Health - Create Acceptance Test Directory [\#87](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/87) ([matty271828](https://github.com/matty271828))
- \[AV-64820\] App Services Terraform Provider Feature Branch [\#84](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/84) ([nidhi07kumar](https://github.com/nidhi07kumar))
- \[AV-67060\] PP Bugs - Pass in sort by parameter  [\#83](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/83) ([matty271828](https://github.com/matty271828))
- \[AV-62009\] Code Health - Pass in Success Codes to Client [\#80](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/80) ([matty271828](https://github.com/matty271828))
- \[AV-66360\] PP Bugs - Make resource object optional [\#76](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/76) ([matty271828](https://github.com/matty271828))
- \[AV-65775\] List APIs - Database Credentials & Buckets Pagination [\#74](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/74) ([matty271828](https://github.com/matty271828))
- \[AV-65679\] List APIs -  Users, ApiKeys & Clusters Pagination [\#73](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/73) ([matty271828](https://github.com/matty271828))
- \[AV-65771\] List APIs - Projects pagination [\#72](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/72) ([matty271828](https://github.com/matty271828))
- \[AV-65452\] CREATE bucket with optional fields [\#68](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/68) ([nidhi07kumar](https://github.com/nidhi07kumar))
- AV-66047 Add PR template [\#67](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/67) ([rijuCB](https://github.com/rijuCB))
- Bump google.golang.org/grpc from 1.57.0 to 1.57.1 [\#66](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/66) ([dependabot[bot]](https://github.com/apps/dependabot))
- \[AV-65773\] Allowlists Datasource - Pagination [\#65](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/65) ([matty271828](https://github.com/matty271828))
- \[AV-65899\] PP Bugs - Move project roles into same array [\#64](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/64) ([matty271828](https://github.com/matty271828))
- \[AV-63471\] User Resource - Update user using patch request [\#62](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/62) ([matty271828](https://github.com/matty271828))
- \[AV-65088\] Better access example for creating db user [\#61](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/61) ([Talina06](https://github.com/Talina06))
- \[AV-65076\] Fixed template.tfvars files for examples [\#60](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/60) ([nidhi07kumar](https://github.com/nidhi07kumar))
- \[AV-64976\] User refresh state for users deleted from UI [\#59](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/59) ([nidhi07kumar](https://github.com/nidhi07kumar))
- Fix added for various resources creation failing in main branch [\#58](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/58) ([aniket-Kumar-c](https://github.com/aniket-Kumar-c))
- AV-65056: acceptance test poc for project [\#54](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/54) ([aniket-Kumar-c](https://github.com/aniket-Kumar-c))
- Bump golang.org/x/net from 0.13.0 to 0.17.0 [\#51](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/51) ([dependabot[bot]](https://github.com/apps/dependabot))
- \[AV-64572\] Code Health - Extract Schema State Validation [\#50](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/50) ([matty271828](https://github.com/matty271828))



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
