package provider

import (
	"context"
	"fmt"
	"terraform-provider-edstem/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &lessonDataSource{}
	_ datasource.DataSourceWithConfigure = &lessonDataSource{}
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &lessonDataSource{}
)

func NewLessonDataSource() datasource.DataSource {
	return &lessonDataSource{}
}

type lessonDataSource struct {
	client *client.Client
	id     int
}

type lessonDataSourceModel struct {
	AvailableAt   types.String              `tfsdk:"available_at"`
	DueAt         types.String              `tfsdk:"due_at"`
	ModuleId      types.Int64               `tfsdk:"module_id"`
	Title         types.String              `tfsdk:"title"`
	ID            types.Int64               `tfsdk:"id"`
	QuizSettings  *lessonQuizSettingsModel  `tfsdk:"quiz_settings"`
	Prerequisites []lessonPrerequisiteModel `tfsdk:"prerequisites"`
}

type lessonQuizSettingsModel struct {
	ActiveStatus        types.String `tfsdk:"active_status"`
	Mode                types.String `tfsdk:"mode"`
	QuestionNumberStyle types.String `tfsdk:"question_number_style"`
}

type lessonPrerequisiteModel struct {
	RequiredLessonId types.Int64 `tfsdk:"required_lesson_id"`
}

// Metadata returns the data source type name.
func (d *lessonDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lesson"
}

// Schema defines the schema for the data source.
func (d *lessonDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"available_at": schema.StringAttribute{
				Computed: true,
			},
			"due_at": schema.StringAttribute{
				Computed: true,
			},
			"module_id": schema.Int64Attribute{
				Computed: true,
				Optional: true,
			},
			"title": schema.StringAttribute{
				Computed: true,
			},
			"quiz_settings": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"active_status": schema.StringAttribute{
						Computed: true,
					},
					"mode": schema.StringAttribute{
						Computed: true,
					},
					"question_number_style": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"prerequisites": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"required_lesson_id": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *lessonDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *lessonDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lessonDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	lesson, err := d.client.GetLesson(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read Ed Lesson with ID %d", d.id),
			err.Error(),
		)
		return
	}

	lesson.AvailableAt.If(func(val string) { state.AvailableAt = types.StringValue(val) })
	lesson.DueAt.If(func(val string) { state.DueAt = types.StringValue(val) })
	lesson.ModuleId.If(func(val int) { state.ModuleId = types.Int64Value(int64(val)) })
	state.Title = types.StringValue(lesson.Title)
	state.QuizSettings = &lessonQuizSettingsModel{}
	state.QuizSettings.ActiveStatus = types.StringValue(lesson.QuizSettings.QuizActiveStatus)
	state.QuizSettings.Mode = types.StringValue(lesson.QuizSettings.QuizMode)
	state.QuizSettings.QuestionNumberStyle = types.StringValue(lesson.QuizSettings.QuizQuestionNumberStyle)
	for _, prereq := range lesson.Prerequisites {
		state.Prerequisites = append(state.Prerequisites, lessonPrerequisiteModel{
			RequiredLessonId: types.Int64Value(int64(prereq.RequiredLessonId)),
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
