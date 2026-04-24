# Code Review

0.  stage all relevant changes (including deletions and renames), run `git add FILEPATH`
1.  run git diff --cached main --name-only -- '*.go' ':!internal/generated/api/openapi.gen.go'

    do not read the openapi.gen.go file as it will consume many tokens and fill up the context window.

2.  verify new datasources and resources use ClientV1.  if not update the code to use ClientV1.

3.  for each go file, look at git diff of the file compared to the one in main
4.  Unexported symbols only need comments if they have side effects, workarounds or constraints not evident from the name and signature.
5.  run goimports -w -local github.com/couchbasecloud/terraform-provider-couchbase-capella on the file
6.  run go vet on the file and fix any errors
7.  build the binary as follows

get the latest git tag and assign to VERSION:

```
VERSION=$(git describe --tags --abbrev=0)
```

then build:

```
go build -ldflags "-s -w -X 'github.com/couchbasecloud/terraform-provider-couchbase-capella/version.ProviderVersion=$VERSION'" -o ./bin/terraform-provider-couchbase-capella
```

fix any errors.

8.  repeat steps 5-7 until there are no errors.  retry up to 5 times.  if errors persist report them.