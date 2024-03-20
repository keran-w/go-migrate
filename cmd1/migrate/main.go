package main

import (
	"github.com/keran-w/go-migrate/docker"
	"log"
	"os/exec"
	"path/filepath"
	"time"
	"bytes"
	"strings"
)

func main() {

	containerName := "ml-container-A"
	container, err := docker.FindContainer(containerName)
	if err != nil {
		log.Fatalf("Error finding container %s: %v", containerName, err)
		return
	}

	log.Printf("Creating checkpoint1 for container %s...\n", containerName)
	// checkpointName := "checkpointA-" + time.Now().Format("MM-DDTHH-mm")
	checkpointName := "checkpointA-1"
	checkpointDir := "/home/ubuntu/go-migrate/checkpoints"

	startTime := time.Now()

	err = container.Checkpoint(checkpointName, checkpointDir, false)
	if err != nil {
		log.Fatalf("Error creating checkpoint for container %s: %v", containerName, err)
		return
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Printf("Time taken for checkpoint1: %v\n", duration)

	cmd := exec.Command("sudo", "chmod", "-R", "777", "./checkpoints")
	if cmd.Run() != nil {
		log.Fatalf("Error changing permissions for checkpoint directory: %v", err)
		return
	} else {
		//log.Printf("Permissions changed for checkpoint directory.\n")
	}
	


	time.Sleep(1 * time.Second)


	log.Printf("Stop old container %s...\n", containerName)

	log.Printf("Creating checkpoint2 for container %s...\n", containerName)
	// checkpointName := "checkpointA-" + time.Now().Format("MM-DDTHH-mm")
	checkpointName = "checkpointA-2"
	checkpointDir = "/home/ubuntu/go-migrate/checkpoints"

	startTime = time.Now()

	err = container.Checkpoint(checkpointName, checkpointDir, true)
	if err != nil {
		log.Fatalf("Error creating checkpoint for container %s: %v", containerName, err)
		return
	}

	endTime = time.Now()
	duration = endTime.Sub(startTime)
	log.Printf("Time taken for checkpoint2: %v\n", duration)

	cmd = exec.Command("sudo", "chmod", "-R", "777", "./checkpoints")
	if cmd.Run() != nil {
		log.Fatalf("Error changing permissions for checkpoint directory: %v", err)
		return
	} else {
		//log.Printf("Permissions changed for checkpoint directory.\n")
	}

	// varName := "CURR"
	// value := container.GetState(varName)
	// log.Printf("Container %s state %s: %s\n", containerName, varName, value)

	//netType := "tcp"
	//host := "localhost"
	//port := "9988"
	//client.ConnectToServer(netType, host, port)



	// Migrate
	imageName := "ml_app"
	containerName = "ml-container-B"
	env := []string{"START=0", "END=3000"}
	container, err = docker.NewContainer(containerName, imageName, env)
	if err != nil {
		log.Fatalf("Error creating container %s: %v", containerName, err)
		return
	}

	// TODO: communications
	// netType := "tcp"
	// host := "localhost"
	// port := "9988"
	// server.StartServer(netType, host, port)

	checkpointID := "checkpointA-1"
	checkpointDir = "/home/ubuntu/go-migrate/checkpoints"
	src := filepath.Join(checkpointDir, checkpointID)

	dst := filepath.Join("/var/lib/docker/containers", container.ID, "checkpoints", checkpointID)

	cmd = exec.Command("sudo", "cp", "-r", src+"/.", dst)
    err = cmd.Run()
    if err != nil {
        log.Fatalf("Error in transmitting checkpoints: %v", err)
        return
    }

	startTime = time.Now()

	newcheckpointID := "checkpointA-2"
	src = filepath.Join(checkpointDir, newcheckpointID)

	err = syncFolders(src, dst)
    if err != nil {
        log.Fatalf("Error syncing folders: %v", err)
        return
    }

	endTime = time.Now()
	duration = endTime.Sub(startTime)
	log.Printf("Time taken for sync: %v\n", duration)

	startTime = time.Now()

	err = container.Restore(checkpointID, dst)
	if err != nil {
		log.Fatalf("Error restoreing from checkpoint %s: %v", checkpointID, err)
	}
	container.Start()

	endTime = time.Now()
	duration = endTime.Sub(startTime)
	log.Printf("Time taken for resuming: %v\n", duration)
}

func syncFolders(src, dst string) error {
	// Step 1: Remove files that exist in dst but not in src
    err := removeExtraFiles(src, dst)
    if err != nil {
        return err
    }

    // Step 2: Check and copy files that do not exist in dst from src to dst 
    err = copyNonexistentFiles(src, dst)
    if err != nil {
        return err
    }
    
    // Step 3: Check and overwrite files in dst with different content from src
    err = copyDifferentFiles(src, dst)
    if err != nil {
        return err
    }

    return nil
}

func removeExtraFiles(src, dst string) error {
    diffCmd := exec.Command("sudo", "diff", "-rq", dst, src)
    var diffOutput, diffErr bytes.Buffer
    diffCmd.Stdout = &diffOutput
    diffCmd.Stderr = &diffErr
    err := diffCmd.Run()
    if err != nil && !bytes.Contains(diffOutput.Bytes(), []byte("differ")) {
        log.Printf("Error running diff command: %v\n", err)
        log.Println("Diff output:", diffOutput.String())
        log.Printf("Diff error: %s\n", diffErr.String())
        return err
    }

    diffFiles := make(map[string]bool)
    diffLines := bytes.Split(diffOutput.Bytes(), []byte("\n"))
    for _, line := range diffLines {
        if bytes.HasPrefix(line, []byte("Only in "+dst)) {
            file := string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("Only in "+dst+":"))))
            diffFiles[file] = true
        }
    }

    for file := range diffFiles {
        dstFilePath := filepath.Join(dst, file)
        rmCmd := exec.Command("sudo", "rm", "-rf", dstFilePath)
        rmOutput, err := rmCmd.CombinedOutput()
        if err != nil {
            log.Printf("Error removeExtraFiles %s: %v\n%s", file, err, rmOutput)
            return err
        }
    }

    return nil
}

func copyNonexistentFiles(src, dst string) error {
    diffCmd := exec.Command("sudo", "diff", "-rq", src, dst)
    var diffOutput, diffErr bytes.Buffer
    diffCmd.Stdout = &diffOutput
    diffCmd.Stderr = &diffErr
    err := diffCmd.Run()
    if err != nil && !bytes.Contains(diffOutput.Bytes(), []byte("differ")) {
        log.Printf("Error running diff command: %v\n", err)
        log.Println("Diff output:", diffOutput.String())
        log.Printf("Diff error: %s\n", diffErr.String())
        return err
    }

    diffFiles := make(map[string]bool)
    diffLines := bytes.Split(diffOutput.Bytes(), []byte("\n"))
    for _, line := range diffLines {
        if bytes.HasPrefix(line, []byte("Only in "+src)) {
            file := string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("Only in "+src+":"))))
            diffFiles[file] = true
        }
    }

    for file := range diffFiles {
		
        srcFilePath := filepath.Join(src, file)
        dstFilePath := filepath.Join(dst, file)
        cpCmd := exec.Command("sudo", "cp", srcFilePath, dstFilePath)
        cpOutput, err := cpCmd.CombinedOutput()
        if err != nil {
            log.Printf("Error copyNonexistentFiles %s: %v\n%s", file, err, cpOutput)
            return err
        }
    }

    return nil
}

func copyDifferentFiles(src, dst string) error {
    diffCmd := exec.Command("sudo", "diff", "-rq", src, dst)
    var diffOutput, diffErr bytes.Buffer
    diffCmd.Stdout = &diffOutput
    diffCmd.Stderr = &diffErr
    err := diffCmd.Run()
    if err != nil && !bytes.Contains(diffOutput.Bytes(), []byte("differ")) {
        log.Printf("Error running diff command: %v\n", err)
        log.Println("Diff output:", diffOutput.String())
        log.Printf("Diff error: %s\n", diffErr.String())
        return err
    }

    diffFiles := make(map[string]bool)
    diffLines := bytes.Split(diffOutput.Bytes(), []byte("\n"))
    for _, line := range diffLines {
        if bytes.HasPrefix(line, []byte("Files ")) {
            file := string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("Files "))))
            diffFiles[file] = true
        }
    }

    for file := range diffFiles {
		file := filepath.Base(file)
		file = strings.TrimSuffix(file, " differ")
		
        srcFilePath := filepath.Join(src, file)
        dstFilePath := filepath.Join(dst, file)
        cpCmd := exec.Command("sudo", "cp", srcFilePath, dstFilePath)
        cpOutput, err := cpCmd.CombinedOutput()
        if err != nil {
            log.Printf("Error copyDifferentFiles %s: %v\n%s", file, err, cpOutput)			
            return err
        }
    }

    return nil
}
