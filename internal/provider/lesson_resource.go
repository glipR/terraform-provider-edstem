package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/resourceclients"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	Id                                   types.Int64  `tfsdk:"id"`
	Attempts                             types.Int64  `tfsdk:"attempts"`
	AvailableAt                          types.String `tfsdk:"available_at"`
	DueAt                                types.String `tfsdk:"due_at"`
	GradePassbackAutoSend                types.Bool   `tfsdk:"grade_passback_auto_send"`
	GradePassbackMode                    types.String `tfsdk:"grade_passback_mode"`
	GradePassbackScaleTo                 types.String `tfsdk:"grade_passback_scale_to"`
	Index                                types.Int64  `tfsdk:"index"`
	IsHidden                             types.Bool   `tfsdk:"is_hidden"`
	IsTimed                              types.Bool   `tfsdk:"is_timed"`
	IsUnlisted                           types.Bool   `tfsdk:"is_unlisted"`
	Kind                                 types.String `tfsdk:"kind"`
	LateSubmissions                      types.Bool   `tfsdk:"late_submissions"`
	LockedAt                             types.String `tfsdk:"locked_at"`
	ModuleId                             types.Int64  `tfsdk:"module_id"`
	Openable                             types.Bool   `tfsdk:"openable"`
	OpenableWithoutAttempt               types.Bool   `tfsdk:"openable_without_attempt"`
	Outline                              types.String `tfsdk:"outline"`
	Password                             types.String `tfsdk:"password"`
	Prerequisites                        types.List   `tfsdk:"prerequisites"`
	ReleaseChallengeSolutions            types.Bool   `tfsdk:"release_challenge_solutions"`
	ReleaseChallengeSolutionsWhileActive types.Bool   `tfsdk:"release_challenge_solutions_while_active"`
	ReleaseFeedback                      types.Bool   `tfsdk:"release_feedback"`
	ReleaseFeedbackWhileActive           types.Bool   `tfsdk:"release_feedback_while_active"`
	ReleaseQuizCorrectnessOnly           types.Bool   `tfsdk:"release_quiz_correctness_only"`
	ReleaseQuizSolutions                 types.Bool   `tfsdk:"release_quiz_solutions"`
	ReOpenSubmissions                    types.Bool   `tfsdk:"reopen_submissions"`
	RequireUserOverride                  types.Bool   `tfsdk:"require_user_override"`
	QuizActiveStatus                     types.String `tfsdk:"quiz_active_status"`
	QuizMode                             types.String `tfsdk:"quiz_mode"`
	QuizQuestionNumberStyle              types.String `tfsdk:"quiz_question_number_style"`
	SolutionsAt                          types.String `tfsdk:"solutions_at"`
	State                                types.String `tfsdk:"state"`
	TimerDuration                        types.Int64  `tfsdk:"timer_duration"`
	TimerExpirationAccess                types.Bool   `tfsdk:"timer_expiration_access"`
	Title                                types.String `tfsdk:"title"`
	TutorialRegex                        types.String `tfsdk:"tutorial_regex"`
	Type                                 types.String `tfsdk:"type"`
	LastUpdated                          types.String `tfsdk:"last_updated"`
}

// Schema defines the schema for the resource.
func (r *lessonResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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
				Default:  stringdefault.StaticString("content"),
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
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (model *lessonResourceModel) MapAPIObj(ctx context.Context) resourceclients.Lesson {
	var obj resourceclients.Lesson
	if !model.Attempts.IsNull() {
		obj.Attempts.Set(int(model.Attempts.ValueInt64()))
	}
	if !model.AvailableAt.IsNull() {
		obj.AvailableAt.Set(model.AvailableAt.ValueString())
	}
	if !model.DueAt.IsNull() {
		obj.DueAt.Set(model.DueAt.ValueString())
	}
	obj.GradePassbackAutoSend = model.GradePassbackAutoSend.ValueBool()
	obj.GradePassbackMode = model.GradePassbackMode.String()
	if !model.GradePassbackScaleTo.IsNull() {
		obj.GradePassbackScaleTo.Set(model.GradePassbackScaleTo.ValueString())
	}
	obj.Id = int(model.Id.ValueInt64())
	if !model.Index.IsNull() {
		obj.Index.Set(int(model.Index.ValueInt64()))
	}
	obj.IsHidden = model.IsHidden.ValueBool()
	obj.IsTimed = model.IsTimed.ValueBool()
	obj.IsUnlisted = model.IsUnlisted.ValueBool()

	obj.Kind = model.Kind.ValueString()
	obj.LateSubmissions = model.LateSubmissions.ValueBool()

	if !model.LockedAt.IsNull() {
		obj.LockedAt.Set(model.LockedAt.ValueString())
	}
	if !model.ModuleId.IsNull() {
		obj.ModuleId.Set(int(model.ModuleId.ValueInt64()))
	}

	obj.Openable = model.Openable.ValueBool()
	obj.OpenableWithoutAttempt = model.OpenableWithoutAttempt.ValueBool()
	obj.Outline = model.Outline.ValueString()
	obj.Password = model.Password.ValueString()
	if !model.Prerequisites.IsNull() {
		obj.Prerequisites = make([]resourceclients.Prerequisite, 0, len(model.Prerequisites.Elements()))
		temp_iterable := make([]types.Int64, 0, len(model.Prerequisites.Elements()))
		// TODO: Error handle
		model.Prerequisites.ElementsAs(ctx, &temp_iterable, false)
		for i := range temp_iterable {
			obj.Prerequisites[i].RequiredLessonId = int(temp_iterable[i].ValueInt64())
		}
	}
	obj.ReleaseChallengeSolutions = model.ReleaseChallengeSolutions.ValueBool()
	obj.ReleaseChallengeSolutionsWhileActive = model.ReleaseChallengeSolutionsWhileActive.ValueBool()
	obj.ReleaseFeedback = model.ReleaseFeedback.ValueBool()
	obj.ReleaseFeedbackWhileActive = model.ReleaseFeedback.ValueBool()
	obj.ReleaseQuizCorrectnessOnly = model.ReleaseQuizCorrectnessOnly.ValueBool()
	obj.ReleaseQuizSolutions = model.ReleaseQuizSolutions.ValueBool()
	obj.ReOpenSubmissions = model.ReOpenSubmissions.ValueBool()
	obj.RequireUserOverride = model.RequireUserOverride.ValueBool()

	obj.QuizSettings = resourceclients.QuizSettings{
		QuizActiveStatus:        model.QuizActiveStatus.ValueString(),
		QuizMode:                model.QuizMode.ValueString(),
		QuizQuestionNumberStyle: model.QuizQuestionNumberStyle.ValueString(),
	}

	if !model.SolutionsAt.IsNull() {
		obj.SolutionsAt.Set(model.SolutionsAt.ValueString())
	}
	obj.State = model.State.ValueString()
	obj.TimerDuration = int(model.TimerDuration.ValueInt64())
	obj.TimerExpirationAccess = model.TimerExpirationAccess.ValueBool()
	obj.Title = model.Title.ValueString()
	obj.TutorialRegex = model.TutorialRegex.ValueString()
	obj.Type = model.Type.ValueString()

	return obj
}

// Create creates the resource and sets the initial Terraform state.
func (r *lessonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan lessonResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api_obj := plan.MapAPIObj(ctx)

	resourceclients.CreateLesson(r.client, &api_obj)

	plan.Id = types.Int64Value(int64(api_obj.Id))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC1123Z))

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

	lesson, err := resourceclients.GetLesson(r.client, int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lesson Object",
			fmt.Sprintf("Could not read Lesson ID %d: %s", state.Id.ValueInt64(), err.Error()),
		)
	}

	lesson.Attempts.If(func(val int) { state.Attempts = types.Int64Value(int64(val)) })
	lesson.AvailableAt.If(func(val string) { state.AvailableAt = types.StringValue(val) })
	lesson.DueAt.If(func(val string) { state.DueAt = types.StringValue(val) })
	state.GradePassbackAutoSend = types.BoolValue(lesson.GradePassbackAutoSend)
	state.GradePassbackMode = types.StringValue(lesson.GradePassbackMode)
	lesson.GradePassbackScaleTo.If(func(val string) { state.GradePassbackScaleTo = types.StringValue(val) })
	lesson.Index.If(func(val int) { state.Index = types.Int64Value(int64(val)) })
	state.IsHidden = types.BoolValue(lesson.IsHidden)
	state.IsTimed = types.BoolValue(lesson.IsTimed)
	state.IsUnlisted = types.BoolValue(lesson.IsUnlisted)
	state.Kind = types.StringValue(lesson.Kind)
	state.LateSubmissions = types.BoolValue(lesson.LateSubmissions)
	lesson.LockedAt.If(func(val string) { state.LockedAt = types.StringValue(val) })
	lesson.ModuleId.If(func(val int) { state.ModuleId = types.Int64Value(int64(val)) })
	state.Openable = types.BoolValue(lesson.Openable)
	state.OpenableWithoutAttempt = types.BoolValue(lesson.OpenableWithoutAttempt)
	state.Outline = types.StringValue(lesson.Outline)
	state.Password = types.StringValue(lesson.Password)
	state.QuizActiveStatus = types.StringValue(lesson.QuizSettings.QuizActiveStatus)
	state.QuizMode = types.StringValue(lesson.QuizSettings.QuizMode)
	state.QuizQuestionNumberStyle = types.StringValue(lesson.QuizSettings.QuizQuestionNumberStyle)
	state.ReOpenSubmissions = types.BoolValue(lesson.ReOpenSubmissions)
	state.ReleaseChallengeSolutions = types.BoolValue(lesson.ReleaseChallengeSolutions)
	state.ReleaseChallengeSolutionsWhileActive = types.BoolValue(lesson.ReleaseChallengeSolutionsWhileActive)
	state.ReleaseFeedback = types.BoolValue(lesson.ReleaseFeedback)
	state.ReleaseFeedbackWhileActive = types.BoolValue(lesson.ReleaseFeedbackWhileActive)
	state.ReleaseQuizCorrectnessOnly = types.BoolValue(lesson.ReleaseQuizCorrectnessOnly)
	state.ReleaseQuizSolutions = types.BoolValue(lesson.ReleaseQuizSolutions)
	state.RequireUserOverride = types.BoolValue(lesson.RequireUserOverride)
	lesson.SolutionsAt.If(func(val string) { state.SolutionsAt = types.StringValue(val) })
	state.State = types.StringValue(lesson.State)
	state.TimerDuration = types.Int64Value(int64(lesson.TimerDuration))
	state.TimerExpirationAccess = types.BoolValue(lesson.TimerExpirationAccess)
	state.Title = types.StringValue(lesson.Title)
	state.TutorialRegex = types.StringValue(lesson.TutorialRegex)
	state.Type = types.StringValue(lesson.Type)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *lessonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan lessonResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api_obj := plan.MapAPIObj(ctx)

	err := resourceclients.UpdateLesson(r.client, &api_obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Lesson Object",
			fmt.Sprintf("Could not update Lesson ID %d: %s", api_obj.Id, err.Error()),
		)
	}

	plan.Id = types.Int64Value(int64(api_obj.Id))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *lessonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *lessonResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	lesson_id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier to be integer. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), lesson_id)...)
}
