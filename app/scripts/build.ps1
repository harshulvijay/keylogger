#requires -version 3.0

# set Golang environment variables
$env:GOARCH="amd64"
# set the C compiler to gcc and CXX compiler to g++
$env:CC="gcc"
$env:CXX="g++"
# enable cgo
$env:CGO_ENABLED="1"
# set optimization flags for cgo
$env:CGO_CFLAGS="-O2 -s"
$env:CGO_CPPFLAGS="-O2 -s"
$env:CGO_CXXFLAGS="-O2 -s"
$env:CGO_LDFLAGS="-O2"

# -ldflags="-s -w" strips the debug info, making it a tiny bit more challenging
# to reverse engineer the binary
# 
# for -H=windowsgui, see https://stackoverflow.com/a/23250506
go build `
  -C "$PSScriptRoot/../src/" `
  -o "$PSScriptRoot/../.out/app.exe" `
  -ldflags="-s -w -H=windowsgui"

