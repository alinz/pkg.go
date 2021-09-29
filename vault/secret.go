package vault

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/alinz/crypto.go"
	"github.com/alinz/hash.go"

	"github.com/alinz/pkg.go/generator"
)

var (
	ErrWrongKey = errors.New("wrong key")
)

type Store interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

type FileStore struct {
	path string
	key  []byte
}

func NewFileStore(path string, passphrase string) *FileStore {
	key := hash.Bytes([]byte(passphrase))
	return &FileStore{
		key:  key,
		path: path,
	}
}

func NewAutoFileStore(path string, passwordTTL time.Duration) *FileStore {
	passphrase := generator.RandString(10, passwordTTL)
	return NewFileStore(path, passphrase)
}

func (f *FileStore) Set(key string, value string) error {
	ciphertext, err := crypto.NewChaCha20().Encrypt([]byte(value), f.key)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(f.path, key))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(ciphertext)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStore) Get(key string) (string, error) {
	file, err := os.Open(filepath.Join(f.path, key))
	if err != nil {
		return "", err
	}
	defer file.Close()

	ciphertext, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	plaintext, err := crypto.NewChaCha20().Decrypt(ciphertext, f.key)
	if errors.Is(err, crypto.ErrWrongKey) {
		return "", ErrWrongKey
	} else if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
