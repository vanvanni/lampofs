package drivers

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3Driver(t *testing.T) {
	driver, err := NewS3Driver(S3Options{
		Region:          "us-east-1",
		BucketName:      "lampofs-test",
		Endpoint:        "http://localhost:9010",
		AccessKeyID:     "local",
		SecretAccessKey: "password123",
		Timeout:         5 * time.Second,
	})

	require.NoError(t, err)
	require.NotNil(t, driver)

	// Test Write
	testData := []byte("Hello, S3 driver test!")
	err = driver.Write("test-file.txt", testData)
	require.NoError(t, err)

	// Test Read
	reader, err := driver.Read("test-file.txt")
	require.NoError(t, err)
	defer reader.Close()

	readData, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, testData, readData)

	// Test Update (append)
	updateData := []byte(" Updated content.")
	err = driver.Update("test-file.txt", updateData, false)
	require.NoError(t, err)

	// Verify update
	reader, err = driver.Read("test-file.txt")
	require.NoError(t, err)
	defer reader.Close()

	expectedData := append(testData, updateData...)
	readData, err = io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, expectedData, readData)

	// Test Delete
	err = driver.Delete("test-file.txt")
	require.NoError(t, err)

	// Verify file
	_, err = driver.Read("test-file.txt")
	assert.Error(t, err)
}

func TestS3DriverConnectionFailure(t *testing.T) {
	_, err := NewS3Driver(S3Options{
		Region:          "us-east-1",
		BucketName:      "non-existent-bucket",
		Endpoint:        "http://non-existent-endpoint:9000",
		AccessKeyID:     "invalid",
		SecretAccessKey: "invalid",
		Timeout:         2 * time.Second,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection test failed")
}
