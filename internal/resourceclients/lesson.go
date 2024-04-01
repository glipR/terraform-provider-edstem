package resourceclients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"terraform-provider-edstem/internal/client"

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

func LessonToTerraform(c *client.Client, lesson_id int, resource_name string, folder_path string) (string, error) {
	lesson, err := GetLesson(c, lesson_id)
	if err != nil {
		return "", err
	}
	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(lesson)
	if err != nil {
		return "", err
	}

	var resource_string = fmt.Sprintf("resource \"edstem_lesson\" %s {\n", resource_name)
	resource_string = resource_string + fmt.Sprintf("\tid = %d\n", lesson.Id)
	lesson.Attempts.If(func(val int) {
		resource_string = resource_string + fmt.Sprintf("\tattempts = %d\n", val)
	})
	lesson.AvailableAt.If(func(val string) {
		resource_string = resource_string + fmt.Sprintf("\tavailable_at = \"%s\"\n", val)
	})
	lesson.DueAt.If(func(val string) {
		resource_string = resource_string + fmt.Sprintf("\tdue_at = \"%s\"\n", val)
	})
	if lesson.GradePassbackAutoSend {
		resource_string = resource_string + fmt.Sprintf("\tgrade_passback_auto_send = %t\n", lesson.GradePassbackAutoSend)
	}
	if lesson.GradePassbackMode != "" {
		resource_string = resource_string + fmt.Sprintf("\tgrade_passback_mode = \"%s\"\n", lesson.GradePassbackMode)
	}
	lesson.GradePassbackScaleTo.If(func(val string) {
		resource_string = resource_string + fmt.Sprintf("\tgrade_passback_scale_to = \"%s\"\n", val)
	})
	lesson.Index.If(func(val int) {
		resource_string = resource_string + fmt.Sprintf("\tindex = %d\n", val)
	})
	resource_string = resource_string + fmt.Sprintf("\tis_hidden = %t\n", lesson.IsHidden)
	if lesson.IsTimed {
		resource_string = resource_string + fmt.Sprintf("\tis_timed = %t\n", lesson.IsTimed)
	}
	if lesson.IsUnlisted {
		resource_string = resource_string + fmt.Sprintf("\tis_unlisted = %t\n", lesson.IsUnlisted)
	}
	resource_string = resource_string + fmt.Sprintf("\tkind = \"%s\"\n", lesson.Kind)
	if lesson.LateSubmissions {
		resource_string = resource_string + fmt.Sprintf("\tlate_submissions = %t\n", lesson.LateSubmissions)
	}
	lesson.LockedAt.If(func(val string) {
		resource_string = resource_string + fmt.Sprintf("\tlocked_at = \"%s\"\n", val)
	})
	lesson.ModuleId.If(func(val int) {
		resource_string = resource_string + fmt.Sprintf("\tindex = %d\n", val)
	})
	if lesson.Openable {
		resource_string = resource_string + fmt.Sprintf("\topenable = %t\n", lesson.Openable)
	}
	if lesson.OpenableWithoutAttempt {
		resource_string = resource_string + fmt.Sprintf("\topenable = %t\n", lesson.OpenableWithoutAttempt)
	}

	if lesson.Outline != "" {
		if strings.Contains(lesson.Outline, "\n") {
			resource_string = resource_string + fmt.Sprintf("\toutline = <<EOT\n%s\nEOT\n", lesson.Outline)
		} else {
			resource_string = resource_string + fmt.Sprintf("\toutline = \"%s\"\n", lesson.Outline)
		}
	}
	if lesson.Password != "" {
		resource_string = resource_string + fmt.Sprintf("\tpassword = \"%s\"\n", lesson.Password)
	}
	// TODO: Prerequisites
	if lesson.ReleaseChallengeSolutions {
		resource_string = resource_string + fmt.Sprintf("\trelease_challenge_solutions = %t\n", lesson.ReleaseChallengeSolutions)
	}
	if lesson.ReleaseChallengeSolutionsWhileActive {
		resource_string = resource_string + fmt.Sprintf("\trelease_challenge_solutions_while_active = %t\n", lesson.ReleaseChallengeSolutionsWhileActive)
	}
	if lesson.ReleaseFeedback {
		resource_string = resource_string + fmt.Sprintf("\trelease_feedback = %t\n", lesson.ReleaseFeedback)
	}
	if lesson.ReleaseFeedbackWhileActive {
		resource_string = resource_string + fmt.Sprintf("\trelease_feedback_while_active = %t\n", lesson.ReleaseFeedbackWhileActive)
	}
	if lesson.ReleaseQuizCorrectnessOnly {
		resource_string = resource_string + fmt.Sprintf("\trelease_quiz_correctness_only = %t\n", lesson.ReleaseQuizCorrectnessOnly)
	}
	if lesson.ReleaseQuizSolutions {
		resource_string = resource_string + fmt.Sprintf("\trelease_quiz_solutions = %t\n", lesson.ReleaseQuizSolutions)
	}
	if lesson.ReOpenSubmissions {
		resource_string = resource_string + fmt.Sprintf("\treopen_submissions = %t\n", lesson.ReOpenSubmissions)
	}
	if lesson.RequireUserOverride {
		resource_string = resource_string + fmt.Sprintf("\trequire_user_override = %t\n", lesson.RequireUserOverride)
	}
	resource_string = resource_string + fmt.Sprintf("\tquiz_active_status = \"%s\"\n", lesson.QuizSettings.QuizActiveStatus)
	resource_string = resource_string + fmt.Sprintf("\tquiz_mode = \"%s\"\n", lesson.QuizSettings.QuizMode)
	if lesson.QuizSettings.QuizQuestionNumberStyle != "" {
		resource_string = resource_string + fmt.Sprintf("\tquiz_question_number_style = \"%s\"\n", lesson.QuizSettings.QuizQuestionNumberStyle)
	}

	lesson.SolutionsAt.If(func(val string) {
		resource_string = resource_string + fmt.Sprintf("\tsolutions_at = \"%s\"\n", val)
	})
	if lesson.State != "active" {
		resource_string = resource_string + fmt.Sprintf("\tstate = \"%s\"\n", lesson.State)
	}
	if lesson.TimerDuration != 60 {
		resource_string = resource_string + fmt.Sprintf("\ttimer_duration = %d\n", lesson.TimerDuration)
	}
	if lesson.TimerExpirationAccess {
		resource_string = resource_string + fmt.Sprintf("\ttimer_expiration_access = %t\n", lesson.TimerExpirationAccess)
	}

	resource_string = resource_string + fmt.Sprintf("\ttitle = \"%s\"\n", lesson.Title)
	if lesson.TutorialRegex != "" {
		resource_string = resource_string + fmt.Sprintf("\ttutorial_regex = \"%s\"\n", lesson.TutorialRegex)
	}
	resource_string = resource_string + fmt.Sprintf("\ttype = \"%s\"\n", lesson.Type)

	resource_string = resource_string + "}"

	slide_ids, e := GetSlideIds(c, lesson_id)
	if e != nil {
		return "", e
	}

	for i := range slide_ids {
		slide_path := path.Join(folder_path, fmt.Sprintf("slide_%d", i))
		e = os.MkdirAll(slide_path, os.ModeDir)
		if e != nil {
			return "", nil
		}
		new_string, slide_err := SlideToTerraform(c, lesson_id, slide_ids[i], fmt.Sprintf("%s_slide_%d", resource_name, i), slide_path, &resource_name)
		if slide_err != nil {
			return "", slide_err
		}
		resource_string = resource_string + "\n\n" + new_string
	}

	return resource_string, nil
}

func CourseToTerraform(c *client.Client, folder_path string) (string, error) {
	lessons, err := GetLessons(c)
	if err != nil {
		return "", err
	}
	lesson_terraform_blocks := make([]string, 0)
	for i, lesson := range lessons {
		lesson_path := fmt.Sprintf("lesson_%d", i)
		res, e := LessonToTerraform(c, lesson.Id, lesson_path, path.Join(folder_path, lesson_path))
		if e != nil {
			return "", e
		}
		lesson_terraform_blocks = append(lesson_terraform_blocks, res)
	}
	return strings.Join(lesson_terraform_blocks, "\n\n\n"), nil
}
