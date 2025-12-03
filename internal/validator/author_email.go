package validator

import (
	"fmt"
	"time"

	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/utils"
	log "github.com/sirupsen/logrus"
)

type AuthorEmailValidator struct {
	authorEmails map[string]string
}

func (m AuthorEmailValidator) Validate(cr *git.CommitRange, g git.Git) (*ValidationResult, error) {
	log.Debugf("Validating if commits has author emails from %s to %s", cr.From, cr.To)
	vr := ValidationResult{
		Validator:  "Author Email",
		Assertions: 0,
		Failures:   0,
	}
	defer func() {
		vr.Duration = utils.Duration(time.Now())
	}()
	var err error
	m.authorEmails, err = g.GetCommitAuthorEmails(cr.From, cr.To)
	if err != nil {
		return nil, err
	}
	vm, err := m.hasAuthorEmail()
	if err != nil {
		return nil, err
	}
	if utils.AllTrue(vm) {
		vr.Valid = true
		vr.Summary = "All commit messages have author email set\n"
	} else {
		vr.Valid = true
		vr.Summary = "Not all commit messages have author email set\n"
	}
	vr.Messages = make(map[string]ResultMessage, len(m.authorEmails))
	for i := range m.authorEmails {
		vr.Assertions++
		if !vm[i] {
			vr.Messages[i] = ResultMessage{false, fmt.Sprintf("Commit Message: \"%s\" has no author email\n", i)}
			vr.Failures++
		} else {
			vr.Messages[i] = ResultMessage{true, fmt.Sprintf("Commit Message: \"%s\" has author email\n", i)}
		}
	}
	return &vr, nil
}

func (m AuthorEmailValidator) hasAuthorEmail() (map[string]bool, error) {
	am := make(map[string]bool, len(m.authorEmails))
	for ch, an := range m.authorEmails {
		if an != "" {
			am[ch] = true
		} else {
			am[ch] = false
		}
	}
	return am, nil
}
