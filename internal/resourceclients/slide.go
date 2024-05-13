package resourceclients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/md2ed"
	"terraform-provider-edstem/internal/tfhelpers"

	"github.com/markphelps/optional"
)

type Slide struct {
	Id          int             `json:"id"`
	Type        string          `json:"type"`
	Title       string          `json:"title"`
	Index       int             `json:"index"`
	IsHidden    bool            `json:"is_hidden"`
	Content     string          `json:"content"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   optional.String `json:"updated_at"`
	LessonId    int             `json:"lesson_id"`
	CourseId    int             `json:"course_id"`
	UserId      int             `json:"user_id"`
	ChallengeId optional.Int    `json:"challenge_id"`
	FileUrl     optional.String `json:"file_url"`
	VideoUrl    optional.String `json:"video_url"`
	Url         optional.String `json:"url"`
	Html        optional.String `json:"html"`
}

type SlideResponse struct {
	Id           int             `json:"id"`
	CourseId     int             `json:"course_id"`
	UserId       int             `json:"user_id"`
	LessonId     int             `json:"lesson_id"`
	Type         string          `json:"type"`
	Title        string          `json:"title"`
	Index        int             `json:"index"`
	IsHidden     bool            `json:"is_hidden"`
	Content      string          `json:"content"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    optional.String `json:"updated_at"`
	RubricPoints optional.Int    `json:"rubric_points"`
	AutoPoints   optional.Int    `json:"auto_points"`
	ChallengeId  optional.Int    `json:"challenge_id"`
	FileUrl      optional.String `json:"file_url"`
	VideoUrl     optional.String `json:"video_url"`
	Url          optional.String `json:"url"`
	Html         optional.String `json:"html"`
}

type SlideCreateRequest struct {
	Type string `json:"type"`
}

type SlideUpdateRequest struct {
	Id           int             `json:"id"`
	CourseId     int             `json:"course_id"`
	UserId       int             `json:"user_id"`
	LessonId     int             `json:"lesson_id"`
	Type         string          `json:"type"`
	Title        string          `json:"title"`
	Index        int             `json:"index"`
	IsHidden     bool            `json:"is_hidden"`
	Content      string          `json:"content"`
	RubricPoints optional.Int    `json:"rubric_points"`
	AutoPoints   optional.Int    `json:"auto_points"`
	VideoUrl     optional.String `json:"video_url"`
	Url          optional.String `json:"url"`
	Html         optional.String `json:"html"`
}

type LessonWithSlidesResponse struct {
	Lesson LessonSlidesObj `json:"lesson"`
}

type LessonSlidesObj struct {
	Slides []SlideResponse `json:"slides"`
}

type SlideEditResponse struct {
	Slide SlideResponse `json:"slide"`
}

func GetSlide(c *client.Client, lesson_id int, slide_id int) (*Slide, error) {
	body, err := c.HTTPRequest(fmt.Sprintf("lessons/%d?view=1", lesson_id), "GET", bytes.Buffer{}, nil)
	if err != nil {
		return nil, err
	}
	resp := &LessonWithSlidesResponse{}
	err = json.NewDecoder(body).Decode(resp)
	if err != nil {
		return nil, err
	}
	slides := resp.Lesson.Slides
	for i := range slides {
		if slides[i].Id == slide_id {
			var slideObj Slide
			tempVar, _ := json.Marshal(slides[i])
			err = json.Unmarshal(tempVar, &slideObj)
			if err != nil {
				return nil, err
			}
			// Not returned from response - intuited from request.
			slideObj.LessonId = lesson_id
			return &slideObj, nil
		}
	}
	return nil, fmt.Errorf("Slide ID %d Not Found", slide_id)
}

func GetSlideIds(c *client.Client, lesson_id int) ([]int, error) {
	body, err := c.HTTPRequest(fmt.Sprintf("lessons/%d?view=1", lesson_id), "GET", bytes.Buffer{}, nil)
	if err != nil {
		return nil, err
	}
	resp := &LessonWithSlidesResponse{}
	err = json.NewDecoder(body).Decode(resp)
	if err != nil {
		return nil, err
	}
	slides := resp.Lesson.Slides

	final := make([]int, len(slides))
	for i := range slides {
		final[i] = slides[i].Id
	}

	return final, nil
}

func UpdateSlide(c *client.Client, slide *Slide) error {
	request := &SlideUpdateRequest{}
	request.Content = slide.Content
	request.Id = slide.Id
	request.Index = slide.Index
	request.IsHidden = slide.IsHidden
	request.Title = slide.Title
	request.Type = slide.Type
	request.CourseId = slide.CourseId
	request.LessonId = slide.LessonId
	request.UserId = slide.UserId
	slide.VideoUrl.If(func(val string) { request.VideoUrl.Set(val) })
	slide.Url.If(func(val string) { request.Url.Set(val) })
	slide.Html.If(func(val string) { request.Html.Set(val) })
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	var boundary string
	var actual_req bytes.Buffer
	if slide.Type != "pdf" {
		boundary = fmt.Sprintf("-----------------------------%s", "28191803313191638583308257490")
		req_text := fmt.Sprintf("--%s\nContent-Disposition: form-data; name=\"slide\"\n\n%s--%s--\n", boundary, buf.String(), boundary)
		actual_req = bytes.Buffer{}
		actual_req.Write([]byte(req_text))
	} else {
		boundary = fmt.Sprintf("-----------------------------%s", "303367121714237365713833509663")
		filename := slide.FileUrl.MustGet()
		filedata, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		req_text := fmt.Sprintf("--%s\nContent-Disposition: form-data; name=\"attachment\"; filename=\"%s\"\nContent-Type: application/pdf\n\n%s\n\n--%s\nContent-Disposition: form-data; name=\"slide\"\n\n%s--%s--\n", boundary, filename, filedata, boundary, buf.String(), boundary)
		actual_req = bytes.Buffer{}
		actual_req.Write([]byte(req_text))
	}

	body, err := c.HTTPRequest(fmt.Sprintf("lessons/slides/%d", slide.Id), "PUT", actual_req, &boundary)
	if err != nil {
		return err
	}
	resp_lesson := &SlideEditResponse{}
	err = json.NewDecoder(body).Decode(resp_lesson)
	if err != nil {
		return err
	}
	slide.Id = resp_lesson.Slide.Id

	// Reordering slides if necessary
	slide_ids, err := GetSlideIds(c, slide.LessonId)
	if err != nil {
		return err
	}

	if (slide.Index - 1) < len(slide_ids) {
		if slide_ids[slide.Index-1] != slide.Id {
			// Reorder
			past_point := 0
			if slide.Index != len(slide_ids) {
				past_point = slide_ids[slide.Index]
			}
			if past_point == slide.Id {
				_, err := c.HTTPRequest(fmt.Sprintf("lessons/slides/%d/reorder/%d", slide.Id, slide_ids[slide.Index-1]), "PUT", bytes.Buffer{}, nil)
				if err != nil {
					return err
				}
			} else {

				// reorder slide_ids[slide.Index-1] to before slide.Id
				_, err := c.HTTPRequest(fmt.Sprintf("lessons/slides/%d/reorder/%d", slide_ids[slide.Index-1], slide.Id), "PUT", bytes.Buffer{}, nil)
				if err != nil {
					return err
				}
				// reorder slide.Id to before past_point
				_, err = c.HTTPRequest(fmt.Sprintf("lessons/slides/%d/reorder/%d", slide.Id, past_point), "PUT", bytes.Buffer{}, nil)
				if err != nil {
					return err
				}
			}
		}
	} // Nothing we can do otherwise - wrong spot.

	return nil
}

// The three lines
//-----------------------------303367121714237365713833509663
//-----------------------------303367121714237365713833509663
//-----------------------------303367121714237365713833509663--
//---------------------------303367121714237365713833509663

func CreateSlide(c *client.Client, slide *Slide) error {
	request := &SlideCreateRequest{}
	request.Type = slide.Type

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	boundary := fmt.Sprintf("-----------------------------%s", "264592028829639346041448524574")
	req_text := fmt.Sprintf("--%s\nContent-Disposition: form-data; name=\"slide\"\n\n%s--%s--\n", boundary, buf.String(), boundary)
	actual_req := bytes.Buffer{}
	actual_req.Write([]byte(req_text))
	fmt.Println(req_text)

	body, err := c.HTTPRequest(fmt.Sprintf("lessons/%d/slides", slide.LessonId), "POST", actual_req, &boundary)
	if err != nil {
		return err
	}
	resp_lesson := &SlideEditResponse{}
	err = json.NewDecoder(body).Decode(resp_lesson)
	if err != nil {
		return err
	}
	slide.Id = resp_lesson.Slide.Id
	slide.CreatedAt = resp_lesson.Slide.CreatedAt
	slide.Index = resp_lesson.Slide.Index
	slide.LessonId = resp_lesson.Slide.LessonId
	slide.CourseId = resp_lesson.Slide.CourseId
	slide.UserId = resp_lesson.Slide.UserId
	return UpdateSlide(c, slide)
}

func SlideToTerraform(c *client.Client, lesson_id int, slide_id int, resource_name string, folder_path string, parent_resource_name *string) (string, []string, error) {
	slide, err := GetSlide(c, lesson_id, slide_id)
	if err != nil {
		return "", []string{}, err
	}
	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(slide)
	if err != nil {
		return "", []string{}, err
	}

	resources := make([]string, 0)
	resources = append(resources, fmt.Sprintf("edstem_slide.%s %d,%d", resource_name, lesson_id, slide_id))

	var resource_string = fmt.Sprintf("resource \"edstem_slide\" %s {\n", resource_name)
	resource_string = resource_string + tfhelpers.TFProp("id", slide.Id, nil)
	resource_string = resource_string + tfhelpers.TFProp("type", slide.Type, nil)
	if parent_resource_name != nil {
		resource_string = resource_string + tfhelpers.TFUnquote("lesson_id", fmt.Sprintf("edstem_lesson.%s.id", *parent_resource_name))
	} else {
		resource_string = resource_string + tfhelpers.TFProp("lesson_id", slide.LessonId, nil)
	}

	resource_string = resource_string + tfhelpers.TFProp("title", slide.Title, "")
	resource_string = resource_string + tfhelpers.TFProp("index", slide.Index, nil)
	resource_string = resource_string + tfhelpers.TFProp("is_hidden", slide.IsHidden, false)
	if slide.Type == "video" {
		resource_string = resource_string + tfhelpers.TFProp("url", slide.VideoUrl, nil)
	} else if slide.Type == "webpage" {
		resource_string = resource_string + tfhelpers.TFProp("url", slide.Url, nil)
	} else if slide.Type == "html" {
		if slide.Html.Present() {
			content_path := path.Join(folder_path, "content.html")
			resource_string = resource_string + tfhelpers.TFFile("content", slide.Html.MustGet(), content_path)
		}
	}
	if slide.Content != "" {
		content_path := path.Join(folder_path, "content.md")
		resource_string = resource_string + tfhelpers.TFFile("content", md2ed.RenderEdToMD(slide.Content, folder_path, true), content_path)
	}
	resource_string = resource_string + "}"

	if slide.Type == "code" {
		s, challenge_resources, e := ChallengeToTerraform(c, lesson_id, slide_id, fmt.Sprintf("%s_challenge", resource_name), folder_path, &resource_name, parent_resource_name)
		if e != nil {
			return "", []string{}, e
		}
		resource_string = resource_string + "\n\n" + s
		resources = append(resources, challenge_resources...)
	}

	return resource_string, resources, nil
}
