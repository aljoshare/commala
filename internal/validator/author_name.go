package validator

import (
	"fmt"
	"time"

	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/utils"
	log "github.com/sirupsen/logrus"
)

type AuthorNameValidator struct {
	authorNames map[string]string
}

func (m AuthorNameValidator) Validate(cr *git.CommitRange, g git.Git) (*ValidationResult, error) {
	log.Debugf("Validating if commits has author names from %s to %s", cr.From, cr.To)
	vr := ValidationResult{
		Validator: "Author Name",
	}
	defer func() {
		vr.Duration = utils.Duration(time.Now())
	}()
	var err error
	m.authorNames, err = g.GetCommitAuthorNames(cr.From, cr.To)
	if err != nil {
		return nil, err
	}
	vm, err := m.hasAuthorName()
	if err != nil {
		return nil, err
	}
	if utils.AllTrue(vm) {
		vr.Valid = true
		vr.Summary = "All commit messages have author names\n"
	} else {
		vr.Valid = false
		vr.Summary = "Not all commit messages have author names\n"
	}
	vr.Messages = make(map[string]ResultMessage, len(m.authorNames))
	for i := range m.authorNames {
		vr.Assertions++
		if !vm[i] {
			vr.Messages[i] = ResultMessage{false, fmt.Sprintf("Commit Message: \"%s\" has no author name\n", i)}
			vr.Failures++
		} else {
			vr.Messages[i] = ResultMessage{true, fmt.Sprintf("Commit Message: \"%s\" has author name\n", i)}
		}
	}
	return &vr, nil
}

func (m AuthorNameValidator) hasAuthorName() (map[string]bool, error) {
	am := make(map[string]bool, len(m.authorNames))
	for ch, an := range m.authorNames {
		if an != "" {
			am[ch] = true
		} else {
			am[ch] = false
		}
	}
	return am, nil
}
