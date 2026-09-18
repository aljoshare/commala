# Reporting Specification

## Purpose

This specification defines the console tabular formatting and JUnit XML report generation capabilities for Commala validation results.

## Requirements

### Requirement: Console Result Table Display

The system SHALL print an interactive, colored summary table of validation results to the terminal upon completion of validation checks.

#### Scenario: Displaying validator outcomes

- **WHEN** validation completes
- **THEN** the system renders a table showing each executed validator, its pass/fail/skip status, assertion count, failure count, duration, and detailed messages for non-compliant commits.

#### Scenario: Distinct visual styling for skipped checks

- **WHEN** a validator check is skipped due to author email, pattern matching, or global branch whitelisting
- **THEN** the result table marks the skipped items distinctively from passes and failures and displays the skip count.

### Requirement: JUnit XML Report Generation

The system SHALL write a structured JUnit XML report containing test suites and test cases corresponding to executed and skipped validators.

#### Scenario: Default JUnit report file destination

- **WHEN** no custom JUnit path is specified
- **THEN** the report is written to "commala-junit.xml" in the current working directory.

#### Scenario: Custom JUnit report path

- **WHEN** a custom path is configured via `report.junit.path`, the `--report-junit-path` flag, or the `COMMALA_REPORT_JUNIT_PATH` environment variable
- **THEN** the JUnit report is written to the specified file path.

### Requirement: JUnit Failure Information

The system SHALL record failed validation assertions as `<failure>` elements inside corresponding `<testcase>` XML tags.

#### Scenario: Validator failure in JUnit report

- **WHEN** a commit fails a validation check (such as non-conventional message or missing sign-off)
- **THEN** the generated `<testcase>` contains a `<failure>` child element with the failure message and validator type attribute.

### Requirement: JUnit Skipped Information with Reason

The system SHALL record skipped validation assertions as `<skipped>` elements inside corresponding `<testcase>` XML tags, preserving the skip reason.

#### Scenario: Skipped check in JUnit report

- **WHEN** a validation check is skipped due to author whitelist, branch pattern, message pattern, or global branch whitelist
- **THEN** the generated `<testcase>` contains a `<skipped>` child element with a `message` attribute detailing the exact reason the check was skipped.
