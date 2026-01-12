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

func (m AuthorNameValidator) Validate(cr *git.CommitRange, g git.Git, whitelist []string) (*ValidationResult, error) {
	log.Debugf("Validating if commits has author names from %s to %s", cr.From, cr.To)
	vr := ValidationResult{
		Validator: "Author Name",
		Skipped:   0,
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

	vr.Messages = make(map[string]ResultMessage, len(m.authorNames))
	for commitHash := range m.authorNames {
		// Check whitelist before validation
		whitelisted, authorEmail, err := IsWhitelisted(commitHash, whitelist, g, cr)
		if err != nil {
			return nil, err
		}

		if whitelisted {
			vr.Skipped++
			vr.Messages[commitHash] = NewSkippedResultMessage(authorEmail)
			continue
		}

		vr.Assertions++
		if !vm[commitHash] {
			vr.Messages[commitHash] = ResultMessage{Valid: false, Message: fmt.Sprintf("Commit Message: \"%s\" has no author name\n", commitHash)}
			vr.Failures++
		} else {
			vr.Messages[commitHash] = ResultMessage{Valid: true, Message: fmt.Sprintf("Commit Message: \"%s\" has author name\n", commitHash)}
		}
	}

	// Update summary with skip count
	if vr.Failures == 0 {
		vr.Valid = true
		if vr.Skipped > 0 {
			vr.Summary = fmt.Sprintf("All author names are present (%d skipped)\n", vr.Skipped)
		} else {
			vr.Summary = "All commit messages have author names\n"
		}
	} else {
		vr.Valid = false
		totalChecked := vr.Assertions + vr.Skipped
		passedCount := vr.Assertions - vr.Failures
		if vr.Skipped > 0 {
			vr.Summary = fmt.Sprintf("%d of %d author names are present (%d skipped)\n", passedCount, totalChecked, vr.Skipped)
		} else {
			vr.Summary = "Not all commit messages have author names\n"
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
