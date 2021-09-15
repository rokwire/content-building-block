package tempstorage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const dir = "./" //we need to set to current dir because of the bin wrapper!

// Adapter struct
type Adapter struct{}

// NewTempStorageAdapter cerates new instance
func NewTempStorageAdapter() *Adapter {
	return &Adapter{}
}

// Save saves a temp file
func (a *Adapter) Save(fileName string, fileType string, fileContent []byte) error {
	log.Printf("Save %s", fileName)

	//create the folder if it does not exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0644)
	}

	newPath := filepath.Join(dir, fileName)

	// write file
	newFile, err := os.Create(newPath)
	if err != nil {
		return errors.New("CANT_WRITE_FILE")
	}
	defer newFile.Close() // idempotent, okay to call twice
	if _, err := newFile.Write(fileContent); err != nil || newFile.Close() != nil {
		return err
	}
	return nil
}

// Delete deletes a file
func (a *Adapter) Delete(fileName string) error {
	log.Printf("Delete %s", fileName)

	path := filepath.Join(dir, fileName)
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("Cannt remove file %s: %s", path, err)
	}
	return nil
}

// Read reads a file
func (a *Adapter) Read(fileName string) (*os.File, error) {
	log.Printf("Read %s", fileName)

	file, err := os.Open(fileName)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return file, nil
}
