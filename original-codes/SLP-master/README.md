# Shorter Linear SLPs for MDS Matrices

SLPs for our ToSC Volume 2017 Issue 4 paper 'Shorter Linear Straight-Line Programs for MDS Matrices'

- The makefile compiles the slp_heuristic C++ program.
- The makefile compiles for each matrix header in the 'paar_header' folder
  two executables: one for Paar1 algorithm and one for Paar2. Each executable prints
  the XOR count for this matrix and the corresponding SLP for implementing the matrix
  with the said XOR count.
- For each matrix there is a second file in the 'bp_format' folder
  that contains the correct format for the slp_heuristic executable. You can
  call the programs on, e.g. the M_4_4 matrix, as follows:
  $ ./slp_heuristic < bp_format/M_4_4.txt
  It will then output the corresponding XOR count and the SLP.
- All matrices and functions for converting them etc. are additionally in the code.sage
  file. This was also used to construct the above headers and text files. The SageMath
  code can be loaded into a SageMath REPL and then used to compute Paar1 or slp_heuristic
  XOR counts (for the latter the program need to be available as a compiled executable).

- The code in the slp_heuristic.cpp is based on the implementation of Boyar and Peralta,
  available at <http://www.imada.sdu.dk/~joan/xor/>. We adjusted the program such that it
  prints the corresponding SLP program.
- Three SLPs where provided by Andrea Visconti through personal communication. The code
  used to generate these SLPs is not available to us.
- Twenty-five SLPs were generated with code kindly provided by René Peralta, based on
  recent work of Joan Boyar, Magnus Gausdal Find, and him. This program is called LinOpt
  in this repository. Its source code is not yet public.

- For every matrix listed in the paper there is a corresponding file in the
  'slp_implementations' directory. These files contain a comment with what program the
  SLP was computed, the XOR count and the SLP itself.
- All SLPs are also contained in the 'check_implementations.sage' file, which can be
  loaded in the SageMath REPL and all SLPs can be verified with the 'run_tests()'
  function.
  
---------------------------------------------------------------------------------------------------------

ToSC Cilt 2017 için SLP'ler Sayı 4 bildiri 'MDS Matrisleri için Daha Kısa Doğrusal Düz Çizgi Programları'

- makefile, slp_heuristic C ++ programını derler.

- makefile 'paar_header' klasöründeki her matris başlığı için iki çalıştırılabilir dosya derler: Biri Paar1 algoritması için ve diğeri Paar2 için.
Her çalıştırılabilir dosya, bu matris için XOR sayısını ve matrisi söz konusu XOR sayısıyla uygulamak için ilgili SLP'yi yazdırır.

- Her matris için 'bp_format' klasöründe slp_heuristic çalıştırılabilir dosyası için doğru formatı içeren ikinci bir dosya bulunur.
Programları örneğin; M_4_4 matrisi için aşağıdaki gibi çalıştırabilirsiniz:
 $ ./slp_heuristic <bp_format / M_4_4.txt
Daha sonra ilgili XOR sayısını ve SLP'yi verecektir.

- Tüm matrisler ve bunların dönüştürülmesi vb. işlemler code.sage dosyasına da eklendi.
Bu aynı zamanda yukarıdaki başlıkları ve metin dosyalarını oluşturmak için de kullanıldı.
SageMath kodu, SageMath REPL'e yüklenebilir ve Paar1 veya slp_heuristic XOR sayısını hesaplamak için kullanılabilir.(ikincisi için programın derlenmiş bir yürütülebilir dosya olması gerekir).

- slp_heuristic.cpp içindeki kod <http://www.imada.sdu.dk/~joan/xor/> adresinde bulunan Boyar ve Peralta'nın uygulanmasına dayanmaktadır. Program, ilgili SLP programını yazdıracak şekilde ayarlanmıştır.

- Andrea Visconti tarafından kişisel iletişim yoluyla sağlanan üç SLP bulunmaktadır. Bu SLP'leri oluşturmak için kullanılan kod bizim için mevcut değil.

- Joan Boyar, Magnus Gausdal Find ve onun son çalışmalarına dayanarak, René Peralta tarafından sağlanan kod ile yirmi beş SLP üretildi.
Bu program bu depoda LinOpt olarak adlandırılır. Kaynak kodu henüz herkese açık değil.

- Makalede listelenen her matris için 'slp_implementations' dizininde ilgili bir dosya var. Bu dosyalar SLP'nin hangi programda hesaplandığı, XOR sayısının ve SLP'nin kendisiyle ilgili bir yorum içerir.

- Tüm SLP'ler, SageMath REPL'e yüklenebilen ve tüm SLP'lerin 'run_tests ()' işleviyle doğrulanabilen 'check_implementations.sage' dosyasında da bulunur.
