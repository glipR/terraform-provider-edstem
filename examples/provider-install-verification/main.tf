terraform {
  required_providers {
    edstem = {
      source = "hashicorp.com/edu/edstem"
    }
  }
}

provider "edstem" {
    course_id = "12108"
}

data "edstem_lesson" "example" {
    id = 36771
}

resource "edstem_lesson" "testing" {
    title = "Terraform Testing"

    timer_duration = 120

    password = "terraform_is_cool"

    index = 3

    due_at = formatdate("YYYY-MM-DD'T'hh:mm:ssZ", "2023-10-28T23:45:54+11:00")
}

resource "edstem_slide" "slide1" {
    type = "document"
    lesson_id = edstem_lesson.testing.id
    title = "Terraform Slide 1 - Document"
    index = 1
    content = file("assets/test.md")
}

resource "edstem_slide" "slide2" {
    type = "quiz"
    lesson_id = edstem_lesson.testing.id
    title = "Terraform Slide 2 - Quiz"
    index = 2
    content = ""
}

resource "edstem_slide" "slide3" {
    type = "code"
    lesson_id = edstem_lesson.testing.id
    title = "Terraform Slide 3 - Code Custom Marking"
    index = 3
    content = "Description content"
}

resource "edstem_challenge" "slide3_code" {
    slide_id = edstem_slide.slide3.id
    lesson_id = edstem_slide.slide3.lesson_id
    folder_path = "assets/code_challenge"
    folder_sha = sha1(join("", [for f in fileset(path.cwd, "assets/code_challenge/**"): filesha1("${path.cwd}/${f}")]))

    type = "custom"
    custom_mark_time_limit_ms = 2500
    custom_run_command = "terraform run"

    feature_anonymous_submissions = true
    feature_manual_completion = false

    criteria = file("assets/criteria.json")
}

resource "edstem_slide" "slide4" {
    type = "code"
    lesson_id = edstem_lesson.testing.id
    title = "Terraform Slide 4 - Code Input Output Marking"
    index = 4
    content = "Description content"
}

resource "edstem_challenge" "slide4_code" {
    slide_id = edstem_slide.slide4.id
    lesson_id = edstem_slide.slide4.lesson_id
    // Just using the same content for code challenge.
    folder_path = "assets/code_challenge"
    folder_sha = sha1(join("", [for f in fileset(path.cwd, "assets/code_challenge/**"): filesha1("${path.cwd}/${f}")]))

    type = "code"
    testcase_json = file("assets/testcases.json")
    testcase_easy = true
    testcase_pty = false
    testcase_overlay_test_files = true

    feature_anonymous_submissions = true
    feature_manual_completion = false
}

resource "edstem_question" "question1" {
    index = 1
    lesson_slide_id = edstem_slide.slide2.id
    type = "multiple-choice"
    answers = ["<document version=\"2.0\"><paragraph>Test1</paragraph></document>", "<document version=\"2.0\"><paragraph>Test2</paragraph></document>"]
    explanation = "<document version=\"2.0\"><paragraph>The answer is not A</paragraph></document>"
    solution = [1]
}

resource "edstem_question" "question2" {
    index = 1
    lesson_slide_id = edstem_slide.slide2.id
    type = "multiple-choice"
    question_document_string = file("assets/question2.md")
}

resource "edstem_slide" "slide5" {
    type = "pdf"
    lesson_id = edstem_lesson.testing.id
    title = "Terraform Slide 5 - PDF"
    index = 5
    file_path = "assets/test.pdf"
}

resource "edstem_slide" "slide6" {
    type = "video"
    lesson_id = edstem_lesson.testing.id
    title = "Terraform Slide 6 - Video"
    index = 6
    url = "https://www.youtube.com/watch?v=feIeCR6oFNM"
}

resource "edstem_slide" "slide7" {
    type = "webpage"
    lesson_id = edstem_lesson.testing.id
    title = "Terraform Slide 7 - Webpage"
    index = 7
    url = "https://www.google.com/"
}

resource "edstem_slide" "slide8" {
    type = "html"
    lesson_id = edstem_lesson.testing.id
    title = "Terraform Slide 8 - HTML"
    index = 8
    content = "<h1>TEST</h1>"
}

output "edu_lessons" {
    value = data.edstem_lesson.example
}
