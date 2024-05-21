---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "edstem_lesson Data Source - terraform-provider-edstem"
subcategory: ""
description: |-
  
---

# edstem_lesson (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `module_id` (Number)
- `prerequisites` (Attributes List) (see [below for nested schema](#nestedatt--prerequisites))

### Read-Only

- `available_at` (String)
- `due_at` (String)
- `id` (Number) The ID of this resource.
- `quiz_settings` (Attributes) (see [below for nested schema](#nestedatt--quiz_settings))
- `title` (String)

<a id="nestedatt--prerequisites"></a>
### Nested Schema for `prerequisites`

Read-Only:

- `required_lesson_id` (Number)


<a id="nestedatt--quiz_settings"></a>
### Nested Schema for `quiz_settings`

Read-Only:

- `active_status` (String)
- `mode` (String)
- `question_number_style` (String)