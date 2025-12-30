package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/proxy"
)

const (
	TorProxyAddress = "127.0.0.1:9150" // Tor Browser portu
	ReportFileName  = "scan_report.log"
	OutputDir       = "output_data"
)

func main() {
	fmt.Println("===========================================")
	fmt.Println("   Siber Tehdit İstihbaratı (CTI) Aracı    ")
	fmt.Println("===========================================")

	// 1. MODÜL: Hedefleri Yükle (Input Handler)
	targets, err := loadTargets("targets.yaml")
	if err != nil {
		log.Fatalf("[!] Giriş dosyası hatası: %v", err)
	}

	// Rapor dosyasını hazırla
	report, _ := os.OpenFile(ReportFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer report.Close()
	report.WriteString(fmt.Sprintf("\n--- Tarama Başlangıcı: %s ---\n", time.Now().Format("2006-01-02 15:04:05")))

	// 2. MODÜL: Tor Proxy Yapılandırması
	torClient, err := setupTorClient()
	if err != nil {
		fmt.Printf("[X] Tor bağlantı hatası: %v\n", err)
		return
	}

	// IP Doğrulama (Ödev Kanıtı)
	verifyTorIP(torClient)

	// Chromedp (Ekran Görüntüsü) Ayarları
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer("socks5://"+TorProxyAddress),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// 3. & 4. MODÜL: Tarama ve Kayıt
	for _, target := range targets {
		processTarget(allocCtx, torClient, target, report)
	}

	fmt.Println("\n[+] Operasyon tamamlandı. Çıktılar '" + OutputDir + "' klasöründe.")
}

func loadTargets(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "url:") {
			parts := strings.SplitN(line, "url:", 2)
			if len(parts) > 1 {
				u := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				if u != "" {
					urls = append(urls, u)
				}
			}
		}
	}
	return urls, nil
}

func setupTorClient() (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", TorProxyAddress, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: &http.Transport{
			Dial:            dialer.Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 60 * time.Second,
	}, nil
}

func verifyTorIP(client *http.Client) {
	fmt.Print("[INFO] Tor bağlantısı kontrol ediliyor... ")
	resp, err := client.Get("https://check.torproject.org/api/ip")
	if err == nil {
		defer resp.Body.Close()
		ip, _ := io.ReadAll(resp.Body)
		fmt.Printf("BAŞARILI. IP: %s\n", string(ip))
	} else {
		fmt.Println("BAŞARISIZ (Proxy aktif mi?)")
	}
}

func processTarget(allocCtx context.Context, client *http.Client, target string, report *os.File) {
	fmt.Printf("[>] Hedef: %s ... ", target)

	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Firefox/115.0")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERİŞİLEMEDİ")
		report.WriteString(fmt.Sprintf("[FAIL] %s - %v\n", target, err))
		return
	}
	defer resp.Body.Close()

	fmt.Println("AKTİF")
	report.WriteString(fmt.Sprintf("[SUCCESS] %s - HTTP %d\n", target, resp.StatusCode))

	// Dosya adını temizle (URL'deki yasaklı karakterleri sil)
	safeDir := strings.ReplaceAll(target, "http://", "")
	safeDir = strings.ReplaceAll(safeDir, "/", "_")
	safeDir = strings.ReplaceAll(safeDir, ":", "_")
	folderPath := filepath.Join(OutputDir, safeDir)
	os.MkdirAll(folderPath, 0755)

	body, _ := io.ReadAll(resp.Body)
	os.WriteFile(filepath.Join(folderPath, "source.html"), body, 0644)

	takeScreenshot(allocCtx, target, folderPath)
}

func takeScreenshot(allocCtx context.Context, url string, folder string) {
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 180*time.Second) // Tor yavaş olduğu için süreyi uzattık
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(30*time.Second), // Sayfanın render edilmesi için bekle
		chromedp.FullScreenshot(&buf, 90),
	)

	if err == nil {
		os.WriteFile(filepath.Join(folder, "screenshot.png"), buf, 0644)
		fmt.Println("   -> Görsel kanıt kaydedildi.")
	} else {
		fmt.Printf("   -> Görsel alınamadı: %v\n", err)
	}
}
