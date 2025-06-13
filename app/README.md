# XOR Optimizasyon - Matris Yönetimi Sistemi

Bu proje, binary matrislerin XOR optimizasyonu için çeşitli algoritmaları (Boyar SLP, Paar, SLP Heuristic) çalıştıran ve sonuçları Mysql veritabanında saklayan bir web uygulamasıdır.

## Özellikler

### 🗄️ Veritabanı Yönetimi
- SQLite veritabanında matris verilerinin saklanması
- Matris hash'leri ile tekrar hesaplama önleme
- Binary ve hex formatında matris depolama
- Algoritma sonuçlarının ve programlarının saklanması

### 🔍 Web Arayüzü
- Modern Bootstrap 5 tabanlı responsive arayüz
- Matris listesi ve detay görüntüleme
- Sayfalama (pagination) ve filtreleme
- Yeni matris ekleme formu
- Algoritmaları yeniden hesaplama özelliği

### ⚡ API Endpoints

#### Orijinal Algoritmalar
- `POST /boyar` - Boyar SLP algoritması
- `POST /paar` - Paar algoritması  
- `POST /slp` - SLP Heuristic algoritması

#### Veritabanı İşlemleri
- `GET /api/matrices` - Matris listesi (sayfalama ve filtreleme ile)
- `POST /api/matrices` - Yeni matris kaydetme
- `GET /api/matrices/{id}` - Matris detayları
- `POST /api/matrices/process` - Matris kaydetme ve tüm algoritmaları çalıştırma
- `POST /api/matrices/recalculate` - Seçili algoritmaları yeniden hesaplama

## Kurulum

### Gereksinimler
- Go 1.21+
- SQLite3

### Bağımlılıklar
```bash
go mod tidy
```

### Çalıştırma
```bash
go run .
```

Uygulama `http://localhost:3000` adresinde çalışacaktır.

## Veritabanı Yapısı

### matrix_records Tablosu
```sql
CREATE TABLE matrix_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    group_name TEXT,                    -- Grup adı (dosya adı)
    matrix_binary TEXT NOT NULL,        -- Binary matris gösterimi
    matrix_hex TEXT NOT NULL,           -- Hex formatında matris
    ham_xor_count INTEGER NOT NULL,     -- Hamming XOR sayısı
    smallest_xor INTEGER,               -- En küçük XOR değeri
    boyar_xor_count INTEGER,            -- Boyar algoritması XOR sayısı
    boyar_depth INTEGER,                -- Boyar algoritması derinlik
    boyar_program TEXT,                 -- Boyar algoritması programı (JSON)
    paar_xor_count INTEGER,             -- Paar algoritması XOR sayısı
    paar_program TEXT,                  -- Paar algoritması programı (JSON)
    slp_xor_count INTEGER,              -- SLP algoritması XOR sayısı
    slp_program TEXT,                   -- SLP algoritması programı (JSON)
    matrix_hash TEXT UNIQUE NOT NULL,   -- Matris hash'i (tekrar önleme)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Kullanım

### Web Arayüzü
1. Tarayıcıda `http://localhost:3000` adresini açın
2. **Matrisler** sekmesinde mevcut matrisleri görüntüleyin
3. **Yeni Matris** sekmesinde yeni matris ekleyin
4. Matris detaylarını görüntülemek için "Detay" butonuna tıklayın
5. "Yeniden Hesapla" butonu ile algoritmaları tekrar çalıştırın

### API Kullanımı

#### Yeni Matris Ekleme
```bash
curl -X POST http://localhost:3000/api/matrices/process \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Matrisi",
    "matrix": [
      ["1", "0", "1"],
      ["0", "1", "0"],
      ["1", "1", "1"]
    ]
  }'
```

#### Matris Listesi
```bash
curl "http://localhost:3000/api/matrices?page=1&limit=10&title=test"
```

#### Yeniden Hesaplama
```bash
curl -X POST http://localhost:3000/api/matrices/recalculate \
  -H "Content-Type: application/json" \
  -d '{
    "matrix_id": 1,
    "algorithms": ["boyar", "paar", "slp"]
  }'
```

### Test Verilerini İçe Aktarma
Test verilerini veritabanına aktarmak için:
```bash
go run test_import.go
```

## Performans Özellikleri

- **Hash Tabanlı Tekrar Önleme**: Aynı matris tekrar hesaplanmaz
- **Sayfalama**: Büyük veri setleri için performanslı listeleme
- **İndeksler**: Hızlı arama için veritabanı indeksleri
- **Asenkron İşlemler**: Web arayüzünde non-blocking işlemler

## Matris Formatları

### JSON Format (API)
```json
[
  ["1", "0", "1"],
  ["0", "1", "0"],
  ["1", "1", "1"]
]
```

### Binary Format (Veritabanı)
```
[1 0 1]
[0 1 0]
[1 1 1]
```

### Hex Format (Veritabanı)
```
5,2,7
```

## Algoritmalar

### Boyar SLP
- Derinlik sınırlı optimizasyon
- XOR sayısı ve derinlik bilgisi
- Program çıktısı

### Paar Algoritması
- Hamming ağırlığı tabanlı optimizasyon
- XOR sayısı optimizasyonu
- Program çıktısı

### SLP Heuristic
- Heuristik tabanlı optimizasyon
- Hızlı hesaplama
- Program çıktısı

## Dosya Yapısı

```
app/
├── main.go              # Ana uygulama ve algoritmalar
├── database.go          # Veritabanı işlemleri
├── api_handlers.go      # API handler'ları
├── test_import.go       # Test verisi import scripti
├── go.mod              # Go modül dosyası
├── matrices.db         # SQLite veritabanı (otomatik oluşur)
└── web/
    ├── index.html      # Ana web arayüzü
    ├── app.js          # JavaScript kodları
    └── test_data.txt   # Test verileri
```

## Geliştirme

### Yeni Algoritma Ekleme
1. `main.go` dosyasına algoritma struct'ını ekleyin
2. `Solve` metodunu implement edin
3. `api_handlers.go` dosyasında yeni algoritma için case ekleyin
4. Veritabanı şemasını güncelleyin

### Yeni API Endpoint Ekleme
1. `api_handlers.go` dosyasına handler fonksiyonu ekleyin
2. `main.go` dosyasında route'u tanımlayın
3. Web arayüzünde gerekli JavaScript kodlarını ekleyin

## Lisans

Bu proje MIT lisansı altında lisanslanmıştır. 