// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
