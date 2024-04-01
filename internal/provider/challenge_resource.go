package provider

import (
	"context"
	"fmt"
	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/resourceclients"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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

	Type types.String `tfsdk:"type"`
	// Points types.Int64  `tfsdk:"points"`

	BuildCommand     types.String `tfsdk:"build_command"`
	RunCommand       types.String `tfsdk:"run_command"`
	TestCommand      types.String `tfsdk:"test_command"`
	TerminalCommand  types.String `tfsdk:"terminal_command"`
	CustomRunCommand types.String `tfsdk:"custom_run_command"`

	// PointLossThreshold types.Int64 `tfsdk:"point_loss_threshold"`
	// PointLossEvery     types.Int64 `tfsdk:"point_loss_every"`
	// PointLossAmount    types.Int64 `tfsdk:"point_loss_amount"`

	PerTestcaseScores types.Bool `tfsdk:"per_testcase_scores"`

	MaxSubmissionsPerInterval types.Int64 `tfsdk:"max_submissions_per_interval"`
	AttemptLimitInterval      types.Int64 `tfsdk:"attempt_limit_interval"`

	OnlyGitSubmission            types.Bool `tfsdk:"only_git_submission"`
	AllowSubmitAfterMarkingLimit types.Bool `tfsdk:"allow_submit_after_marking_limit"`

	PassbackScoringMode       types.String  `tfsdk:"passback_scoring_mode"`
	PassbackMaxAutomaticScore types.Float64 `tfsdk:"passback_max_automatic_score"`
	PassbackScaleTo           types.Float64 `tfsdk:"passback_scale_to"`

	Run                  types.Bool `tfsdk:"feature_run"`
	Check                types.Bool `tfsdk:"feature_check"`
	Mark                 types.Bool `tfsdk:"feature_mark"`
	Terminal             types.Bool `tfsdk:"feature_terminal"`
	Connect              types.Bool `tfsdk:"feature_connect"`
	Feedback             types.Bool `tfsdk:"feature_feedback"`
	ManualCompletion     types.Bool `tfsdk:"feature_manual_completion"`
	AnonymousSubmissions types.Bool `tfsdk:"feature_anonymous_submissions"`
	Arguments            types.Bool `tfsdk:"feature_arguments"`
	ConfirmSubmit        types.Bool `tfsdk:"feature_confirm_submit"`
	RunBeforeSubmit      types.Bool `tfsdk:"feature_run_before_submit"`
	GitSubmission        types.Bool `tfsdk:"feature_git_submission"`
	Editor               types.Bool `tfsdk:"feature_editor"`
	RemoteDesktop        types.Bool `tfsdk:"feature_remote_desktop"`
	IntermediateFiles    types.Bool `tfsdk:"feature_intermediate_files"`

	CustomMarkTimeLimitMS types.Int64  `tfsdk:"custom_mark_time_limit_ms"`
	TestcaseJSON          types.String `tfsdk:"testcase_json"`

	// TODO:
	// * Rubric
	// * Test cases for the standard marking procedure
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
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"folder_path": schema.StringAttribute{
				Required: true,
			},
			"folder_sha": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Default:  stringdefault.StaticString("none"),
				Optional: true,
				Computed: true,
			},
			"build_command": schema.StringAttribute{
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"run_command": schema.StringAttribute{
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"test_command": schema.StringAttribute{
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"terminal_command": schema.StringAttribute{
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"custom_run_command": schema.StringAttribute{
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"per_testcase_scores": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"max_submissions_per_interval": schema.Int64Attribute{
				Optional: true,
			},
			"attempt_limit_interval": schema.Int64Attribute{
				Optional: true,
			},
			"only_git_submission": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"allow_submit_after_marking_limit": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"passback_scoring_mode": schema.StringAttribute{
				Optional: true,
			},
			"passback_max_automatic_score": schema.Float64Attribute{
				Optional: true,
			},
			"passback_scale_to": schema.Float64Attribute{
				Optional: true,
			},
			"feature_run": schema.BoolAttribute{
				Default:  booldefault.StaticBool(true),
				Optional: true,
				Computed: true,
			},
			"feature_check": schema.BoolAttribute{
				Default:  booldefault.StaticBool(true),
				Optional: true,
				Computed: true,
			},
			"feature_mark": schema.BoolAttribute{
				Default:  booldefault.StaticBool(true),
				Optional: true,
				Computed: true,
			},
			"feature_terminal": schema.BoolAttribute{
				Default:  booldefault.StaticBool(true),
				Optional: true,
				Computed: true,
			},
			"feature_connect": schema.BoolAttribute{
				Default:  booldefault.StaticBool(true),
				Optional: true,
				Computed: true,
			},
			"feature_feedback": schema.BoolAttribute{
				Default:  booldefault.StaticBool(true),
				Optional: true,
				Computed: true,
			},
			"feature_manual_completion": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"feature_anonymous_submissions": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"feature_arguments": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"feature_confirm_submit": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"feature_run_before_submit": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"feature_git_submission": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"feature_editor": schema.BoolAttribute{
				Default:  booldefault.StaticBool(true),
				Optional: true,
				Computed: true,
			},
			"feature_remote_desktop": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"feature_intermediate_files": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"custom_mark_time_limit_ms": schema.Int64Attribute{
				Optional: true,
			},
			"testcase_json": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (model *challengeResourceModel) MapAPIObj(ctx context.Context, client *client.Client) (*resourceclients.Challenge, error) {

	lesson_id := model.LessonId.ValueInt64()
	slide_id := model.SlideId.ValueInt64()

	chal, err := resourceclients.GetChallenge(client, int(lesson_id), int(slide_id))
	if err != nil {
		return nil, err
	}

	chal.Type = model.Type.ValueString()
	chal.Explanation = model.Explanation.ValueString()

	chal.Features.Run = model.Run.ValueBool()
	chal.Features.Check = model.Check.ValueBool()
	chal.Features.Mark = model.Mark.ValueBool()
	chal.Features.Connect = model.Connect.ValueBool()
	chal.Features.Terminal = model.Terminal.ValueBool()
	chal.Features.Feedback = model.Feedback.ValueBool()
	chal.Features.ManualCompletion = model.ManualCompletion.ValueBool()
	chal.Features.AnonymousSubmissions = model.AnonymousSubmissions.ValueBool()
	chal.Features.Arguments = model.Arguments.ValueBool()
	chal.Features.ConfirmSubmit = model.ConfirmSubmit.ValueBool()
	chal.Features.RunBeforeSubmit = model.RunBeforeSubmit.ValueBool()
	chal.Features.GitSubmission = model.GitSubmission.ValueBool()
	chal.Features.Editor = model.Editor.ValueBool()
	chal.Features.RemoteDesktop = model.RemoteDesktop.ValueBool()
	chal.Features.IntermediateFiles = model.IntermediateFiles.ValueBool()

	chal.Settings.BuildCommand = model.BuildCommand.ValueString()
	chal.Settings.CheckCommand = model.TestCommand.ValueString()
	chal.Settings.RunCommand = model.RunCommand.ValueString()
	chal.Settings.TerminalCommand = model.TerminalCommand.ValueString()
	chal.Settings.OnlyGitSubmission = model.OnlyGitSubmission.ValueBool()
	chal.Settings.AttemptLimitInterval = int(model.AttemptLimitInterval.ValueInt64())
	chal.Settings.AllowSubmitAfterMarkingLimit = model.AllowSubmitAfterMarkingLimit.ValueBool()
	chal.Settings.MaxSubmissionsPerInterval = int(model.MaxSubmissionsPerInterval.ValueInt64())
	chal.Settings.Passback.MaxAutomaticScore = model.PassbackMaxAutomaticScore.ValueFloat64()
	chal.Settings.Passback.ScoringMode = model.PassbackScoringMode.ValueString()
	chal.Settings.Passback.ScaleTo = model.PassbackScaleTo.ValueFloat64()
	chal.Settings.PerTestCaseScores = model.PerTestcaseScores.ValueBool()

	chal.Tickets.MarkUnit.BuildCommand = model.BuildCommand.ValueString()
	chal.Tickets.MarkCustom.BuildCommand = model.BuildCommand.ValueString()
	chal.Tickets.MarkStandard.BuildCommand = model.BuildCommand.ValueString()

	if chal.Type == "none" {
		chal.Tickets.RunStandard.RunCommand = model.RunCommand.ValueString()
		chal.Tickets.RunStandard.BuildCommand = model.BuildCommand.ValueString()
	} else if chal.Type == "code" {
		chal.Tickets.RunStandard.RunCommand = model.RunCommand.ValueString()
		chal.Tickets.RunStandard.BuildCommand = model.BuildCommand.ValueString()

		// TODO
		// chal.Tickets.MarkStandard.Testcases = ...
	} else if chal.Type == "custom" {
		chal.Tickets.MarkCustom.RunCommand = model.CustomRunCommand.ValueString()
		if !model.CustomMarkTimeLimitMS.IsNull() {
			chal.Tickets.MarkCustom.RunLimit.CpuTime.Set(model.CustomMarkTimeLimitMS.ValueInt64())
			chal.Tickets.MarkCustom.RunLimit.WallTime.Set(model.CustomMarkTimeLimitMS.ValueInt64())
		}
	}

	return chal, nil
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

	resourceclients.UpdateChallenge(r.client, plan.FolderPath.ValueString(), api_obj)

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

	resourceclients.UpdateChallenge(r.client, plan.FolderPath.ValueString(), api_obj)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *challengeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
