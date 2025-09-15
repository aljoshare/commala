package utils

import (
	"io/ioutil"

	"github.com/charmbracelet/log"
)

func WriteFile(path string, content []byte) error {
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		log.Errorf("Can't write file %s: %v", path, err)
		return err
	}
	return nil
}
