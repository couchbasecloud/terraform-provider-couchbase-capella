# Code Review


1.  git diff main --name-only -- '*.go' ':!internal/generated/api/openapi.gen.go'

    do not read the openapi.gen.go file as it will consume many tokens

2.  for each go file, look at git diff of the file compared to the one in main
3.  all functions, structs and global var/const must have a comment
4.  run goimports -w -local github.com/couchbasecloud/terraform-provider-couchbase-capella on the file
5.  run go vet on the file and fix any errors
6.  build the binary as follows

get the latest git tag and assign to VERSION:

VERSION=$(git describe --tags --abbrev=0)

then build:

go build -ldflags "-s -w -X 'github.com/couchbasecloud/terraform-provider-couchbase-capella/version.ProviderVersion=$VERSION'" -o ./bin/terraform-provider-couchbase-capella

fix any errors.

7.  repeat steps 4-6 until there are no errors.  retry up to 5 times.  if errors persist report them.