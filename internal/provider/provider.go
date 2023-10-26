// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"terraform-provider-edstem/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &edstemProvider{}

type edstemProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &edstemProvider{
			version: version,
		}
	}
}

func (p *edstemProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "edstem"
	resp.Version = p.version
}

type edstemProviderModel struct {
	CourseId types.String `tfsdk:"course_id"`
	Token    types.String `tfsdk:"token"`
}

func (p *edstemProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"course_id": schema.StringAttribute{
				Required: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *edstemProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config edstemProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.CourseId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("course_id"),
			"Unknown Edstem Course ID",
			"The provider cannot create the Edstem API client as these is an unknown configuration value for the Edstem API course_id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDSTEM_COURSE_ID environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Edstem Access Token",
			"The provider cannot create the Edstem API client as these is an unknown configuration value for the Edstem API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDSTEM_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	course_id := os.Getenv("EDSTEM_COURSE_ID")
	token := os.Getenv("EDSTEM_TOKEN")

	if !config.CourseId.IsNull() {
		course_id = config.CourseId.ValueString()
	}
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if course_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("course_id"),
			"Missing Edstem API Course ID",
			"The provider cannot create the Edstem API client as there is a missing or empty value for the Edstem API course_id. "+
				"Set the Course ID value in the configuration or use the EDSTEM_COURSE_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Edstem API Token",
			"The provider cannot create the Edstem API client as there is a missing or empty value for the Edstem API token. "+
				"Set the Token value in the configuration or use the EDSTEM_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := client.NewClient(&course_id, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Edstem API Client",
			"An unexpected error occurred when creating the Edstem API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Edstem Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

}

func (p *edstemProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewLessonDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *edstemProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLessonResource,
		NewSlideResource,
	}
}
