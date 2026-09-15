package git

import (
	"testing"

	"github.com/go-git/go-git/v6/plumbing"
)

func clearCiEnv(t *testing.T) {
	t.Helper()
	ciEnvVars := []string{
		"COMMALA_BRANCH_NAME",
		"GITHUB_HEAD_REF",
		"CI_MERGE_REQUEST_SOURCE_BRANCH_NAME",
		"CI_COMMIT_REF_NAME",
		"CI_COMMIT_BRANCH",
		"BITBUCKET_BRANCH",
	}
	for _, envVar := range ciEnvVars {
		t.Setenv(envVar, "")
	}
}

func TestGetBranchName_CommalaBranchNameEnv(t *testing.T) {
	clearCiEnv(t)
	t.Setenv("COMMALA_BRANCH_NAME", "my-custom-branch")
	r := RealGit{}
	branch, err := r.GetBranchName()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch != "my-custom-branch" {
		t.Errorf("expected branch my-custom-branch, got %s", branch)
	}
}

func TestResolveBranchName_NormalBranch(t *testing.T) {
	clearCiEnv(t)
	ref := plumbing.NewReferenceFromStrings("refs/heads/feature/login", "0123456789abcdef0123456789abcdef01234567")
	branch, err := resolveBranchName(ref, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch != "feature/login" {
		t.Errorf("expected branch feature/login, got %s", branch)
	}
}

func TestResolveBranchName_DetachedHead_GithubHeadRef(t *testing.T) {
	clearCiEnv(t)
	t.Setenv("GITHUB_HEAD_REF", "release-please--branches--main")
	ref := plumbing.NewReferenceFromStrings("HEAD", "0123456789abcdef0123456789abcdef01234567")
	branch, err := resolveBranchName(ref, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch != "release-please--branches--main" {
		t.Errorf("expected release-please--branches--main, got %s", branch)
	}
}

func TestResolveBranchName_DetachedHead_GitlabCi(t *testing.T) {
	clearCiEnv(t)
	t.Setenv("CI_MERGE_REQUEST_SOURCE_BRANCH_NAME", "feature/gitlab-mr")
	ref := plumbing.NewReferenceFromStrings("HEAD", "0123456789abcdef0123456789abcdef01234567")
	branch, err := resolveBranchName(ref, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch != "feature/gitlab-mr" {
		t.Errorf("expected feature/gitlab-mr, got %s", branch)
	}
}

func TestResolveBranchName_DetachedHead_NoCiEnv(t *testing.T) {
	clearCiEnv(t)
	ref := plumbing.NewReferenceFromStrings("HEAD", "0123456789abcdef0123456789abcdef01234567")
	branch, err := resolveBranchName(ref, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch != "HEAD" {
		t.Errorf("expected HEAD, got %s", branch)
	}
}

func TestMockGit_BranchName(t *testing.T) {
	mDefault := MockGit{}
	name, err := mDefault.GetBranchName()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "feature/mock-branch" {
		t.Errorf("expected feature/mock-branch, got %s", name)
	}

	mCustom := MockGit{BranchName: "release-please--branches--main"}
	name, err = mCustom.GetBranchName()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "release-please--branches--main" {
		t.Errorf("expected release-please--branches--main, got %s", name)
	}
}
