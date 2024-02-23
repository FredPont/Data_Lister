echo "Start compilation for windows 64bits"

#CC=x86_64-w64-mingw32-gcc GOOS=windows CGO_ENABLED=1 go build -o Win_DataLister.exe .
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -o Win_DataLister.exe .


echo "Start compilation for Mac 64bits"
#GOOS=darwin GOARCH=amd64 go build -o MacX64_DataLister.bin .
#GOOS=darwin GOARCH=arm64 go build -o MacARM64_DataLister.bin .

echo "Start compilation for linux 64bits"
go build -o Linux_DataLister.bin .
