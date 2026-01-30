@echo off
set CGO_ENABLED=0

echo === 1. Build Frontend ===
if exist "frontend" (
    cd frontend
) else (
    echo Error: frontend directory not found!
    pause
    exit /b 1
)

echo Installing dependencies...
call pnpm install
if %errorlevel% neq 0 (
    echo Error during pnpm install
    pause
    exit /b %errorlevel%
)

echo Building frontend assets...
call pnpm build
if %errorlevel% neq 0 (
    echo Error during pnpm build
    pause
    exit /b %errorlevel%
)

cd ..

echo === 2. Build Backend (Go) ===
echo Building filebrowser.exe...
"C:\Users\anhnt\dev\go\bin\go.exe" build -ldflags="-s -w" -o filebrowser.exe
if %errorlevel% neq 0 (
    echo Error during go build
    pause
    exit /b %errorlevel%
)

echo.
echo ==========================================
echo       BUILD SUCCESSFUL
echo ==========================================
echo filebrowser.exe has been updated.
@REM pause
