package resourceclients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/wshelpers"

	"github.com/markphelps/optional"
)

type Challenge struct {
	Id          int            `json:"id"`
	CourseId    int            `json:"course_id"`
	LessonId    optional.Int64 `json:"lesson_id"`
	SlideId     optional.Int64 `json:"slide_id"`
	Type        string         `json:"type"`
	Explanation string         `json:"explanation"`

	Features ChallengeFeatures `json:"features"`
	Settings ChallengeSettings `json:"settings"`
	Tickets  ChallengeTickets  `json:"tickets"`
}

type ChallengeFeatures struct {
	Run                  bool `json:"run"`
	Check                bool `json:"check"`
	Mark                 bool `json:"mark"`
	Terminal             bool `json:"terminal"`
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
}

type MarkStandardTicket struct {
	BuildCommand string `json:"build_command"`
	RunCommand   string `json:"run_command"`
	// Testcases []any `json:"testcases"`
	/*
		{
			"testcases": [
				{
					"name": "Test1",
					"description": "",
					"hidden": true,
					"private": false,
					"score": 0,
					"max_score": 0,
					"skip": false,
					"run_limit": {
						"cpu_time": 3000,
						"wall_time": 3000,
						"pty_size": {
							"rows": 0,
							"cols": 0
						}
					},
					"run_command": "special_command",
					"stdin_path": "1.in",
					"extra_paths": null,
					"checks": [
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
					],
					"output_files": []
				}
			]
		}
	*/
	MarkAll bool `json:"mark_all"`
}

type PassbackSettings struct {
	MaxAutomaticScore float64 `json:"max_automatic_score"`
	ScaleTo           float64 `json:"scale_to"`
	ScoringMode       string  `json:"scoring_mode"`
}

type ChallegeResponseJSON struct {
	Challenge Challenge `json:"challenge"`
}

type ChallengeResource struct {
	Id       int `json:"id"`
	CourseId int `json:"course_id"`
	SlideId  int `json:"slide_id"`
	LessonId int `json:"lesson_id"`

	Explanation string `json:"explanation"`
	FolderPath  string `json:"folder_path"`
	FolderSha   string `json:"folder_sha"`

	Type string `json:"type"`
	// Points int  `json:"points"`

	BuildCommand     string `json:"build_command"`
	RunCommand       string `json:"run_command"`
	TestCommand      string `json:"test_command"`
	TerminalCommand  string `json:"terminal_command"`
	CustomRunCommand string `json:"custom_run_command"`

	// PointLossThreshold int `json:"point_loss_threshold"`
	// PointLossEvery     int `json:"point_loss_every"`
	// PointLossAmount    int `json:"point_loss_amount"`

	PerTestcaseScores bool `json:"per_testcase_scores"`

	MaxSubmissionsPerInterval int `json:"max_submissions_per_interval"`
	AttemptLimitInterval      int `json:"attempt_limit_interval"`

	OnlyGitSubmission            bool `json:"only_git_submission"`
	AllowSubmitAfterMarkingLimit bool `json:"allow_submit_after_marking_limit"`

	PassbackScoringMode       string  `json:"passback_scoring_mode"`
	PassbackMaxAutomaticScore float64 `json:"passback_max_automatic_score"`
	PassbackScaleTo           float64 `json:"passback_scale_to"`

	Run                  bool `json:"feature_run"`
	Check                bool `json:"feature_check"`
	Mark                 bool `json:"feature_mark"`
	Terminal             bool `json:"feature_terminal"`
	Feedback             bool `json:"feature_feedback"`
	ManualCompletion     bool `json:"feature_manual_completion"`
	AnonymousSubmissions bool `json:"feature_anonymous_submissions"`
	Arguments            bool `json:"feature_arguments"`
	ConfirmSubmit        bool `json:"feature_confirm_submit"`
	RunBeforeSubmit      bool `json:"feature_run_before_submit"`
	GitSubmission        bool `json:"feature_git_submission"`
	Editor               bool `json:"feature_editor"`
	RemoteDesktop        bool `json:"feature_remote_desktop"`
	IntermediateFiles    bool `json:"feature_intermediate_files"`

	// TODO:
	// * Rubric
	// * Test cases for the standard marking procedure
}

type TicketResponse struct {
	Ticket string `json:"ticket"`
}

func GetChallenge(c *client.Client, lesson_id int, slide_id int) (*Challenge, error) {
	body, err := c.HTTPRequest(fmt.Sprintf("lessons/%d?view=1", lesson_id), "GET", bytes.Buffer{}, nil)
	if err != nil {
		return nil, err
	}
	resp := &LessonWithSlidesResponse{}
	err = json.NewDecoder(body).Decode(resp)
	if err != nil {
		return nil, err
	}
	challenge_id := 0
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
			if !slideObj.ChallengeId.Present() {
				return nil, fmt.Errorf("Challenge for Slide %d Not Found", slideObj.Id)
			}
			challenge_id = slideObj.ChallengeId.MustGet()
			body, err = c.HTTPRequest(fmt.Sprintf("challenges/%d?view=1", challenge_id), "GET", bytes.Buffer{}, nil)
			if err != nil {
				return nil, err
			}
			resp := &ChallegeResponseJSON{}
			err = json.NewDecoder(body).Decode(resp)
			if err != nil {
				return nil, err
			}
			return &resp.Challenge, nil
		}
	}

	return nil, fmt.Errorf("Challenge for Slide %d Not Found", slide_id)
}

func UpdateChallenge(conn *client.Client, folder_path string, challenge *Challenge) error {
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
	f, f_err := os.Create("data2.json")
	if f_err != nil {
		return f_err
	}
	f.Write(buf.Bytes())
	f.Close()
	_, patch_err := conn.HTTPRequest(fmt.Sprintf("challenges/%d", challenge.Id), "PATCH", buf, nil)
	if patch_err != nil {
		return patch_err
	}

	return nil
}

func ChallengeToTerraform(c *client.Client, lesson_id int, slide_id int, resource_name string, folder_path string) (string, error) {
	chal, err := GetChallenge(c, lesson_id, slide_id)
	if err != nil {
		return "", err
	}
	var resource_string = fmt.Sprintf("resource \"edstem_challenge\" %s {\n", resource_name)
	resource_string = resource_string + fmt.Sprintf("\tslide_id = %d\n", slide_id)
	resource_string = resource_string + fmt.Sprintf("\tlesson_id = %d\n", lesson_id)

	if chal.Explanation != "" {
		if strings.Contains(chal.Explanation, "\n") {
			resource_string = resource_string + fmt.Sprintf("\texplanation = <<EOT\n%s\nEOT\n", chal.Explanation)
		} else {
			resource_string = resource_string + fmt.Sprintf("\texplanation = \"%s\"\n", chal.Explanation)
		}
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

	// TODO: Rubrics and Test cases
	resource_string = resource_string + "}"

	return resource_string, nil
}
