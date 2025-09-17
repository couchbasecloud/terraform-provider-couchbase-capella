//go:generate oapi-codegen -package api -generate types,client -response-type-suffix Resp -o ./openapi.gen.go ../../../openapi.generated.yaml

package api
