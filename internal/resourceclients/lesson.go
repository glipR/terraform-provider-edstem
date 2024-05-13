package resourceclients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/tfhelpers"

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

func GetLessons(c *client.Client) ([]Lesson, error) {
	body, err := c.HTTPRequest(fmt.Sprintf("courses/%s/lessons", c.CourseID), "GET", bytes.Buffer{}, nil)
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

func GetLesson(c *client.Client, lesson_id int) (*Lesson, error) {
	body, err := c.HTTPRequest(fmt.Sprintf("lessons/%d?view=1", lesson_id), "GET", bytes.Buffer{}, nil)
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

func UpdateLesson(c *client.Client, lesson *Lesson) error {
	request := &LessonUpdateRequest{LessonObj: *lesson}
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	body, err := c.HTTPRequest(fmt.Sprintf("lessons/%d", lesson.Id), "PUT", buf, nil)
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

func CreateLesson(c *client.Client, lesson *Lesson) error {
	lesson_request := &NewLessonRequest{Kind: lesson.Kind}
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(lesson_request)
	if err != nil {
		return err
	}
	body, err := c.HTTPRequest(fmt.Sprintf("courses/%s/lessons", c.CourseID), "POST", buf, nil)
	if err != nil {
		return err
	}
	resp_lesson := &LessonResponse{}
	err = json.NewDecoder(body).Decode(resp_lesson)
	if err != nil {
		return err
	}
	lesson.Id = resp_lesson.LessonObj.Id
	return UpdateLesson(c, lesson)
}

func LessonToTerraform(c *client.Client, lesson_id int, resource_name string, folder_path string) (string, []string, error) {
	lesson, err := GetLesson(c, lesson_id)
	if err != nil {
		return "", []string{}, err
	}
	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(lesson)
	if err != nil {
		return "", []string{}, err
	}

	resources := make([]string, 0)
	resources = append(resources, fmt.Sprintf("edstem_lesson.%s %d", resource_name, lesson_id))

	var resource_string = fmt.Sprintf("resource \"edstem_lesson\" %s {\n", resource_name)

	resource_string = resource_string + tfhelpers.TFProp("id", lesson.Id, nil)
	resource_string = resource_string + tfhelpers.TFProp("attempts", lesson.Attempts, nil)

	resource_string = resource_string + tfhelpers.TFProp("available_at", lesson.AvailableAt, nil)
	resource_string = resource_string + tfhelpers.TFProp("due_at", lesson.DueAt, nil)
	resource_string = resource_string + tfhelpers.TFProp("locked_at", lesson.LockedAt, nil)
	resource_string = resource_string + tfhelpers.TFProp("solutions_at", lesson.SolutionsAt, nil)

	resource_string = resource_string + tfhelpers.TFProp("grade_passback_auto_send", lesson.GradePassbackAutoSend, false)
	resource_string = resource_string + tfhelpers.TFProp("grade_passback_mode", lesson.GradePassbackMode, "")
	resource_string = resource_string + tfhelpers.TFProp("grade_passback_scale_to", lesson.GradePassbackScaleTo, nil)

	resource_string = resource_string + tfhelpers.TFProp("index", lesson.Index, nil)
	resource_string = resource_string + tfhelpers.TFProp("is_hidden", lesson.IsHidden, nil)
	resource_string = resource_string + tfhelpers.TFProp("is_timed", lesson.IsTimed, false)
	resource_string = resource_string + tfhelpers.TFProp("is_unlisted", lesson.IsUnlisted, false)
	resource_string = resource_string + tfhelpers.TFProp("late_submissions", lesson.LateSubmissions, false)
	resource_string = resource_string + tfhelpers.TFProp("openable", lesson.Openable, false)
	resource_string = resource_string + tfhelpers.TFProp("openable_without_attempt", lesson.OpenableWithoutAttempt, false)

	resource_string = resource_string + tfhelpers.TFProp("kind", lesson.Kind, "")
	resource_string = resource_string + tfhelpers.TFProp("module_id", lesson.ModuleId, nil)

	if lesson.Outline != "" {
		if strings.Contains(lesson.Outline, "\n") {
			resource_string = resource_string + tfhelpers.TFUnquote("outline", fmt.Sprintf("<<EOT\n%s\nEOT", lesson.Outline))
		} else {
			resource_string = resource_string + tfhelpers.TFProp("outline", lesson.Outline, "")
		}
	}
	resource_string = resource_string + tfhelpers.TFProp("password", lesson.Password, "")
	resource_string = resource_string + tfhelpers.TFProp("state", lesson.State, "active")
	resource_string = resource_string + tfhelpers.TFProp("title", lesson.Title, "")
	resource_string = resource_string + tfhelpers.TFProp("tutorial_regex", lesson.TutorialRegex, "")
	resource_string = resource_string + tfhelpers.TFProp("type", lesson.Type, "")

	// TODO: Prerequisites

	resource_string = resource_string + tfhelpers.TFProp("release_challenge_solutions", lesson.ReleaseChallengeSolutions, false)
	resource_string = resource_string + tfhelpers.TFProp("release_challenge_solutions_while_active", lesson.ReleaseChallengeSolutionsWhileActive, false)
	resource_string = resource_string + tfhelpers.TFProp("release_feedback", lesson.ReleaseFeedback, false)
	resource_string = resource_string + tfhelpers.TFProp("release_feedback_while_active", lesson.ReleaseFeedbackWhileActive, false)
	resource_string = resource_string + tfhelpers.TFProp("release_quiz_correctness_only", lesson.ReleaseQuizCorrectnessOnly, false)
	resource_string = resource_string + tfhelpers.TFProp("release_quiz_solutions", lesson.ReleaseQuizSolutions, false)
	resource_string = resource_string + tfhelpers.TFProp("reopen_submissions", lesson.ReOpenSubmissions, false)
	resource_string = resource_string + tfhelpers.TFProp("require_user_override", lesson.RequireUserOverride, false)

	resource_string = resource_string + tfhelpers.TFProp("quiz_active_status", lesson.QuizSettings.QuizActiveStatus, "")
	resource_string = resource_string + tfhelpers.TFProp("quiz_mode", lesson.QuizSettings.QuizMode, "")
	resource_string = resource_string + tfhelpers.TFProp("quiz_question_number_style", lesson.QuizSettings.QuizQuestionNumberStyle, "")

	resource_string = resource_string + tfhelpers.TFProp("timer_duration", lesson.TimerDuration, 60)
	resource_string = resource_string + tfhelpers.TFProp("timer_expiration_access", lesson.TimerExpirationAccess, false)

	resource_string = resource_string + "}"

	slide_ids, e := GetSlideIds(c, lesson_id)
	if e != nil {
		return "", []string{}, e
	}

	for i := range slide_ids {
		slide_path := path.Join(folder_path, fmt.Sprintf("slide_%d", i))
		e = os.MkdirAll(slide_path, 0777)
		if e != nil {
			return "", []string{}, nil
		}
		new_string, slide_resources, slide_err := SlideToTerraform(c, lesson_id, slide_ids[i], fmt.Sprintf("%s_slide_%d", resource_name, i), slide_path, &resource_name)
		if slide_err != nil {
			return "", []string{}, slide_err
		}
		resource_string = resource_string + "\n\n" + new_string
		resources = append(resources, slide_resources...)
	}

	return resource_string, resources, nil
}

func CourseToTerraform(c *client.Client, folder_path string) (string, []string, error) {
	lessons, err := GetLessons(c)
	if err != nil {
		return "", []string{}, err
	}
	lesson_terraform_blocks := make([]string, 0)
	resources := make([]string, 0)
	for i, lesson := range lessons {
		lesson_path := fmt.Sprintf("lesson_%d", i)
		res, lesson_resources, e := LessonToTerraform(c, lesson.Id, lesson_path, path.Join(folder_path, lesson_path))
		if e != nil {
			return "", []string{}, e
		}
		lesson_terraform_blocks = append(lesson_terraform_blocks, res)
		resources = append(resources, lesson_resources...)
	}
	return strings.Join(lesson_terraform_blocks, "\n\n\n"), resources, nil
}
