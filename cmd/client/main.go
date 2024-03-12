package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/keran-w/go-migrate/docker"
)

func cmp(file1, file2 string) bool {
	// Read the contents of the first file
	data1, err := os.ReadFile(file1)
	if err != nil {
		fmt.Println("Error reading file:", file1, err)
		return false
	}

	// Read the contents of the second file
	data2, err := os.ReadFile(file2)
	if err != nil {
		fmt.Println("Error reading file:", file2, err)
		return false
	}

	// Compare the contents of the two files
	return bytes.Equal(data1, data2)
}

func getAllFilenames(dirPath string) (map[string]struct{}, error) {
	files := make(map[string]struct{})
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files[filepath.Base(path)] = struct{}{}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func sendFile(src string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(dst, data, 0777)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func signal() error {
	signalDst := "/home/ubuntu/test/server/signal.txt"
	signalSrc := "/home/ubuntu/test/client/signal.txt"
	err := sendFile(signalSrc, signalDst)
	return err
}

func main() {

	act := os.Args[1]
	switch act {
	// ./client migrate containerName srcDir dstDir
	case "migrate":
		containerName := os.Args[2]
		checkpointDir := os.Args[3]
		dstDir := os.Args[4]

		container, err := docker.FindContainer(containerName)
		if err != nil {
			log.Fatalf("Error finding container %s: %v", containerName, err)
			return
		}

		log.Printf("Creating checkpoint for container %s...\n", containerName)
		// checkpointName := "checkpointA-" + time.Now().Format("MM-DDTHH-mm")
		maxIter := 2
		for i := 0; i < maxIter; i++ {
			checkpointName := containerName + strconv.Itoa(i)
			err = container.Checkpoint(checkpointName, checkpointDir, false)
			if err != nil {
				log.Fatalf("Error creating checkpoint for container %s: %v", containerName, err)
				return
			}
			curDir := checkpointDir + "/" + checkpointName
			curFilenames, err := getAllFilenames(curDir)
			if err != nil {
				log.Fatalf("Error fetching the checkpoint files %s: %v", curDir, err)
				return
			}
			if i == 0 {
				for filename := range curFilenames {
					dstFile := filepath.Join(dstDir, filename)
					srcFile := filepath.Join(curDir, filename)
					err := sendFile(srcFile, dstFile)
					if err != nil {
						return
					}
				}
				continue
			}
			prevDir := checkpointDir + "/" + containerName + strconv.Itoa(i-1)
			prevFilenames, err := getAllFilenames(prevDir)
			if err != nil {
				log.Fatalf("Error fetching the checkpoint files %s: %v", prevDir, err)
				return
			}

			for filename := range curFilenames {
				if _, ok := prevFilenames[filename]; ok {
					//do something here
					if cmp(prevDir+"/"+filename, curDir+"/"+filename) {
						continue
					} else {
						fmt.Println("Sending Different Files:", filename)
						dstFile := filepath.Join(dstDir, filename)
						srcFile := filepath.Join(curDir, filename)
						err := sendFile(srcFile, dstFile)
						if err != nil {
							return
						}
					}
				} else {
					fmt.Println("Sending New Files:", filename)
					dstFile := filepath.Join(dstDir, filename)
					srcFile := filepath.Join(curDir, filename)
					err := sendFile(srcFile, dstFile)
					if err != nil {
						return
					}
				}
			}

		}
		err = signal()
		if err != nil {
			log.Fatalf("Error %v", err)
		}
		container.Stop()
		fmt.Println("Finished")
	}

	// cmd := exec.Command("sudo", "chmod", "-R", "777", "./checkpoints")
	// if cmd.Run() != nil {
	// 	log.Fatalf("Error changing permissions for checkpoint directory: %v", err)
	// 	return
	// } else {
	// 	log.Printf("Permissions changed for checkpoint directory.\n")
	// }

	// varName := "CURR"
	// value := container.GetState(varName)
	// log.Printf("Container %s state %s: %s\n", containerName, varName, value)

	//netType := "tcp"
	//host := "localhost"
	//port := "9988"
	//client.ConnectToServer(netType, host, port)
}
