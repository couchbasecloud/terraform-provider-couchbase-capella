# Changelog

## [v1.5.3](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.5.3) (2025-08-06)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.5.2...v1.5.3)

**Fixed bugs:**

- \[AV-107267\] Restrict setting IOPS for GCP clusters [\#369](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/369) ([l0n3star](https://github.com/l0n3star))
- \[AV-107190\] Return API errors [\#367](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/367) ([l0n3star](https://github.com/l0n3star))
- \[AV-107162\] Do not allow setting autoexpansion for AWS/GCP [\#366](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/366) ([l0n3star](https://github.com/l0n3star))

## [v1.5.3](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.5.3) (2025-07-31)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.5.2...v1.5.3)

**Fixed bugs:**

- \[AV-107190\] Return API errors [\#367](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/367) ([l0n3star](https://github.com/l0n3star))
- \[AV-107162\] Do not allow setting autoexpansion for AWS/GCP [\#366](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/366) ([l0n3star](https://github.com/l0n3star))

## [v1.5.2](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.5.2) (2025-07-22)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.5.1...v1.5.2)

**Implemented enhancements:**

- \[AV-104869\] Implement GCP Get Private Endpoint Command support [\#359](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/359) ([akhilravuri-cb](https://github.com/akhilravuri-cb))
- \[AV-104294\] Consume 504 gateway error [\#357](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/357) ([l0n3star](https://github.com/l0n3star))
- \[AV-102981\] Implement App services Allowed CIDR [\#353](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/353) ([mohammed-madi](https://github.com/mohammed-madi))
- \[AV-103003\] Create common app service for acc tests [\#351](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/351) ([l0n3star](https://github.com/l0n3star))

**Fixed bugs:**

- \[AV-102723\] Initialize num\_replica to null [\#350](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/350) ([l0n3star](https://github.com/l0n3star))

**Closed issues:**

- Provider resource for App Endpoints [\#346](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues/346)

**Merged pull requests:**

- \[AV-102370\] Add PR title and description checker script [\#344](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/344) ([PaulomeeCb](https://github.com/PaulomeeCb))

## [v1.5.1](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.5.1) (2025-05-26)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.5.0...v1.5.1)

**Fixed bugs:**

- AV-102419: Fix PFT cluster schema [\#345](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/345) ([SaicharanCB](https://github.com/SaicharanCB))

## [v1.5.0](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.5.0) (2025-05-21)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.4.1...v1.5.0)

**Implemented enhancements:**

- \[AV-98659\] Implement Free Tier On Off  [\#286](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/286) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-98401\] Add Free Tier App Service Resource [\#283](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/283) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-98308\] Add Free Tier Bucket Resource [\#282](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/282) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-84484\] Implement Free Tier Cluster [\#264](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/264) ([SaicharanCB](https://github.com/SaicharanCB))

**Fixed bugs:**

- \[AV-99812\] Set Provider Type [\#287](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/287) ([l0n3star](https://github.com/l0n3star))
- \[AV-97596\] Set Computed Value for Replica Correctly [\#278](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/278) ([l0n3star](https://github.com/l0n3star))
- \[AV-97308\] Clean Up Resources on Setup Fail [\#276](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/276) ([l0n3star](https://github.com/l0n3star))
- \[AV-97306\] Set indexName for Secondary Index [\#275](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/275) ([l0n3star](https://github.com/l0n3star))
- \[AV-97171\] Fix Github Files to Use ubuntu-latest and upgrade sdk [\#273](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/273) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-97053\] Get Default Primary Index Name Correctly [\#269](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/269) ([l0n3star](https://github.com/l0n3star))

**Closed issues:**

- Improve Documentation [\#289](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues/289)
- VPC peering resource is recreated on each apply [\#284](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues/284)
- Problem creating primary index with one replica [\#277](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues/277)
- Provider resource for cluster backups [\#268](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues/268)
- Docs: Provider docs wrong variable name [\#261](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues/261)
- How to create a primary index using the couchbase-capella\_query\_indexes [\#259](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues/259)
- Does provider support Alert Integration [\#257](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues/257)
- Database Injection Attack vulnerability [\#256](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues/256)

**Merged pull requests:**

- Bump github.com/hashicorp/terraform-plugin-framework from 1.14.1 to 1.15.0 [\#342](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/342) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-testing from 1.12.0 to 1.13.0 [\#341](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/341) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-go from 0.26.0 to 0.27.0 [\#340](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/340) ([dependabot[bot]](https://github.com/apps/dependabot))
- \[AV-102211\] Address Final Docs Review Comments [\#339](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/339) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-101713\] Address Resources Docs Review Comments [\#338](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/338) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-102119\] Add Documentation for Scopes and Collections Resources and Datasources [\#337](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/337) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-101712\] Address Datasources Docs Review Comments [\#336](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/336) ([PaulomeeCb](https://github.com/PaulomeeCb))
- Bump github.com/hashicorp/terraform-plugin-framework-validators from 0.17.0 to 0.18.0 [\#335](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/335) ([dependabot[bot]](https://github.com/apps/dependabot))
- \[AV-101826\] Automate Addition of PR labels [\#334](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/334) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-101710\] Update Documentation for All Resources [\#331](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/331) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-101592\] Update index.md and 1.5.0-upgrade-guide.md Files [\#330](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/330) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-101448\] Update Documentation for All Datasources [\#328](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/328) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100853\] Update Readme File in Examples [\#327](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/327) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100854\] Add documentation for GSI Resources and Datasources [\#326](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/326) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100850\] Add Documentation for Free-Tier Cluster [\#325](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/325) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-100851\] Add Documentation for Free Tier On/Off [\#324](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/324) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-100882\] Add Documentation for Cluster On/Off Schedule [\#323](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/323) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-100881\] Add Documentation for  Cluster On/Off On Demand [\#322](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/322) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-100880\] Add Documentation for App Service On/Off [\#321](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/321) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-100879\] Add documentation for Certificate Datasource [\#320](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/320) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-100878\] Add documentation for Backup Schedule Resource and Datasource [\#319](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/319) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-100877\] Add documentation for Backup Resource and Datasource [\#318](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/318) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-100866\] Add documentation for Sample Bucket Resource [\#317](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/317) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100867\] Add documentation for Project Resource and Datasource \(Part 2\) [\#316](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/316) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100868\] Add documentation for Project Event Datasource \(Part 1\) [\#315](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/315) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100869\] Add documentation for Project Resource and Datasource [\#314](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/314) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100870\] Add documentation for Bucket Flush Datasource [\#313](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/313) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100871\] Add documentation for Event Datasource \(Part 2\) [\#312](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/312) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100872\] Add documentation for Event Datasource \(Part 1\)  [\#311](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/311) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100873\] Add documentation for Database Credentials Resource and Datasource [\#310](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/310) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100874\] Add documentation for App Services Resource and Datasource [\#309](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/309) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100875\] Add documentation for API Keys Resource and Datasource [\#308](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/308) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100855\] Add documentation for Network Peering Resource and Datasource [\#307](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/307) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100856\] Add documentation for Azure VNet Peering Datasource [\#306](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/306) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100857\] Add documentation for Organization Datasource [\#305](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/305) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100858\] Add documentation for Private Endpoints Resource and Datasource \(Part 3\) [\#304](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/304) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100859\] Add documentation for Private Endpoints Resource and Datasource \(Part 2\) [\#303](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/303) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100860\] Add documentation for Private Endpoints Resource and Datasource \(Part 1\) [\#302](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/302) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100861\] Add documentation for Audit log Export Resource and Datasource [\#301](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/301) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100862\] Add documentation for Audit log Settings [\#300](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/300) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100863\] Add documentation for Audit log Event IDs Datasource [\#299](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/299) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100864\] Add documentation for Cluster Resource and Datasource [\#298](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/298) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-101023\] Add import code snippet to documentation [\#297](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/297) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100893\] Add markdown description to fields using Generics [\#295](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/295) ([PaulomeeCb](https://github.com/PaulomeeCb))
- \[AV-100865\] Add documentation for User Resource and Datasource [\#293](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/293) ([karanjain-ops](https://github.com/karanjain-ops))
- \[AV-100849\] Add documentation for Free Tier Buckets [\#292](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/292) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-100830\] Add documentation for free-tier resources [\#291](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/291) ([SaicharanCB](https://github.com/SaicharanCB))
- \[AV-100782\] Add documentation for the Allowlist Resource [\#290](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/290) ([PaulomeeCb](https://github.com/PaulomeeCb))
- Bump golang.org/x/net from 0.37.0 to 0.38.0 [\#288](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/288) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-testing from 1.11.0 to 1.12.0 [\#281](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/281) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump golang.org/x/time from 0.10.0 to 0.11.0 [\#279](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/279) ([dependabot[bot]](https://github.com/apps/dependabot))
- \[AV-97194\] Update Backup Documentation to add the term Bucket [\#274](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/274) ([l0n3star](https://github.com/l0n3star))
- \[AV-97063\] Tie down the versions for golangci-lint and ubuntu in Github workflows [\#270](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/270) ([Talina06](https://github.com/Talina06))
- Bump github.com/hashicorp/terraform-plugin-framework from 1.13.0 to 1.14.1 [\#267](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/267) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-framework-validators from 0.16.0 to 0.17.0 [\#265](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/265) ([dependabot[bot]](https://github.com/apps/dependabot))
- \[Docs\] Update documentation to use the correct environment variable [\#263](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/263) ([cdsre](https://github.com/cdsre))
- Bump golang.org/x/time from 0.8.0 to 0.10.0 [\#260](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/260) ([dependabot[bot]](https://github.com/apps/dependabot))
- \[AV-92992\] Refactor and Optimize Acceptance Tests [\#255](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/255) ([l0n3star](https://github.com/l0n3star))
- Bump github.com/hashicorp/terraform-plugin-testing from 1.7.0 to 1.11.0 [\#253](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/253) ([dependabot[bot]](https://github.com/apps/dependabot))



## [v1.4.1](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/v1.4.1) (2024-12-20)

[Full Changelog](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/compare/v1.4.0...v1.4.1)

**Implemented enhancements:**
- Allows the provider `host` and `authentication_token` to be set through environment variables prefixed with `CAPELLA_HOST` and `CAPELLA_AUTHENTICATION_TOKEN`. [\#239](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/239) ([cdsre](https://github.com/cdsre))

**Fixed bugs:**
- \[AV-92771\] Fix GSI documentation link  [\#248](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/248) ([a-atri](https://github.com/a-atri))

**Documentation Enhancement:**
- \[AV-92775\] Update the Getting Started GSI example to create a non-deferred secondary index [\#247](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/247) ([l0n3star](https://github.com/l0n3star))

**Merged pull requests:**
- Bump github.com/stretchr/testify from 1.9.0 to 1.10.0 [#251](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/251) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump golang.org/x/crypto from 0.30.0 to 0.31.0 [#249](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/249) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/couchbase/tools-common/functional from 1.2.0 to 1.3.1 [#246](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/246) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-framework-validators from 0.12.0 to 0.16.0 [#237](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/237) ([dependabot[bot]](https://github.com/apps/dependabot))



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
