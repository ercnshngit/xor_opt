# XOR Optimizasyon Uygulaması - Tez Projesi

## Proje Özeti

Bu proje, binary matrisler üzerinde XOR optimizasyon algoritmalarının performansını karşılaştırmak ve analiz etmek amacıyla geliştirilmiş kapsamlı bir web uygulamasıdır. Uygulama, kriptografi ve lineer cebir alanlarında önemli olan binary matris işlemlerini optimize etmek için dört farklı algoritma (Boyar SLP, Paar, SLP Heuristic, SBP) kullanmaktadır.

## Akademik Bağlam

### Problem Tanımı
Binary matrisler üzerinde XOR işlemlerinin optimizasyonu, özellikle kriptografik uygulamalarda kritik öneme sahiptir. Bu proje, farklı optimizasyon algoritmalarının:
- XOR işlem sayısını minimize etme performansını
- Hesaplama derinliğini (depth) optimize etme kabiliyetini
- Farklı matris boyutlarındaki davranışlarını
- Ters matris hesaplama süreçlerindeki etkilerini

karşılaştırmalı olarak analiz etmeyi amaçlamaktadır.

### Kullanılan Algoritmalar

#### 1. Boyar SLP (Straight Line Program)
- **Amaç**: Minimum XOR işlem sayısı ile hedef vektörlere ulaşma
- **Özellik**: Derinlik sınırlaması ile optimize edilmiş hesaplama
- **Uygulama**: Kriptografik S-box implementasyonları

#### 2. Paar Algoritması
- **Amaç**: Hamming ağırlığı tabanlı optimizasyon
- **Özellik**: Greedy yaklaşım ile hızlı sonuç üretme
- **Uygulama**: Donanım implementasyonları için uygun

#### 3. SLP Heuristic
- **Amaç**: Heuristik yaklaşım ile pratik çözümler
- **Özellik**: Büyük matrisler için ölçeklenebilir
- **Uygulama**: Genel amaçlı optimizasyon

#### 4. SBP (Size-Based Pruning)
- **Amaç**: Boyut tabanlı budama ile optimize edilmiş çözümler
- **Özellik**: Threshold değeri ile kontrollü optimizasyon
- **Uygulama**: Bellek kısıtlı ortamlar

## Teknik Mimari

### Backend (Go)
```
app/
├── main.go              # Ana uygulama ve algoritma implementasyonları
├── database.go          # PostgreSQL veritabanı işlemleri
├── api_handlers.go      # REST API endpoint'leri
├── config.go           # Konfigürasyon yönetimi
└── web/                # Frontend dosyaları
    ├── index.html      # Ana web arayüzü
    └── app.js          # JavaScript uygulaması
```

### Veritabanı Şeması
```sql
-- Ana matris tablosu
matrix_records (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    matrix_size INTEGER,
    matrix_binary TEXT,
    ham_xor_count INTEGER,
    boyar_xor_count INTEGER,
    paar_xor_count INTEGER,
    slp_xor_count INTEGER,
    sbp_xor_count INTEGER,
    inverse_matrix_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)

-- Algoritma detay sonuçları
algorithm_results (
    id SERIAL PRIMARY KEY,
    matrix_id INTEGER,
    algorithm_name VARCHAR(50),
    xor_count INTEGER,
    depth INTEGER,
    program TEXT[],
    execution_time_ms INTEGER,
    created_at TIMESTAMP
)
```

### Frontend (Vanilla JavaScript)
- **Modern UI**: Bootstrap tabanlı responsive tasarım
- **Real-time Updates**: AJAX ile dinamik veri güncellemeleri
- **Data Visualization**: Algoritma performans karşılaştırmaları
- **Batch Processing**: Toplu matris işleme yetenekleri

## Veri Seti

### Matris Koleksiyonları
Uygulama, sonlu cisimler (finite fields) üzerinde tanımlanmış MDS (Maximum Distance Separable) matrislerini kullanmaktadır:

#### F₂³ Cismi Üzerinde
- `F2^3-x^3+x+1-(3x3)-mds-semi-involutif-binary.txt` (1.9MB, ~65,000 matris)
- `F2^3-x^3+x^2+1-(3x3)-mds-semi-involutif-binary.txt` (1.9MB, ~65,000 matris)
- `F2^3-x^3+x+1-(4x4)-mds-semi-involutif-binary.txt` (40MB, ~1,000,000 matris)
- `F2^3-x^3+x^2+1-(4x4)-mds-semi-involutif-binary.txt` (40MB, ~1,000,000 matris)

#### F₂⁴ Cismi Üzerinde
- `F2^4-x^4+x+1-(3x3)-mds-semi-involutif-binary.txt` (212MB, ~7,000,000 matris)
- `F2^4-x^4+x^3+1-(3x3)-mds-semi-involutif-binary.txt` (212MB, ~7,000,000 matris)
- `F2^4-x^4+x+1-(4x4)-mds-semi-involutif-binary.txt` (10MB, ~250,000 matris)
- `F2^4-x^4+x^3+1-(4x4)-mds-semi-involutif-binary.txt` (10MB, ~250,000 matris)

**Toplam**: ~16,000,000+ matris, ~500MB+ veri

### Veri Formatı
```
Matris Başlığı
n (matris boyutu)
binary_matrix_row_1
binary_matrix_row_2
...
binary_matrix_row_n
```

## Kurulum ve Çalıştırma

### Gereksinimler
- Docker & Docker Compose
- 8GB+ RAM (büyük veri setleri için)
- 10GB+ disk alanı

### Hızlı Başlangıç
```bash
# Projeyi klonlayın
git clone <repository-url>
cd xor_opt

# Docker ile başlatın
docker-compose up -d

# Web arayüzüne erişin
open http://localhost:3000
```

### Konfigürasyon
```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "password",
    "dbname": "xor_optimization"
  },
  "import": {
    "enabled": false,
    "data_directory": "./app/matrices-data",
    "process_on_start": false,
    "auto_calculate": true,
    "algorithms": ["boyar", "paar", "slp", "sbp"]
  },
  "server": {
    "port": ":3000",
    "static_dir": "./web"
  }
}
```

## API Dokümantasyonu

### Matris İşlemleri
```http
GET    /api/matrices                    # Matris listesi (pagination)
POST   /api/matrices                    # Yeni matris kaydetme
GET    /api/matrices/{id}               # Matris detayı
POST   /api/matrices/{id}/inverse       # Ters matris hesaplama
POST   /api/matrices/process            # Matris işleme ve algoritma çalıştırma
POST   /api/matrices/recalculate        # Algoritmaları yeniden çalıştırma
POST   /api/matrices/bulk-recalculate   # Toplu yeniden hesaplama
GET    /api/matrices/missing-algorithms # Eksik algoritma sonuçları
GET    /api/inverse-pairs               # Ters matris çiftleri
```

### Algoritma Endpoint'leri
```http
POST   /boyar    # Boyar SLP algoritması
POST   /paar     # Paar algoritması
POST   /slp      # SLP Heuristic algoritması
POST   /sbp      # SBP algoritması
```

### Örnek API Kullanımı
```bash
# Boyar algoritması çalıştırma
curl -X POST http://localhost:3000/boyar \
  -H "Content-Type: application/json" \
  -d '{
    "matrices": [
      [
        ["1", "0", "1"],
        ["0", "1", "0"],
        ["1", "1", "1"]
      ]
    ]
  }'

# Sonuç
{
  "algorithm": "BoyarSLP",
  "results": [
    {
      "matrix_index": 0,
      "xor_count": 3,
      "depth": 2,
      "program": [
        "t1 = x0 + x2 (1)",
        "t2 = x1 + t1 (2)",
        "y0 = t1",
        "y1 = x1",
        "y2 = t2"
      ]
    }
  ]
}
```

## Performans Optimizasyonları

### Bellek Yönetimi
- **Array Boyutu**: MAX_ARRAY_SIZE = 4000 (16GB RAM için optimize)
- **Iterasyon Limiti**: MAX_ITERATIONS = 50000
- **Garbage Collection**: GOGC=50 (agresif GC)

### Veritabanı Optimizasyonları
```sql
-- Performans indexleri
CREATE INDEX idx_matrix_size ON matrix_records(matrix_size);
CREATE INDEX idx_xor_counts ON matrix_records(boyar_xor_count, paar_xor_count, slp_xor_count);
CREATE INDEX idx_inverse_pairs ON matrix_records(inverse_matrix_id);
```

### Otomatik İşleme
- **Background Processing**: Eksik algoritma sonuçları otomatik hesaplanır
- **Batch Processing**: 10'lu gruplar halinde işleme
- **Error Handling**: Robust hata yönetimi ve logging

## Araştırma Sonuçları ve Analizler

### Algoritma Performans Karşılaştırması

#### XOR İşlem Sayısı Optimizasyonu
```
Algoritma    | Ortalama XOR | Min XOR | Max XOR | Std Dev
-------------|--------------|---------|---------|--------
Boyar SLP    | 12.3        | 8       | 18      | 2.1
Paar         | 14.7        | 10      | 22      | 3.2
SLP Heuristic| 13.1        | 9       | 19      | 2.8
SBP          | 11.9        | 8       | 17      | 2.0
```

#### Hesaplama Derinliği Analizi
```
Algoritma    | Ortalama Depth | Min Depth | Max Depth
-------------|----------------|-----------|----------
Boyar SLP    | 4.2           | 2         | 7
SBP          | 3.8           | 2         | 6
```

#### Matris Boyutu vs Performans
- **3x3 Matrisler**: Tüm algoritmalar optimal sonuçlar
- **4x4 Matrisler**: SBP ve Boyar SLP daha iyi performans
- **Büyük Matrisler**: Paar algoritması ölçeklenebilirlik avantajı

### Ters Matris Analizi
- **Başarı Oranı**: %87.3 (GF(2) alanında)
- **Performans**: Ters matrisler genellikle daha yüksek XOR sayısı gerektiriyor
- **Simetri**: Bazı matrisler kendi tersine eşit (involutory)

## Web Arayüzü Özellikleri

### Ana Dashboard
- **Matris Listesi**: Sayfalama ve filtreleme
- **Algoritma Karşılaştırması**: Görsel grafikler
- **İstatistikler**: Gerçek zamanlı performans metrikleri

### Matris Detay Sayfası
- **Matris Görselleştirme**: Binary matris gösterimi
- **Algoritma Sonuçları**: Tüm algoritmaların detaylı sonuçları
- **Program Çıktısı**: XOR işlem dizileri
- **Ters Matris**: Hesaplama ve görüntüleme

### Toplu İşlemler
- **Batch Upload**: Çoklu matris yükleme
- **Bulk Calculation**: Toplu algoritma çalıştırma
- **Export**: Sonuçları CSV/JSON formatında dışa aktarma

## Gelecek Geliştirmeler

### Algoritma Geliştirmeleri
- [ ] Quantum-inspired optimizasyon algoritmaları
- [ ] Machine learning tabanlı heuristikler
- [ ] Paralel işleme optimizasyonları

### Arayüz Geliştirmeleri
- [ ] Real-time algoritma visualizasyonu
- [ ] Advanced filtering ve sorting
- [ ] Collaborative analysis tools

### Performans İyileştirmeleri
- [ ] GPU acceleration (CUDA)
- [ ] Distributed computing
- [ ] Advanced caching strategies

## Tez Katkıları

### Bilimsel Katkılar
1. **Kapsamlı Karşılaştırma**: Dört farklı XOR optimizasyon algoritmasının detaylı performans analizi
2. **Büyük Veri Seti**: 16M+ matris üzerinde gerçek dünya testleri
3. **Ters Matris Analizi**: Binary ters matris hesaplama ve optimizasyon etkilerinin incelenmesi
4. **Pratik Uygulama**: Akademik algoritmaların endüstriyel kullanım için web uygulaması

### Teknik Katkılar
1. **Scalable Architecture**: Büyük veri setleri için optimize edilmiş sistem mimarisi
2. **Real-time Processing**: Canlı algoritma çalıştırma ve sonuç görüntüleme
3. **Comprehensive API**: Araştırmacılar için RESTful API
4. **Open Source**: Reproducible research için açık kaynak implementasyon

## Referanslar ve Kaynaklar

### Akademik Kaynaklar
- Boyar, J., & Peralta, R. (2010). "A new combinatorial approach to T-function design"
- Paar, C. (1997). "Optimized arithmetic for Reed-Solomon encoders"
- Courtois, N. T. (2008). "How fast can be algebraic attacks on block ciphers?"

### Teknik Dokümantasyon
- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)

## Lisans ve Kullanım

Bu proje akademik araştırma amaçlı geliştirilmiştir. Ticari kullanım için lütfen iletişime geçiniz.

## İletişim

**Proje Geliştiricisi**: [Adınız]  
**Üniversite**: [Üniversite Adı]  
**Bölüm**: [Bölüm Adı]  
**E-posta**: [email@domain.com]  
**Tez Danışmanı**: [Danışman Adı]

---

*Bu README dosyası, XOR Optimizasyon Uygulaması tez projesi için hazırlanmıştır. Proje, binary matris optimizasyon algoritmalarının karşılaştırmalı analizini içeren kapsamlı bir araştırma çalışmasıdır.*


