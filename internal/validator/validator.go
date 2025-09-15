package validator

import (
	"sort"

	"github.com/aljoshare/commala/internal/config"
	"github.com/aljoshare/commala/internal/git"
)

func Validate(cr *git.CommitRange, g git.Git, c config.Config) ([]*ValidationResult, error) {
	var results []*ValidationResult

	validators := []struct {
		enabled  bool
		validate func() (*ValidationResult, error)
	}{
		{c.BranchEnabled, func() (*ValidationResult, error) {
			return BranchValidator{}.Validate(g)
		}},
		{c.MessageEnabled, func() (*ValidationResult, error) {
			return MessageValidator{}.Validate(cr, g)
		}},
		{c.SignOffEnabled, func() (*ValidationResult, error) {
			return SignOffValidator{}.Validate(cr, g)
		}},
		{c.AuthorNameEnabled, func() (*ValidationResult, error) {
			return AuthorNameValidator{}.Validate(cr, g)
		}},
		{c.AuthorEmailEnabled, func() (*ValidationResult, error) {
			return AuthorEmailValidator{}.Validate(cr, g)
		}},
	}

	type result struct {
		vr  *ValidationResult
		err error
	}
	ch := make(chan result, len(validators))
	var tasks int

	for _, v := range validators {
		if v.enabled {
			tasks++
			go func(validate func() (*ValidationResult, error)) {
				vr, err := validate()
				ch <- result{vr, err}
			}(v.validate)
		}
	}

	for i := 0; i < tasks; i++ {
		r := <-ch
		if r.err != nil {
			return nil, r.err
		}
		results = append(results, r.vr)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Validator < results[j].Validator
	})
	return results, nil
}
