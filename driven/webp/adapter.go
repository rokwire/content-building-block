package webp

import (
	"content/core/model"
	"fmt"
	"log"
	"os/exec"
)

// Adapter struct
type Adapter struct{}

// NewWebpAdapter cerates new instance
func NewWebpAdapter() *Adapter {
	return &Adapter{}
}

// Convert converts an image
func (a *Adapter) Convert(inputFileName string, outputFileName string, spec model.ImageSpec) error {
	log.Println("Convert")

	args := []string{}
	var cmd *exec.Cmd

	if spec.Height > 0 || spec.Width > 0 {
		args = append(args, "-resize", fmt.Sprintf("%d", spec.Width), fmt.Sprintf("%d", spec.Height))
	}
	if spec.Quality > 0 {
		args = append(args, "-q", fmt.Sprintf("%d", spec.Quality))
	}
	args = append(args, inputFileName, "-o", outputFileName)
	cmd = exec.Command("cwebp", args...)

	err := cmd.Run()
	if err != nil {
		log.Printf("Error on convertion:%s", err.Error())
	}

	return err
}
