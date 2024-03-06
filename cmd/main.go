package main

import (
	"bufio"
	"dbf-column-filler/internal/text"
	"fmt"
	"github.com/Riznyk01/dbf"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile(text.LogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Printf("%s: %v", text.FailedToOpen, err)
		os.Exit(1)
	}
	logger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	var wg sync.WaitGroup
	var params []string

	fmt.Printf(text.StartMessage)

	checkForDroppedFiles(os.Args)
	checkForOtherFormats(os.Args)

	scanner := bufio.NewScanner(os.Stdin)

	var line string
	for {
		scanner.Scan()
		line = scanner.Text()

		params = strings.Split(line, " ")
		if len(params)%2 != 0 {
			fmt.Printf(text.EnterEven)
		} else if len(params) == 0 {
			fmt.Printf(text.DidntEnter)
		} else {
			break
		}

	}
	wg.Add(len(os.Args) - 1)
	go func() {
		for _, filePath := range os.Args[1:] {
			go processDBFFile(filePath, params, &wg)
		}
	}()

	wg.Wait()
	fmt.Println(text.SuccessMessage)
	<-time.After(15 * time.Second)
}

func processDBFFile(filePath string, par []string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%s %s\n", text.Working, filePath)
	_, fileName := filepath.Split(filePath)
	dirPath := filepath.Dir(filePath)
	changedFilesDir := filepath.Join(dirPath, text.OutputFolder)
	pathForTheChangedFiles := filepath.Join(dirPath, text.OutputFolder, fileName)

	t, err := dbf.LoadFile(filePath)
	if err != nil {
		logger.Println(err)
	}

	for i := 0; i < t.NumRecords(); i++ {
		defer func() {
			if r := recover(); r != nil {
				logger.Printf(text.PanicMessage, r)
			}
		}()
		for j := 0; j < len(par); j += 2 {
			t.SetFieldValueByName(i, par[j], par[j+1])
		}
	}

	_, err = os.Stat(changedFilesDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(changedFilesDir, os.ModePerm)
		if err != nil {
			logger.Println(text.CreatingFoldersError, err)
			return
		}
	}

	err = t.SaveFile(pathForTheChangedFiles)
	if err != nil {
		logger.Println(err)
	}
	fmt.Printf(text.FileSavedMessage, pathForTheChangedFiles)
}
func checkForDroppedFiles(files []string) {
	for len(files) < 2 {
		fmt.Printf("%s", text.DropTheFiles)
		<-time.After(2 * time.Second)
		os.Exit(0)
	}
}
func checkForOtherFormats(files []string) {
	for _, filePath := range files[1:] {
		if !strings.HasSuffix(filePath, text.FileExt) {
			fmt.Printf("%s", text.DropDBF)
			<-time.After(2 * time.Second)
			os.Exit(0)
		}
	}
}
