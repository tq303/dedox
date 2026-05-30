@echo off
setlocal

set ARCH=amd64
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" set ARCH=arm64

set URL=https://github.com/tq303/dedox/releases/latest/download/ddx-windows-%ARCH%.exe
set DEST=%SystemRoot%\System32\ddx.exe

echo Installing ddx for windows/%ARCH%...
curl -sL "%URL%" -o "%DEST%"
echo Installed to %DEST%
