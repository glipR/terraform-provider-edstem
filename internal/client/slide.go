package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/markphelps/optional"
)

type Slide struct {
	Id        int             `json:"id"`
	Type      string          `json:"type"`
	Title     string          `json:"title"`
	Index     int             `json:"index"`
	IsHidden  bool            `json:"is_hidden"`
	Content   string          `json:"content"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt optional.String `json:"updated_at"`
	LessonId  int             `json:"lesson_id"`
}

type SlideResponse struct {
	Id        int             `json:"id"`
	Type      string          `json:"type"`
	Title     string          `json:"title"`
	Index     int             `json:"index"`
	IsHidden  bool            `json:"is_hidden"`
	Content   string          `json:"content"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt optional.String `json:"updated_at"`
}

type SlideCreateRequest struct {
	Type string `json:"type"`
}

type LessonWithSlidesResponse struct {
	Lesson LessonSlidesObj `json:"lesson"`
}

type LessonSlidesObj struct {
	Slides []SlideResponse `json:"slides"`
}

func (c *Client) GetSlide(lesson_id int, slide_id int) (*Slide, error) {
	body, err := c.httpRequest(fmt.Sprintf("lessons/%d?view=1", lesson_id), "GET", bytes.Buffer{}, nil)
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

func (c *Client) UpdateSlide(slide *Slide) error {
	request := &SlideResponse{}
	request.Content = slide.Content
	request.CreatedAt = slide.CreatedAt
	request.Id = slide.Id
	request.Index = slide.Index
	request.IsHidden = slide.IsHidden
	request.Title = slide.Title
	request.Type = slide.Type
	slide.UpdatedAt.If(func(s string) { request.UpdatedAt.Set(s) })
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	boundary := fmt.Sprintf("-----------------------------%s", "264592028829639346041448524574")
	req_text := fmt.Sprintf("%s\nContent-Disposition: form-data; name=\"slide\"\n\n%s%s--\n", boundary, buf.String(), boundary)
	fmt.Printf(req_text)
	actual_req := bytes.Buffer{}
	actual_req.Write([]byte(req_text))

	body, err := c.httpRequest(fmt.Sprintf("lessons/%d/slides/%d", slide.LessonId, slide.Id), "PUT", actual_req, &boundary)
	if err != nil {
		return err
	}
	resp_lesson := &LessonResponse{}
	err = json.NewDecoder(body).Decode(resp_lesson)
	if err != nil {
		return err
	}
	slide.Id = resp_lesson.LessonObj.Id
	return err
}

func (c *Client) CreateSlide(slide *Slide) error {
	request := &SlideCreateRequest{}
	request.Type = slide.Type

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	boundary := fmt.Sprintf("-----------------------------%s", "264592028829639346041448524574")
	req_text := fmt.Sprintf("%s\nContent-Disposition: form-data; name=\"slide\"\n\n%s%s--\n", boundary, buf.String(), boundary)
	actual_req := bytes.Buffer{}
	actual_req.Write([]byte(req_text))
	fmt.Println(req_text)

	body, err := c.httpRequest(fmt.Sprintf("lessons/%d/slides", slide.LessonId), "POST", actual_req, &boundary)
	if err != nil {
		return err
	}
	resp_lesson := &LessonResponse{}
	err = json.NewDecoder(body).Decode(resp_lesson)
	if err != nil {
		return err
	}
	slide.Id = resp_lesson.LessonObj.Id
	return c.UpdateSlide(slide)
}
