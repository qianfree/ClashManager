@echo off
echo Building ClashDoc...

REM 1. Build frontend
echo Step 1: Building frontend...
cd web
call pnpm install
call pnpm run build
if %errorlevel% neq 0 (
    echo Frontend build failed!
    exit /b 1
)
cd ..


REM 2. Build Go application
echo Step 2: Building Go application...
go build -o clash-manager.exe ./cmd/server
if %errorlevel% neq 0 (
    echo Go build failed!
    exit /b 1
)


echo.
echo Build completed successfully!
echo Output: clash-manager.exe
