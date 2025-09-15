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
	messages []string
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
	if utils.AllTrueSlice(vs) {
		vr.Valid = true
		vr.Message = append(vr.Message, "All commit messages are signed off\n")
		return &vr, nil
	}
	for _, v := range vs {
		if !v {
			for i, m := range m.messages {
				if !vs[i] {
					vr.Message = append(vr.Message, fmt.Sprintf("Commit Message: \"%s\" is not signed off\n", m))
				}
			}
		}
	}
	return &vr, nil
}

func (m SignOffValidator) isSignedOff() ([]bool, error) {
	var l []bool
	l = make([]bool, 0, len(m.messages))
	re, err := regexp.Compile(`(?m)^Signed-off-by: .+$`)
	if err != nil {
		return nil, fmt.Errorf("can't compile regexp: %w", err)
	}
	for _, cm := range m.messages {
		matched := re.Match([]byte(cm))
		l = append(l, matched)
	}
	return l, nil
}
