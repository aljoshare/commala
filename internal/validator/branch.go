package validator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/utils"
	log "github.com/sirupsen/logrus"
)

type BranchValidator struct {
	branchName string
}

func (b BranchValidator) Validate(g git.Git) (*ValidationResult, error) {
	log.Debug("Validating Branch")
	vr := ValidationResult{
		Validator: "Branch",
	}
	defer func() {
		vr.Duration = utils.Duration(time.Now())
	}()
	var err error
	b.branchName, err = g.GetBranchName()
	if err != nil {
		return nil, err
	}
	valid, err := b.isConventional()
	if valid {
		vr.Valid = true
		vr.Message = append(vr.Message, fmt.Sprintf("Branch name \"%s\" is conventional", b.branchName))
	} else {
		vr.Valid = false
		vr.Message = append(vr.Message, fmt.Sprintf("Branch name \"%s\" is not conventional", b.branchName))
	}
	if err != nil {
		return nil, err
	}
	return &vr, nil
}

func (b BranchValidator) isConventional() (bool, error) {
	var s []bool
	s = make([]bool, 0, 5)
	matched, err := regexp.Match(`^(feature|feat|bugfix|fix|hotfix|release|chore|docs|test|refactor)/[a-zA-Z0-9._-]+$`, []byte(b.branchName))
	if err != nil {
		return false, fmt.Errorf("can't match branch name")
	}
	s = append(s, matched)
	return utils.AllTrueSlice(s), nil
}
