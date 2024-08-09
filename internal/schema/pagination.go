package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Cursor struct {
	Hrefs *Hrefs `tfsdk:"hrefs"`
	Pages *Pages `tfsdk:"pages"`
}

type Hrefs struct {
	First    types.String `tfsdk:"first"`
	Last     types.String `tfsdk:"last"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type Pages struct {
	Last       types.Int64 `tfsdk:"last"`
	Next       types.Int64 `tfsdk:"next"`
	Page       types.Int64 `tfsdk:"page"`
	PerPage    types.Int64 `tfsdk:"per_page"`
	Previous   types.Int64 `tfsdk:"previous"`
	TotalItems types.Int64 `tfsdk:"total_items"`
}
