# terraform-provider-edstem

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## How do I use?

Include the following snippet in the `required_providers` block of your terraform code:

```
terraform {
  required_providers {
    edstem = {
      source = "hashicorp.com/edu/edstem"
    }
  }
}
```

Unfortunately `-parallelism=1` must be used with this provider because we can't have multiple slides being applied at the same time.

## How do I import existing Ed lessons etc. into my terraform?

You'll need to invoke this module yourself (TODO: Add what this script is for people installing the package)

```
// Make a new directory for terraform
mkdir my_course
// Specify what resources you'd like to bring through (and everything more granular)
// This one just grabs lesson 36778 from course 12108
go run main.go import_tf lesson my_course -c 12108 -l 36778
// This one grabs all lessons from course 12108
go run main.go import_tf course my_course -c 12108
```

(TODO: Make the import script have the (default) option to also fill the tfstate file)

## Currently not functional components

* Terraform Actions:
    * Destroying objects
    * Reading current state for diff
* Slides
    * Survey, SQL Challenge, RStudio Challenge, Jupyter Challenge, Web Challenge
    * Question types other than Multi-Choice
    * Code Challenges that aren't `none`, `custom` or `code`.
* Misc
    * Images in content


## Cautionary areas

* Ed/MD rendering hasn't been rigorously tested
* Some minor elements of the challenges api aren't fully understood, so some minor differences may occur when importing/re-applying.

## Development Notes

Once you've written your provider, you'll want to [publish it on the Terraform Registry](https://developer.hashicorp.com/terraform/registry/providers/publishing) so that others can use it.

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
