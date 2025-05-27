# XOR Optimization - Configuration System

Bu proje artık otomatik data import işlemleri için kapsamlı bir config sistemi içeriyor.

## Config Dosyası

Uygulama başlatıldığında `config.json` dosyasını arar. Dosya yoksa varsayılan ayarlarla otomatik oluşturur.

### Config Yapısı

```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres", 
    "password": "password",
    "dbname": "xor_optimization",
    "sslmode": "disable"
  },
  "import": {
    "enabled": true,
    "data_directory": "./matrices-data",
    "file_extensions": [".txt", ".csv", ".json"],
    "max_file_size_mb": 500,
    "process_on_start": true,
    "watch_directory": false,
    "skip_existing": true,
    "batch_size": 10,
    "auto_calculate": true,
    "algorithms": ["boyar", "paar", "slp"]
  },
  "server": {
    "port": ":3000",
    "host": "localhost",
    "enable_cors": true,
    "log_level": "info",
    "static_dir": "./web"
  }
}
```

## Import Ayarları

### `enabled` (bool)
- Otomatik import işlemini etkinleştirir/devre dışı bırakır
- Varsayılan: `true`

### `data_directory` (string)
- Matrix dosyalarının bulunduğu dizin
- Varsayılan: `"./matrices-data"`

### `file_extensions` ([]string)
- Desteklenen dosya uzantıları
- Varsayılan: `[".txt", ".csv", ".json"]`

### `max_file_size_mb` (int64)
- Maksimum dosya boyutu (MB)
- Varsayılan: `500`

### `process_on_start` (bool)
- Uygulama başlatıldığında otomatik import yapılsın mı
- Varsayılan: `true`

### `watch_directory` (bool)
- Dizini sürekli izleyip yeni dosyaları otomatik import etsin mi
- Varsayılan: `false` (henüz implement edilmedi)

### `skip_existing` (bool)
- Zaten var olan matrisleri atlasın mı
- Varsayılan: `true`

### `batch_size` (int)
- Toplu işlem boyutu
- Varsayılan: `10`

### `auto_calculate` (bool)
- Import edilen matrisler için algoritmaları otomatik hesaplasın mı
- Varsayılan: `true`

### `algorithms` ([]string)
- Otomatik hesaplanacak algoritmalar
- Seçenekler: `["boyar", "paar", "slp"]`
- Varsayılan: `["boyar", "paar", "slp"]`

## Desteklenen Dosya Formatları

### 1. Text Format (.txt)
```
1 0 1
0 1 0
1 1 1
```

### 2. CSV Format (.csv)
```
1,0,1
0,1,0
1,1,1
```

### 3. JSON Format (.json)
```json
{
  "matrix": [
    ["1", "0", "1"],
    ["0", "1", "0"],
    ["1", "1", "1"]
  ]
}
```

veya

```json
{
  "matrices": [
    [
      ["1", "0", "1"],
      ["0", "1", "0"]
    ],
    [
      ["1", "1"],
      ["0", "1"]
    ]
  ]
}
```

## API Endpoints

### Config Endpoints

#### GET /api/config
Mevcut konfigürasyonu döndürür.

#### POST /api/config/import
Manuel import işlemini tetikler.

```bash
curl -X POST http://localhost:3000/api/config/import
```

## Kullanım

### 1. Config Dosyasını Düzenleme
```bash
nano config.json
```

### 2. Data Dizinini Ayarlama
```json
{
  "import": {
    "data_directory": "/path/to/your/matrices"
  }
}
```

### 3. Sadece Belirli Algoritmaları Çalıştırma
```json
{
  "import": {
    "auto_calculate": true,
    "algorithms": ["paar"]
  }
}
```

### 4. Büyük Dosyaları İşleme
```json
{
  "import": {
    "max_file_size_mb": 1000,
    "batch_size": 5
  }
}
```

## Örnek Kullanım Senaryoları

### Senaryo 1: Sadece Import, Hesaplama Yok
```json
{
  "import": {
    "enabled": true,
    "auto_calculate": false
  }
}
```

### Senaryo 2: Sadece PAAR Algoritması
```json
{
  "import": {
    "auto_calculate": true,
    "algorithms": ["paar"]
  }
}
```

### Senaryo 3: Farklı Data Dizini
```json
{
  "import": {
    "data_directory": "/home/user/research/matrices",
    "file_extensions": [".txt"]
  }
}
```

## Log Mesajları

Import işlemi sırasında detaylı log mesajları görürsünüz:

```
[INFO] Config yüklendi
[INFO] Otomatik import başlatılıyor: ./matrices-data
[INFO] Dosya import ediliyor: ./matrices-data/test.txt
[INFO] Matrix kaydedildi: test_matrix_1 (ID: 123)
[INFO] Algoritma hesaplandı: test_matrix_1 -> boyar
[INFO] Otomatik import tamamlandı
```

## Hata Durumları

- **Dosya çok büyük**: `max_file_size_mb` ayarını artırın
- **Desteklenmeyen format**: `file_extensions` listesini kontrol edin
- **Veritabanı hatası**: Database ayarlarını kontrol edin
- **Dizin bulunamadı**: `data_directory` yolunu kontrol edin

## Performans İpuçları

1. **Büyük dosyalar için**: `batch_size` değerini düşürün
2. **Hızlı import için**: `auto_calculate` false yapın, sonra manuel hesaplatın
3. **Bellek tasarrufu için**: `max_file_size_mb` sınırını ayarlayın
4. **Sadece yeni dosyalar için**: `skip_existing` true bırakın 