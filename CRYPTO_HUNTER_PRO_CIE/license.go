package main

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// ============================================================================
// SISTEMA DE LICENCA - NTP + HTTPS TIME VERIFICATION (ANTI-BYPASS)
// ============================================================================
//
// A variavel expirationDate e injetada em tempo de compilacao via:
//   go build -ldflags="-X main.expirationDate=20260401"
//
// Se vazia = versao vitalicia (sem expiracao)
// Se preenchida = versao trial com data limite no formato YYYYMMDD
//
// PROTECOES:
// 1. NTP via IP direto (nao usa DNS = arquivo hosts nao afeta)
// 2. Validacao cruzada entre servidores (detecta NTP falso)
// 3. HTTPS como backup (certificado SSL impede falsificacao)
// ============================================================================

// expirationDate - injetada via -ldflags na compilacao
// Formato: YYYYMMDD (ex: "20260401" = 01/Abril/2026)
// Vazio = versao vitalicia
var expirationDate string

// ============================================================================
// NTP CLIENT - IPs DIRETOS (anti-hosts file bypass)
// ============================================================================

// ntpServerIPs - IPs diretos dos servidores NTP (nao usa DNS)
// Arquivo hosts do Windows so intercepta nomes de dominio, nao IPs
var ntpServerIPs = []string{
	"216.239.35.0:123",   // time1.google.com
	"216.239.35.4:123",   // time2.google.com
	"216.239.35.8:123",   // time3.google.com
	"216.239.35.12:123",  // time4.google.com
	"162.159.200.1:123",  // time.cloudflare.com
	"162.159.200.123:123", // time.cloudflare.com (secundario)
	"132.163.97.6:123",   // time.nist.gov (NIST)
	"17.253.20.45:123",   // time.apple.com
}

// queryNTP faz uma consulta NTP a um servidor especifico via IP direto
func queryNTP(serverIP string) (time.Time, error) {
	conn, err := net.DialTimeout("udp", serverIP, 5*time.Second)
	if err != nil {
		return time.Time{}, fmt.Errorf("conexao falhou: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Montar pacote NTP (48 bytes)
	// LI=0, VN=4, Mode=3 (client) => primeiro byte = 0x23
	req := make([]byte, 48)
	req[0] = 0x23

	_, err = conn.Write(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("envio falhou: %v", err)
	}

	resp := make([]byte, 48)
	_, err = conn.Read(resp)
	if err != nil {
		return time.Time{}, fmt.Errorf("leitura falhou: %v", err)
	}

	// Extrair timestamp de transmissao (bytes 40-47)
	seconds := binary.BigEndian.Uint32(resp[40:44])
	fraction := binary.BigEndian.Uint32(resp[44:48])

	// Converter NTP epoch (1900) para Unix epoch (1970)
	const ntpEpochOffset = 2208988800

	if uint64(seconds) < ntpEpochOffset {
		return time.Time{}, fmt.Errorf("timestamp NTP invalido")
	}

	unixSeconds := int64(seconds) - int64(ntpEpochOffset)
	nanoseconds := int64(fraction) * 1e9 / (1 << 32)

	ntpTime := time.Unix(unixSeconds, nanoseconds).UTC()

	// Validacao basica: hora deve ser razoavel (apos 2024)
	if ntpTime.Year() < 2024 {
		return time.Time{}, fmt.Errorf("hora NTP fora do esperado: %v", ntpTime)
	}

	return ntpTime, nil
}

// ============================================================================
// HTTPS TIME - Backup via HTTPS (anti-NTP fake)
// ============================================================================

// getHTTPSTime obtem a hora via header HTTP Date de servidores HTTPS
// O certificado SSL garante que a resposta vem do servidor real
func getHTTPSTime() (time.Time, error) {
	httpsURLs := []string{
		"https://www.google.com",
		"https://www.cloudflare.com",
		"https://www.microsoft.com",
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	for _, url := range httpsURLs {
		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		dateHeader := resp.Header.Get("Date")
		if dateHeader == "" {
			continue
		}

		// Parse do header Date (formato RFC 1123)
		httpTime, err := time.Parse(time.RFC1123, dateHeader)
		if err != nil {
			// Tentar formato alternativo
			httpTime, err = time.Parse("Mon, 02 Jan 2006 15:04:05 MST", dateHeader)
			if err != nil {
				continue
			}
		}

		httpTime = httpTime.UTC()

		if httpTime.Year() < 2024 {
			continue
		}

		return httpTime, nil
	}

	return time.Time{}, fmt.Errorf("nenhum servidor HTTPS retornou hora valida")
}

// ============================================================================
// VERIFICACAO DE TEMPO COM VALIDACAO CRUZADA
// ============================================================================

// getVerifiedTime obtem a hora real com validacao cruzada
// Precisa de pelo menos 2 fontes concordando (diferenca < 2 minutos)
func getVerifiedTime() (time.Time, error) {
	var times []time.Time

	// Fase 1: Coletar horas de servidores NTP (por IP direto)
	for _, server := range ntpServerIPs {
		t, err := queryNTP(server)
		if err == nil {
			times = append(times, t)
		}
		// Se ja temos 3 respostas NTP, suficiente
		if len(times) >= 3 {
			break
		}
	}

	// Fase 2: Coletar hora via HTTPS como fonte adicional
	httpsTime, err := getHTTPSTime()
	if err == nil {
		times = append(times, httpsTime)
	}

	// Precisamos de pelo menos 2 fontes
	if len(times) < 2 {
		if len(times) == 1 {
			// Se so temos 1 fonte, nao podemos validar cruzado
			// Mas ainda e melhor que nada - aceitar com cautela
			return times[0], nil
		}
		return time.Time{}, fmt.Errorf("fontes de tempo insuficientes (%d)", len(times))
	}

	// Fase 3: Validacao cruzada - verificar se as fontes concordam
	// Tolerancia: 2 minutos de diferenca maxima entre quaisquer 2 fontes
	const maxDrift = 2 * time.Minute

	for i := 0; i < len(times); i++ {
		for j := i + 1; j < len(times); j++ {
			diff := times[i].Sub(times[j])
			if diff < 0 {
				diff = -diff
			}
			if diff > maxDrift {
				return time.Time{}, fmt.Errorf("divergencia detectada entre fontes de tempo: %v", diff)
			}
		}
	}

	// Todas as fontes concordam - usar a mediana
	medianTime := getMedianTime(times)
	return medianTime, nil
}

// getMedianTime retorna o tempo mediano de uma lista de tempos
func getMedianTime(times []time.Time) time.Time {
	n := len(times)
	if n == 0 {
		return time.Time{}
	}

	// Converter para unix timestamps para ordenar
	stamps := make([]int64, n)
	for i, t := range times {
		stamps[i] = t.UnixNano()
	}

	// Bubble sort simples (lista pequena)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if stamps[j] < stamps[i] {
				stamps[i], stamps[j] = stamps[j], stamps[i]
			}
		}
	}

	// Mediana
	var medianNano int64
	if n%2 == 0 {
		medianNano = (stamps[n/2-1] + stamps[n/2]) / 2
	} else {
		medianNano = stamps[n/2]
	}

	return time.Unix(0, medianNano).UTC()
}

// ============================================================================
// VERIFICACAO DE LICENCA
// ============================================================================

// checkLicense verifica se a licenca e valida
// Retorna true se pode continuar, false se expirado ou erro
func checkLicense() bool {
	// Versao vitalicia - sem data de expiracao
	if expirationDate == "" {
		return true
	}

	// Parse da data de expiracao
	expDate, err := time.Parse("20060102", expirationDate)
	if err != nil {
		fmt.Println()
		fmt.Println("================================================================")
		fmt.Println("  [ERRO] Data de expiracao invalida no binario.")
		fmt.Println("  [ERROR] Invalid expiration date in binary.")
		fmt.Println("================================================================")
		fmt.Println()
		os.Exit(1)
		return false
	}

	// A expiracao e no final do dia (23:59:59 UTC)
	expDate = expDate.Add(24*time.Hour - time.Second)

	// Obter hora real com validacao cruzada
	fmt.Println()
	fmt.Println("  Verificando licenca / Checking license...")

	verifiedTime, err := getVerifiedTime()
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "divergencia") {
			// Detectou manipulacao de NTP
			fmt.Println()
			fmt.Println("================================================================")
			fmt.Println("  [ERRO] Manipulacao de tempo detectada!")
			fmt.Println("  [ERROR] Time manipulation detected!")
			fmt.Println("================================================================")
			fmt.Println()
			fmt.Println("  As fontes de tempo nao concordam entre si.")
			fmt.Println("  The time sources do not agree with each other.")
			fmt.Println()
			fmt.Println("  Verifique se seu sistema nao esta com NTP adulterado.")
			fmt.Println("  Check if your system does not have tampered NTP.")
			fmt.Println()
		} else {
			// Sem internet
			fmt.Println()
			fmt.Println("================================================================")
			fmt.Println("  [ERRO] Nao foi possivel verificar a licenca!")
			fmt.Println("  [ERROR] Could not verify license!")
			fmt.Println("================================================================")
			fmt.Println()
			fmt.Println("  E necessaria conexao com a internet para verificar a licenca.")
			fmt.Println("  Internet connection is required to verify the license.")
			fmt.Println()
			fmt.Println("  Verifique sua conexao e tente novamente.")
			fmt.Println("  Check your connection and try again.")
			fmt.Println()
		}

		waitExit()
		return false
	}

	// Verificar se expirou
	if verifiedTime.After(expDate) {
		fmt.Println()
		fmt.Println("================================================================")
		fmt.Println("  [EXPIRADO] Sua licenca expirou!")
		fmt.Println("  [EXPIRED] Your license has expired!")
		fmt.Println("================================================================")
		fmt.Println()
		fmt.Printf("  Data de expiracao / Expiration date: %s\n", expDate.Format("02/01/2006"))
		fmt.Printf("  Data atual / Current date: %s\n", verifiedTime.Format("02/01/2006 15:04:05 UTC"))
		fmt.Println()
		fmt.Println("  Entre em contato com a Equipe Crypto Hunter Pro para renovar.")
		fmt.Println("  Contact the Crypto Hunter Pro Team to renew.")
		fmt.Println()
		waitExit()
		return false
	}

	// Licenca valida - mostrar dias restantes
	remaining := expDate.Sub(verifiedTime)
	days := int(math.Ceil(remaining.Hours() / 24))

	fmt.Printf("  [OK] Licenca valida - expira em %d dia(s) (%s)\n",
		days, expDate.Format("02/01/2006"))
	fmt.Printf("  [OK] License valid - expires in %d day(s) (%s)\n",
		days, expDate.Format("02/01/2006"))

	return true
}

// waitExit aguarda Enter e encerra o programa
func waitExit() {
	fmt.Print("  Pressione Enter para sair / Press Enter to exit...")
	var input string
	fmt.Scanln(&input)
	os.Exit(1)
}
