package drivers

import (
	"github.com/vanvanni/lampofs/errors"
	"io"
	"os"
	"path/filepath"
)

type LocalDriver struct {
	rootPath string
}

func NewLocalDriver(rootPath string) (*LocalDriver, error) {
	if err := os.MkdirAll(rootPath, 0755); err != nil {
		return nil, err
	}

	return &LocalDriver{
		rootPath: rootPath,
	}, nil
}

func (d *LocalDriver) Read(path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(d.rootPath, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, errors.ErrFileNotFound
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (d *LocalDriver) Write(path string, data []byte) error {
	fullPath := filepath.Join(d.rootPath, path)

	if _, err := os.Stat(fullPath); err == nil {
		return errors.ErrFileExists
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (d *LocalDriver) Put(path string, data []byte) error {
	fullPath := filepath.Join(d.rootPath, path)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (d *LocalDriver) Delete(path string) error {
	fullPath := filepath.Join(d.rootPath, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return errors.ErrFileNotFound
	}

	return os.Remove(fullPath)
}

func (d *LocalDriver) Update(path string, data []byte, prepend bool) error {
	fullPath := filepath.Join(d.rootPath, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// If file doesn't exist, create it with the new data
		return d.Put(path, data)
	}

	if prepend {
		return d.prependToFile(fullPath, data)
	}

	return d.appendToFile(fullPath, data)
}

func (d *LocalDriver) appendToFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (d *LocalDriver) prependToFile(path string, data []byte) error {
	existingData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return err
	}

	_, err = file.Write(existingData)
	return err
}
