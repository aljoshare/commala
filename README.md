![Logo for commala](assets/logo.png)

# commala - A commit linter with a lot of rice

![GitHub Release](https://img.shields.io/github/v/release/aljoshare/commala?style=flat&logo=github&label=release&color=eca13d)
![GitHub Release Date](https://img.shields.io/github/release-date/aljoshare/commala?display_date=published_at&style=flat&logo=github&label=release%20date&color=eca13d)

![Static Badge](https://img.shields.io/badge/language-grey?logo=go)
![Static Badge](https://img.shields.io/badge/platform-linux-eca13d?logo=docker)
![Static Badge](https://img.shields.io/badge/arch-amd64-eca13d?logo=docker)
![Static Badge](https://img.shields.io/badge/arch-arm64-eca13d?logo=docker)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/aljoshare/commala/badge)](https://scorecard.dev/viewer/?uri=github.com/aljoshare/commala)

> “Go then, there are other commits than these.”  
> — commala, probably

commala is a commit linting tool that ensures that certain standards are met before you merge to keep your git history clean and consistent. commala is part of your Ka-tet, when you walk through the Wastelands of software development.

![Example of a commala workflow](assets/commala.gif)

## Validators

- Author Name: Check if the name of the author is set
- Author Email: Check if the email address of the author is set
- Branch: Check if the branch follows the [conventional branch specification](https://conventional-branch.github.io/)
- Message: Check if the commit message follows the [conventional commit specification](https://www.conventionalcommits.org)
- Sign-off: Check if the commit is signed off

## Getting started

If you want to use it on Github, try out the [Github Action](https://github.com/aljoshare/commala-action). You can find an example workflow [here](examples/github/example.yml). For Gitlab CI/CD, you can copy [this example](examples/gitlab/.gitlab-ci.yml) and modify it to your needs.

### Configuration

You can configure commala via `.commala.yml` or command line parameters:

```yaml
report:
  junit:
    path: commala-junit.xml # --report-junit-path
validate:
  author:
    name:
      enabled: true # --author-name-enabled
      whitelist: [] # --author-name-whitelist
    email:
      enabled: true # --author-email-enabled
      whitelist: [] # --author-email-whitelist
  branch:
    enabled: true # --branch-enabled
    whitelist: [] # --branch-whitelist
  message:
    enabled: true # --message-enabled
    whitelist: [] # --message-whitelist
  signoff:
    enabled: true # --signoff-enabled
    whitelist: [] # --signoff-whitelist
```

### Contributor Whitelists

Commala supports whitelisting specific contributors (by email) to skip validation for their commits. This is useful for automated bot accounts like Dependabot or Renovate that may not follow conventional commit standards.

Each validator can have its own whitelist configured via `.commala.yml`:

```yaml
validate:
  branch:
    enabled: true
    whitelist:
      - "dependabot[bot]@users.noreply.github.com"
      - "renovate[bot]@users.noreply.github.com"
  message:
    enabled: true
    whitelist:
      - "dependabot[bot]@users.noreply.github.com"
```

Or via CLI flags:

```bash
commala check HEAD~5 \
  --branch-whitelist="dependabot[bot]@users.noreply.github.com" \
  --message-whitelist="dependabot[bot]@users.noreply.github.com"
```

**How It Works:**

- Commits from whitelisted authors are marked as "skipped" during validation
- Skipped commits are clearly marked in console output (gray color)
- JUnit reports include `<skipped>` elements for whitelisted commits
- Skipped commits don't count as failures
- Whitelist matching uses exact email comparison (case-sensitive)

### CLI

The commala command is pretty easy. You can run the checks on all commits like this:

```shell
commala check
```

If you want to check all commits, just pass two dots:

```shell
commala check ..
```

If you want to specify the commit to start and check until HEAD, just specify the commit hash followed by two dots:

```shell
commala check a1b2c3d4e5f67890abcdef1234567890abcdef12..
```

If you want to specify the commit to end the check and start from the beginning, just specify the commit hash preceded by two dots:

```shell
commala check ..a1b2c3d4e5f67890abcdef1234567890abcdef12
```

If you want to specify a commit range, just specify two commit hashes with two dots between them:

```shell
commala check f725bf88adb76df5c8c576b514def199e20fc6a0..a1b2c3d4e5f67890abcdef1234567890abcdef12
```

If you want to specify a negative index, just use the swung dash notation:

```shell
commala check HEAD~3
```

### Result

To make it easy to use commala as part of a CI/CD job, it will output the result on the command line but also writes the result in JUnit XML format, so that it can be picked up by the source code versioning system of your choice. If one of the checks fails, commala will exit with a non-zero status.
