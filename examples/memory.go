package main

import (
	"fmt"
	"github.com/vanvanni/lampofs"
	"github.com/vanvanni/lampofs/drivers"
	"io"
	"log"
)

func main() {
	fmt.Println("\n>Memory Driver Example")

	// Create the driver
	driver := drivers.NewMemoryDriver()

	// Create a Lampo instance
	lampo := lampofs.NewLampo(driver)

	// Write
	err := lampo.Write("memory-example.txt", []byte("Hello, Memory World!"))
	if err != nil {
		log.Fatal(err)
	}

	// Read
	reader, err := lampo.Read("memory-example.txt")
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Read from memory: %s\n", string(data))

	// Update (as prepend)
	err = lampo.Update("memory-example.txt", []byte("Prepended - "), true)
	if err != nil {
		log.Fatal(err)
	}

	// Read as test
	reader, err = lampo.Read("memory-example.txt")
	if err != nil {
		log.Fatal(err)
	}
	data, err = io.ReadAll(reader)
	reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Read updated memory file: %s\n", string(data))

	// Delete
	err = lampo.Delete("memory-example.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Memory file deleted successfully")
}
