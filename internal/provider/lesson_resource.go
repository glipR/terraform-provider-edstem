package provider

import (
	"context"
	"fmt"

	"terraform-provider-edstem/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
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
	reopen_submissions                       types.Bool   `tfsdk:"reopen_submissions"`
	require_user_override                    types.Bool   `tfsdk:"require_user_override"`
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
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"grade_passback_mode": schema.StringAttribute{
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"grade_passback_scale_to": schema.StringAttribute{
				Optional: true,
			},
			"index": schema.Int64Attribute{
				Optional: true,
			},
			"is_hidden": schema.BoolAttribute{
				Default:  booldefault.StaticBool(true),
				Optional: true,
				Computed: true,
			},
			"is_timed": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"is_unlisted": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"kind": schema.StringAttribute{
				Default:  stringdefault.StaticString("legacy"),
				Optional: true,
				Computed: true,
			},
			"late_submissions": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"locked_at": schema.StringAttribute{
				Optional: true,
			},
			"module_id": schema.Int64Attribute{
				Optional: true,
			},
			"openable": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"openable_without_attempt": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"outline": schema.StringAttribute{
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"password": schema.StringAttribute{
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"prerequisites": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
			},
			"release_challenge_solutions": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"release_challenge_solutions_while_active": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"release_feedback": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"release_feedback_while_active": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"release_quiz_correctness_only": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"release_quiz_solutions": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"reopen_submissions": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"require_user_override": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"quiz_active_status": schema.StringAttribute{
				Default:  stringdefault.StaticString("active"),
				Optional: true,
				Computed: true,
			},
			"quiz_mode": schema.StringAttribute{
				Default:  stringdefault.StaticString("multiple-attempts"),
				Optional: true,
				Computed: true,
			},
			"quiz_question_number_style": schema.StringAttribute{
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"solutions_at": schema.StringAttribute{
				Optional: true,
			},
			"state": schema.StringAttribute{
				Default:  stringdefault.StaticString("active"),
				Optional: true,
				Computed: true,
			},
			"timer_duration": schema.Int64Attribute{
				Default:  int64default.StaticInt64(60),
				Optional: true,
				Computed: true,
			},
			"timer_expiration_access": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"title": schema.StringAttribute{
				Required: true,
			},
			"tutorial_regex": schema.StringAttribute{
				Default:  stringdefault.StaticString(""),
				Optional: true,
				Computed: true,
			},
			"type": schema.StringAttribute{
				Default:  stringdefault.StaticString("general"),
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (model *lessonResourceModel) MapAPIObj(ctx context.Context) client.Lesson {
	var obj client.Lesson
	if !model.attempts.IsNull() {
		obj.Attempts.Set(int(model.attempts.ValueInt64()))
	}
	if !model.available_at.IsNull() {
		obj.AvailableAt.Set(model.available_at.ValueString())
	}
	if !model.due_at.IsNull() {
		obj.DueAt.Set(model.due_at.ValueString())
	}
	obj.GradePassbackAutoSend = model.grade_passback_auto_send.ValueBool()
	obj.GradePassbackMode = model.grade_passback_mode.String()
	if !model.grade_passback_scale_to.IsNull() {
		obj.GradePassbackScaleTo.Set(model.grade_passback_scale_to.ValueString())
	}
	obj.Id = int(model.id.ValueInt64())
	if !model.index.IsNull() {
		obj.Index.Set(int(model.index.ValueInt64()))
	}
	obj.IsHidden = model.is_hidden.ValueBool()
	obj.IsTimed = model.is_timed.ValueBool()
	obj.IsUnlisted = model.is_unlisted.ValueBool()

	obj.Kind = model.kind.ValueString()
	obj.LateSubmissions = model.late_submissions.ValueBool()

	if !model.locked_at.IsNull() {
		obj.LockedAt.Set(model.locked_at.ValueString())
	}
	if !model.module_id.IsNull() {
		obj.ModuleId.Set(int(model.module_id.ValueInt64()))
	}

	obj.Openable = model.openable.ValueBool()
	obj.OpenableWithoutAttempt = model.openable_without_attempt.ValueBool()
	obj.Outline = model.outline.ValueString()
	obj.Password = model.password.ValueString()
	if !model.prerequisites.IsNull() {
		obj.Prerequisites = make([]client.Prerequisite, 0, len(model.prerequisites.Elements()))
		temp_iterable := make([]types.Int64, 0, len(model.prerequisites.Elements()))
		// TODO: Error handle
		model.prerequisites.ElementsAs(ctx, &temp_iterable, false)
		for i := range temp_iterable {
			obj.Prerequisites[i].RequiredLessonId = int(temp_iterable[i].ValueInt64())
		}
	}
	obj.ReleaseChallengeSolutions = model.release_challenge_solutions.ValueBool()
	obj.ReleaseChallengeSolutionsWhileActive = model.release_challenge_solutions_while_active.ValueBool()
	obj.ReleaseFeedback = model.release_feedback.ValueBool()
	obj.ReleaseFeedbackWhileActive = model.release_feedback_while_active.ValueBool()
	obj.ReleaseQuizCorrectnessOnly = model.release_quiz_correctness_only.ValueBool()
	obj.ReleaseQuizSolutions = model.release_quiz_solutions.ValueBool()
	obj.ReOpenSubmissions = model.reopen_submissions.ValueBool()
	obj.RequireUserOverride = model.require_user_override.ValueBool()

	obj.QuizSettings = client.QuizSettings{
		QuizActiveStatus:        model.quiz_active_status.ValueString(),
		QuizMode:                model.quiz_mode.ValueString(),
		QuizQuestionNumberStyle: model.quiz_question_number_style.ValueString(),
	}

	if !model.solutions_at.IsNull() {
		obj.SolutionsAt.Set(model.solutions_at.ValueString())
	}
	obj.State = model.state.ValueString()
	obj.TimerDuration = int(model.timer_duration.ValueInt64())
	obj.TimerExpirationAccess = model.timer_expiration_access.ValueBool()
	obj.Title = model.title.ValueString()
	obj.TutorialRegex = model.tutorial_regex.ValueString()
	obj.Type = model.lesson_type.ValueString()

	return obj
}

// Create creates the resource and sets the initial Terraform state.
func (r *lessonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var test types.String
	diags := req.Plan.GetAttribute(ctx, path.Root("title"), &test)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fmt.Println(test.ValueString())

	// Retrieve values from plan
	var plan lessonResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api_obj := plan.MapAPIObj(ctx)

	r.client.CreateLesson(&api_obj)

	plan.id = types.Int64Value(int64(api_obj.Id))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *lessonResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state lessonResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.GetLesson(int(state.id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lesson Object",
			fmt.Sprintf("Could not read Lesson ID %d: %s", state.id.ValueInt64(), err.Error()),
		)
	}

	// TODO: For now, nothing happens with the read elements. Should update state.

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *lessonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *lessonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
