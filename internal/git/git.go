package git

import (
	"fmt"
	"strconv"
	"strings"

	gogit "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	log "github.com/sirupsen/logrus"
)

type CommitRange struct {
	From string
	To   string
}

type Git interface {
	GetBranchName() (string, error)
	GetCommitMessages(from string, to string) (map[string]string, error)
	GetCommitAuthorNames(from string, to string) (map[string]string, error)
	GetCommitAuthorEmails(from string, to string) (map[string]string, error)
	GetCommitHash(hash string) (*object.Commit, error)
	GetInitialCommit() (string, error)
	GetLatestCommit() (string, error)
	ParseCommitRange(commitrange string) (*CommitRange, error)
}

type RealGit struct{}

func (r RealGit) getRepo() (*gogit.Repository, error) {
	repo, err := gogit.PlainOpen(".")
	if err != nil {
		return nil, fmt.Errorf("can't get git repo")
	}
	return repo, nil
}

func (r RealGit) getHead() (*plumbing.Reference, error) {
	repo, err := r.getRepo()
	if err != nil {
		return nil, err
	}
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("can't get git head")
	}
	return ref, nil
}

func (r RealGit) GetBranchName() (string, error) {
	h, err := r.getHead()
	if err != nil {
		return "", fmt.Errorf("can't get branch name")
	}
	bn := h.Name().Short()
	return bn, nil
}

func (r RealGit) GetCommitMessages(from string, to string) (map[string]string, error) {
	var m = make(map[string]string)
	repo, err := r.getRepo()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	fromHash, err := r.GetCommitHash(from)
	if err != nil {
		return nil, err
	}
	toHash, err := r.GetCommitHash(to)
	if err != nil {
		return nil, err
	}
	commitIter, err := repo.Log(&gogit.LogOptions{From: fromHash.Hash, To: toHash.Hash})
	if err != nil {
		return nil, fmt.Errorf("can't get git log")
	}

	err = commitIter.ForEach(func(commit *object.Commit) error {
		m[commit.Hash.String()] = commit.Message
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("can't iterate git log")
	}
	return m, nil
}

func (r RealGit) GetCommitAuthorNames(from string, to string) (map[string]string, error) {
	var m = make(map[string]string)
	repo, err := r.getRepo()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	fromHash, err := r.GetCommitHash(from)
	if err != nil {
		return nil, err
	}
	toHash, err := r.GetCommitHash(to)
	if err != nil {
		return nil, err
	}
	commitIter, err := repo.Log(&gogit.LogOptions{From: fromHash.Hash, To: toHash.Hash})
	if err != nil {
		return nil, fmt.Errorf("can't get git log")
	}

	err = commitIter.ForEach(func(commit *object.Commit) error {
		m[commit.Hash.String()] = commit.Author.Name
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("can't iterate git log")
	}
	return m, nil
}

func (r RealGit) GetCommitAuthorEmails(from string, to string) (map[string]string, error) {
	var m = make(map[string]string)
	repo, err := r.getRepo()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	fromHash, err := r.GetCommitHash(from)
	if err != nil {
		return nil, err
	}
	toHash, err := r.GetCommitHash(to)
	if err != nil {
		return nil, err
	}
	commitIter, err := repo.Log(&gogit.LogOptions{From: fromHash.Hash, To: toHash.Hash})
	if err != nil {
		return nil, fmt.Errorf("can't get git log")
	}

	err = commitIter.ForEach(func(commit *object.Commit) error {
		m[commit.Hash.String()] = commit.Author.Email
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("can't iterate git log")
	}
	return m, nil
}

func (r RealGit) GetCommitHash(hash string) (*object.Commit, error) {
	repo, err := r.getRepo()
	if err != nil {
		return nil, err
	}
	commitHash := plumbing.NewHash(hash)
	commit, err := repo.CommitObject(commitHash)
	if err != nil {
		return nil, fmt.Errorf("can't find commit %s", hash)
	}
	return commit, nil
}

func (r RealGit) GetInitialCommit() (string, error) {
	repo, err := r.getRepo()
	if err != nil {
		return "", err
	}
	commitIter, err := repo.Log(&gogit.LogOptions{Order: gogit.LogOrderCommitterTime})
	if err != nil {
		return "", fmt.Errorf("can't get git log")
	}

	var initialCommit string
	commitIter.ForEach(func(c *object.Commit) error {
		initialCommit = c.Hash.String()
		return nil
	})
	if initialCommit == "" {
		log.Errorf("Can't find initial commit")
		return "", fmt.Errorf("can't find initial commit")
	} else {
		log.Warnf("No initial commit defined! Starting with %s", initialCommit)
		return initialCommit, nil
	}
}

func (r RealGit) GetLatestCommit() (string, error) {
	head, err := r.getHead()
	if err != nil {
		return "", err
	}
	latestCommit := head.Hash().String()
	log.Debugf("No latest commit defined! Ending with %s", latestCommit)
	return latestCommit, nil
}

func (r RealGit) ParseCommitRange(commitrange string) (*CommitRange, error) {
	cr := CommitRange{}
	var err error
	if strings.Compare(commitrange, "..") == 0 {
		cr.To, err = r.GetInitialCommit()
		if err != nil {
			return nil, fmt.Errorf("can't find initial commit")
		}
		cr.From, err = r.GetLatestCommit()
		if err != nil {
			return nil, fmt.Errorf("can't find latest commit")
		}
		return &cr, nil
	}
	if strings.HasPrefix(commitrange, "..") || !strings.Contains(commitrange, "..") {
		cr.To, err = r.GetInitialCommit()
		if err != nil {
			return nil, fmt.Errorf("can't find initial commit")
		}
		cr.To = strings.Replace(commitrange, "..", "", 1)
		return &cr, nil
	}
	if strings.HasSuffix(commitrange, "..") {
		cr.From, err = r.GetLatestCommit()
		cr.To = strings.Replace(commitrange, "..", "", 1)
		if err != nil {
			return nil, fmt.Errorf("can't find latest commit")
		}
		return &cr, nil
	}
	s := strings.Split(commitrange, "..")
	cr.To = s[0]
	cr.From = s[1]
	return &cr, nil
}

func (r RealGit) ParseNegativeIndex(commitrange string) (*CommitRange, error) {
	cr := CommitRange{}
	if strings.HasPrefix(commitrange, "HEAD~") {
		negindex := strings.Replace(commitrange, "HEAD~", "", 1)
		n, err := strconv.Atoi(negindex)
		if err != nil {
			return nil, fmt.Errorf("can't parse negative index %s", negindex)
		}
		head, err := r.getHead()
		if err != nil {
			return nil, fmt.Errorf("can't get git head")
		}
		fromCommit := head.Hash()
		repo, err := r.getRepo()
		if err != nil {
			return nil, err
		}
		commitIter, err := repo.Log(&gogit.LogOptions{From: fromCommit, Order: gogit.LogOrderCommitterTime})
		if err != nil {
			return nil, fmt.Errorf("can't get git log")
		}
		i := 0
		var toCommit *object.Commit
		err = commitIter.ForEach(func(commit *object.Commit) error {
			if i == n {
				toCommit = commit
				return fmt.Errorf("stop iteration")
			}
			i++
			return nil
		})
		if err != nil && err.Error() != "stop iteration" {
			return nil, fmt.Errorf("can't iterate git log")
		}
		if toCommit == nil {
			return nil, fmt.Errorf("can't find commit at negative index %d", n)
		}
		cr.From = fromCommit.String()
		cr.To = toCommit.Hash.String()
		log.Debugf("Parsed negative index %s to commit range %s..%s", commitrange, cr.From, cr.To)
		return &cr, nil
	} else {
		return nil, fmt.Errorf("can't parse negative index from %s", commitrange)
	}
}
