//go:generate oapi-codegen -package apigen -generate types,client -response-type-suffix Resp -o ./openapi.gen.go ../../openapi.generated.yaml

package apigen
