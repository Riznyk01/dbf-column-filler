package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for len(os.Args) < 2 {
		fmt.Println("please drop files on the EXE file")
		<-time.After(2 * time.Second)
		os.Exit(0)
	}

	for _, filePath := range os.Args[1:] {
		if !strings.HasSuffix(filePath, "DBF") {
			fmt.Println("please drop only DBF files, try again")
			<-time.After(2 * time.Second)
			os.Exit(0)
		}
	}

	wg.Add(len(os.Args) - 1)
	go func() {
		for _, filePath := range os.Args[1:] {
			go processDBFFile(filePath, &wg)
		}
	}()

	wg.Wait()
	<-time.After(3 * time.Second)
}

func processDBFFile(filePath string, wg *sync.WaitGroup) {
	fmt.Printf("Processing file: %s\n", filePath)
	wg.Done()
}
