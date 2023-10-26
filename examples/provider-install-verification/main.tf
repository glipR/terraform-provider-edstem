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

    timer_duration = 30
}

output "edu_lessons" {
    value = data.edstem_lesson.example
}
