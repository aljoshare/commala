package validator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/utils"
	log "github.com/sirupsen/logrus"
)

type SignOffValidator struct {
	messages map[string]string
}

func (m SignOffValidator) Validate(cr *git.CommitRange, g git.Git, whitelist []string) (*ValidationResult, error) {
	log.Debugf("Validating if commit messages are signed off from %s to %s", cr.From, cr.To)
	vr := ValidationResult{
		Validator: "Signed-Off Message",
		Skipped:   0,
	}
	defer func() {
		vr.Duration = utils.Duration(time.Now())
	}()
	var err error
	m.messages, err = g.GetCommitMessages(cr.From, cr.To)
	if err != nil {
		return nil, err
	}
	vs, err := m.isSignedOff()
	if err != nil {
		return nil, err
	}
	if utils.AllTrue(vs) {
		vr.Valid = true
		vr.Summary = "All commit messages are signed off\n"
		return &vr, nil
	} else {
		vr.Valid = false
		vr.Summary = "Not all commit messages are signed off\n"
	}
	vr.Messages = make(map[string]ResultMessage, len(m.messages))
	for commitHash := range m.messages {
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
		if !vs[commitHash] {
			vr.Messages[commitHash] = ResultMessage{Valid: false, Message: fmt.Sprintf("Commit Message: \"%s\" is not signed off\n", commitHash)}
			vr.Failures++
		} else {
			vr.Messages[commitHash] = ResultMessage{Valid: true, Message: fmt.Sprintf("Commit Message: \"%s\" is signed off\n", commitHash)}
		}
	}

	// Update summary with skip count
	if vr.Failures == 0 {
		vr.Valid = true
		if vr.Skipped > 0 {
			vr.Summary = fmt.Sprintf("All commits are signed off (%d skipped)\n", vr.Skipped)
		} else {
			vr.Summary = "All commit messages are signed off\n"
		}
		return &vr, nil
	}
	vr.Valid = false
	totalChecked := vr.Assertions + vr.Skipped
	passedCount := vr.Assertions - vr.Failures
	if vr.Skipped > 0 {
		vr.Summary = fmt.Sprintf("%d of %d commits are signed off (%d skipped)\n", passedCount, totalChecked, vr.Skipped)
	} else {
		vr.Summary = "Not all commit messages are signed off\n"
	}
	return &vr, nil
}

func (m SignOffValidator) isSignedOff() (map[string]bool, error) {
	var l = make(map[string]bool, len(m.messages))
	re, err := regexp.Compile(`(?m)^Signed-off-by: .+$`)
	if err != nil {
		return nil, fmt.Errorf("can't compile regexp: %w", err)
	}
	for i, cm := range m.messages {
		matched := re.Match([]byte(cm))
		l[i] = matched
	}
	return l, nil
}
