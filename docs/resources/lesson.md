---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "edstem_lesson Resource - terraform-provider-edstem"
subcategory: ""
description: |-
  
---

# edstem_lesson (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `title` (String) Lesson title.

### Optional

- `attempts` (Number) The number of attempts that the user can submit
- `available_at` (String) The timestamp the lesson becomes available.
- `due_at` (String) The timestamp the lesson is due.
- `grade_passback_auto_send` (Boolean) Whether to automatically do grade passback.
- `grade_passback_mode` (String)
- `grade_passback_scale_to` (String)
- `index` (Number) The index of the lesson within the lessons listing.
- `is_hidden` (Boolean) Whether to hide the lesson from students (inaccessible even with link).
- `is_timed` (Boolean) Whether this lesson is timed.
- `is_unlisted` (Boolean) Whether this lesson is only accessible to students via a link.
- `kind` (String) What type of lesson this is. Defaults to content.
- `late_submissions` (Boolean) Whether to allow late submissions.
- `locked_at` (String) Timestamp of when this lesson will no longer accept any submissions.
- `module_id` (Number) ID of the module within the lesson listing.
- `openable` (Boolean)
- `openable_without_attempt` (Boolean)
- `outline` (String) Short description of the lesson to show in the lesson listing view.
- `password` (String) Blocks user access without entering the password.
- `prerequisites` (List of Number) List of Lesson IDs that need to be completed before this lesson can be commenced.
- `quiz_active_status` (String)
- `quiz_mode` (String)
- `quiz_question_number_style` (String)
- `release_challenge_solutions` (Boolean) Release challenge solutions once due date has passed.
- `release_challenge_solutions_while_active` (Boolean) Release challenge solutions while the lesson is active.
- `release_feedback` (Boolean) Release marking feedback once due date has passed.
- `release_feedback_while_active` (Boolean) Release marking feedback while the lesson is active.
- `release_quiz_correctness_only` (Boolean) Release whether a students quiz answer was right/wrong on completion.
- `release_quiz_solutions` (Boolean) Release the correct quiz answer on completion.
- `reopen_submissions` (Boolean)
- `require_user_override` (Boolean)
- `solutions_at` (String) The timestamp the lesson solutions becomes available.
- `state` (String)
- `timer_duration` (Number) For timed lessons, how long in minutes the duration of the lesson lasts.
- `timer_expiration_access` (Boolean)
- `tutorial_regex` (String) Restrict access to students whose tutorial group match the regular expression.
- `type` (String)

### Read-Only

- `id` (Number) Integer ID identifying the Lesson. This can be found in the URL of a lesson. For example, `https://edstem.org/au/courses/<course_id>/lessons/<lesson_id>/slides/<slide_id>`. Here we want the lesson_id.
- `last_updated` (String)