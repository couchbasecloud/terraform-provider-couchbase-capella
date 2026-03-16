# Code Review


1.  git diff main --name-only -- '*.go' ':!internal/generated/api/openapi.gen.go'

    do not read the openapi.gen.go file as it will consume many tokens and fill up the context window.

2.  verify new datasources and resources use ClientV1.  if not update the code to use ClientV1.

3.  for each go file, look at git diff of the file compared to the one in main
4.  all functions, structs and global var/const must have a comment
5.  run goimports -w -local github.com/couchbasecloud/terraform-provider-couchbase-capella on the file
6.  run go vet on the file and fix any errors
7.  build the binary as follows

get the latest git tag and assign to VERSION:

VERSION=$(git describe --tags --abbrev=0)

then build:

go build -ldflags "-s -w -X 'github.com/couchbasecloud/terraform-provider-couchbase-capella/version.ProviderVersion=$VERSION'" -o ./bin/terraform-provider-couchbase-capella

fix any errors.

8.  repeat steps 5-7 until there are no errors.  retry up to 5 times.  if errors persist report them.