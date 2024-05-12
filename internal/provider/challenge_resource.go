package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/md2ed"
	"terraform-provider-edstem/internal/resourceclients"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
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

	CustomMarkTimeLimitMS    types.Int64  `tfsdk:"custom_mark_time_limit_ms"`
	TestcaseJSON             types.String `tfsdk:"testcase_json"`
	TestcasePty              types.Bool   `tfsdk:"testcase_pty"`
	TestcaseEasy             types.Bool   `tfsdk:"testcase_easy"`
	TestcaseMarkAll          types.Bool   `tfsdk:"testcase_mark_all"`
	TestcaseOverlayTestFiles types.Bool   `tfsdk:"testcase_overlay_test_files"`

	Criteria     types.String `tfsdk:"criteria"`
	Rubric       types.String `tfsdk:"rubric"`
	RubricPoints types.Int64  `tfsdk:"rubric_points"`
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
				Default:  int64default.StaticInt64(0),
				Optional: true,
				Computed: true,
			},
			"attempt_limit_interval": schema.Int64Attribute{
				Default:  int64default.StaticInt64(0),
				Optional: true,
				Computed: true,
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
				Default:  float64default.StaticFloat64(0),
				Optional: true,
				Computed: true,
			},
			"passback_scale_to": schema.Float64Attribute{
				Default:  float64default.StaticFloat64(0),
				Optional: true,
				Computed: true,
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
				Default:  stringdefault.StaticString("[]"),
				Optional: true,
				Computed: true,
			},
			"testcase_pty": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"testcase_easy": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"testcase_mark_all": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"testcase_overlay_test_files": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"criteria": schema.StringAttribute{
				Default:  stringdefault.StaticString("[]"),
				Optional: true,
				Computed: true,
			},
			"rubric": schema.StringAttribute{
				Default:  stringdefault.StaticString("{}"),
				Optional: true,
				Computed: true,
			},
			"rubric_points": schema.Int64Attribute{
				Optional: true,
			},
		},
	}
}

func (model *challengeResourceModel) MapAPIObj(ctx context.Context, client *client.Client) (*resourceclients.Challenge, *resourceclients.Rubric, error) {

	lesson_id := model.LessonId.ValueInt64()
	slide_id := model.SlideId.ValueInt64()

	chal, rubric, err := resourceclients.GetChallengeAndRubric(client, int(lesson_id), int(slide_id))
	if err != nil {
		return nil, nil, err
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

		testcases := model.TestcaseJSON.ValueString()
		if testcases != "" {
			resp := &[]resourceclients.TestCase{}
			err = json.NewDecoder(strings.NewReader(testcases)).Decode(resp)
			if err != nil {
				return nil, nil, err
			}
			chal.Tickets.MarkStandard.Testcases = *resp
		}
		chal.Tickets.MarkStandard.RunLimit.Pty.Set(model.TestcasePty.ValueBool())
		chal.Tickets.MarkStandard.Easy = model.TestcaseEasy.ValueBool()
		chal.Tickets.MarkStandard.MarkAll = model.TestcaseMarkAll.ValueBool()
		chal.Tickets.MarkStandard.Overlay = model.TestcaseOverlayTestFiles.ValueBool()
	} else if chal.Type == "custom" {
		chal.Tickets.MarkCustom.RunCommand = model.CustomRunCommand.ValueString()
		if !model.CustomMarkTimeLimitMS.IsNull() {
			chal.Tickets.MarkCustom.RunLimit.CpuTime.Set(model.CustomMarkTimeLimitMS.ValueInt64())
			chal.Tickets.MarkCustom.RunLimit.WallTime.Set(model.CustomMarkTimeLimitMS.ValueInt64())
		}
	}

	if !model.Criteria.IsNull() {
		var crit []resourceclients.Criteria
		err = json.NewDecoder(strings.NewReader(model.Criteria.ValueString())).Decode(&crit)
		if err != nil {
			return nil, nil, err
		}
		chal.Settings.Criteria = crit
	}

	if !model.Rubric.IsNull() {
		rubric_data := &resourceclients.Rubric{}
		err = json.NewDecoder(strings.NewReader(model.Rubric.ValueString())).Decode(&rubric_data)
		if err != nil {
			return nil, nil, err
		}
		// TODO: Update but keep ids
		for i, section := range rubric_data.Sections {
			for j, item := range section.Items {
				rubric_data.Sections[i].Items[j].Title = md2ed.RenderMDToEd(item.Title)
			}
		}
		for i, item := range rubric_data.UnsectionedItems {
			rubric_data.UnsectionedItems[i].Title = md2ed.RenderMDToEd(item.Title)
		}
		rubric = rubric_data

		if !model.RubricPoints.IsNull() {
			chal.RubricPoints.Set(int(model.RubricPoints.ValueInt64()))
		}
	}

	return chal, rubric, nil
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

	api_obj, rubric, err := plan.MapAPIObj(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Challenge Object",
			fmt.Sprintf("Could not create Challenge: %s", err.Error()),
		)
		return
	}

	resourceclients.UpdateChallenge(r.client, plan.FolderPath.ValueString(), api_obj, rubric)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func compareCriteria(crit1 []resourceclients.Criteria, crit2 []resourceclients.Criteria) bool {
	if len(crit1) != len(crit2) {
		return false
	}
	for i := range crit1 {
		if crit1[i].Name != crit2[i].Name {
			return false
		}
		if len(crit1[i].Levels) != len(crit2[i].Levels) {
			return false
		}
		for j := range crit1[i].Levels {
			if crit1[i].Levels[j].Description != crit2[i].Levels[j].Description {
				return false
			}
			if crit1[i].Levels[j].Mark != crit2[i].Levels[j].Mark {
				return false
			}
		}
	}
	return true
}

func compareTestCase(tc1 []resourceclients.TestCase, tc2 []resourceclients.TestCase) bool {
	if len(tc1) != len(tc2) {
		return false
	}
	for i := range tc1 {
		if tc1[i].Description != tc2[i].Description ||
			tc1[i].Hidden != tc2[i].Hidden ||
			tc1[i].MaxScore != tc2[i].MaxScore ||
			tc1[i].Name != tc2[i].Name ||
			tc1[i].Private != tc2[i].Private ||
			tc1[i].RunCommand.OrElse("") != tc2[i].RunCommand.OrElse("") ||
			tc1[i].Score != tc2[i].Score ||
			tc1[i].Skip != tc2[i].Skip ||
			tc1[i].StdinPath != tc2[i].StdinPath ||
			tc1[i].RunLimit.CpuTime.OrElse(0) != tc2[i].RunLimit.CpuTime.OrElse(0) ||
			tc1[i].RunLimit.Pty.OrElse(false) != tc2[i].RunLimit.Pty.OrElse(false) {
			return false
		}
		if len(tc1[i].Checks) != len(tc2[i].Checks) {
			return false
		}
		if len(tc1[i].OutputFiles) != len(tc2[i].OutputFiles) {
			return false
		}
		for j := range tc1[i].Checks {
			if tc1[i].Checks[j].ExpectPath != tc2[i].Checks[j].ExpectPath ||
				tc1[i].Checks[j].Markdown != tc2[i].Checks[j].Markdown ||
				tc1[i].Checks[j].Name != tc2[i].Checks[j].Name ||
				tc1[i].Checks[j].Type != tc2[i].Checks[j].Type {
				return false
			}
		}
		for j := range tc1[i].OutputFiles {
			if tc1[i].OutputFiles[j] != tc2[i].OutputFiles[j] {
				return false
			}
		}
	}
	return true
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

	challenge, _, err := resourceclients.GetChallengeAndRubric(r.client, int(state.LessonId.ValueInt64()), int(state.SlideId.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Challenge Object",
			fmt.Sprintf("Could not read Challenge from Slide ID %d: %s", state.SlideId.ValueInt64(), err.Error()),
		)
		return
	}

	state.AllowSubmitAfterMarkingLimit = types.BoolValue(challenge.Settings.AllowSubmitAfterMarkingLimit)
	state.AnonymousSubmissions = types.BoolValue(challenge.Features.AnonymousSubmissions)
	state.Arguments = types.BoolValue(challenge.Features.Arguments)
	state.AttemptLimitInterval = types.Int64Value(int64(challenge.Settings.AttemptLimitInterval))
	state.BuildCommand = types.StringValue(challenge.Settings.BuildCommand)
	state.Check = types.BoolValue(challenge.Features.Check)
	state.ConfirmSubmit = types.BoolValue(challenge.Features.ConfirmSubmit)
	state.Connect = types.BoolValue(challenge.Features.Connect)
	var cur_state []resourceclients.Criteria
	json.NewDecoder(strings.NewReader(state.Criteria.ValueString())).Decode(&cur_state)
	if !compareCriteria(cur_state, challenge.Settings.Criteria) {
		// Criteria are different, set the state.
		crit, err := json.MarshalIndent(challenge.Settings.Criteria, "", "  ")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Criteria Object",
				fmt.Sprintf("Could not read Criteria from Slide ID %d: %s", state.SlideId.ValueInt64(), err.Error()),
			)
		}
		state.Criteria = types.StringValue(string(crit))
	}
	challenge.Tickets.MarkCustom.RunLimit.CpuTime.If(func(val int64) { state.CustomMarkTimeLimitMS = types.Int64Value(val) })
	state.CustomRunCommand = types.StringValue(challenge.Tickets.MarkCustom.RunCommand)
	state.Editor = types.BoolValue(challenge.Features.Editor)
	state.Explanation = types.StringValue(challenge.Explanation)
	state.Feedback = types.BoolValue(challenge.Features.Feedback)
	// TODO: Check the challenge folder contents and compare.
	state.GitSubmission = types.BoolValue(challenge.Features.GitSubmission)
	state.IntermediateFiles = types.BoolValue(challenge.Features.IntermediateFiles)
	state.ManualCompletion = types.BoolValue(challenge.Features.ManualCompletion)
	state.Mark = types.BoolValue(challenge.Features.Mark)
	state.MaxSubmissionsPerInterval = types.Int64Value(int64(challenge.Settings.MaxSubmissionsPerInterval))
	state.OnlyGitSubmission = types.BoolValue(challenge.Settings.OnlyGitSubmission)
	state.PassbackMaxAutomaticScore = types.Float64Value(challenge.Settings.Passback.MaxAutomaticScore)
	state.PassbackScaleTo = types.Float64Value(challenge.Settings.Passback.ScaleTo)
	state.PassbackScoringMode = types.StringValue(challenge.Settings.Passback.ScoringMode)
	state.PerTestcaseScores = types.BoolValue(challenge.Settings.PerTestCaseScores)
	state.RemoteDesktop = types.BoolValue(challenge.Features.RemoteDesktop)
	state.Run = types.BoolValue(challenge.Features.Run)
	state.RunBeforeSubmit = types.BoolValue(challenge.Features.RunBeforeSubmit)
	state.RunCommand = types.StringValue(challenge.Settings.RunCommand)
	state.TerminalCommand = types.StringValue(challenge.Settings.TerminalCommand)
	state.TestCommand = types.StringValue(challenge.Settings.CheckCommand)
	state.Terminal = types.BoolValue(challenge.Features.Terminal)
	var cur_tests []resourceclients.TestCase
	json.NewDecoder(strings.NewReader(state.TestcaseJSON.ValueString())).Decode(&cur_tests)
	if !compareTestCase(cur_tests, challenge.Tickets.MarkStandard.Testcases) {
		// Mismatching test case data.
		testcase, err := json.MarshalIndent(challenge.Tickets.MarkStandard.Testcases, "", "  ")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Testcases Object",
				fmt.Sprintf("Could not read Test cases from Slide ID %d: %s", state.SlideId.ValueInt64(), err.Error()),
			)
		}
		state.TestcaseJSON = types.StringValue(string(testcase))
	}
	state.TestcaseEasy = types.BoolValue(challenge.Tickets.MarkStandard.Easy)
	state.TestcaseMarkAll = types.BoolValue(challenge.Tickets.MarkStandard.MarkAll)
	state.TestcaseOverlayTestFiles = types.BoolValue(challenge.Tickets.MarkStandard.Overlay)
	challenge.Tickets.MarkStandard.RunLimit.Pty.If(func(val bool) { state.TestcasePty = types.BoolValue(val) })
	state.Type = types.StringValue(challenge.Type)

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

	api_obj, rubric, err := plan.MapAPIObj(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Challenge Object",
			fmt.Sprintf("Could not update Challenge: %s", err.Error()),
		)
		return
	}

	resourceclients.UpdateChallenge(r.client, plan.FolderPath.ValueString(), api_obj, rubric)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *challengeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *challengeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier to be two space separated integers: lesson_id,slide_id. Got: %q", req.ID),
		)
		return
	}
	lesson_id, err := strconv.Atoi(idParts[0])
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier to be two space separated integers: lesson_id,slide_id. Got: %q", req.ID),
		)
		return
	}
	slide_id, err := strconv.Atoi(idParts[1])
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier to be two space separated integers: lesson_id,slide_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("lesson_id"), lesson_id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("slide_id"), slide_id)...)
}
