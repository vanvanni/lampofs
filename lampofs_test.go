package lampofs

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockDriver struct {
	readFunc   func(path string) (io.ReadCloser, error)
	writeFunc  func(path string, data []byte) error
	putFunc    func(path string, data []byte) error
	deleteFunc func(path string) error
	updateFunc func(path string, data []byte, prepend bool) error
}

func (m *mockDriver) Read(path string) (io.ReadCloser, error) {
	if m.readFunc != nil {
		return m.readFunc(path)
	}
	return nil, nil
}

func (m *mockDriver) Write(path string, data []byte) error {
	if m.writeFunc != nil {
		return m.writeFunc(path, data)
	}
	return nil
}

func (m *mockDriver) Put(path string, data []byte) error {
	if m.putFunc != nil {
		return m.putFunc(path, data)
	}
	return nil
}

func (m *mockDriver) Delete(path string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(path)
	}
	return nil
}

func (m *mockDriver) Update(path string, data []byte, prepend bool) error {
	if m.updateFunc != nil {
		return m.updateFunc(path, data, prepend)
	}
	return nil
}

func TestNewLampo(t *testing.T) {
	driver := &mockDriver{}
	lampo := NewLampo(driver)

	assert.NotNil(t, lampo)
	assert.Equal(t, driver, lampo.driver)
}

func TestLampoRead(t *testing.T) {
	expectedData := "test data"
	driver := &mockDriver{
		readFunc: func(path string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewBufferString(expectedData)), nil
		},
	}

	lampo := NewLampo(driver)
	reader, err := lampo.Read("test.txt")

	assert.NoError(t, err)
	assert.NotNil(t, reader)

	data, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, string(data))
}

func TestLampoWrite(t *testing.T) {
	driver := &mockDriver{
		writeFunc: func(path string, data []byte) error {
			return nil
		},
	}

	lampo := NewLampo(driver)
	testData := []byte("test data")
	err := lampo.Write("test.txt", testData)

	assert.NoError(t, err)
}

func TestLampoPut(t *testing.T) {
	driver := &mockDriver{
		putFunc: func(path string, data []byte) error {
			return nil
		},
	}

	lampo := NewLampo(driver)
	testData := []byte("test data")
	err := lampo.Put("test.txt", testData)

	assert.NoError(t, err)
}

func TestLampoDelete(t *testing.T) {
	driver := &mockDriver{
		deleteFunc: func(path string) error {
			return nil
		},
	}

	lampo := NewLampo(driver)
	err := lampo.Delete("test.txt")

	assert.NoError(t, err)
}

func TestLampoUpdate(t *testing.T) {
	driver := &mockDriver{
		updateFunc: func(path string, data []byte, prepend bool) error {
			return nil
		},
	}

	lampo := NewLampo(driver)
	testData := []byte("test data")
	err := lampo.Update("test.txt", testData, true)

	assert.NoError(t, err)
}

func TestLampoOn(t *testing.T) {
	driver := &mockDriver{}
	lampo := NewLampo(driver)

	eventReceived := false
	lampo.On(func(event LampEvent) {
		eventReceived = true
		assert.Equal(t, "WRITE", event.Type)
		assert.Equal(t, "test.txt", event.Path)
	})

	driver.writeFunc = func(path string, data []byte) error {
		return nil
	}

	err := lampo.Write("test.txt", []byte("test"))
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	assert.True(t, eventReceived)
}
