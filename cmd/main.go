package main

import (
	"fmt"
	"github.com/Riznyk01/dbf"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	startMessage         = "Hello 👋 \nThis program is designed to modify column values in either a single file or a group of files that you drag and drop onto the executable file. \nSimply drag and drop the files onto the executable, enter the field name, and specify the new value. \nThe modified files will be saved in the 'changed' directory within the same folder where the original files are located.\n"
	outputFolder         = "changed"
	logFileName          = "error_log.txt"
	fileExt              = "DBF"
	failedToOpen         = "Failed to open log file"
	enterColumn          = "Please, enter a column name to change:\n"
	enterValue           = "Please, enter a value and press enter:\n"
	dropDBF              = "Please, drop only DBF files, try again"
	enteredColumnValue   = "You entered column: %s, value: %s\n"
	dropTheFiles         = "Please drop files on the executable file"
	working              = "Processing file:"
	creatingFoldersError = "Error occurred while creating folders"
	fileSavedMessage     = "File %s has been successfully saved.\n"
	panicMessage         = "Panic occurred:"
	successMessage       = "Program completed successfully."
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Printf("%s: %v", failedToOpen, err)
		os.Exit(1)
	}
	logger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	var wg sync.WaitGroup
	var columnValue, column string

	fmt.Printf(startMessage)

	checkForDroppedFiles(os.Args)
	checkForOtherFormats(os.Args)

	fmt.Printf("%s", enterColumn)
	_, err := fmt.Scan(&column)
	if err != nil {
		logger.Println(err)
	}
	fmt.Printf("%s", enterValue)
	_, err = fmt.Scan(&columnValue)
	if err != nil {
		logger.Println(err)
	}
	fmt.Printf(enteredColumnValue, column, columnValue)
	wg.Add(len(os.Args) - 1)
	go func() {
		for _, filePath := range os.Args[1:] {
			go processDBFFile(filePath, column, columnValue, &wg)
		}
	}()

	wg.Wait()
	fmt.Println(successMessage)
	<-time.After(15 * time.Second)
}

func processDBFFile(filePath string, col string, val string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%s %s\n", working, filePath)
	_, fileName := filepath.Split(filePath)
	dirPath := filepath.Dir(filePath)
	changedFilesDir := filepath.Join(dirPath, outputFolder)
	pathForTheChangedFiles := filepath.Join(dirPath, outputFolder, fileName)

	t, err := dbf.LoadFile(filePath)
	if err != nil {
		logger.Println(err)
	}

	for i := 0; i < t.NumRecords(); i++ {
		defer func() {
			if r := recover(); r != nil {
				logger.Printf(panicMessage, r)
			}
		}()
		t.SetFieldValueByName(i, col, val)
	}

	_, err = os.Stat(changedFilesDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(changedFilesDir, os.ModePerm)
		if err != nil {
			logger.Println(creatingFoldersError, err)
			return
		}
	}

	err = t.SaveFile(pathForTheChangedFiles)
	if err != nil {
		logger.Println(err)
	}
	fmt.Printf(fileSavedMessage, pathForTheChangedFiles)
}
func checkForDroppedFiles(files []string) {
	for len(files) < 2 {
		fmt.Printf("%s", dropTheFiles)
		<-time.After(2 * time.Second)
		os.Exit(0)
	}
}
func checkForOtherFormats(files []string) {
	for _, filePath := range files[1:] {
		if !strings.HasSuffix(filePath, fileExt) {
			fmt.Printf("%s", dropDBF)
			<-time.After(2 * time.Second)
			os.Exit(0)
		}
	}
}
