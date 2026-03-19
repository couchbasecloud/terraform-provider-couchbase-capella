---
name: tf-datasource-gen
description: generate terraform datasources based on openapi spec.
---

# Terraform Datasource Generator

## Instructions

1.  datasource code should be in internal/datasources/
2.  schema for datasource should be in its own file with format <feature>_schema.go

    add validation for organization_id, project_id and cluster_id if present.  for example with organization_id

    capellaschema.AddAttr(attrs, "organization_id", snapshotBackupBuilder, requiredStringWithValidator())

    func requiredStringWithValidator() *schema.StringAttribute {
    	return &schema.StringAttribute{
    		Required:   true,
    		Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
    	}
    }

3.  implement one or two datasources depending on the spec provided.

    - the first is to get a specific resource.  use the get endpoint.  if there is no get endpoint then skip this implementation.
    for example if the feature is Buckets then need bucket.go to get a specific bucket.

    - the second datasource is to list all resources.  use the list endpoint. if there is no list endpoint then skip this implementation.
    the file name should have a plural resource name.
    for example if the feature is Buckets then need buckets.go to list all buckets.

4.  create struct with feature name that embeds Data struct.  for example if the feature is Buckets then need this struct

    type Buckets struct {
    	*providerschema.Data
    }
5.  need New function.  for example if feature is Buckets then need this function

    func NewBuckets() datasource.DataSource {
    	return &Buckets{}
    }

6.  type should implement interfaces datasource.DataSource and datasource.DataSourceWithConfigure.
    must use type conversion of nil to assert that the type implements the interfaces.

    for example for Buckets

    var (
    	_ datasource.DataSource              = (*Buckets)(nil)
    	_ datasource.DataSourceWithConfigure = (*Buckets)(nil)
    )
7.  need Metadata function.  for example with Buckets

    func (d *Buckets) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    	resp.TypeName = req.ProviderTypeName + "_buckets"
    }

8.  need Configure function.  for example with Buckets

    func (d *Buckets) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    	if req.ProviderData == nil {
    		return
    	}

    	data, ok := req.ProviderData.(*providerschema.Data)
    	if !ok {
    		resp.Diagnostics.AddError(
    			"Unexpected Data Source Configure Type",
    			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    		)

    		return
    	}

    	d.Data = data
    }

9.  generate necessary structs to handle API response.  put structs in internal/api/
    use ClientV1 struct to make API calls with retry logic.  for example:

    response, err := s.ClientV1.ExecuteWithRetry

10.  register the datasource in internal/provider/provider.go in func (p *capellaProvider) DataSources

    for example with buckets need datasources.NewBuckets,

11.  create acceptance tests for both datasources in acceptance_tests/ with format <feature>_test.go.
     for example if feature is Buckets then need buckets_test.go

12.  acceptance tests should run in parallel.  that is use resource.ParallelTest()


