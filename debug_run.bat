@REM  User 'admin' initialized with randomly generated password: vxy4B5UT18doPv96
@echo off
setlocal EnableExtensions
title File Browser Debug Run

echo ==========================================
echo    File Browser - Debug Run Tool
echo ==========================================

:: 1. Check for Go
echo [1/3] Checking Go...
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go is not found in your PATH. 
    echo Please install Go and try again.
    pause
    exit /b 1
)
echo - Go: Found.

:: 2. Check for pnpm
echo [2/3] Checking pnpm...
where pnpm >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] pnpm is not found in your PATH.
    echo Please run: npm install -g pnpm
    pause
    exit /b 1
)
echo - pnpm: Found.

:: 3. Choose mode
echo.
echo [3/3] Choose run mode:
echo [1] Dev mode (Vite + Go Backend) - FAST REFRESH
echo [2] Integration mode (Build + Go Backend) - FINAL TEST
echo.
set /p mode="Select mode [1/2, default=1]: "
if "%mode%"=="" set mode=1

if "%mode%"=="1" (
    echo Starting Dev Mode...
    start "FB-Backend" cmd /c "go run main.go -p 8080 -d filebrowser.db"
    timeout /t 2 > nul
    cd frontend
    start http://localhost:5173
    call pnpm dev
) else (
    echo Starting Integration Mode...
    cd frontend
    echo Building...
    call pnpm build
    if %errorlevel% neq 0 (
        echo [ERROR] Build failed!
        pause
        exit /b 1
    )
    cd ..
    echo Running...
    start http://localhost:8080
    go run main.go -p 8080 -d filebrowser.db
)

pause
