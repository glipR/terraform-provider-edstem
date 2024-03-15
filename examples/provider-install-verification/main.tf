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

    due_at = formatdate("YYYY-MM-DD'T'hh:mm:ss.000Z", "2023-10-28T23:45:54+11:00")
}

resource "edstem_slide" "slide1" {
    type = "document"
    lesson_id = edstem_lesson.testing.id
    title = "Terraform Slide - Document"
    index = 1
    content = file("assets/test.md")
}

resource "edstem_slide" "slide2" {
    type = "quiz"
    lesson_id = edstem_lesson.testing.id
    title = "Terraform Slide - Quiz"
    index = 2
    content = ""
}

resource "edstem_slide" "slide3" {
    type = "code"
    lesson_id = edstem_lesson.testing.id
    title = "Terrafor Slide - Code"
    index = 3
    content = "Description content"
}

resource "edstem_challenge" "slide3_code" {
    slide_id = edstem_slide.slide3.id
    lesson_id = edstem_slide.slide3.lesson_id
    folder_path = "assets/code_challenge"
    folder_sha = sha1(join("", [for f in fileset(path.cwd, "assets/code_challenge/**"): filesha1("${path.cwd}/${f}")]))
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

output "edu_lessons" {
    value = data.edstem_lesson.example
}
