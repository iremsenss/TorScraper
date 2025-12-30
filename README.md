# 🛡️ TorScraper: Automated CTI Intelligence Collector

**TorScraper**, Siber Tehdit İstihbaratı (CTI) süreçlerinin "Toplama" (Collection) aşamasını otomatize etmek amacıyla geliştirilmiş, Go tabanlı bir anonim veri kazıma aracıdır. Tor ağı üzerindeki `.onion` uzantılı adresleri güvenli ve anonim bir şekilde analiz ederek metinsel ve görsel kanıt toplar.

---

## 🚀 Öne Çıkan Özellikler

* **🕵️ Tam Anonimlik:** Tüm trafik SOCKS5 üzerinden Tor ağına yönlendirilir, IP sızıntısı önlenir.
* **📸 Görsel İstihbarat:** `chromedp` kullanarak web sitelerinin tam sayfa ekran görüntülerini (screenshot) otomatik olarak alır.
* **📄 Veri Arşivleme:** Hedef sitelerin HTML kaynak kodlarını analiz için yerel diskte depolar.
* **🛠️ Hata Toleransı:** Zaman aşımı (timeout) ve bağlantı hatalarını yöneterek tarama sürecini kesintisiz sürdürür.
* **📊 Dinamik Raporlama:** Süreçle ilgili tüm detayları (SUCCESS/FAIL) zaman damgalı bir log dosyasına kaydeder.
* **⚡ Performans:** Go'nun eşzamanlılık avantajlarını kullanarak optimize edilmiş tarama süreci.

---

## 🏗️ Teknik Mimari

Proje 4 ana modül üzerine inşa edilmiştir:
1.  **Input Handler:** `targets.yaml` dosyasından hedefleri temizleyerek okur.
2.  **Tor Proxy Client:** `net/http` ve `proxy` paketleri ile Tor tünellemesi yapar.
3.  **Scraper & Screenshot:** Veri toplama ve görsel kayıt işlemlerini yürütür.
4.  **Output Writer:** Elde edilen verileri hiyerarşik bir klasör yapısında arşivler.

---

## 📋 Gereksinimler

* **Go:** v1.18+
* **Tor Service:** Arka planda çalışıyor olmalıdır (Port: `9150` veya `9050`).
* **Tarayıcı:** Chrome veya Chromium (Screenshot özelliği için).

---

## 🔧 Kurulum ve Kullanım

### 1. Bağımlılıkları Yükleyin
```bash
git clone https://github.com/iremsenss/TorScraper.git
cd TorScraper
go mod download
go run main.go



