compiler=go
buildArg=build
outputArg=-o
binFile=bcp
sourceFileMain=main.go
sourceFileFunction=CalcFunctions.go  FileFunctions.go


all: build

build:
	$(compiler) $(buildArg) $(outputArg) $(binFile) $(sourceFileMain) $(sourceFileFunction)

clear:
	rm $(binFile)