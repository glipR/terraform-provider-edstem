package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/markphelps/optional"
)

type QuizSettings struct {
	QuizActiveStatus        string `json:"quiz_active_status"`
	QuizMode                string `json:"quiz_mode"`
	QuizQuestionNumberStyle string `json:"quiz_question_number_style"`
}

type Prerequisite struct {
	RequiredLessonId int `json:"required_lesson_id"`
}

type Lesson struct {
	Attempts                             optional.Int    `json:"attempts"`
	AvailableAt                          optional.String `json:"available_at"`
	DueAt                                optional.String `json:"due_at"`
	GradePassbackAutoSend                bool            `json:"grade_passback_auto_send"`
	GradePassbackMode                    string          `json:"grade_passback_mode"`
	GradePassbackScaleTo                 optional.String `json:"grade_passback_scale_to"`
	Id                                   int             `json:"id"`
	Index                                optional.Int    `json:"index"`
	IsHidden                             bool            `json:"is_hidden"`
	IsTimed                              bool            `json:"is_timed"`
	IsUnlisted                           bool            `json:"is_unlisted"`
	Kind                                 string          `json:"kind"`
	LateSubmissions                      bool            `json:"late_submissions"`
	LockedAt                             optional.String `json:"locked_at"`
	ModuleId                             optional.Int    `json:"module_id"`
	Openable                             bool            `json:"openable"`
	OpenableWithoutAttempt               bool            `json:"openable_without_attempt"`
	Outline                              string          `json:"outline"`
	Password                             string          `json:"password"`
	Prerequisites                        []Prerequisite  `json:"prerequisites"`
	ReleaseChallengeSolutions            bool            `json:"release_challenge_solutions"`
	ReleaseChallengeSolutionsWhileActive bool            `json:"release_challenge_solutions_while_active"`
	ReleaseFeedback                      bool            `json:"release_feedback"`
	ReleaseFeedbackWhileActive           bool            `json:"release_feedback_while_active"`
	ReleaseQuizCorrectnessOnly           bool            `json:"release_quiz_correctness_only"`
	ReleaseQuizSolutions                 bool            `json:"release_quiz_solutions"`
	ReOpenSubmissions                    bool            `json:"reopen_submissions"`
	RequireUserOverride                  bool            `json:"require_user_override"`
	QuizSettings                         QuizSettings    `json:"settings"`
	SolutionsAt                          optional.String `json:"solutions_at"`
	State                                string          `json:"state"`
	TimerDuration                        int             `json:"timer_duration"`
	TimerExpirationAccess                bool            `json:"timer_expiration_access"`
	Title                                string          `json:"title"`
	TutorialRegex                        string          `json:"tutorial_regex"`
	Type                                 string          `json:"type"`
}

type LessonResponse struct {
	LessonObj Lesson `json:"lesson"`
}

type CourseResponse struct {
	LessonList []Lesson `json:"lessons"`
}

type NewLessonRequest struct {
	Kind string `json:"kind"`
}

type LessonUpdateRequest struct {
	LessonObj Lesson `json:"lesson"`
}

func (c *Client) GetLessons() ([]Lesson, error) {
	body, err := c.httpRequest(fmt.Sprintf("courses/%s/lessons", c.CourseID), "GET", bytes.Buffer{}, nil)
	if err != nil {
		return nil, err
	}
	response := &CourseResponse{}
	err = json.NewDecoder(body).Decode(response)
	if err != nil {
		return nil, err
	}
	return response.LessonList, nil
}

func (c *Client) GetLesson(lesson_id int) (*Lesson, error) {
	body, err := c.httpRequest(fmt.Sprintf("lessons/%d?view=1", lesson_id), "GET", bytes.Buffer{}, nil)
	if err != nil {
		return nil, err
	}
	lesson := &LessonResponse{}
	err = json.NewDecoder(body).Decode(lesson)
	if err != nil {
		return nil, err
	}
	return &lesson.LessonObj, nil
}

func (c *Client) UpdateLesson(lesson *Lesson) error {
	request := &LessonUpdateRequest{LessonObj: *lesson}
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	body, err := c.httpRequest(fmt.Sprintf("lessons/%d", lesson.Id), "PUT", buf, nil)
	if err != nil {
		return err
	}
	resp_lesson := &LessonResponse{}
	err = json.NewDecoder(body).Decode(resp_lesson)
	if err != nil {
		return err
	}
	lesson.Id = resp_lesson.LessonObj.Id
	return err
}

func (c *Client) CreateLesson(lesson *Lesson) error {
	lesson_request := &NewLessonRequest{Kind: lesson.Kind}
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(lesson_request)
	if err != nil {
		return err
	}
	body, err := c.httpRequest(fmt.Sprintf("courses/%s/lessons", c.CourseID), "POST", buf, nil)
	if err != nil {
		return err
	}
	resp_lesson := &LessonResponse{}
	err = json.NewDecoder(body).Decode(resp_lesson)
	if err != nil {
		return err
	}
	return c.UpdateLesson(lesson)
}
