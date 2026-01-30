@echo off
set CGO_ENABLED=0

@REM echo === Build frontend ===
@REM cd frontend
@REM pnpm install
@REM pnpm build
@REM cd ..

echo === Build backend + embed ===
"C:\Users\anhnt\dev\go\bin\go.exe" build -ldflags="-s -w" -o filebrowser.exe

echo DONE: filebrowser.exe
