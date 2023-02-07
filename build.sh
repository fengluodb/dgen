#! /bin/sh
if [ "$(uname)" == "Darwin" ]; then # Mac OS X
    go build -o dgen main.go
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then # Linux
    go build -o dgen main.go
elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ];then # Windows
    go build -o dgen.exe main.go
fi