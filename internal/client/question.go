package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/markphelps/optional"
)

type Question struct {
	Id            int64          `json:"id"`
	Index         optional.Int64 `json:"index"`
	LessonSlideId int64          `json:"lesson_slide_id"`
	AutoPoints    int64          `json:"auto_points"`

	Type        string          `json:"type"`
	Answers     []string        `json:"answers"`
	Content     optional.String `json:"content"`
	Explanation optional.String `json:"explanation"`
	Solution    []int           `json:"solution"`

	Formatted         bool `json:"formatted"`
	MultipleSelection bool `json:"multiple_selection"`
}

type MultiChoiceQuestionResponse struct {
	Id            int64                   `json:"id"`
	Index         optional.Int64          `json:"index"`
	LessonSlideId int64                   `json:"lesson_slide_id"`
	AutoPoints    int64                   `json:"auto_points"`
	Data          MultiChoiceQuestionData `json:"data"`
}

type MultiChoiceQuestionRequest struct {
	Id            optional.Int64          `json:"id"`
	Index         optional.Int64          `json:"index"`
	LessonSlideId int64                   `json:"lesson_slide_id"`
	AutoPoints    int64                   `json:"auto_points"`
	Data          MultiChoiceQuestionData `json:"data"`
}

type MultiChoiceQuestionData struct {
	Answers           []string `json:"answers"`
	Assessed          bool     `json:"assessed"`
	Content           string   `json:"content"`
	Explanation       string   `json:"explanation"`
	Formatted         bool     `json:"formatted"`
	MultipleSelection bool     `json:"multiple_selection"`
	Solution          []int    `json:"solution"`
	Type              string   `json:"type"`
}

type MultiChoiceQuestionRequestActual struct {
	Question MultiChoiceQuestionRequest `json:"question"`
}
type MultiChoiceQuestionResponseActual struct {
	Question MultiChoiceQuestionResponse `json:"question"`
}

type QuestionReadResponse struct {
	Questions []MultiChoiceQuestionResponse `json:"questions"`
}

func (c *Client) GetQuestion(lesson_slide_id int, question_id int) (*Question, error) {
	body, err := c.httpRequest(fmt.Sprintf("lessons/slides/%d/questions", lesson_slide_id), "GET", bytes.Buffer{}, nil)
	if err != nil {
		return nil, err
	}
	resp := &QuestionReadResponse{}
	err = json.NewDecoder(body).Decode(resp)
	if err != nil {
		return nil, err
	}
	questions := resp.Questions
	for i := range questions {
		if questions[i].Id == int64(question_id) {
			var questionObj Question
			tempVar, _ := json.Marshal(questions[i])
			err = json.Unmarshal(tempVar, &questionObj)
			if err != nil {
				return nil, err
			}
			return &questionObj, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Question ID %d Not Found", question_id))
}

func (c *Client) UpdateMultichoiceQuestion(question *Question) error {
	request := &MultiChoiceQuestionRequest{}
	request.Id.Set(question.Id)
	request.Index = question.Index
	request.LessonSlideId = question.LessonSlideId
	request.AutoPoints = question.AutoPoints
	request.Data.Answers = question.Answers
	request.Data.Content = question.Content.OrElse("")
	request.Data.Explanation = question.Explanation.OrElse("")
	request.Data.Formatted = question.Formatted
	request.Data.MultipleSelection = question.MultipleSelection
	request.Data.Solution = question.Solution
	request.Data.Type = question.Type
	request_actual := &MultiChoiceQuestionRequestActual{}
	request_actual.Question = *request
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(request_actual)
	if err != nil {
		return err
	}
	body, err := c.httpRequest(fmt.Sprintf("lessons/slides/questions/%d", question.Id), "PUT", buf, nil)
	if err != nil {
		return err
	}
	resp_lesson := &MultiChoiceQuestionResponseActual{}
	err = json.NewDecoder(body).Decode(resp_lesson)
	if err != nil {
		return err
	}
	question.Id = resp_lesson.Question.Id
	return err
}

func (c *Client) CreateQuestion(question *Question) error {
	request := &MultiChoiceQuestionRequest{}
	request.Index = question.Index
	request.LessonSlideId = question.LessonSlideId
	request.AutoPoints = question.AutoPoints
	request.Data.Answers = question.Answers
	request.Data.Content = question.Content.OrElse("")
	request.Data.Explanation = question.Explanation.OrElse("")
	request.Data.Formatted = question.Formatted
	request.Data.MultipleSelection = question.MultipleSelection
	request.Data.Solution = question.Solution
	request.Data.Type = question.Type
	request_actual := &MultiChoiceQuestionRequestActual{}
	request_actual.Question = *request
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(request_actual)
	if err != nil {
		return err
	}
	body, err := c.httpRequest(fmt.Sprintf("lessons/slides/%d/questions", question.LessonSlideId), "POST", buf, nil)
	if err != nil {
		return err
	}
	resp_lesson := &MultiChoiceQuestionResponseActual{}
	err = json.NewDecoder(body).Decode(resp_lesson)
	if err != nil {
		return err
	}
	question.Id = resp_lesson.Question.Id
	return err
}
