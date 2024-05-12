package resourceclients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/md2ed"
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

func ChallengeToTerraform(c *client.Client, lesson_id int, slide_id int, resource_name string, folder_path string, slide_resource_name *string, lesson_resource_name *string) (string, error) {
	chal, rubric, err := GetChallengeAndRubric(c, lesson_id, slide_id)
	if err != nil {
		return "", err
	}
	var resource_string = fmt.Sprintf("resource \"edstem_challenge\" %s {\n", resource_name)
	if slide_resource_name != nil {
		resource_string = resource_string + fmt.Sprintf("\tslide_id = edstem_slide.%s.id\n", *slide_resource_name)
	} else {
		resource_string = resource_string + fmt.Sprintf("\tslide_id = %d\n", slide_id)
	}
	if lesson_resource_name != nil {
		resource_string = resource_string + fmt.Sprintf("\tlesson_id = edstem_lesson.%s.id\n", *lesson_resource_name)
	} else {
		resource_string = resource_string + fmt.Sprintf("\tlesson_id = %d\n", lesson_id)
	}

	if chal.Explanation != "" {
		content_path := path.Join(folder_path, "explanation.md")
		f, e := os.Create(content_path)
		if e != nil {
			return "", e
		}
		f.WriteString(md2ed.RenderEdToMD(chal.Explanation, folder_path, true))
		resource_string = resource_string + fmt.Sprintf("\tcontent = file(\"%s\")\n", content_path)
	}

	var repos = []string{"scaffold", "solution", "testbase"}

	for _, repo := range repos {
		err = wshelpers.ReadChallengeRepo(c, chal.Id, folder_path, repo)
		if err != nil {
			return "", err
		}
	}

	resource_string = resource_string + fmt.Sprintf("\tfolder_path = \"%s\"\n", folder_path)
	resource_string = resource_string + fmt.Sprintf("\tfolder_sha = sha1(join(\"\", [for f in fileset(path.cwd, \"%s/**\"): filesha1(\"${path.cwd}/${f}\")]))\n", folder_path)
	resource_string = resource_string + fmt.Sprintf("\ttype = \"%s\"\n", chal.Type)

	if chal.Settings.BuildCommand != "" {
		resource_string = resource_string + fmt.Sprintf("\tbuild_command = \"%s\"\n", chal.Settings.BuildCommand)
	}
	if chal.Settings.RunCommand != "" {
		resource_string = resource_string + fmt.Sprintf("\trun_command = \"%s\"\n", chal.Settings.RunCommand)
	}
	if chal.Settings.CheckCommand != "" {
		resource_string = resource_string + fmt.Sprintf("\ttest_command = \"%s\"\n", chal.Settings.CheckCommand)
	}
	if chal.Settings.TerminalCommand != "" {
		resource_string = resource_string + fmt.Sprintf("\tterminal_command = \"%s\"\n", chal.Settings.TerminalCommand)
	}

	if chal.Tickets.MarkCustom.RunCommand != "" {
		resource_string = resource_string + fmt.Sprintf("\tcustom_run_command = \"%s\"\n", chal.Tickets.MarkCustom.RunCommand)
	}

	if chal.Settings.PerTestCaseScores {
		resource_string = resource_string + fmt.Sprintf("\tper_testcase_scores = %t\n", chal.Settings.PerTestCaseScores)
	}
	if chal.Settings.MaxSubmissionsPerInterval != 0 {
		resource_string = resource_string + fmt.Sprintf("\tmax_submissions_per_interval = %d\n", chal.Settings.MaxSubmissionsPerInterval)
	}
	if chal.Settings.AttemptLimitInterval != 0 {
		resource_string = resource_string + fmt.Sprintf("\tattempt_limit_interval = %d\n", chal.Settings.AttemptLimitInterval)
	}
	if chal.Settings.OnlyGitSubmission {
		resource_string = resource_string + fmt.Sprintf("\tonly_git_submission = %t\n", chal.Settings.OnlyGitSubmission)
	}
	if chal.Settings.AllowSubmitAfterMarkingLimit {
		resource_string = resource_string + fmt.Sprintf("\tallow_submit_after_marking_limit = %t\n", chal.Settings.AllowSubmitAfterMarkingLimit)
	}
	if chal.Settings.Passback.ScoringMode != "" {
		resource_string = resource_string + fmt.Sprintf("\tpassback_scoring_mode = \"%s\"\n", chal.Settings.Passback.ScoringMode)
	}
	if chal.Settings.Passback.MaxAutomaticScore != 0 {
		resource_string = resource_string + fmt.Sprintf("\tpassback_max_automatic_score = %f\n", chal.Settings.Passback.MaxAutomaticScore)
	}
	if chal.Settings.Passback.ScaleTo != 0 {
		resource_string = resource_string + fmt.Sprintf("\tpassback_max_automatic_score = %f\n", chal.Settings.Passback.ScaleTo)
	}

	if !chal.Features.Run {
		resource_string = resource_string + fmt.Sprintf("\tfeature_run = %t\n", chal.Features.Run)
	}
	if !chal.Features.Check {
		resource_string = resource_string + fmt.Sprintf("\tfeature_check = %t\n", chal.Features.Check)
	}
	if !chal.Features.Mark {
		resource_string = resource_string + fmt.Sprintf("\tfeature_mark = %t\n", chal.Features.Mark)
	}
	if !chal.Features.Connect {
		resource_string = resource_string + fmt.Sprintf("\tfeature_connect = %t\n", chal.Features.Connect)
	}
	if !chal.Features.Terminal {
		resource_string = resource_string + fmt.Sprintf("\tfeature_terminal = %t\n", chal.Features.Terminal)
	}
	if !chal.Features.Feedback {
		resource_string = resource_string + fmt.Sprintf("\tfeature_feedback = %t\n", chal.Features.Feedback)
	}
	if chal.Features.ManualCompletion {
		resource_string = resource_string + fmt.Sprintf("\tfeature_manual_completion = %t\n", chal.Features.ManualCompletion)
	}
	if chal.Features.AnonymousSubmissions {
		resource_string = resource_string + fmt.Sprintf("\tfeature_anonymous_submissions = %t\n", chal.Features.AnonymousSubmissions)
	}
	if chal.Features.Arguments {
		resource_string = resource_string + fmt.Sprintf("\tfeature_arguments = %t\n", chal.Features.Arguments)
	}
	if chal.Features.ConfirmSubmit {
		resource_string = resource_string + fmt.Sprintf("\tfeature_confirm_submit = %t\n", chal.Features.ConfirmSubmit)
	}
	if chal.Features.RunBeforeSubmit {
		resource_string = resource_string + fmt.Sprintf("\tfeature_run_before_submit = %t\n", chal.Features.RunBeforeSubmit)
	}
	if chal.Features.GitSubmission {
		resource_string = resource_string + fmt.Sprintf("\tfeature_git_submission = %t\n", chal.Features.GitSubmission)
	}
	if !chal.Features.Editor {
		resource_string = resource_string + fmt.Sprintf("\tfeature_editor = %t\n", chal.Features.Editor)
	}
	if chal.Features.RemoteDesktop {
		resource_string = resource_string + fmt.Sprintf("\tfeature_remote_desktop = %t\n", chal.Features.RemoteDesktop)
	}
	if chal.Features.IntermediateFiles {
		resource_string = resource_string + fmt.Sprintf("\tfeature_intermediate_files = %t\n", chal.Features.IntermediateFiles)
	}
	chal.Tickets.MarkCustom.RunLimit.CpuTime.If(func(val int64) {
		resource_string = resource_string + fmt.Sprintf("\tcustom_mark_time_limit_ms = %d\n", val)
	})

	if len(chal.Settings.Criteria) > 0 {
		var buf = bytes.Buffer{}
		err = json.NewEncoder(&buf).Encode(chal.Settings.Criteria)
		if err != nil {
			return "", err
		}
		content_path := path.Join(folder_path, "criteria.json")
		f, e := os.Create(content_path)
		if e != nil {
			return "", e
		}
		f.WriteString(buf.String())
		resource_string = resource_string + fmt.Sprintf("\tcriteria = file(\"%s\")\n", content_path)
	}
	if rubric != nil {
		var buf = bytes.Buffer{}
		err = json.NewEncoder(&buf).Encode(rubric)
		if err != nil {
			return "", err
		}
		content_path := path.Join(folder_path, "rubric.json")
		f, e := os.Create(content_path)
		if e != nil {
			return "", e
		}
		f.WriteString(buf.String())
		resource_string = resource_string + fmt.Sprintf("\trubric = file(\"%s\")\n", content_path)
	}

	if len(chal.Tickets.MarkStandard.Testcases) > 0 {
		res, err := json.MarshalIndent(chal.Tickets.MarkStandard.Testcases, "", "  ")
		if err != nil {
			return "", err
		}
		content_path := path.Join(folder_path, "testcases.json")
		f, e := os.Create(content_path)
		if e != nil {
			return "", e
		}
		f.Write(res)
		resource_string = resource_string + fmt.Sprintf("\ttestcase_json = file(\"%s\")\n", content_path)
	}

	chal.Tickets.MarkStandard.RunLimit.Pty.If(func(val bool) {
		resource_string = resource_string + fmt.Sprintf("\ttestcase_pty = %t\n", val)
	})
	if chal.Tickets.MarkStandard.Easy {
		resource_string = resource_string + fmt.Sprintf("\ttestcase_easy = %t\n", chal.Tickets.MarkStandard.Easy)
	}
	if chal.Tickets.MarkStandard.MarkAll {
		resource_string = resource_string + fmt.Sprintf("\ttestcase_mark_all = %t\n", chal.Tickets.MarkStandard.MarkAll)
	}
	if chal.Tickets.MarkStandard.Overlay {
		resource_string = resource_string + fmt.Sprintf("\ttestcase_overlay_test_files = %t\n", chal.Tickets.MarkStandard.Overlay)
	}

	// TODO: Test cases
	resource_string = resource_string + "}"

	return resource_string, nil
}
