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
	"terraform-provider-edstem/internal/wshelpers"

	"github.com/markphelps/optional"
)

type Challenge struct {
	Id           int            `json:"id"`
	CourseId     int            `json:"course_id"`
	LessonId     optional.Int64 `json:"lesson_id"`
	SlideId      optional.Int64 `json:"slide_id"`
	Type         string         `json:"type"`
	Explanation  string         `json:"explanation"`
	RubricId     optional.Int   `json:"rubric_id"`
	RubricPoints optional.Int   `json:"rubric_points"`

	Features ChallengeFeatures `json:"features"`
	Settings ChallengeSettings `json:"settings"`
	Tickets  ChallengeTickets  `json:"tickets"`
}

type RubricResponse struct {
	Rubric Rubric `json:"rubric"`
}

type Rubric struct {
	PositiveGrading  bool            `json:"positive_grading"`
	Id               optional.Int    `json:"id"`
	Sections         []RubricSection `json:"sections"`
	UnsectionedItems []RubricItem    `json:"unsectioned_items"`
}

type RubricSection struct {
	Id        optional.Int   `json:"id"`
	SelectOne bool           `json:"bool"`
	MarkClamp optional.Int64 `json:"mark_clamp"`
	Title     string         `json:"title"`
	Index     int            `json:"index"`
	Items     []RubricItem   `json:"items"`
}

type RubricItem struct {
	Id               optional.Int `json:"id"`
	Points           int          `json:"points"`
	Title            string       `json:"title"`
	StaffDescription string       `json:"staff_description"`
	Index            int          `json:"index"`
}

type ChallengeFeatures struct {
	Run                  bool `json:"run"`
	Check                bool `json:"check"`
	Mark                 bool `json:"mark"`
	Terminal             bool `json:"terminal"`
	Connect              bool `json:"connect"`
	Feedback             bool `json:"feedback"`
	ManualCompletion     bool `json:"manual_completion"`
	AnonymousSubmissions bool `json:"anonymous_submissions"`
	Arguments            bool `json:"arguments"`
	ConfirmSubmit        bool `json:"confirm_submit"`
	RunBeforeSubmit      bool `json:"run_before_submit"`
	GitSubmission        bool `json:"git_submission"`
	Editor               bool `json:"editor"`
	RemoteDesktop        bool `json:"remote_desktop"`
	IntermediateFiles    bool `json:"intermediate_files"`
}

type ChallengeSettings struct {
	BuildCommand                        string           `json:"build_command"`
	CheckCommand                        string           `json:"check_command"`
	RunCommand                          string           `json:"run_command"`
	TerminalCommand                     string           `json:"terminal_command"`
	OnlyGitSubmission                   bool             `json:"only_git_submission"`
	AttemptLimitInterval                int              `json:"attempt_limit_interval"`
	AllowSubmitAfterMarkingLimit        bool             `json:"allow_submit_after_marking_limit"`
	MaxSubmissionsPerInterval           int              `json:"max_submissions_per_interval"`
	MaxSubmissionsWithIntermediateFiles int              `json:"max_submissions_with_intermediate_files"`
	Passback                            PassbackSettings `json:"passback"`
	PerTestCaseScores                   bool             `json:"per_testcase_scores"`
	Criteria                            []Criteria       `json:"criteria"`
}

type Criteria struct {
	Name   string          `json:"name"`
	Levels []CriteriaLevel `json:"levels"`
}

type CriteriaLevel struct {
	Mark        string `json:"mark"`
	Description string `json:"description"`
}

type ChallengeTickets struct {
	RunUnit     SimpleRunTicket `json:"run_unit"`
	RunCustom   SimpleRunTicket `json:"run_custom"`
	RunStandard SimpleRunTicket `json:"run_standard"`

	MarkUnit     MarkUnitTicket     `json:"mark_unit"`
	MarkCustom   MarkCustomTicket   `json:"mark_custom"`
	MarkStandard MarkStandardTicket `json:"mark_standard"`
}

type SimpleRunTicket struct {
	RunCommand   string `json:"run_command"`
	TestCommand  string `json:"test_command"`
	BuildCommand string `json:"build_command"`
}

type MarkUnitTicket struct {
	BuildCommand        string `json:"build_command"`
	TestcasePath        string `json:"testcase_path"`
	AdditionalClasspath string `json:"additional_classpath"`
}

type MarkCustomTicket struct {
	BuildCommand string         `json:"build_command"`
	RunCommand   string         `json:"run_command"`
	RunLimit     RunLimitConfig `json:"run_limit"`
}

type RunLimitConfig struct {
	CpuTime  optional.Int64 `json:"cpu_time"`
	WallTime optional.Int64 `json:"wall_time"`
	Pty      optional.Bool  `json:"pty"`
}

type MarkStandardTicket struct {
	BuildCommand string         `json:"build_command"`
	RunCommand   string         `json:"run_command"`
	Testcases    []TestCase     `json:"testcases"`
	Easy         bool           `json:"easy"`
	MarkAll      bool           `json:"mark_all"`
	RunLimit     RunLimitConfig `json:"run_limit"`
	Overlay      bool           `json:"overlay"`
}

type TestCase struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Hidden      bool            `json:"hidden"`
	Private     bool            `json:"private"`
	Score       int             `json:"score"`
	MaxScore    int             `json:"max_score"`
	Skip        bool            `json:"skip"`
	RunCommand  optional.String `json:"run_command"`
	StdinPath   string          `json:"stdin_path"`
	OutputFiles []string        `json:"output_files"`
	Checks      []TestCaseCheck `json:"checks"`
	RunLimit    RunLimitConfig  `json:"run_limit"`
}

type TestCaseCheck struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	ExpectPath string `json:"expect_path"`
	Markdown   bool   `json:"markdown"`
	// TODO: More fields
	/*
		{
					"name": "",
					"type": "check_diff",
					"source": {
						"type": "source_mixed",
						"file": ""
					},
					"transforms": [],
					"expect_path": "1.out",
					"acceptable_line_error_rate": 0,
					"acceptable_char_error_rate": 0,
					"acceptable_line_errors": 0,
					"acceptable_char_errors": 0,
					"regex_match": "",
					"run_limit": {
						"pty_size": {
							"rows": 0,
							"cols": 0
						}
					},
					"run_command": "",
					"markdown": false
				}
	*/
}

type PassbackSettings struct {
	MaxAutomaticScore float64 `json:"max_automatic_score"`
	ScaleTo           float64 `json:"scale_to"`
	ScoringMode       string  `json:"scoring_mode"`
}

type ChallegeResponseJSON struct {
	Challenge Challenge `json:"challenge"`
}

type TicketResponse struct {
	Ticket string `json:"ticket"`
}

func GetChallengeAndRubric(c *client.Client, lesson_id int, slide_id int) (*Challenge, *Rubric, error) {
	body, err := c.HTTPRequest(fmt.Sprintf("lessons/%d?view=1", lesson_id), "GET", bytes.Buffer{}, nil)
	if err != nil {
		return nil, nil, err
	}
	resp := &LessonWithSlidesResponse{}
	err = json.NewDecoder(body).Decode(resp)
	if err != nil {
		return nil, nil, err
	}
	challenge_id := 0
	slides := resp.Lesson.Slides
	for i := range slides {
		if slides[i].Id == slide_id {
			var slideObj Slide
			tempVar, _ := json.Marshal(slides[i])
			err = json.Unmarshal(tempVar, &slideObj)
			if err != nil {
				return nil, nil, err
			}
			// Not returned from response - intuited from request.
			slideObj.LessonId = lesson_id
			if !slideObj.ChallengeId.Present() {
				return nil, nil, fmt.Errorf("Challenge for Slide %d Not Found", slideObj.Id)
			}
			challenge_id = slideObj.ChallengeId.MustGet()
			body, err = c.HTTPRequest(fmt.Sprintf("challenges/%d?view=1", challenge_id), "GET", bytes.Buffer{}, nil)
			if err != nil {
				return nil, nil, err
			}
			resp := &ChallegeResponseJSON{}
			err = json.NewDecoder(body).Decode(resp)
			if err != nil {
				return nil, nil, err
			}
			resp.Challenge.LessonId.Set(int64(lesson_id))
			resp.Challenge.SlideId.Set(int64(slide_id))
			// Rubric Data
			var rubric *RubricResponse
			resp.Challenge.RubricId.If(func(val int) {
				body, err = c.HTTPRequest(fmt.Sprintf("rubrics/%d", val), "GET", bytes.Buffer{}, nil)
				rubric = &RubricResponse{}
				err = json.NewDecoder(body).Decode(&rubric)
			})
			if err != nil {
				return nil, nil, err
			}
			if rubric == nil {
				return &resp.Challenge, nil, nil
			}
			return &resp.Challenge, &rubric.Rubric, nil
		}
	}

	return nil, nil, fmt.Errorf("Challenge for Slide %d Not Found", slide_id)
}

func UpdateChallenge(conn *client.Client, folder_path string, challenge *Challenge, rubric *Rubric) error {
	dir_entries, err := os.ReadDir(folder_path)
	if err != nil {
		return err
	}

	for _, subdir := range dir_entries {
		wshelpers.UpdateChallengeRepo(conn, challenge.Id, folder_path, subdir.Name())
	}

	var request = &ChallegeResponseJSON{}
	request.Challenge = *challenge

	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	body, patch_err := conn.HTTPRequest(fmt.Sprintf("challenges/%d", challenge.Id), "PATCH", buf, nil)
	if patch_err != nil {
		return patch_err
	}

	resp := &ChallegeResponseJSON{}
	err = json.NewDecoder(body).Decode(resp)
	if err != nil {
		return err
	}

	if rubric != nil {
		rubric.Id.Set(resp.Challenge.Id)
		challenge.RubricId = resp.Challenge.RubricId
		if challenge.RubricId.Present() {
			// Update
			var request = &RubricResponse{}
			request.Rubric = *rubric
			buf := bytes.Buffer{}
			err = json.NewEncoder(&buf).Encode(request)
			if err != nil {
				return err
			}
			_, err := conn.HTTPRequest(fmt.Sprintf("rubrics/%d", challenge.RubricId.MustGet()), "PUT", buf, nil)
			if err != nil {
				return err
			}
		} else {
			// Create
			var request = &RubricResponse{}
			request.Rubric = *rubric
			buf := bytes.Buffer{}
			err = json.NewEncoder(&buf).Encode(request)
			if err != nil {
				return err
			}
			_, err := conn.HTTPRequest(fmt.Sprintf("markable/%d/rubric?replace=false", challenge.LessonId.MustGet()), "PUT", buf, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ChallengeToTerraform(c *client.Client, lesson_id int, slide_id int, resource_name string, folder_path string, slide_resource_name *string, lesson_resource_name *string) (string, []string, error) {
	chal, rubric, err := GetChallengeAndRubric(c, lesson_id, slide_id)
	if err != nil {
		return "", []string{}, err
	}
	resources := make([]string, 0)
	resources = append(resources, fmt.Sprintf("edstem_challenge.%s %d,%d", resource_name, lesson_id, slide_id))
	var resource_string = fmt.Sprintf("resource \"edstem_challenge\" %s {\n", resource_name)
	if slide_resource_name != nil {
		resource_string = resource_string + tfhelpers.TFUnquote("slide_id", fmt.Sprintf("edstem_slide.%s.id\n", *slide_resource_name))
	} else {
		resource_string = resource_string + tfhelpers.TFProp("slide_id", slide_id, "")
	}
	if lesson_resource_name != nil {
		resource_string = resource_string + tfhelpers.TFUnquote("lesson_id", fmt.Sprintf("edstem_lesson.%s.id\n", *lesson_resource_name))
	} else {
		resource_string = resource_string + tfhelpers.TFProp("lesson_id", lesson_id, "")
	}

	if chal.Explanation != "" {
		content_path := path.Join(folder_path, "explanation.md")
		resource_string = resource_string + tfhelpers.TFFile("content", md2ed.RenderEdToMD(chal.Explanation, folder_path, true), content_path)
	}

	var repos = []string{"scaffold", "solution", "testbase"}

	for _, repo := range repos {
		err = wshelpers.ReadChallengeRepo(c, chal.Id, folder_path, repo)
		if err != nil {
			return "", []string{}, err
		}
	}

	resource_string = resource_string + tfhelpers.TFProp("folder_path", folder_path, "")
	resource_string = resource_string + tfhelpers.TFUnquote("folder_sha", fmt.Sprintf("sha1(join(\"\", [for f in fileset(path.cwd, \"%s/**\"): filesha1(\"${path.cwd}/${f}\")]))\n", folder_path))
	resource_string = resource_string + tfhelpers.TFProp("type", chal.Type, "")

	resource_string = resource_string + tfhelpers.TFProp("build_command", chal.Settings.BuildCommand, "")
	resource_string = resource_string + tfhelpers.TFProp("run_command", chal.Settings.RunCommand, "")
	resource_string = resource_string + tfhelpers.TFProp("test_command", chal.Settings.CheckCommand, "")
	resource_string = resource_string + tfhelpers.TFProp("terminal_command", chal.Settings.TerminalCommand, "")
	resource_string = resource_string + tfhelpers.TFProp("custom_run_command", chal.Tickets.MarkCustom.RunCommand, "")

	resource_string = resource_string + tfhelpers.TFProp("per_testcase_scores", chal.Settings.PerTestCaseScores, false)
	resource_string = resource_string + tfhelpers.TFProp("max_submissions_per_interval", chal.Settings.MaxSubmissionsPerInterval, 0)
	resource_string = resource_string + tfhelpers.TFProp("attempt_limit_interval", chal.Settings.AttemptLimitInterval, 0)
	resource_string = resource_string + tfhelpers.TFProp("only_git_submission", chal.Settings.OnlyGitSubmission, false)
	resource_string = resource_string + tfhelpers.TFProp("allow_submit_after_marking_limit", chal.Settings.AllowSubmitAfterMarkingLimit, false)

	resource_string = resource_string + tfhelpers.TFProp("passback_scoring_mode", chal.Settings.Passback.ScoringMode, "")
	resource_string = resource_string + tfhelpers.TFProp("passback_max_automatic_score", chal.Settings.Passback.MaxAutomaticScore, float64(0))
	resource_string = resource_string + tfhelpers.TFProp("passback_scale_to", chal.Settings.Passback.ScaleTo, float64(0))

	resource_string = resource_string + tfhelpers.TFProp("feature_run", chal.Features.Run, true)
	resource_string = resource_string + tfhelpers.TFProp("feature_check", chal.Features.Check, true)
	resource_string = resource_string + tfhelpers.TFProp("feature_mark", chal.Features.Mark, true)
	resource_string = resource_string + tfhelpers.TFProp("feature_connect", chal.Features.Connect, true)
	resource_string = resource_string + tfhelpers.TFProp("feature_terminal", chal.Features.Terminal, true)
	resource_string = resource_string + tfhelpers.TFProp("feature_feedback", chal.Features.Feedback, true)
	resource_string = resource_string + tfhelpers.TFProp("feature_editor", chal.Features.Editor, true)

	resource_string = resource_string + tfhelpers.TFProp("feature_manual_completion", chal.Features.ManualCompletion, false)
	resource_string = resource_string + tfhelpers.TFProp("feature_anonymous_submissions", chal.Features.AnonymousSubmissions, false)
	resource_string = resource_string + tfhelpers.TFProp("feature_arguments", chal.Features.Arguments, false)
	resource_string = resource_string + tfhelpers.TFProp("feature_confirm_submit", chal.Features.ConfirmSubmit, false)
	resource_string = resource_string + tfhelpers.TFProp("feature_run_before_submit", chal.Features.RunBeforeSubmit, false)
	resource_string = resource_string + tfhelpers.TFProp("feature_git_submission", chal.Features.GitSubmission, false)
	resource_string = resource_string + tfhelpers.TFProp("feature_remote_desktop", chal.Features.RemoteDesktop, false)
	resource_string = resource_string + tfhelpers.TFProp("feature_intermediate_files", chal.Features.IntermediateFiles, false)

	resource_string = resource_string + tfhelpers.TFProp("custom_mark_time_limit_ms", chal.Tickets.MarkCustom.RunLimit.CpuTime, optional.Int64{})

	if len(chal.Settings.Criteria) > 0 {
		res, err := json.MarshalIndent(chal.Settings.Criteria, "", "  ")
		if err != nil {
			return "", []string{}, err
		}
		content_path := path.Join(folder_path, "criteria.json")
		resource_string = resource_string + tfhelpers.TFFile("criteria", string(res), content_path)
	}
	if rubric != nil {
		res, err := json.MarshalIndent(rubric, "", "  ")
		if err != nil {
			return "", []string{}, err
		}
		content_path := path.Join(folder_path, "rubric.json")
		resource_string = resource_string + tfhelpers.TFFile("rubric", string(res), content_path)
	}

	if len(chal.Tickets.MarkStandard.Testcases) > 0 {
		res, err := json.MarshalIndent(chal.Tickets.MarkStandard.Testcases, "", "  ")
		if err != nil {
			return "", []string{}, err
		}
		content_path := path.Join(folder_path, "testcases.json")
		resource_string = resource_string + tfhelpers.TFFile("testcases", string(res), content_path)
	}

	tfhelpers.TFProp("testcase_pty", chal.Tickets.MarkStandard.RunLimit.Pty, nil)
	tfhelpers.TFProp("testcase_easy", chal.Tickets.MarkStandard.Easy, false)
	tfhelpers.TFProp("testcase_mark_all", chal.Tickets.MarkStandard.MarkAll, false)
	tfhelpers.TFProp("testcase_overlay_test_files", chal.Tickets.MarkStandard.Overlay, false)

	resource_string = resource_string + "}"

	return resource_string, resources, nil
}
