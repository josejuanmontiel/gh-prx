# gh-prx

<img src="https://github.com/ilaif/gh-prx/raw/main/assets/logo.png" width="200">

A GitHub (`gh`) CLI extension to automate the daily work with **branches**, **commits** and **pull requests**.

## Usage

[Jump to installation](#installation)

1. Checking out to an automatically generated branch:

    ```sh
    gh prx checkout-new 1234 # Where 1234 is the issue's number/code
    ```

2. Creating a new PR with automatically generated title/body and checklist prompt:

    ```sh
    gh prx create
    ```

> Explore further by running `gh prx --help`

## Why?

As developers, we rely heavily on git and our git provider (in this case, GitHub).

Many of us find the terminal and CLI applications as our main toolkit.

`gh-prx` helps with automating and standardizing how we work with git and GitHub to create a faster and more streamlined workflow for individuals and teams.

## Features

* Automatically creating new branches named based on issues fetched from project management tools (GitHub, and more in the future...)
* Extended PR creation:
  * Automatically push branch to origin
  * Parse branch names by a pattern into a customized PR title and description template
  * Add labels based on issue types
  * Filter commits and display them in the PR description
  * Interactively answer PR checklists before creating the PR
  * All `gh pr create` original flags are extended into the tool

> `gh-prx` is an early-stage project. Got a new feature in mind? Open a pull request or a [feature request](https://github.com/ilaif/gh-prx/issues/new) 🙏

## Configuration

Configuration is provided from `.github/.gh-prx.yaml` and is advised to be committed to git to maintain standardization across the team.

The default values for `.gh-prx.yaml` are:

```yaml
branch:
   template: "{{.Type}}/{{.Issue}}-{{.Description}}" # Branch name template
   patterns: # A map of patterns to match for each template variable
      Type: "fix|feat|chore|docs|refactor|test|style|build|ci|perf|revert"
      Issue: "#[0-9]+"
      Description: ".*"
   token_separators: ["-", "_"] # Characters used to separate branch name into a human-readable string
pr:
   title: "{{.Type}}({{.Issue}}): {{ humanize .Description}}" # PR title template
   ignore_commits_patterns: ["^wip"] # Patterns that when matched, filters out a commit message from the {{.Commits}} variable
   answer_checklist: true # Whether to prompt the user to answer PR description checklists. Possible answers: yes, no, skip (remove the item)
   push_to_remote: true # Whether to push the local changes to remote before creating the PR
issue:
   provider: github # The provider to use for fetching issue details
```

### PR Description (Body)

The PR description is based on the repo's `.github/pull_request_template.md`. If this file does not exist, a default template is used:

   ```markdown
   {{with .Issue}}Closes {{.Issue}}.

   {{end}}## Description

   {{ humanize .Description}}

   ## PR Checklist

   - [ ] Tests are included
   - [ ] Documentation is changed or added
   ```

### Templating

The templates are based on [Go text template](https://pkg.go.dev/text/template).

Additional template functions:

#### `humanize`

Humanizes a string by separating it into tokens (words) based on `branch.token_separators`.

Example:

Given:

* `token_separators: ["-"]`
* `{{.Description}}`: "my-dashed-string"
* Template:

   ```go-template
   This is "{{ humanize .Description}}"
   ```

Result:

```txt
This is "my dashed string"
```

Special template variable names:

* `{{.Type}}` - Used to interpret GitHub labels to add to the PR and issue type to add the branch name.
* `{{.Issue}}` - Used as a placeholder for the issue number/code when creating a new branch.
* `{{.Description}}` - Used as a placeholder for the issue title when creating a new branch.
* `{{.Commits}}` - Used as a placeholder in a PR description (body) to iterate over filtered commits.

## Installation

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation)

   _Installation requires a minimum version (2.0.0) of the the GitHub CLI that supports extensions._

2. Install this extension:

   ```sh
   gh extension install ilaif/gh-prx
   ```

<details>
   <summary><strong>Installing Manually</strong></summary>

> If you want to install this extension **manually**, follow these steps:

1. Clone the repo

   ```bash
   # git
   git clone https://github.com/ilaif/gh-prx
   # GitHub CLI
   gh repo clone ilaif/gh-prx
   ```

2. Cd into it

   ```bash
   cd gh-prx
   ```

3. Install it locally

   ```bash
   gh extension install .
   ```

</details>

## Questions, bug reporting and feature requests

You're more than welcome to [Create a new issue](https://github.com/ilaif/gh-prx/issues/new) or contribute.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. See [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## License

gh-prx is licensed under the [MIT](https://choosealicense.com/licenses/mit/) license. For more information, please see the [LICENSE](LICENSE) file.
