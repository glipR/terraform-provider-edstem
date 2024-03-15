package provider

import (
	"context"
	"fmt"
	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/resourceclients"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &challengeResource{}
	_ resource.ResourceWithConfigure = &challengeResource{}
)

// NewChallengeResource is a helper function to simplify the provider implementation.
func NewChallengeResource() resource.Resource {
	return &challengeResource{}
}

// challengeResource is the resource implementation.
type challengeResource struct {
	client *client.Client
}

// Configure adds the provider configured client to the resource.
func (r *challengeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *challengeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_challenge"
}

type challengeResourceModel struct {
	SlideId  types.Int64 `tfsdk:"slide_id"`
	LessonId types.Int64 `tfsdk:"lesson_id"`

	Explanation types.String `tfsdk:"explanation"`
	FolderPath  types.String `tfsdk:"folder_path"`
	FolderSha   types.String `tfsdk:"folder_sha"`
}

// Schema defines the schema for the resource.
func (r *challengeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"slide_id": schema.Int64Attribute{
				Required: true,
			},
			"lesson_id": schema.Int64Attribute{
				Required: true,
			},
			"explanation": schema.StringAttribute{
				Optional: true,
			},
			"folder_path": schema.StringAttribute{
				Required: true,
			},
			"folder_sha": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (model *challengeResourceModel) MapAPIObj(ctx context.Context, client *client.Client) (*resourceclients.Challenge, error) {

	lesson_id := model.LessonId.ValueInt64()
	slide_id := model.SlideId.ValueInt64()

	return resourceclients.GetChallenge(client, int(lesson_id), int(slide_id))
}

// Create creates the resource and sets the initial Terraform state.
func (r *challengeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan challengeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api_obj, err := plan.MapAPIObj(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Challenge Object",
			fmt.Sprintf("Could not create Challenge: %s", err.Error()),
		)
		return
	}

	var chal resourceclients.ChallengeResource
	chal.Id = api_obj.Id
	chal.CourseId = api_obj.CourseId
	chal.FolderPath = plan.FolderPath.ValueString()

	resourceclients.UpdateChallenge(r.client, &chal)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *challengeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state challengeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := resourceclients.GetChallenge(r.client, int(state.LessonId.ValueInt64()), int(state.SlideId.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Challenge Object",
			fmt.Sprintf("Could not read Challenge from Slide ID %d: %s", state.SlideId.ValueInt64(), err.Error()),
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
func (r *challengeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan challengeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api_obj, err := plan.MapAPIObj(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Challenge Object",
			fmt.Sprintf("Could not update Challenge: %s", err.Error()),
		)
	}
	var chal resourceclients.ChallengeResource
	chal.Id = api_obj.Id
	chal.CourseId = api_obj.CourseId
	chal.FolderPath = plan.FolderPath.ValueString()

	resourceclients.UpdateChallenge(r.client, &chal)

	// TODO Actually update

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *challengeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
