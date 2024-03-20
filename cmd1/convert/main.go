package main

import (
	"fmt"
	"io/ioutil"
)

// FileInfo represents basic metadata about a file in the checkpoint directory.
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"` // Size in bytes
	Type string `json:"type"`
}

func main() {
	dirPath := "checkpoints/checkpointA-1"
	fileName := "pages-1.img"
	filePath := fmt.Sprintf("%s/%s", dirPath, fileName)
	// read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}
	// print the file data
	fmt.Printf("File %s data: %s\n", filePath, string(data))

}
