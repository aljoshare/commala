package git

import (
	"github.com/go-git/go-git/v6/plumbing/object"
)

type MockGit struct{}

func (m MockGit) GetBranchName() (string, error) {
	return "feature/mock-branch", nil
}

func (m MockGit) GetCommitMessages(from string, to string) ([]string, error) {
	var l []string
	l = make([]string, 0, 5)
	l = append(l, "feat: add new feature\nSigned-off-by: Johen Doe <john@doe.com>")
	l = append(l, "fix: fix a bug\nSigned-off-by: Johen Doe <john@doe.com>")
	return l, nil
}

func (m MockGit) GetCommitAuthorNames(from string, to string) (map[string]string, error) {
	var authorMap = make(map[string]string)
	authorMap["commit1"] = "Test User"
	authorMap["commit2"] = "Another User"
	return authorMap, nil
}

func (m MockGit) GetCommitAuthorEmails(from string, to string) (map[string]string, error) {
	var emailMap = make(map[string]string)
	emailMap["commit1"] = "test@user.de"
	emailMap["commit2"] = "anothertest@user.de"
	return emailMap, nil
}

func (m MockGit) GetCommitHash(ref string) (*object.Commit, error) {
	c := object.Commit{}
	return &c, nil
}
func (m MockGit) GetInitialCommit() (string, error) {
	return "initialCommitHash", nil
}
func (m MockGit) GetLatestCommit() (string, error) {
	return "latestCommitHash", nil
}
func (m MockGit) ParseCommitRange(commitrange string) (*CommitRange, error) {
	cr := CommitRange{
		From: "initialCommitHash",
		To:   "latestCommitHash",
	}
	return &cr, nil
}
