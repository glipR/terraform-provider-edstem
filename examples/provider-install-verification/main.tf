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

output "edu_lessons" {
    value = data.edstem_lesson.example
}
