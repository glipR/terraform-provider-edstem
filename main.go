// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"terraform-provider-edstem/internal/client"
	"terraform-provider-edstem/internal/md2ed"
	"terraform-provider-edstem/internal/provider"
	"terraform-provider-edstem/internal/resourceclients"

	"github.com/akamensky/argparse"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

type ImportArgs struct {
	CourseId string
	LessonId *string
	SlideId  *string

	ResourceName *string
	FolderPath   string
}

func import_tf(object_type string, args ImportArgs) error {
	var token = os.Getenv("EDSTEM_TOKEN")
	if token == "" {
		return fmt.Errorf("Please provide the EDSTEM_TOKEN environment variable")
	}
	var client, err = client.NewClient(&args.CourseId, &token)
	if err != nil {
		return err
	}

	var tf string

	if object_type == "course" {
		tf, err = resourceclients.CourseToTerraform(client, args.FolderPath)
		if err != nil {
			return err
		}
	} else if object_type == "lesson" {
		if args.LessonId == nil {
			return fmt.Errorf("LessonId must be provided")
		}
		var lesson_id, err = strconv.Atoi(*args.LessonId)
		if err != nil {
			return err
		}
		var resource_name string
		if args.ResourceName == nil {
			resource_name = "my_lesson"
		} else {
			resource_name = *args.ResourceName
		}
		tf, err = resourceclients.LessonToTerraform(client, lesson_id, resource_name, args.FolderPath)
		if err != nil {
			return err
		}
	} else if object_type == "slide" {
		if args.LessonId == nil {
			return fmt.Errorf("LessonId must be provided")
		}
		var lesson_id, err_lesson = strconv.Atoi(*args.LessonId)
		if err_lesson != nil {
			return err_lesson
		}
		if args.SlideId == nil {
			return fmt.Errorf("SlideId must be provided")
		}
		var slide_id, err_slide = strconv.Atoi(*args.SlideId)
		if err_slide != nil {
			return err_slide
		}
		var resource_name string
		if args.ResourceName == nil {
			resource_name = "my_slide"
		} else {
			resource_name = *args.ResourceName
		}
		tf, err = resourceclients.SlideToTerraform(client, lesson_id, slide_id, resource_name, args.FolderPath, nil)
		if err != nil {
			return err
		}
	} else if object_type == "challenge" {
		if args.LessonId == nil {
			return fmt.Errorf("LessonId must be provided")
		}
		var lesson_id, err_lesson = strconv.Atoi(*args.LessonId)
		if err_lesson != nil {
			return err_lesson
		}
		if args.SlideId == nil {
			return fmt.Errorf("SlideId must be provided")
		}
		var slide_id, err_slide = strconv.Atoi(*args.SlideId)
		if err_slide != nil {
			return err_slide
		}
		var resource_name string
		if args.ResourceName == nil {
			resource_name = "my_challenge"
		} else {
			resource_name = *args.ResourceName
		}
		tf, err = resourceclients.ChallengeToTerraform(client, lesson_id, slide_id, resource_name, args.FolderPath, nil, nil)
		if err != nil {
			return err
		}
	}

	const preamble = `terraform {
  required_providers {
	edstem = {
	  source = "hashicorp.com/edu/edstem"
	}
  }
}

provider "edstem" {
  course_id = "12108"
}`

	f, e := os.Create(path.Join(args.FolderPath, "main.tf"))
	if e != nil {
		return e
	}
	f.WriteString(preamble)
	f.WriteString("\n\n")
	f.WriteString(tf)
	return nil
}

func sysargs() {
	// terraform plan/apply
	// go run main.go import_tf lesson temp -c 12108 -l 36778
	// go run main.go render_ed examples/provider-install-verification/assets/test.md
	if len(os.Args) == 1 {
		// No args.
		terraform()
		return
	}
	if os.Args[1] == "import_tf" {
		parser := argparse.NewParser("print", "Prints provided string to stdout.\nExample: go run main.go import_tf lesson temp -c 12108 -l 36778")
		parser.SelectorPositional([]string{"import_tf"}, nil)
		resource_type := parser.SelectorPositional([]string{"course", "lesson", "slide", "challenge"}, nil)
		folder_path := parser.StringPositional(nil)
		course_id := parser.String("c", "course_id", &argparse.Options{Required: true, Help: "Course ID"})
		lesson_id := parser.String("l", "lesson_id", &argparse.Options{Required: false, Help: "Lesson ID"})
		slide_id := parser.String("s", "slide_id", &argparse.Options{Required: false, Help: "Slide ID"})
		resource_name := parser.String("r", "resource_name", &argparse.Options{Required: false, Help: "Resource Name"})

		err := parser.Parse(os.Args)
		if err != nil {
			fmt.Print(parser.Usage(err))
		}

		if *course_id == "" {
			course_id = nil
		}
		if *lesson_id == "" {
			lesson_id = nil
		}
		if *slide_id == "" {
			slide_id = nil
		}
		if *resource_name == "" {
			resource_name = nil
		}

		args := ImportArgs{
			CourseId:     *course_id,
			LessonId:     lesson_id,
			SlideId:      slide_id,
			ResourceName: resource_name,
			FolderPath:   *folder_path,
		}
		err = import_tf(*resource_type, args)
		if err != nil {
			fmt.Println("An error occurred: ", err)
		}
	} else if os.Args[1] == "render_ed" {
		fpath := os.Args[2]
		content, _ := os.ReadFile(fpath)
		fmt.Println(md2ed.RenderMDToEd(string(content)))
	} else if os.Args[1] == "render_md" {
		fpath := os.Args[2]
		content, _ := os.ReadFile(fpath)
		fmt.Println(md2ed.RenderEdToMD(string(content), "", false))
	}
}

func main() {
	sysargs()
}

func terraform() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/hashicorp/hashicups. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "hashicorp.com/edu/edstem",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
