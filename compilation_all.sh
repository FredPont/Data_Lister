echo "Start compilation for windows 64bits"

CC=x86_64-w64-mingw32-gcc GOOS=windows CGO_ENABLED=1 go build -o Win_DataLister.exe .


echo "Start compilation for Mac 64bits"
GOOS=darwin go build -o Mac_DataLister.bin .


echo "Start compilation for linux 64bits"
go build -o Linux_DataLister.bin .
