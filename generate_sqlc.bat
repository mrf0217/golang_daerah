@echo off
REM SQLC Code Generation Script for Windows
REM This script generates Go code from SQL queries using sqlc

echo Generating sqlc code...

REM Check if sqlc is installed
where sqlc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: sqlc is not installed
    echo Install it with: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    exit /b 1
)

REM Generate code for all databases
sqlc generate

if %ERRORLEVEL% EQU 0 (
    echo [SUCCESS] sqlc code generated successfully!
    echo Generated files are in: internal/repository/generated/
) else (
    echo [ERROR] sqlc generation failed
    exit /b 1
)

