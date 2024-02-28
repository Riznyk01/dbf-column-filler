package main

import (
	"fmt"
	"github.com/tadvi/dbf"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var columnValue, column string

	checkForDroppedFiles(os.Args)
	checkForOtherFormats(os.Args)

	fmt.Printf("please enter a column name to change:\n")
	_, err := fmt.Scan(&column)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("please enter a value and press enter:\n")
	_, err = fmt.Scan(&columnValue)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("you entered column:%s, value:%s\n", column, columnValue)

	wg.Add(len(os.Args) - 1)
	go func() {
		for _, filePath := range os.Args[1:] {
			go processDBFFile(filePath, column, columnValue, &wg)
		}
	}()

	wg.Wait()
	<-time.After(15 * time.Second)
}

func processDBFFile(filePath string, col string, val string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("processing file: %s\n", filePath)
	_, fileName := filepath.Split(filePath)
	dirPath := filepath.Dir(filePath)
	changedFilesDir := filepath.Join(dirPath, "changed")
	pathForTheChangedFiles := filepath.Join(dirPath, "changed", fileName)

	t, err := dbf.LoadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < t.NumRecords(); i++ {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("panic occurred:", r)
			}
		}()
		t.SetFieldValueByName(i, col, val)
	}

	_, err = os.Stat(changedFilesDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(changedFilesDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating folders:", err)
			return
		}
	}

	err = t.SaveFile(pathForTheChangedFiles)
	if err != nil {
		fmt.Println(err)
	}
}
func checkForDroppedFiles(files []string) {
	for len(files) < 2 {
		fmt.Println("please drop files on the EXE file")
		<-time.After(2 * time.Second)
		os.Exit(0)
	}
}
func checkForOtherFormats(files []string) {
	for _, filePath := range files[1:] {
		if !strings.HasSuffix(filePath, "DBF") {
			fmt.Println("please drop only DBF files, try again")
			<-time.After(2 * time.Second)
			os.Exit(0)
		}
	}
}
