package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func CheckIfDir(FileName string) (int, error) {
	GetStat, GetStatError := os.Stat(FileName)
	if GetStatError != nil {
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

func GetFileSize(FileName string) int64 {
	FileStat, _ := os.Stat(FileName)
	return FileStat.Size()
}

func SeparateFileName(FileName string) string {
	FileNameLen := len(FileName)
	var GetFileName, RevFileName []string
	var result string

	for i := FileNameLen - 1; ; i-- {
		if strings.Compare(string(FileName[i]), "/") == 0 {
			break
		}
		GetFileName = append(GetFileName, string(FileName[i]))
	}

	for i := len(GetFileName) - 1; i >= 0; i-- {
		RevFileName = append(RevFileName, GetFileName[i])
	}
	result = strings.Join(RevFileName, "")
	return result
}

func OpenFiles(SourceFile, DestinationPath string) (*os.File, *os.File, error) {
	var DestinyFilePath string
	if strings.Contains(SourceFile, "/") {
		DestinyFilePath = DestinationPath + SeparateFileName(SourceFile)
	} else {
		DestinyFilePath = DestinationPath + SourceFile
	}
	_, DstFileCreateStatus := os.Create(DestinyFilePath)
	if DstFileCreateStatus != nil {
		return nil, nil, DstFileCreateStatus
	}

	DstFileOpen, DstFileOpenStatus := os.OpenFile(DestinyFilePath, os.O_CREATE|os.O_RDWR, 0444)
	if DstFileOpenStatus != nil {
		return nil, nil, DstFileOpenStatus
	}

	SrcFileOpen, SrcFileOpenStatus := os.Open(SourceFile)
	if SrcFileOpenStatus != nil {
		return nil, nil, SrcFileOpenStatus
	}
	return SrcFileOpen, DstFileOpen, nil
}

func CopyProcess(SourceFile, DestinationPath *os.File, FileSize int64, SrcFileName, DstFileName string) error {
	var BytesCopied int64
	var WrittenBytesCounter int64

	for {
		BlockSize := CalcBlockSize(FileSize, BytesCopied)
		WrittenBytes, CopyNError := io.CopyN(DestinationPath, SourceFile, BlockSize)
		WrittenBytesCounter += WrittenBytes
		PercentProgress := (float64(WrittenBytesCounter) / float64(FileSize)) * 100
		fmt.Printf("%s[*] %sFile:%s Destination:%s Progress: %2.f%%\r", BLUE, WHITE, SrcFileName, DstFileName, PercentProgress)
		if CopyNError == io.EOF {
			fmt.Printf("\n")
			return nil
		}
	}
}
