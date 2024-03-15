package resourceclients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"terraform-provider-edstem/internal/client"

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
}

type SlideCreateRequest struct {
	Type string `json:"type"`
}

type SlideUpdateRequest struct {
	Id           int          `json:"id"`
	CourseId     int          `json:"course_id"`
	UserId       int          `json:"user_id"`
	LessonId     int          `json:"lesson_id"`
	Type         string       `json:"type"`
	Title        string       `json:"title"`
	Index        int          `json:"index"`
	IsHidden     bool         `json:"is_hidden"`
	Content      string       `json:"content"`
	RubricPoints optional.Int `json:"rubric_points"`
	AutoPoints   optional.Int `json:"auto_points"`
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
	return nil, errors.New(fmt.Sprintf("Slide ID %d Not Found", slide_id))
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
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	boundary := fmt.Sprintf("-----------------------------%s", "28191803313191638583308257490")
	req_text := fmt.Sprintf("--%s\nContent-Disposition: form-data; name=\"slide\"\n\n%s--%s--\n", boundary, buf.String(), boundary)
	fmt.Printf(req_text)
	actual_req := bytes.Buffer{}
	actual_req.Write([]byte(req_text))

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
	return err
}

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
