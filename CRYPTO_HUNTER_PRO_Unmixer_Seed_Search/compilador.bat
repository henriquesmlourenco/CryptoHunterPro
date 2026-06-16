@echo off
chcp 65001 >nul 2>&1
cd /d "%~dp0"
setlocal enabledelayedexpansion

echo ================================================================
echo   CRYPTO HUNTER PRO - UNMIXER SEED SEARCH
echo   Compilador para Windows (x64 e x86)
echo ================================================================
echo.

REM ================================================================
REM  Verificar Go instalado
REM ================================================================
echo [1/6] Verificando Go...
go version >nul 2>&1
if errorlevel 1 (
    echo.
    echo [ERRO] Go nao encontrado!
    echo.
    echo Instale o Go em: https://go.dev/dl/
    echo Baixe a versao "go1.24.4.windows-amd64.msi" ou superior.
    echo Apos instalar, REINICIE o computador e execute novamente.
    echo.
    pause
    exit /b 1
)
for /f "tokens=3" %%v in ('go version') do set GO_VERSION=%%v
echo [OK] Go encontrado: %GO_VERSION%
echo.

REM ================================================================
REM  Selecionar tipo de licenca
REM ================================================================
echo ================================================================
echo   TIPO DE LICENCA
echo ================================================================
echo.
echo   [1] Vitalicio (sem expiracao - versao completa)
echo   [2] Com expiracao (trial - versao com prazo)
echo.

:ask_license
set /p LICENSE_TYPE="  Escolha (1 ou 2): "

if "%LICENSE_TYPE%"=="1" goto lifetime
if "%LICENSE_TYPE%"=="2" goto trial

echo   [!] Opcao invalida. Digite 1 ou 2.
goto ask_license

REM ================================================================
REM  MODO VITALICIO
REM ================================================================
:lifetime
echo.
echo   [OK] Modo VITALICIO selecionado (sem expiracao).
echo.

set LDFLAGS=-s -w
set EXE_NAME_64=unmixer_seed_search_x64.exe
set EXE_NAME_86=unmixer_seed_search_x86.exe

goto compile

REM ================================================================
REM  MODO TRIAL (COM EXPIRACAO)
REM ================================================================
:trial
echo.

:ask_days
set /p TRIAL_DAYS="  Quantos dias de trial? (ex: 5, 15, 30): "

REM Validar que e um numero
set "VALID="
for /f "delims=0123456789" %%i in ("%TRIAL_DAYS%") do set VALID=%%i
if defined VALID (
    echo   [!] Digite apenas numeros.
    goto ask_days
)
if "%TRIAL_DAYS%"=="" (
    echo   [!] Digite um numero de dias.
    goto ask_days
)
if %TRIAL_DAYS% LEQ 0 (
    echo   [!] O numero de dias deve ser maior que zero.
    goto ask_days
)

REM Calcular data de expiracao usando PowerShell
for /f "usebackq" %%d in (`powershell -NoProfile -Command "(Get-Date).AddDays(%TRIAL_DAYS%).ToString('yyyyMMdd')"`) do set EXP_DATE=%%d
for /f "usebackq" %%d in (`powershell -NoProfile -Command "(Get-Date).AddDays(%TRIAL_DAYS%).ToString('dd/MM/yyyy')"`) do set EXP_DATE_DISPLAY=%%d

echo.
echo   [OK] Modo TRIAL selecionado.
echo   [OK] Expiracao em %TRIAL_DAYS% dia(s): %EXP_DATE_DISPLAY%
echo   [OK] Data embutida no binario: %EXP_DATE%
echo.

set LDFLAGS=-s -w -X main.expirationDate=%EXP_DATE%
set EXE_NAME_64=unmixer_seed_search_x64_Expires_in_%EXP_DATE%.exe
set EXE_NAME_86=unmixer_seed_search_x86_Expires_in_%EXP_DATE%.exe

goto compile

REM ================================================================
REM  COMPILACAO
REM ================================================================
:compile

REM ================================================================
REM  Baixar dependencias
REM ================================================================
echo [2/6] Baixando dependencias...
go mod tidy >nul 2>&1
go mod download >nul 2>&1
echo [OK] Dependencias prontas.
echo.

REM ================================================================
REM  Compilar para Windows 64-bit
REM ================================================================
echo [3/6] Compilando versao Windows 64-bit...
echo.

set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

go build -ldflags="!LDFLAGS!" -trimpath -o "!EXE_NAME_64!" .

if errorlevel 1 (
    echo.
    echo [ERRO] Falha na compilacao 64-bit!
    echo Verifique os arquivos .go e execute: go vet ./...
    echo.
    pause
    exit /b 1
)

if not exist "%~dp0!EXE_NAME_64!" (
    echo [ERRO] Executavel 64-bit nao foi gerado!
    pause
    exit /b 1
)

for %%A in ("%~dp0!EXE_NAME_64!") do set SIZE64=%%~zA
echo [OK] !EXE_NAME_64! (%SIZE64% bytes)
echo.

REM ================================================================
REM  Compilar para Windows 32-bit
REM ================================================================
echo [4/6] Compilando versao Windows 32-bit...
echo.

set CGO_ENABLED=0
set GOOS=windows
set GOARCH=386

go build -ldflags="!LDFLAGS!" -trimpath -o "!EXE_NAME_86!" .

if errorlevel 1 (
    echo.
    echo [ERRO] Falha na compilacao 32-bit!
    echo Verifique os arquivos .go e execute: go vet ./...
    echo.
    pause
    exit /b 1
)

if not exist "%~dp0!EXE_NAME_86!" (
    echo [ERRO] Executavel 32-bit nao foi gerado!
    pause
    exit /b 1
)

for %%A in ("%~dp0!EXE_NAME_86!") do set SIZE86=%%~zA
echo [OK] !EXE_NAME_86! (%SIZE86% bytes)
echo.

REM ================================================================
REM  Resultado final
REM ================================================================
echo [5/6] Verificacao final...
echo.
echo ================================================================
echo   COMPILACAO CONCLUIDA COM SUCESSO!
echo ================================================================
echo.
echo   Executaveis gerados:
echo.
echo   !EXE_NAME_64!  -  Windows 64-bit (Windows 7/8/10/11)
echo   !EXE_NAME_86!  -  Windows 32-bit (Windows 7/8/10/11)
echo.

if "%LICENSE_TYPE%"=="1" (
    echo   Tipo: VITALICIO (sem expiracao)
) else (
    echo   Tipo: TRIAL (expira em %EXP_DATE_DISPLAY%)
    echo   Data embutida no binario: %EXP_DATE%
    echo   A data NAO pode ser alterada sem recompilar.
)

echo.
echo   Ambos compilados com linkagem estatica (CGO_ENABLED=0).
echo   Nao precisam de Go ou outras dependencias para executar.
echo.
echo   Qual usar?
echo   - Se o Windows for 64-bit (maioria): !EXE_NAME_64!
echo   - Se o Windows for 32-bit (antigo):  !EXE_NAME_86!
echo   - Na duvida, use o x64. Se nao abrir, tente o x86.
echo.

if "%LICENSE_TYPE%"=="2" (
    echo   SEGURANCA DA LICENCA:
    echo   - A data de expiracao esta embutida dentro do .exe
    echo   - O programa verifica a hora real via NTP (internet)
    echo   - Nao e possivel burlar alterando o relogio do Windows
    echo   - Sem internet, o programa nao abre (protecao anti-bypass)
    echo   - Para renovar, recompile com nova data de expiracao
    echo.
)

pause
