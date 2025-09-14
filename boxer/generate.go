//go:build generate
// +build generate

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mholt/archiver/v4"
)

func deleteIfExists(filename string) {
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return
	} else if err != nil {
		log.Fatalf("failed to check if %s exists: %w", filename, err)
	}
}

func boxRepositoryFiles() []archiver.File {
	// Read the box code repository folder from disk
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		"../box": "box",
	})
	if err != nil {
		panic(fmt.Errorf("failed to read box repository folder from disk: %w", err))
	}
	return files
}

func createGzipFromFiles(files []archiver.File, archiveFilename string) {
	// Create the archive file
	out, err := os.Create(archiveFilename)
	if err != nil {
		panic(fmt.Errorf("failed to create archive of box repository folder: %w", err))
	}
	defer out.Close()

	// Compress/archive into the file
	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}
	err = format.Archive(context.Background(), out, files)
	if err != nil {
		panic(fmt.Errorf("failed to compress/archive box repository folder to gzip: %w", err))
	}
}

func createArchiveOfBoxRepositoryFolder(archiveFilename, boxRepositoryFolder string) {
	files := boxRepositoryFiles()
	createGzipFromFiles(files, archiveFilename)
	fmt.Printf("Successfully created archive file: %s\n", archiveFilename)
}

func goGenerateOnBoxRepository(boxRepositoryFolder string) {
	cmd := exec.Command("go", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = boxRepositoryFolder
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to run go generate on box repository folder: %v", err)
	}
}

func main() {
	boxRepositoryFolder := "../box"
	archiveFilename := "embedded_box.tar.gz"
	deleteIfExists(archiveFilename)
	createArchiveOfBoxRepositoryFolder(archiveFilename, boxRepositoryFolder)
}
