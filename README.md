# XOR Optimizasyon Uygulaması

Bu uygulama, binary matrisler üzerinde XOR optimizasyon algoritmalarını (Boyar SLP, Paar, SLP Heuristic) çalıştıran bir web uygulamasıdır.

## Özellikler

- **Web Arayüzü**: Modern ve kullanıcı dostu web arayüzü
- **Algoritma Desteği**: Boyar SLP, Paar ve SLP Heuristic algoritmaları
- **Veritabanı**: PostgreSQL ile matris verilerinin saklanması
- **Otomatik Import**: Uygulama başlatıldığında matrices-data klasöründeki dosyaların otomatik olarak veritabanına import edilmesi
- **Filtreleme**: XOR sayılarına göre filtreleme
- **Pagination**: Büyük veri setleri için sayfalama
- **Toplu Yükleme**: Birden fazla matrisin aynı anda yüklenmesi

## Kurulum

### Gereksinimler

- Docker
- Docker Compose

### Kurulum Adımları

1. **Projeyi klonlayın:**
   ```bash
   git clone <repository-url>
   cd xor_opt
   ```

2. **Docker Compose ile başlatın:**
   ```bash
   docker-compose up -d
   ```

3. **Uygulamaya erişin:**
   - Web Arayüzü: http://localhost:3000
   - PostgreSQL: localhost:5432

### İlk Başlatma

Uygulama ilk kez başlatıldığında:

1. PostgreSQL veritabanı otomatik olarak oluşturulur
2. Gerekli tablolar ve indexler oluşturulur
3. `matrices-data` klasöründeki 4 dosya otomatik olarak taranır
4. Veritabanında eksik matrisler varsa, dosyalardan otomatik import edilir
5. Bu işlem background'da çalışır ve uygulamanın başlamasını engellemez

### Matrices-Data Dosyaları

Aşağıdaki dosyalar otomatik olarak import edilir:
- `F2^3-x^3+x^2+1-(3x3)-mds-semi-involutif-binary.txt`
- `F2^3-x^3+x+1-(3x3)-mds-semi-involutif-binary.txt`
- `F2^4-x^4+x^3+1-(3x3)-mds-semi-involutif-binary.txt`
- `F2^4-x^4+x+1-(3x3)-mds-semi-involutif-binary.txt`

## Kullanım

### Web Arayüzü

1. **Matris Listesi**: Ana sayfada tüm matrisler listelenir
2. **Filtreleme**: XOR sayılarına göre filtreleme yapabilirsiniz
3. **Matris Detayı**: Herhangi bir matrise tıklayarak detaylarını görüntüleyebilirsiniz
4. **Yeni Matris**: Yeni matris ekleyebilirsiniz
5. **Toplu Yükleme**: Birden fazla matrisin aynı anda yüklenmesi

### API Endpoints

#### Matris İşlemleri
- `GET /api/matrices` - Matris listesi (pagination ve filtreleme ile)
- `POST /api/matrices` - Yeni matris kaydetme
- `GET /api/matrices/{id}` - Matris detayı
- `POST /api/matrices/process` - Matris kaydetme ve algoritmaları çalıştırma
- `POST /api/matrices/recalculate` - Algoritmaları yeniden çalıştırma

#### Algoritma Endpoints
- `POST /boyar` - Boyar SLP algoritması
- `POST /paar` - Paar algoritması
- `POST /slp` - SLP Heuristic algoritması

### Örnek API Kullanımı

```bash
# Matris listesi
curl "http://localhost:3000/api/matrices?page=1&limit=10"

# Yeni matris kaydetme
curl -X POST http://localhost:3000/api/matrices \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Matrix",
    "matrix": [
      ["1", "0", "1"],
      ["0", "1", "0"],
      ["1", "1", "1"]
    ]
  }'

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
```

## Geliştirme

### Yerel Geliştirme

```bash
# Veritabanını başlatın
docker-compose up -d db

# Go uygulamasını yerel olarak çalıştırın
cd app
go mod download
go run .
```

### Veritabanı Bağlantısı

Uygulama aşağıdaki environment variable'ları kullanır:

- `DB_HOST`: PostgreSQL host (varsayılan: localhost)
- `DB_PORT`: PostgreSQL port (varsayılan: 5432)
- `DB_NAME`: Veritabanı adı (varsayılan: xor_opt)
- `DB_USER`: Kullanıcı adı (varsayılan: xor_user)
- `DB_PASSWORD`: Şifre (varsayılan: xor_password)
- `DB_SSLMODE`: SSL modu (varsayılan: disable)
- `MATRICES_DATA_PATH`: Matris dosyalarının yolu (varsayılan: ./matrices-data)

## Durdurma

```bash
docker-compose down
```

Veritabanı verilerini de silmek için:
```bash
docker-compose down -v
```

## Loglar

```bash
# Tüm servislerin logları
docker-compose logs -f

# Sadece uygulama logları
docker-compose logs -f app

# Sadece veritabanı logları
docker-compose logs -f db
```

## Sorun Giderme

### Veritabanı Bağlantı Sorunu
- PostgreSQL container'ının çalıştığından emin olun: `docker-compose ps`
- Logları kontrol edin: `docker-compose logs db`

### Import Sorunu
- Matrices-data dosyalarının mevcut olduğundan emin olun
- Uygulama loglarını kontrol edin: `docker-compose logs app`

### Port Çakışması
- 3000 veya 5432 portları kullanımda ise docker-compose.yml dosyasında port numaralarını değiştirin 



all 3x3 makalesine atif yapip involutifleri nasil buldugumuzu yazalim
kanıtı yaz  65 sayfa 4.bolum hocanin doktora tezindeki(atifi hocanin yolladigi makaleden yapacagim, ) involutif elde edilmesini yazcaz. about semi involutary pdf komple turkceye cevrilip yazilacak

2.bolumde magmadan bahset(cikis yeri vs) magma nasil kullanilir ufak ornekler,   slp ve paar den bahset

4bbolumde 3x3 ve 4x4 ornekler sunulacak 4x4 icin aciklama yazilacak  xor sayisi lseq 80 seklinde bakildi 2 uzeri 4 ornekleri


5.bolum literatur karsilastirmasi ve sonuclar

involutif mds -> ayni cismin alfa  uzeri 1 ile carpilmis hali -> cikan matris ->


semilerin tersi kendisine esit olmadigi icin sifreleme desifreleme esit olmadigi icin yan kanal saldirilarina acik olabilir bu sonuc bolumunde belirtilmeli



az xor a sahip olanlarin terslerinin de xor maliyetlerine bakmak lazim

programi nasil yazdigini anlat json format

