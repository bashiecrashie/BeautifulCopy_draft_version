package main

import (
	"fmt"
	"os"
	"io"
	"strings"
)

//Colors
const WHITE = "\033[1;37m"
const GREEN = "\033[1;32m"
const RED = "\033[1;31m"
const BLUE = "\033[1;34m"

const KBYTE = 1000

func main(){
	if len(os.Args) < 3{
		fmt.Printf("%s[-] %sUsage: bcp file1 file2 fileN /path/to/copy/\n", RED, WHITE)
		os.Exit(1)
	}
	GetFileNames := os.Args[1:len(os.Args)-1]
	GetPath := os.Args[len(os.Args)-1]

	for index := 0; index < len(GetFileNames); index++{
		CheckIfDirStatus, CheckIfDirError := CheckIfDir(GetFileNames[index])
		if CheckIfDirError != nil{
			fmt.Printf("%s[-] %s%v\n", RED, WHITE, CheckIfDirError)
		}

		if CheckIfDirStatus == 0{
			SrcFile, DstFile, OpenFilesStatus := OpenFiles(GetFileNames[index], GetPath)
			if OpenFilesStatus != nil{
				fmt.Printf("\n%s[-] %sFile:%s Destination:%s Error:%v\n", RED, WHITE, GetFileNames[index], GetPath, OpenFilesStatus)
			}else{
				FileSize := GetFileSize(GetFileNames[index])
				CopyProcessStatus := CopyProcess(SrcFile, DstFile, FileSize, GetFileNames[index], GetPath)
				if CopyProcessStatus != nil{
					fmt.Printf("%sFAIL Error:%v\n", RED, WHITE, GetFileNames[index], GetPath, CopyProcessStatus)
				}else{
					continue
				}
			}
		}else{
			continue
		}
	}
}

func CheckIfDir(FileName string)(int, error){
	GetStat, GetStatError:= os.Stat(FileName)
	if GetStatError != nil{
		return 3, GetStatError
	}

	switch GetMode := GetStat.Mode(); {
		case GetMode.IsDir():
			return 1, nil
		case GetMode.IsRegular():
			return 0, nil
		default:
			return 3, nil
	}
}

func GetFileSize(FileName string)(int64){
	FileStat, _ := os.Stat(FileName)
	return FileStat.Size()
}

func SeparateFileName(FileName string)(string){
	FileNameLen := len(FileName)

	var GetFileName, RevFileName []string
    
    var result string

  	for i := FileNameLen-1; ;i--{
    	    if strings.Compare(string(FileName[i]), "/") == 0{
        	           break
           	}
           	GetFileName = append(GetFileName, string(FileName[i]))
   	}

   for i := len(GetFileName)-1; i >= 0; i--{
           RevFileName = append(RevFileName, GetFileName[i])
   }

   result = strings.Join(RevFileName, "")

   return result
}

func OpenFiles(SourceFile, DestinationPath string)(*os.File, *os.File, error){
	
	var DestinyFilePath string
	if strings.Contains(SourceFile, "/") == true{
		DestinyFilePath = DestinationPath + SeparateFileName(SourceFile)
	}
	DestinyFilePath = DestinationPath + SourceFile

	_, DstFileCreateStatus := os.Create(DestinyFilePath)
	if DstFileCreateStatus != nil{
		return nil, nil, DstFileCreateStatus
	}

	DstFileOpen, DstFileOpenStatus := os.OpenFile(DestinyFilePath, os.O_CREATE|os.O_RDWR, 444)
	if DstFileOpenStatus != nil{
		return nil, nil, DstFileOpenStatus
	}

	SrcFileOpen, SrcFileOpenStatus := os.Open(SourceFile)
	if SrcFileOpenStatus != nil{
		return nil, nil, SrcFileOpenStatus
	}
	return SrcFileOpen, DstFileOpen, nil
}

func CalcBlockSize(FileSize, BytesCopied int64) (int64){
	var BlockSize int64 = 64000
	remain := FileSize - BytesCopied
	if remain < BlockSize{
		return remain
	}
	return BlockSize
}

func CopyProcess(SourceFile, DestinationPath *os.File, FileSize int64, SrcFileName, DstFileName string)(error){
	var BytesCopied int64
	var WrittenBytesCounter int64
	for{
		BlockSize := CalcBlockSize(FileSize, BytesCopied)
		WrittenBytes, CopyNError := io.CopyN(DestinationPath, SourceFile, BlockSize)
		WrittenBytesCounter += WrittenBytes
		PercentProgress := (float64(WrittenBytesCounter) / float64(FileSize)) * 100
		fmt.Printf("%s[*] %sFile:%s Destination:%s Progress: %2.f%%\r", BLUE, WHITE, SrcFileName, DstFileName, PercentProgress)
		if CopyNError == io.EOF{
			fmt.Printf("\n")
			return nil
		}
	}	
	return nil
}