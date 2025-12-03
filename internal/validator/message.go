package validator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aljoshare/commala/internal/git"
	"github.com/aljoshare/commala/internal/utils"
	log "github.com/sirupsen/logrus"
)

type MessageValidator struct {
	messages map[string]string
}

func (m MessageValidator) Validate(cr *git.CommitRange, g git.Git) (*ValidationResult, error) {
	log.Debugf("Validating if commit messages are conventional from %s to %s", cr.From, cr.To)
	vr := ValidationResult{
		Validator: "Conventional Message",
	}
	defer func() {
		vr.Duration = utils.Duration(time.Now())
	}()
	var err error
	m.messages, err = g.GetCommitMessages(cr.From, cr.To)
	if err != nil {
		return nil, err
	}
	vs, err := m.isConventional()
	if err != nil {
		return nil, err
	}
	if utils.AllTrue(vs) {
		vr.Valid = true
		vr.Summary = "All commit messages are conventional\n"
		return &vr, nil
	} else {
		vr.Valid = false
		vr.Summary = "Not all commit messages are conventional\n"
	}
	vr.Messages = make(map[string]ResultMessage, len(m.messages))
	for i := range m.messages {
		vr.Assertions++
		if !vs[i] {
			vr.Messages[i] = ResultMessage{false, fmt.Sprintf("Commit Message: \"%s\" is not conventional\n", i)}
			vr.Failures++
		} else {
			vr.Messages[i] = ResultMessage{true, fmt.Sprintf("Commit Message: \"%s\" is conventional\n", i)}
		}
	}
	return &vr, nil
}

func (m MessageValidator) isConventional() (map[string]bool, error) {
	var l = make(map[string]bool, len(m.messages))
	re := regexp.MustCompile(`^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test){1}(\([\w\-\.]+\))?(!)?: ([\w ])+([\s\S]*)`)
	for i, cm := range m.messages {
		matched := re.Match([]byte(cm))
		l[i] = matched
	}
	return l, nil
}
