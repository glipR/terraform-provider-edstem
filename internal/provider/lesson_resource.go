package provider

import (
	"context"
	"fmt"

	"terraform-provider-edstem/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &lessonResource{}
	_ resource.ResourceWithConfigure = &lessonResource{}
)

// NewLessonResource is a helper function to simplify the provider implementation.
func NewLessonResource() resource.Resource {
	return &lessonResource{}
}

// lessonResource is the resource implementation.
type lessonResource struct {
	client *client.Client
}

// Configure adds the provider configured client to the resource.
func (r *lessonResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *lessonResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lesson"
}

type lessonResourceModel struct {
	id                                       types.Int64  `tfsdk:"id"`
	attempts                                 types.Int64  `tfsdk:"attempts"`
	available_at                             types.String `tfsdk:"available_at"`
	due_at                                   types.String `tfsdk:"due_at"`
	grade_passback_auto_send                 types.Bool   `tfsdk:"grade_passback_auto_send"`
	grade_passback_mode                      types.String `tfsdk:"grade_passback_mode"`
	grade_passback_scale_to                  types.String `tfsdk:"grade_passback_scale_to"`
	index                                    types.Int64  `tfsdk:"index"`
	is_hidden                                types.Bool   `tfsdk:"is_hidden"`
	is_timed                                 types.Bool   `tfsdk:"is_timed"`
	is_unlisted                              types.Bool   `tfsdk:"is_unlisted"`
	kind                                     types.String `tfsdk:"kind"`
	late_submissions                         types.Bool   `tfsdk:"late_submissions"`
	locked_at                                types.String `tfsdk:"locked_at"`
	module_id                                types.Int64  `tfsdk:"module_id"`
	openable                                 types.Bool   `tfsdk:"openable"`
	openable_without_attempt                 types.Bool   `tfsdk:"openable_without_attempt"`
	outline                                  types.String `tfsdk:"outline"`
	password                                 types.String `tfsdk:"password"`
	prerequisites                            types.List   `tfsdk:"prerequisites"`
	release_challenge_solutions              types.Bool   `tfsdk:"release_challenge_solutions"`
	release_challenge_solutions_while_active types.Bool   `tfsdk:"release_challenge_solutions"`
	release_feedback                         types.Bool   `tfsdk:"release_feedback"`
	release_feedback_while_active            types.Bool   `tfsdk:"release_feedback_while_active"`
	release_quiz_correctness_only            types.Bool   `tfsdk:"release_quiz_correctness_only"`
	release_quiz_solutions                   types.Bool   `tfsdk:"release_quiz_solutions"`
	quiz_active_status                       types.String `tfsdk:"quiz_active_status"`
	quiz_mode                                types.String `tfsdk:"quiz_mode"`
	quiz_question_number_style               types.String `tfsdk:"quiz_question_number_style"`
	solutions_at                             types.String `tfsdk:"solutions_at"`
	state                                    types.String `tfsdk:"state"`
	timer_duration                           types.Int64  `tfsdk:"timer_duration"`
	timer_expiration_access                  types.Bool   `tfsdk:"timer_expiration_access"`
	title                                    types.String `tfsdk:"title"`
	tutorial_regex                           types.String `tfsdk:"tutorial_regex"`
	lesson_type                              types.String `tfsdk:"type"`
}

// Schema defines the schema for the resource.
func (r *lessonResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"attempts": schema.Int64Attribute{
				Optional: true,
			},
			"available_at": schema.StringAttribute{
				Optional: true,
			},
			"due_at": schema.StringAttribute{
				Optional: true,
			},
			"grade_passback_auto_send": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"grade_passback_mode": schema.StringAttribute{
				Default: stringdefault.StaticString(""),
			},
			"grade_passback_scale_to": schema.StringAttribute{
				Optional: true,
			},
			"index": schema.Int64Attribute{
				Optional: true,
			},
			"is_hidden": schema.BoolAttribute{
				Default: booldefault.StaticBool(true),
			},
			"is_timed": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"is_unlisted": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"kind": schema.StringAttribute{
				Default: stringdefault.StaticString("legacy"),
			},
			"late_submissions": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"locked_at": schema.StringAttribute{
				Optional: true,
			},
			"module_id": schema.Int64Attribute{
				Optional: true,
			},
			"openable": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"openable_without_attempt": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"outline": schema.StringAttribute{
				Default: stringdefault.StaticString(""),
			},
			"password": schema.StringAttribute{
				Default: stringdefault.StaticString(""),
			},
			"prerequisites": schema.ListAttribute{
				ElementType: types.StringType,
			},
			"release_challenge_solutions": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"release_challenge_solutions_while_active": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"release_feedback": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"release_feedback_while_active": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"release_quiz_correctness_only": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"release_quiz_solutions": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"reopen_submissions": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"quiz_active_status": schema.StringAttribute{
				Default: stringdefault.StaticString("active"),
			},
			"quiz_mode": schema.StringAttribute{
				Default: stringdefault.StaticString("multiple-attempts"),
			},
			"quiz_question_number_style": schema.StringAttribute{
				Default: stringdefault.StaticString(""),
			},
			"solutions_at": schema.StringAttribute{
				Optional: true,
			},
			"state": schema.StringAttribute{
				Default: stringdefault.StaticString("active"),
			},
			"timer_duration": schema.Int64Attribute{
				Default: int64default.StaticInt64(60),
			},
			"timer_expiration_access": schema.BoolAttribute{
				Default: booldefault.StaticBool(false),
			},
			"title": schema.StringAttribute{
				Required: true,
			},
			"tutorial_regex": schema.StringAttribute{
				Default: stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				Default: stringdefault.StaticString("general"),
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *lessonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read refreshes the Terraform state with the latest data.
func (r *lessonResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *lessonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *lessonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
