package drivers

import (
	"bytes"
	"github.com/vanvanni/lampofs/errors"
	"io"
	"sync"
	"time"
)

type MemoryDriver struct {
	files map[string]*memoryFile
	mutex sync.RWMutex
}

type memoryFile struct {
	data      []byte
	createdAt time.Time
	updatedAt time.Time
}

func NewMemoryDriver() *MemoryDriver {
	return &MemoryDriver{
		files: make(map[string]*memoryFile),
	}
}

func (d *MemoryDriver) Read(path string) (io.ReadCloser, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	file, exists := d.files[path]
	if !exists {
		return nil, errors.ErrFileNotFound
	}

	dataCopy := make([]byte, len(file.data))
	copy(dataCopy, file.data)

	return io.NopCloser(bytes.NewReader(dataCopy)), nil
}

func (d *MemoryDriver) Write(path string, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, exists := d.files[path]; exists {
		return errors.ErrFileExists
	}

	now := time.Now()
	d.files[path] = &memoryFile{
		data:      data,
		createdAt: now,
		updatedAt: now,
	}

	return nil
}

func (d *MemoryDriver) Put(path string, data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	now := time.Now()

	if file, exists := d.files[path]; exists {
		file.data = data
		file.updatedAt = now
	} else {
		d.files[path] = &memoryFile{
			data:      data,
			createdAt: now,
			updatedAt: now,
		}
	}

	return nil
}

func (d *MemoryDriver) Delete(path string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, exists := d.files[path]; !exists {
		return errors.ErrFileNotFound
	}

	delete(d.files, path)
	return nil
}

func (d *MemoryDriver) Update(path string, data []byte, prepend bool) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	file, exists := d.files[path]
	if !exists {
		now := time.Now()
		d.files[path] = &memoryFile{
			data:      data,
			createdAt: now,
			updatedAt: now,
		}
		return nil
	}

	if prepend {
		newData := make([]byte, len(data)+len(file.data))
		copy(newData, data)
		copy(newData[len(data):], file.data)
		file.data = newData
	} else {
		file.data = append(file.data, data...)
	}

	file.updatedAt = time.Now()
	return nil
}
