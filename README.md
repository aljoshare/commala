![Logo for commala](assets/logo.png)

# commala - A commit linter with a lot of rice

![GitHub Release](https://img.shields.io/github/v/release/aljoshare/commala?style=flat&logo=github&label=release&color=eca13d)
![GitHub Release Date](https://img.shields.io/github/release-date/aljoshare/commala?display_date=published_at&style=flat&logo=github&label=release%20date&color=eca13d)

![Static Badge](https://img.shields.io/badge/language-grey?logo=go)
![Static Badge](https://img.shields.io/badge/platform-linux-eca13d?logo=docker)
![Static Badge](https://img.shields.io/badge/arch-amd64-eca13d?logo=docker)
![Static Badge](https://img.shields.io/badge/arch-arm64-eca13d?logo=docker)

> “Go then, there are other commits than these.”  
> — commala, probably

commala is a commit linting tool that ensures that certain standards are met before you merge to keep your git history clean and consistent. Inspired by the The Dark Tower, commala is part of your Ka-tet, when you walk through the Wastelands of software development.

## Getting started

If you want to use it on Github, try out the Github Action. For Gitlab CI/CD, you can copy this example and modify it to your needs.

### Configuration

You can configure commala via `.commala.yml` or command line parameters:

```yaml
report:
  junit:
    path: commala-junit.xml # --report-junit-path
validators:
  author:
    name:
      enabled: true # --author-name-enabled
    email:
      enabled: true # --author-email-enabled
  branch:
    enabled: true # --branch-enabled
  message:
    enabled: true # --message-enabled
  signoff:
    enabled: true # --signoff-enabled
```

### CLI

The commala command is pretty easy. You can run the checks on all commits like this:

```shell
commala check
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

### Result

To make it easy to use commala as part of a CI/CD job, it will output the result on the command line but also writes the result in JUnit XML format, so that it can be picked up by the source code versioning system of your choice. If one of the checks fails, commala will exit with a non-zero status.
