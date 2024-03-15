package provider

import (
	"context"
	"fmt"
	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/md2ed"
	"terraform-provider-edstem/internal/resourceclients"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &slideResource{}
	_ resource.ResourceWithConfigure = &slideResource{}
)

// NewSlideResource is a helper function to simplify the provider implementation.
func NewSlideResource() resource.Resource {
	return &slideResource{}
}

// slideResource is the resource implementation.
type slideResource struct {
	client *client.Client
}

// Configure adds the provider configured client to the resource.
func (r *slideResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *slideResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_slide"
}

type slideResourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	LessonId    types.Int64  `tfsdk:"lesson_id"`
	Title       types.String `tfsdk:"title"`
	Index       types.Int64  `tfsdk:"index"`
	IsHidden    types.Bool   `tfsdk:"is_hidden"`
	Content     types.String `tfsdk:"content"`
	ContentType types.String `tfsdk:"content_type"`
}

// Schema defines the schema for the resource.
func (r *slideResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				// TODO: Validate
			},
			"lesson_id": schema.Int64Attribute{
				Required: true,
			},
			"title": schema.StringAttribute{
				Required: true,
			},
			"index": schema.Int64Attribute{
				Required: true,
			},
			"is_hidden": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"content": schema.StringAttribute{
				Required: true,
			},
			"content_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("md"),
			},
		},
	}
}

func (model *slideResourceModel) MapAPIObj(ctx context.Context) (*resourceclients.Slide, error) {
	var obj resourceclients.Slide

	obj.Id = int(model.Id.ValueInt64())
	obj.Type = model.Type.ValueString()
	obj.LessonId = int(model.LessonId.ValueInt64())
	obj.Title = model.Title.ValueString()
	obj.Index = int(model.Index.ValueInt64())
	obj.IsHidden = model.IsHidden.ValueBool()
	obj.Content = model.Content.ValueString()
	if model.ContentType.ValueString() == "md" {
		obj.Content = md2ed.RenderMDToEd(obj.Content)
		fmt.Print(obj.Content)
	}
	return &obj, nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *slideResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan slideResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api_obj, err := plan.MapAPIObj(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Slide Object",
			fmt.Sprintf("Could not create Slide: %s", err.Error()),
		)
	}

	err = resourceclients.CreateSlide(r.client, api_obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Slide Object",
			fmt.Sprintf("Could not create Slide for Lesson ID %d: %s", api_obj.LessonId, err.Error()),
		)
	}

	plan.Id = types.Int64Value(int64(api_obj.Id))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *slideResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state slideResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := resourceclients.GetSlide(r.client, int(state.LessonId.ValueInt64()), int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Slide Object",
			fmt.Sprintf("Could not read Slide ID %d: %s", state.Id.ValueInt64(), err.Error()),
		)
	}

	// TODO: For now, nothing happens with the read elements. Should update state to confirm any changes necessary.

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *slideResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan slideResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api_obj, err := plan.MapAPIObj(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Slide Object",
			fmt.Sprintf("Could not update Slide: %s", err.Error()),
		)
	}

	err = resourceclients.UpdateSlide(r.client, api_obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Slide Object",
			fmt.Sprintf("Could not update Lesson ID %d: %s", api_obj.Id, err.Error()),
		)
	}

	plan.Id = types.Int64Value(int64(api_obj.Id))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *slideResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
