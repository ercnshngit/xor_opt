# XOR Optimizasyon Web Arayüzü

Bu web arayüzü, backend Go API'sini kullanarak XOR optimizasyon algoritmalarını çalıştırır.

## Özellikler

- **SLP Heuristic Algoritması** - Sezgisel yaklaşım
- **Boyar SLP Algoritması** - Derinlik sınırlı optimizasyon
- **PAAR Algoritması** - Hamming ağırlığı tabanlı optimizasyon

## Kurulum ve Çalıştırma

### 1. Backend API'sini Başlatın

```bash
cd ../backend
go run main.go
```

Backend 8080 portunda çalışacaktır.

### 2. Web Sunucusunu Başlatın

```bash
cd web
python3 -m http.server 3000
```

Web arayüzü http://localhost:3000 adresinde erişilebilir olacaktır.

## Kullanım

1. **Matris Girişi**: Matrislerinizi metin alanına yapıştırın veya dosya yükleyin
2. **Algoritma Seçimi**: İstediğiniz algoritmayı seçin
3. **Sonuçları Görüntüleme**: Sonuçlar accordion formatında gösterilir

## Desteklenen Formatlar

### Manuel Giriş
```
x^2 ile carpilmis A1 matrisi (binary):
[1 0 1]
[0 1 0]
[1 1 1]
HamXOR Sayisi:
5
-----
```

### Dosya Yükleme
- `.txt` dosyaları desteklenir
- Büyük dosyalar için satır sayısı belirtebilirsiniz

## API Endpoints

Web arayüzü aşağıdaki backend endpoint'lerini kullanır:

- `POST http://localhost:8080/slp` - SLP Heuristic
- `POST http://localhost:8080/boyar` - Boyar SLP  
- `POST http://localhost:8080/paar` - PAAR Algoritması

## Dosya Yapısı

```
web/
├── index.html       # Ana web sayfası
├── main.js          # JavaScript kodu
├── test_data.txt    # Test verisi
└── README.md        # Bu dosya
```

## Özellikler

- **Responsive Tasarım** - Mobil ve masaüstü uyumlu
- **Progress Tracking** - İşlem durumu takibi
- **Accordion Sonuçlar** - Kolay görüntüleme
- **Dosya Yükleme** - Büyük dosya desteği
- **Hata Yönetimi** - Detaylı hata mesajları 