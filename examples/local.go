package main

import (
	"fmt"
	"github.com/vanvanni/lampofs"
	"github.com/vanvanni/lampofs/drivers"
	"io"
	"log"
)

func main() {
	fmt.Println("\n> Local Driver Example")

	// Create the driver
	driver, err := drivers.NewLocalDriver("./data")
	if err != nil {
		log.Fatal(err)
	}

	// Create a Lampo instance
	lampo := lampofs.NewLampo(driver)

	// Write
	err = lampo.Write("example.txt", []byte("Hello, Local World!"))
	if err != nil {
		log.Fatal(err)
	}

	// Read
	reader, err := lampo.Read("example.txt")
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Read from file: %s\n", string(data))

	// Update (as append)
	err = lampo.Update("example.txt", []byte(" - Appended"), false)
	if err != nil {
		log.Fatal(err)
	}

	// Read as test
	reader, err = lampo.Read("example.txt")
	if err != nil {
		log.Fatal(err)
	}
	data, err = io.ReadAll(reader)
	reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Read updated file: %s\n", string(data))

	// Put (overwrite)
	err = lampo.Put("example.txt", []byte("New content"))
	if err != nil {
		log.Fatal(err)
	}

	// Read (overwrite)
	reader, err = lampo.Read("example.txt")
	if err != nil {
		log.Fatal(err)
	}
	data, err = io.ReadAll(reader)
	reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Read overwritten file: %s\n", string(data))

	// Delete
	err = lampo.Delete("example.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("File deleted successfully\n\n")
}
