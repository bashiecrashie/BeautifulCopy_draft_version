package main

func CalcBlockSize(FileSize, BytesCopied int64) int64 {
	var BlockSize int64 = 64000
	remain := FileSize - BytesCopied
	if remain < BlockSize {
		return remain
	}
	return BlockSize
}
