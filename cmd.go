package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

type RestoreFlag string

const (
	RestoreCopyFlag    RestoreFlag = "setup"
	RestoreReleaseFlag RestoreFlag = "release"
)

func listFilesInDirectory(path string) (map[string]struct{}, error) {
	out := make(map[string]struct{})
	err := filepath.Walk(path, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//ignore root directory
		if filePath == path {
			return nil
		}
		if IsDir(filePath) {
			return nil
		}

		//remove root path
		filePath = strings.Replace(filePath, path, "", 1)

		//ignore files
		for i := range config.Ignore {
			if strings.HasPrefix(filePath, config.Ignore[i]) {
				return nil
			}
		}
		out[filePath] = struct{}{}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
func move(sourceFiles, setupFiles map[string]struct{}) error {
	// todo: implement
	// 1. filter for candidates
	movedFiles, _ := filterForCandidates(setupFiles)

	// 2. move setup files to source
	for i := range movedFiles {
		sourceFile := config.Source + movedFiles[i]
		setupFile := config.Setup + movedFiles[i]
		if _, ok := sourceFiles[movedFiles[i]]; ok {
			// 3. if source is exists, move to backup
			//move source to backup
			backupFile := config.Backup + movedFiles[i]
			err := DeepCopyFile(sourceFile, backupFile)
			if err != nil {
				log.Printf("copy backup file failed, err: %v", err)
				return err
			}
			log.Printf("backup file: %v", movedFiles[i])
		}
		//move setup files to source
		err := DeepCopyFile(setupFile, sourceFile)
		if err != nil {
			log.Printf("copy setup file failed, err: %v", err)
			return err
		}
		log.Printf("setup file: %v", movedFiles[i])
	}

	return nil
}
func release(sourceFiles, backupFiles map[string]struct{}) error {
	//filter for candidates
	moveFiles, removeFiles := filterForCandidates(sourceFiles)

	//1.remove files
	for i := range removeFiles {
		//remove
		err := os.Remove(config.Source + removeFiles[i])
		if err != nil {
			log.Printf("Remove file failed, file: %v, err: %v", removeFiles[i], err)
			return err
		}
		log.Printf("removed file: %v", removeFiles[i])
	}

	//2.move files to setup
	for i := range moveFiles {
		sourceFile := config.Source + moveFiles[i]
		setupFile := config.Setup + moveFiles[i]
		err := DeepCopyFile(sourceFile, setupFile)
		if err != nil {
			log.Printf("copy source file failed, err: %v", err)
			return err
		}

		//remove sourceFile
		os.RemoveAll(sourceFile)
		log.Printf("release file: %v", moveFiles[i])
	}

	//3.restore backups
	for path := range backupFiles {
		sourceFile := config.Source + path
		backupFile := config.Backup + path
		err := DeepCopyFile(backupFile, sourceFile)
		if err != nil {
			log.Printf("copy backup file failed, path: %v, source: %v, target: %v, err: %v", path, backupFile, sourceFile, err)
			return err
		}
		log.Printf("restore file: %v", path)
	}

	return nil
}
func backup(path, target string) error {
	// cleanBackup()
	return copy.Copy(path, target)
}

func restore(restorePath, sourcePath string) {
	err := os.RemoveAll(sourcePath)
	if err != nil {
		log.Printf("Remove sourcePath file failed, file: %v, err: %v", sourcePath, err)
		return
	}
	err = backup(restorePath, sourcePath)
	if err != nil {
		log.Printf("move files failed, err: %v", err)
	}

	err = os.RemoveAll(restorePath)
	if err != nil {
		log.Printf("Remove restorePath file failed, file: %v, err: %v", restorePath, err)
		return
	}
}
func cleanDirectory(path string) error {
	return os.RemoveAll(path)
}

func filterForCandidates(source map[string]struct{}) ([]string, []string) {
	//find moved & removed
	moveFiles, removeFiles := make([]string, 0), make([]string, 0)
	for filePath := range source {
		for i := range moveReg {
			if moveReg[i].MatchString(filePath) {
				moveFiles = append(moveFiles, filePath)
				break
			}
		}
		for i := range removeReg {
			if removeReg[i].MatchString(filePath) {
				removeFiles = append(removeFiles, filePath)
				break
			}
		}
	}
	return moveFiles, removeFiles
}

func Setup() {
	//back up source path
	log.Println("backup source files...")
	err := cleanDirectory(config.MoveCpy)
	if err != nil {
		log.Fatalf("clean movecpy files failed, err: %v", err)
	}
	err = backup(config.Source, config.MoveCpy)
	if err != nil {
		log.Fatalf("move files failed, err: %v", err)
	}

	//clean backup
	log.Println("clean up backup files...")
	err = cleanDirectory(config.Backup)
	if err != nil {
		log.Fatalf("clean backup files failed, err: %v", err)
	}

	log.Println("collecting files...")
	//list setup files
	setupFiles, err := listFilesInDirectory(config.Setup)
	if err != nil {
		log.Fatalf("list setup files failed, path: %v, err: %v", config.Setup, err)
	}

	//list source files
	sourceFiles, err := listFilesInDirectory(config.Source)
	if err != nil {
		log.Fatalf("list src files failed, path: %v, err: %v", config.Source, err)
	}

	//do move
	log.Println("start setup...")
	err = move(sourceFiles, setupFiles)
	if err != nil {
		//restore
		log.Printf("move files failed, err: %v", err)
		log.Println("restore files")
		restore(config.MoveCpy, config.Source)
		// log.Fatalf("move files failed, err: %v", err)
		return
	}
	log.Println("setup successfully")
}

func Release() {
	//back up source path
	log.Println("backup source files...")
	err := cleanDirectory(config.ReleaseCpy)
	if err != nil {
		log.Fatalf("clean releaseCpy files failed, err: %v", err)
	}
	err = backup(config.Source, config.ReleaseCpy)
	if err != nil {
		log.Fatalf("move files failed, err: %v", err)
	}

	//clean setup
	log.Println("clean up setup files...")
	err = cleanDirectory(config.Setup)
	if err != nil {
		log.Fatalf("clean setup files failed, err: %v", err)
	}

	//list source files
	log.Println("collecting files...")
	sourceFiles, err := listFilesInDirectory(config.Source)
	if err != nil {
		log.Fatalf("list src files failed, path: %v, err: %v", config.Source, err)
	}

	//list backup files
	backupFiles, err := listFilesInDirectory(config.Backup)
	if err != nil {
		log.Fatalf("list backup files failed, path: %v, err: %v", config.Backup, err)
	}

	log.Println("start releasing...")
	err = release(sourceFiles, backupFiles)
	if err != nil {
		//restore
		log.Printf("release files failed, err: %v", err)
		log.Println("restore files")
		restore(config.ReleaseCpy, config.Source)
		return
	}
	log.Println("release successfully")
}

func Clean() {
	//remove cpy files
	os.RemoveAll(config.MoveCpy)
	os.RemoveAll(config.ReleaseCpy)
}

func Restore(flag RestoreFlag) {
	switch flag {
	case RestoreCopyFlag:
		log.Println("restore setup backup")
		restore(config.MoveCpy, config.Source)
	case RestoreReleaseFlag:
		log.Println("restore release backup")
		restore(config.ReleaseCpy, config.Source)
	default:
		log.Fatal("unknown flag")
	}
}
