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

func (m SignOffValidator) Validate(cr *git.CommitRange, g git.Git) (*ValidationResult, error) {
	log.Debugf("Validating if commit messages are signed off from %s to %s", cr.From, cr.To)
	vr := ValidationResult{
		Validator: "Signed-Off Message",
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
	for i := range m.messages {
		vr.Assertions++
		if !vs[i] {
			vr.Messages[i] = ResultMessage{false, fmt.Sprintf("Commit Message: \"%s\" is not signed off\n", i)}
			vr.Failures++
		} else {
			vr.Messages[i] = ResultMessage{true, fmt.Sprintf("Commit Message: \"%s\" is signed off\n", i)}
		}
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
