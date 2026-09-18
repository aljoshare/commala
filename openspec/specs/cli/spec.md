# CLI and Configuration Specification

## Purpose

This specification defines the command-line interface, configuration loading hierarchy, commit range parsing, and exit code semantics of the Commala tool.

## Requirements

### Requirement: Commit Range Parsing

The system SHALL accept and parse Git commit ranges specified as standard two-dot range syntax or negative relative index notation.

#### Scenario: Two-dot commit range notation

- **WHEN** the user provides a commit range argument of the form "<from-hash>..<to-hash>"
- **THEN** the system resolves the commits between the two hashes for validation.

#### Scenario: Open-ended start commit range

- **WHEN** the user provides an argument of the form "..<to-hash>"
- **THEN** the system resolves commits starting from the repository initial commit up to the specified commit hash.

#### Scenario: Open-ended end commit range

- **WHEN** the user provides an argument of the form "<from-hash>.."
- **THEN** the system resolves commits starting from the specified commit hash up to the repository latest commit (HEAD).

#### Scenario: Full history range

- **WHEN** the user provides the argument ".."
- **THEN** the system resolves all commits from the repository's initial commit up to the latest commit (HEAD).

#### Scenario: Negative index notation

- **WHEN** the user provides an argument of the form "HEAD~<n>" (e.g., "HEAD~5")
- **THEN** the system resolves the commit range from HEAD back through the previous n commits.

#### Scenario: Invalid range argument syntax

- **WHEN** the user provides an argument that does not contain ".." and does not start with "HEAD~"
- **THEN** the system reports an error stating that the argument must be a commit range or a negative index and exits with code 1.

### Requirement: Configuration Discovery and Custom Path

The system SHALL discover configuration from standard file locations or load from a custom path specified via CLI flag or environment variable.

#### Scenario: Default configuration discovery

- **WHEN** no custom configuration path is specified
- **THEN** the system searches for `.commala.yml` in the current working directory, falling back to `$HOME/.commala/.commala.yaml`.

#### Scenario: Custom configuration path via flag

- **WHEN** the `--config` flag is supplied with a path to a configuration file
- **THEN** the system loads configuration exclusively from the specified file path.

#### Scenario: Custom configuration path via environment variable

- **WHEN** the `COMMALA_CONFIG` environment variable is set
- **THEN** the system loads configuration from the path specified by the environment variable.

### Requirement: Configuration Override Precedence

The system SHALL allow CLI flags and environment variables (`COMMALA_*`) to override values defined in the configuration file.

#### Scenario: CLI flag overrides configuration file value

- **WHEN** a validator is enabled in `.commala.yml` but the corresponding CLI flag (e.g., `--branch-enabled=false`) is passed
- **THEN** the CLI flag value takes precedence and disables the validator.

#### Scenario: Environment variable overrides configuration file value

- **WHEN** an environment variable (e.g., `COMMALA_VALIDATE_BRANCH_ENABLED=false`) is set
- **THEN** the environment variable value takes precedence over `.commala.yml`.

### Requirement: Process Exit Codes

The system SHALL return appropriate process exit codes based on validation outcomes and operational errors.

#### Scenario: All validations pass or are skipped

- **WHEN** all executed validators report valid status true with zero failures
- **THEN** the application terminates with exit code 0.

#### Scenario: At least one validation fails

- **WHEN** any validator reports valid status false with one or more failures
- **THEN** the application terminates with exit code 1.

#### Scenario: Operational or Git error occurs

- **WHEN** an error occurs during argument parsing, configuration reading, or Git repository operations
- **THEN** the error is reported to standard error and the application terminates with exit code 1.
