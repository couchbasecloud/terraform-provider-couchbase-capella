package tf

// NOTE:
// 1) Generates a Provider Code Spec JSON from the OpenAPI spec limited to project resource/datasource.
// 2) Generates Terraform Plugin Framework schemas/stubs from that spec.
// Run via `make gen-tf`.

//go:generate sh -c "tfplugingen-openapi generate -config config.yaml -output spec.json ../../../openapi.generated.yaml"

// generate project data sources and resources from OpenAPI-derived spec
//go:generate sh -c "tfplugingen-framework generate data-sources --input spec.json --output ."
//go:generate sh -c "tfplugingen-framework generate resources --input spec.json --output ."
