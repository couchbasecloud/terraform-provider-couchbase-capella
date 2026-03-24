# Code Review

-  for each new or modified file, run git add <file>
-  do not run git commit
-  run git diff --cached main --name-only -- '*.go' ':!internal/generated/api/openapi.gen.go'
-  do not read the openapi.gen.go file as it will consume many tokens and fill up the context window.
-  verify new datasources and resources use ClientV1.  if not update the code to use ClientV1.
-  unexported symbols only need comments if they have side effects, workarounds or constraints not evident from the name and signature.
-  format and build as follows:
    -   run goimports -w -local github.com/couchbasecloud/terraform-provider-couchbase-capella on the file
    -   run go vet on the file and fix any errors
    -   get the latest git tag and assign to VERSION:
            VERSION=$(git describe --tags --abbrev=0)
    -   build with this command:
        go build -ldflags "-s -w -X 'github.com/couchbasecloud/terraform-provider-couchbase-capella/version.ProviderVersion=$VERSION'" -o ./bin/terraform-provider-couchbase-capella

        fix any errors.
    -   repeat format and build steps until there are no errors.  retry up to 5 times.  if errors persist report them.