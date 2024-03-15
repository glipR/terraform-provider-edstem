package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/md2ed"
	"terraform-provider-edstem/internal/resourceclients"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &questionResource{}
	_ resource.ResourceWithConfigure = &questionResource{}
)

// NewQuestionResource is a helper function to simplify the provider implementation.
func NewQuestionResource() resource.Resource {
	return &questionResource{}
}

// questionResource is the resource implementation.
type questionResource struct {
	client *client.Client
}

// Configure adds the provider configured client to the resource.
func (r *questionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *questionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_question"
}

type questionResourceModel struct {
	Id            types.Int64 `tfsdk:"id"`
	Index         types.Int64 `tfsdk:"index"`
	LessonSlideId types.Int64 `tfsdk:"lesson_slide_id"`
	AutoPoints    types.Int64 `tfsdk:"auto_points"`

	Type types.String `tfsdk:"type"`
	// Option 1 - Specify everything in Terraform
	Answers     types.List   `tfsdk:"answers"`
	Content     types.String `tfsdk:"content"`
	Explanation types.String `tfsdk:"explanation"`
	Solution    types.List   `tfsdk:"solution"`

	// Option 2 - Specify everything in a single string by starting lines with !content, !answer-1, ...
	QuestionDocumentString types.String `tfsdk:"question_document_string"`

	Formatted         types.Bool `tfsdk:"formatted"`
	MultipleSelection types.Bool `tfsdk:"multiple_selection"`
}

// Schema defines the schema for the resource.
func (r *questionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"index": schema.Int64Attribute{
				Required: true,
			},
			"lesson_slide_id": schema.Int64Attribute{
				Required: true,
			},
			"auto_points": schema.Int64Attribute{
				Default:  int64default.StaticInt64(1),
				Optional: true,
				Computed: true,
			},
			"type": schema.StringAttribute{
				Required: true,
				// TODO: Validate
			},
			"answers": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"content": schema.StringAttribute{
				Optional: true,
			},
			"explanation": schema.StringAttribute{
				Optional: true,
			},
			"solution": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
			},
			"question_document_string": schema.StringAttribute{
				Optional: true,
			},
			"formatted": schema.BoolAttribute{
				Default:  booldefault.StaticBool(true),
				Optional: true,
				Computed: true,
			},
			"multiple_selection": schema.BoolAttribute{
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (model *questionResourceModel) MapAPIObj(ctx context.Context) (*resourceclients.Question, error) {
	var obj resourceclients.Question

	obj.Id = model.Id.ValueInt64()
	if !model.Index.IsNull() {
		obj.Index.Set(model.Index.ValueInt64())
	}
	obj.LessonSlideId = model.LessonSlideId.ValueInt64()
	obj.AutoPoints = model.AutoPoints.ValueInt64()

	obj.Type = model.Type.ValueString()

	if model.QuestionDocumentString.IsNull() {

		tmp := make([]types.String, 0, len(model.Answers.Elements()))
		model.Answers.ElementsAs(ctx, &tmp, false)
		obj.Answers = make([]string, len(tmp), len(tmp))
		for i := range tmp {
			obj.Answers[i] = tmp[i].ValueString()
		}

		obj.Content.Set(model.Content.ValueString())
		obj.Explanation.Set(model.Explanation.ValueString())

		tmp2 := make([]types.Int64, 0, len(model.Solution.Elements()))
		model.Solution.ElementsAs(ctx, &tmp2, false)
		obj.Solution = make([]int, len(tmp2), len(tmp2))
		for i := range tmp2 {
			obj.Solution[i] = int(tmp2[i].ValueInt64())
		}
	} else {

		docstring := strings.ReplaceAll(model.QuestionDocumentString.ValueString(), "\r", "")
		docstring = "\n!nothing\n" + docstring

		split := strings.Split(docstring, "\n!")

		answer_counter := 0

		for i, s := range split {
			if i <= 1 {
				continue
			}
			sep := strings.SplitAfterN(s, "\n", 2)
			if sep[0] == "content\n" {
				obj.Content.Set(md2ed.RenderMDToEd(strings.TrimSpace(sep[1])))
			} else if sep[0] == "explanation\n" {
				obj.Explanation.Set(md2ed.RenderMDToEd(strings.TrimSpace(sep[1])))
			} else if strings.HasPrefix(sep[0], "answer") {
				answer_split := strings.Split(sep[0], "-")
				obj.Answers = append(obj.Answers, md2ed.RenderMDToEd(strings.TrimSpace(sep[1])))
				if len(answer_split) > 1 {
					obj.Solution = append(obj.Solution, answer_counter)
				}
				answer_counter++
			} else {
				return nil, errors.New(fmt.Sprintf("Unmatched exclamation line: %s:%s", sep[0], sep[1]))
			}
		}
	}

	obj.Formatted = model.Formatted.ValueBool()
	obj.MultipleSelection = model.MultipleSelection.ValueBool()

	return &obj, nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *questionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan questionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api_obj, err := plan.MapAPIObj(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Question Object",
			fmt.Sprintf("Could not create Question: %s", err.Error()),
		)
		return
	}

	err = resourceclients.CreateQuestion(r.client, api_obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Question Object",
			fmt.Sprintf("Could not create Question for Lesson Slide ID %d: %s", api_obj.LessonSlideId, err.Error()),
		)
		return
	}

	plan.Id = types.Int64Value(int64(api_obj.Id))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *questionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state questionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := resourceclients.GetQuestion(r.client, int(state.LessonSlideId.ValueInt64()), int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Question Object",
			fmt.Sprintf("Could not read Question ID %d: %s", state.Id.ValueInt64(), err.Error()),
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
func (r *questionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan questionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api_obj, err := plan.MapAPIObj(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Question Object",
			fmt.Sprintf("Could not update Question: %s", err.Error()),
		)
	}

	err = resourceclients.UpdateMultichoiceQuestion(r.client, api_obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Question Object",
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
func (r *questionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
