package main

import (
	"fmt"
	"os"
)

const WHITE = "\033[1;37m"
const GREEN = "\033[1;32m"
const RED = "\033[1;31m"
const BLUE = "\033[1;34m"

const KBYTE = 1000

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("%s[-] %sUsage: bcp file1 file2 fileN /path/to/copy/\n", RED, WHITE)
		os.Exit(1)
	}
	GetFileNames := os.Args[1 : len(os.Args)-1]
	FilePath := os.Args[len(os.Args)-1]

	var GetPath string

	if FilePath[len(FilePath)-1] != '/' {
		GetPath = FilePath + "/"
	} else {
		GetPath = FilePath
	}

	for index := 0; index < len(GetFileNames); index++ {
		CheckIfDirStatus, CheckIfDirError := CheckIfDir(GetFileNames[index])
		if CheckIfDirError != nil {
			fmt.Printf("%s[-] %s%v\n", RED, WHITE, CheckIfDirError)
		}

		if CheckIfDirStatus == 0 {
			SrcFile, DstFile, OpenFilesStatus := OpenFiles(GetFileNames[index], GetPath)
			if OpenFilesStatus != nil {
				fmt.Printf("\n%s[-] %sFile:%s Destination:%s Error:%v\n", RED, WHITE, GetFileNames[index], GetPath, OpenFilesStatus)
			} else {
				FileSize := GetFileSize(GetFileNames[index])
				CopyProcessStatus := CopyProcess(SrcFile, DstFile, FileSize, GetFileNames[index], GetPath)
				if CopyProcessStatus != nil {
					fmt.Printf("%s[-] %sFAIL Error:%v\n", RED, WHITE, CopyProcessStatus)
				} else {
					continue
				}
			}
		} else {
			continue
		}
	}
}
