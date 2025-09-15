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
		Validator: "Author Email",
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
	if utils.AllTrueMap(vm) {
		vr.Valid = true
		vr.Message = append(vr.Message, "All commit messages have author email set\n")
		return &vr, nil
	}
	for _, v := range vm {
		if !v {
			for i, m := range m.authorEmails {
				if !vm[i] {
					vr.Message = append(vr.Message, fmt.Sprintf("Commit Message: \"%s\" has no author email\n", m))
				}
			}
		}
	}
	return &vr, nil
}

func (m AuthorEmailValidator) hasAuthorEmail() (map[string]bool, error) {
	am := make(map[string]bool)
	for ch, an := range m.authorEmails {
		if an != "" {
			am[ch] = true
		} else {
			am[ch] = false
		}
	}
	return am, nil
}
