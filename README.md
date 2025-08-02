# LampoFS - Go Filesystem Library

LampoFS is a Go library that makes working with different storage systems easy and simple. It gives you one simple way to handle files whether they're on your local computer, in cloud storage like S3, or just in memory. No need to learn different APIs for each storage type!

## Simple Reference

### Lampo Methods

- `Read(path string) (io.ReadCloser, error)` - Read a file
- `Write(path string, data []byte) error` - Write a new file (fails if file exists)
- `Put(path string, data []byte) error` - Create or overwrite a file
- `Delete(path string) error` - Delete a file
- `Update(path string, data []byte, prepend bool) error` - Append or prepend to a file
- `On(handler func(event LampEvent))` - Register an event listener

### LampEvent

Events fired by the filesystem operations:

- `Type`: "READ", "WRITE", "PUT", "DELETE", "APPEND", "PREPEND"
- `Path`: Path of the file
- `Timestamp`: Unix timestamp of the event
- `Data`: Additional data (size of data for write operations)

## Testing
Run tests with:

```bash
# Run tests with gotestsum (beautiful output)
make test
```

## License
MIT
