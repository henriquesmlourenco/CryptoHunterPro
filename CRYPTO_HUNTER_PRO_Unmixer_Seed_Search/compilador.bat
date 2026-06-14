@echo off
chcp 65001 >nul 2>nul
setlocal enabledelayedexpansion

echo ================================================================
echo CRYPTO HUNTER PRO - UNMIXER SEED SEARCH v1.6
echo Compilador para Windows (x64 e x86)
echo ================================================================
echo.

REM Verificar se Go esta instalado
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [ERRO] Go nao esta instalado!
    echo.
    echo Por favor, instale Go em: https://go.dev/dl/
    echo.
    pause
    exit /b 1
)

echo [OK] Go encontrado!
go version
echo.

echo Baixando dependencias...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo [ERRO] Falha ao baixar dependencias!
    pause
    exit /b 1
)
echo [OK] Dependencias baixadas!
echo.

echo ================================================================
echo TIPO DE LICENCA
echo ================================================================
echo.
echo [1] Vitalicio (sem data de expiracao)
echo [2] Com data de expiracao (trial)
echo.

:ESCOLHA_LICENCA
set /p LICENCA_TIPO="Escolha (1 ou 2): "
if "%LICENCA_TIPO%"=="1" goto VITALICIO
if "%LICENCA_TIPO%"=="2" goto TRIAL
echo [X] Opcao invalida! Escolha 1 ou 2.
goto ESCOLHA_LICENCA

:VITALICIO
set "LDFLAGS="
set "NOME_X64=unmixer_seed_search_x64.exe"
set "NOME_X86=unmixer_seed_search_x86.exe"
set "LICENCA_INFO=Vitalicio (sem expiracao)"
echo.
echo [OK] Licenca vitalicia selecionada.
echo.
goto COMPILAR

:TRIAL
echo.
:PEDIR_DIAS
set /p DIAS="Quantos dias de validade? (ex: 30): "
if "%DIAS%"=="" (
    echo [X] Digite um numero valido!
    goto PEDIR_DIAS
)

REM Calcular data de expiracao usando PowerShell
for /f "tokens=*" %%a in ('powershell -Command "(Get-Date).AddDays(%DIAS%).ToString('yyyyMMdd')"') do set EXPDATE=%%a

if "%EXPDATE%"=="" (
    echo [ERRO] Falha ao calcular data de expiracao!
    pause
    exit /b 1
)

REM Formatar data para exibicao DD/MM/YYYY
set "EXP_ANO=%EXPDATE:~0,4%"
set "EXP_MES=%EXPDATE:~4,2%"
set "EXP_DIA=%EXPDATE:~6,2%"

set "LDFLAGS=-X main.expirationDate=%EXPDATE%"
set "NOME_X64=unmixer_seed_search_x64_Expires_in_%EXPDATE%.exe"
set "NOME_X86=unmixer_seed_search_x86_Expires_in_%EXPDATE%.exe"
set "LICENCA_INFO=Trial - Expira em %EXP_DIA%/%EXP_MES%/%EXP_ANO% (%DIAS% dias)"

echo.
echo [OK] Licenca trial selecionada.
echo     Expira em: %EXP_DIA%/%EXP_MES%/%EXP_ANO% (%DIAS% dias)
echo.

:COMPILAR
echo ================================================================
echo COMPILANDO...
echo ================================================================
echo.

REM Compilar x64
echo Compilando x64 (64-bit)...
if "%LDFLAGS%"=="" (
    set CGO_ENABLED=0
    set GOOS=windows
    set GOARCH=amd64
    go build -o "%NOME_X64%" .
) else (
    set CGO_ENABLED=0
    set GOOS=windows
    set GOARCH=amd64
    go build -ldflags "%LDFLAGS%" -o "%NOME_X64%" .
)
if %ERRORLEVEL% NEQ 0 (
    echo [ERRO] Falha na compilacao x64!
    pause
    exit /b 1
)
echo [OK] x64 compilado: %NOME_X64%
echo.

REM Compilar x86
echo Compilando x86 (32-bit)...
if "%LDFLAGS%"=="" (
    set CGO_ENABLED=0
    set GOOS=windows
    set GOARCH=386
    go build -o "%NOME_X86%" .
) else (
    set CGO_ENABLED=0
    set GOOS=windows
    set GOARCH=386
    go build -ldflags "%LDFLAGS%" -o "%NOME_X86%" .
)
if %ERRORLEVEL% NEQ 0 (
    echo [ERRO] Falha na compilacao x86!
    pause
    exit /b 1
)
echo [OK] x86 compilado: %NOME_X86%
echo.

echo ================================================================
echo COMPILACAO CONCLUIDA COM SUCESSO!
echo ================================================================
echo.
echo Tipo de licenca: %LICENCA_INFO%
echo.
echo Executaveis criados:
echo   x64: %NOME_X64%
echo   x86: %NOME_X86%
echo.
echo Tamanhos:
dir /b "%NOME_X64%" 2>nul && for %%F in ("%NOME_X64%") do echo   x64: %%~zF bytes
dir /b "%NOME_X86%" 2>nul && for %%F in ("%NOME_X86%") do echo   x86: %%~zF bytes
echo.
echo ================================================================
echo SEGURANCA:
echo   - A data de expiracao e verificada via NTP (anti-manipulacao)
echo   - Servidores: Google, Cloudflare, pool.ntp.org, Windows, Apple
echo   - Sem internet = programa nao abre (protecao anti-bypass)
echo ================================================================
echo.
pause
