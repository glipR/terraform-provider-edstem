package resourceclients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/wshelpers"
)

type Challenge struct {
	Id       int `json:"id"`
	CourseId int `json:"course_id"`
	LessonId int `json:"lesson_id"`
	SlideId  int `json:"slide_id"`
}

type ChallengeResponse struct {
	Id       int `json:"id"`
	CourseId int `json:"course_id"`
}

type ChallegeResponseJSON struct {
	Challenge ChallengeResponse `json:"challenge"`
}

type ChallengeResource struct {
	Id         int    `json:"id"`
	CourseId   int    `json:"course_id"`
	FolderPath string `json:"folder_path"`
}

type TicketResponse struct {
	Ticket string `json:"ticket"`
}

type ChallengePatch struct {
	Id       int    `json:"id"`
	CourseID int    `json:"course_id"`
	Type     string `json:"type"`
	Kind     string `json:"kind"`

	Content     string `json:"content"`
	Explanation string `json:"explanation"`
	Outline     string `json:"outline"`
	Password    string `json:"password"`

	IsActive                 bool `json:"is_active"`
	IsHidden                 bool `json:"is_hidden"`
	IsExam                   bool `json:"is_exam"`
	IsFeedbackVisible        bool `json:"is_feedback_visible"`
	IsGradePassbackAvailable bool `json:"is_grade_passback_available"`

	Settings ChallengeSettings `json:"settings"`
	Features ChallengeFeatures `json:"features"`

	// TODO: Rubric stuff

	ScaffoldHash string `json:"scaffold_hash"`
	SolutionHash string `json:"solution_hash"`
	TestbaseHash string `json:"testbase_hash"`
}

type ChallengeSettings struct {
	BuildCommand                        string           `json:"build_command"`
	CheckCommand                        string           `json:"check_command"`
	RunCommand                          string           `json:"run_command"`
	TerminalCommand                     string           `json:"terminal_command"`
	MaxSubmissionsPerInterval           int              `json:"max_submissions_per_interval"`
	MaxSubmissionsWithIntermediateFiles int              `json:"max_submissions_with_intermediate_files"`
	Passback                            PassbackSettings `json:"passback"`
	PerTestCaseScores                   bool             `json:"per_testcase_scores"`
}

type ChallengeFeatures struct {
	AnonymousSubmissions  bool `json:"anonymous_submissions"`
	Arguments             bool `json:"arguments"`
	Check                 bool `json:"check"`
	ConfirmSubmit         bool `json:"confirm_submit"`
	Connect               bool `json:"connect"`
	Editor                bool `json:"editor"`
	Feedback              bool `json:"feedback"`
	Full                  bool `json:"full"`
	GitSubmissions        bool `json:"git_submissions"`
	IntermediateFiles     bool `json:"Intermediate_files"`
	Internet              bool `json:"internet"`
	ManualCompletion      bool `json:"manual_completion"`
	Mark                  bool `json:"mark"`
	Network               bool `json:"network"`
	RemoteDesktop         bool `json:"remote_desktop"`
	Run                   bool `json:"run"`
	RunBeforeSubmit       bool `json:"run_before_submit"`
	Terminal              bool `json:"terminal"`
	Treeview              bool `json:"treeview"`
	TreeviewInitiallyOpen bool `json:"treeview_initially_open"`
}

type PassbackSettings struct {
	MaxAutomaticScore int    `json:"max_automatic_score"`
	ScaleTo           int    `json:"scale_to"`
	ScoringMode       string `json:"scoring_mode"`
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
			ret := &Challenge{}
			ret.Id = resp.Challenge.Id
			ret.CourseId = resp.Challenge.CourseId
			ret.LessonId = lesson_id
			ret.SlideId = slide_id
			return ret, nil
		}
	}

	return nil, fmt.Errorf("Challenge for Slide %d Not Found", slide_id)
}

func UpdateChallenge(conn *client.Client, challenge *ChallengeResource) error {
	// TODO

	dir_entries, err := os.ReadDir(challenge.FolderPath)
	if err != nil {
		return err
	}

	for _, subdir := range dir_entries {
		wshelpers.UpdateChallengeRepo(conn, challenge.Id, challenge.FolderPath, subdir.Name())

	}

	return nil
}
