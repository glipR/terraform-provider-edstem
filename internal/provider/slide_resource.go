package provider

import (
	"context"
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
	FilePath    types.String `tfsdk:"file_path"`
	Url         types.String `tfsdk:"url"`
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
				MarkdownDescription: "Integer ID identifying the Slide. This can be found in the URL of a slide. For example, `https://edstem.org/au/courses/<course_id>/lessons/<lesson_id>/slides/<slide_id>`. Here we want the slide_id.",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "String identifying the type of slide. Options are `document`, `quiz`, `code`, `pdf`, `video`, `webpage`, `html`.",
				// TODO: Validate
			},
			"lesson_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Integer ID identifying the Lesson containing this slide. This can be found in the URL of a slide. For example, `https://edstem.org/au/courses/<course_id>/lessons/<lesson_id>/slides/<slide_id>`. Here we want the lesson_id.",
			},
			"title": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Title of the slide.",
			},
			"index": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Where this slide should slot within the slide list. 1 = first slide, 2 = second slide...",
			},
			"is_hidden": schema.BoolAttribute{
				Default:             booldefault.StaticBool(false),
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether this slide should be hidden from students.",
			},
			"content": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Content of the slide.",
			},
			"content_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("md"),
				MarkdownDescription: "Format of the slide content. Defaults to `md` (Markdown). Set to `ed` if you want to enter in the xml directly.",
			},
			"file_path": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The path for certain slide types to load content (like `video` or `pdf`)",
			},
			"url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The path for webpage slides to load from.",
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
	obj.FileUrl.Set(model.FilePath.ValueString())
	if model.Type.ValueString() == "video" {
		if !model.Url.IsNull() {
			obj.VideoUrl.Set(model.Url.ValueString())
		}
	} else if model.Type.ValueString() == "webpage" {
		if !model.Url.IsNull() {
			obj.Url.Set(model.Url.ValueString())
		}
	} else if model.Type.ValueString() == "html" {
		obj.Html.Set(model.Content.ValueString())
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

	slide, err := resourceclients.GetSlide(r.client, int(state.LessonId.ValueInt64()), int(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Slide Object",
			fmt.Sprintf("Could not read Slide ID %d: %s", state.Id.ValueInt64(), err.Error()),
		)
		return
	}

	state.Content = types.StringValue(slide.Content)
	if state.ContentType.ValueString() == "md" {
		state.Content = types.StringValue(md2ed.RenderEdToMD(state.Content.ValueString(), "", false))
	}
	// The index reported by the slide endpoint is wrong.
	// Should infer from the ordering in the lesson response instead.
	slide_ids, err := resourceclients.GetSlideIds(r.client, int(state.LessonId.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Slide Indexes",
			fmt.Sprintf("Could not read Lesson ID %d: %s", state.LessonId.ValueInt64(), err.Error()),
		)
	}
	for index, slide_id := range slide_ids {
		if slide_id == int(state.Id.ValueInt64()) {
			state.Index = types.Int64Value(int64(index + 1))
		}
	}

	state.IsHidden = types.BoolValue(slide.IsHidden)
	state.Title = types.StringValue(slide.Title)
	state.Type = types.StringValue(slide.Type)
	if slide.Type == "html" {
		slide.Html.If(func(val string) { state.Content = types.StringValue(val) })
	} else if slide.Type == "video" {
		slide.VideoUrl.If(func(val string) { state.Url = types.StringValue(val) })
	} else if slide.Type == "webpage" {
		slide.Url.If(func(val string) { state.Url = types.StringValue(val) })
	}

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

func (r *slideResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), slide_id)...)
}
