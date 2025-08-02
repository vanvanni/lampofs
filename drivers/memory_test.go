package drivers

import (
	"github.com/vanvanni/lampofs/errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryDriver(t *testing.T) {
	driver := NewMemoryDriver()
	assert.NotNil(t, driver)

	// Test Write
	testData := []byte("Hello, Memory!")
	err := driver.Write("test.txt", testData)
	assert.NoError(t, err)

	// Test Read
	reader, err := driver.Read("test.txt")
	assert.NoError(t, err)
	assert.NotNil(t, reader)

	data, err := io.ReadAll(reader)
	reader.Close()
	assert.NoError(t, err)
	assert.Equal(t, testData, data)

	// Test Put (overwrite)
	newData := []byte("New data")
	err = driver.Put("test.txt", newData)
	assert.NoError(t, err)

	reader, err = driver.Read("test.txt")
	assert.NoError(t, err)
	assert.NotNil(t, reader)

	data, err = io.ReadAll(reader)
	reader.Close()
	assert.NoError(t, err)
	assert.Equal(t, newData, data)

	// Test Update (append)
	appendData := []byte(" Appended")
	err = driver.Update("test.txt", appendData, false)
	assert.NoError(t, err)

	reader, err = driver.Read("test.txt")
	assert.NoError(t, err)
	assert.NotNil(t, reader)

	data, err = io.ReadAll(reader)
	reader.Close()
	assert.NoError(t, err)

	expected := append(newData, appendData...)
	assert.Equal(t, expected, data)

	// Test Update (prepend)
	prependData := []byte("Prepended ")
	err = driver.Update("test.txt", prependData, true)
	assert.NoError(t, err)

	reader, err = driver.Read("test.txt")
	assert.NoError(t, err)
	assert.NotNil(t, reader)

	data, err = io.ReadAll(reader)
	reader.Close()
	assert.NoError(t, err)

	expected = append(prependData, expected...)
	assert.Equal(t, expected, data)

	// Test Delete
	err = driver.Delete("test.txt")
	assert.NoError(t, err)

	// Verify file
	_, err = driver.Read("test.txt")
	assert.Error(t, err)
	assert.Equal(t, errors.ErrFileNotFound, err)
}

func TestMemoryDriverErrors(t *testing.T) {
	driver := NewMemoryDriver()
	assert.NotNil(t, driver)

	// Reading non-existent file
	_, err := driver.Read("nonexistent.txt")
	assert.Error(t, err)
	assert.Equal(t, errors.ErrFileNotFound, err)

	// Writing to existing file (should fail)
	testData := []byte("test")
	err = driver.Write("existing.txt", testData)
	assert.NoError(t, err)

	err = driver.Write("existing.txt", testData)
	assert.Error(t, err)
	assert.Equal(t, errors.ErrFileExists, err)

	// Deleting non-existent file
	err = driver.Delete("nonexistent.txt")
	assert.Error(t, err)
	assert.Equal(t, errors.ErrFileNotFound, err)
}
