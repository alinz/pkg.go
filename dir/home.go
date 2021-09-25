package dir

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

func Home() (string, error) {
	return homedir.Dir()
}

func Expand(path string) (string, error) {
	return homedir.Expand(path)
}

func Create(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
