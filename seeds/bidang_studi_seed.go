package seeds

import (
	"fmt"
	"pustaka-backend/models"
	"gorm.io/gorm"
)

// BidangStudiSeeder seeds the bidang_studi table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func BidangStudiSeeder(db *gorm.DB) error {
	fmt.Println("📝 Seeding bidang_studi table...")

	// Helper function to create string pointer
	// strPtr := func(s string) *string {
	// 	return &s
	// }

	bidangStudiData := []models.BidangStudi{
  {
    Code: "A.BYA",
    Name: "AKUNTANSI BIAYA",
  },
  {
    Code: "A.KEU",
    Name: "AKUNTANSI KEUANGAN",
  },
  {
    Code: "A.PRB",
    Name: "AKUNTANSI PERBANKAN",
  },
  {
    Code: "LPEA",
    Name: "EKO BISNIS DAN ADM UMUM LAYANANAN PERBANKAN",
  },
  {
    Code: "ANTRO",
    Name: "ANTROPOLOGI",
  },
  {
    Code: "BINDO",
    Name: "BAHASA INDONESIA",
  },
  {
    Code: "BIGTK",
    Name: "BAHASA INGGRIS TEHNIK",
  },
  {
    Code: "B.IND",
    Name: "BAHASA  INDONESIA",
  },
  {
    Code: "B.ING",
    Name: "BAHASA INGGRIS",
  },
  {
    Code: "B.SID",
    Name: "BAHASA DAN SASTRA INDONESIA",
  },
  {
    Code: "BIO",
    Name: "BIOLOGI",
  },
  {
    Code: "BIOPA",
    Name: "BIOLOGI ( KEAHLIAN KESEHATAN )",
  },
  {
    Code: "BK",
    Name: "BIMBINGAN KONSELING",
  },
  {
    Code: "EKO",
    Name: "EKONOMI",
  },
  {
    Code: "EKOPS",
    Name: "EKONOMI KOPERASI & PENGELOLAAN SDM ( PENGETAHUAN SOSIAL )",
  },
  {
    Code: "F.PA",
    Name: "FISIKA ( PENGETAHUAN ALAM )",
  },
  {
    Code: "FIS",
    Name: "FISIKA",
  },
  {
    Code: "GAB",
    Name: "GABUNGAN",
  },
  {
    Code: "GEO",
    Name: "GEOGRAFI",
  },
  {
    Code: "GEOPS",
    Name: "GEOGRAFI ( PENGETAHUAN SOSIAL )",
  },
  {
    Code: "IPA",
    Name: "ILMU PENGETAHUAN ALAM",
  },
  {
    Code: "IPAT",
    Name: "ILMU PENGETAHUAN ALAM TERPADU",
  },
  {
    Code: "IPS",
    Name: "ILMU PENGETAHUAN SOSIAL",
  },
  {
    Code: "IPST",
    Name: "ILMU PENGETAHUAN SOSIAL TERPADU",
  },
  {
    Code: "KIM",
    Name: "KIMIA",
  },
  {
    Code: "KOMPT",
    Name: "KOMPUTER",
  },
  {
    Code: "KTK",
    Name: "PENDIDIKAN KESENIAN DAN KETRAMPILAN",
  },
  {
    Code: "SSKRS",
    Name: "SISTEM SASIS KENDARAAN  RINGAN SETAHUN",
  },
  {
    Code: "MTK",
    Name: "MATEMATIKA",
  },
  {
    Code: "MTKA",
    Name: "MATEMATIKA IPA",
  },
  {
    Code: "MTKS",
    Name: "MATEMATIKA IPS",
  },
  {
    Code: "MTKT",
    Name: "MATEMATIKA TEHNIK",
  },
  {
    Code: "MUK",
    Name: "MEMBUKA USAHA KECIL",
  },
  {
    Code: "P.P",
    Name: "PELAYANAN PRIMA",
  },
  {
    Code: "PKMP",
    Name: "PENGELOLAAN KEARSIPAN MP",
  },
  {
    Code: "PJKES",
    Name: "PENDIDIKAN JASMANI DAN KESEHATAN",
  },
  {
    Code: "PJS",
    Name: "PENDIDIKAN JASMANI",
  },
  {
    Code: "PKES",
    Name: "PENDIDIKAN KESENIAN",
  },
  {
    Code: "PKPS",
    Name: "PEND. KEWARGANEGARAAN & ILMU PENGETAHUAN SOSIAL",
  },
  {
    Code: "PKS",
    Name: "PENDIDIKAN KEWARGANEGARAAN & SEJARAH",
  },
  {
    Code: "PPKN",
    Name: "PENDIDIKAN KEWARGANEGARAAN",
  },
  {
    Code: "S.AKT",
    Name: "SIKLUS AKUNTANSI",
  },
  {
    Code: "PRMP",
    Name: "PENGELOLAAN RAPAT MP",
  },
  {
    Code: "SEJ1",
    Name: "SEJARAH",
  },
  {
    Code: "SEJPA",
    Name: "SEJARAH ( IPA)",
  },
  {
    Code: "SEJPS",
    Name: "SEJARAH ( PENGETAHUAN SOSIAL )",
  },
  {
    Code: "SMART",
    Name: "MAJALAH",
  },
  {
    Code: "SNK",
    Name: "SURAT NIAGA DAN KEARSIPAN",
  },
  {
    Code: "SOS",
    Name: "SOSIOLOGI",
  },
  {
    Code: "T.I.K",
    Name: "TEKHNOLOGI INFORMASI & KOMPUTER",
  },
  {
    Code: "TN",
    Name: "TATA NEGARA",
  },
  {
    Code: "X.AKT",
    Name: "AKUNTANSI",
  },
  {
    Code: "X.BIO",
    Name: "BIOLOGI X",
  },
  {
    Code: "X.EKO",
    Name: "EKONOMI X",
  },
  {
    Code: "X.FIS",
    Name: "FISIKA X",
  },
  {
    Code: "X.GEO",
    Name: "GEOGRAFI X",
  },
  {
    Code: "XINDL",
    Name: "BAHASA INDONESIA (CL) X",
  },
  {
    Code: "XINGL",
    Name: "BAHASA INGGRIS (CL) X",
  },
  {
    Code: "XIPAL",
    Name: "IPA (CL) X",
  },
  {
    Code: "XIPSL",
    Name: "IPS (CL) X",
  },
  {
    Code: "X.KIM",
    Name: "KIMIA X",
  },
  {
    Code: "X.KOM",
    Name: "KOMPUTER X",
  },
  {
    Code: "X.KTK",
    Name: "KTK X",
  },
  {
    Code: "X.KWU",
    Name: "KEWIRAUSAHAAN X",
  },
  {
    Code: "XMTKL",
    Name: "MATEMATIKA (CL) X",
  },
  {
    Code: "XMTKA",
    Name: "MATEMATIKA IPA X",
  },
  {
    Code: "X.PAI",
    Name: "PENDIDIKAN AGAMA ISLAM X",
  },
  {
    Code: "X.PJS",
    Name: "PENJAS X",
  },
  {
    Code: "X.SEJ",
    Name: "SEJARAH X",
  },
  {
    Code: "X.SOS",
    Name: "SOSIOLOGI X",
  },
  {
    Code: "X.TN",
    Name: "TATA NEGARA X",
  },
  {
    Code: "XIGT",
    Name: "BAHASA INGGRIS TEHNIK X",
  },
  {
    Code: "XMTT",
    Name: "MATEMATIKA TEHNIK X",
  },
  {
    Code: "XPPKN",
    Name: "PENDIDIKAN PANCASILA & KEWARGANEGARAAN (CL) X",
  },
  {
    Code: "XSBUD",
    Name: "SENI BUDAYA X",
  },
  {
    Code: "XTIK",
    Name: "TEKNOLOGI INFORMASI & KOMPUTER X",
  },
  {
    Code: "PK",
    Name: "PENGALAMANKU",
  },
  {
    Code: "BHT",
    Name: "BENDA, HEWAN & TANAMAN",
  },
  {
    Code: "PHL",
    Name: "PAHLAWANKU",
  },
  {
    Code: "CC",
    Name: "CITA-CITAKU",
  },
  {
    Code: "MSB",
    Name: "MAKANAN SEHAT & BERGIZI",
  },
  {
    Code: "PRA",
    Name: "PRAKARYA",
  },
  {
    Code: "FTEKR",
    Name: "FISIKA ( KEAHLIAN TEKNOLOGI & REKAYASA )",
  },
  {
    Code: "F.KES",
    Name: "FISIKA ( KEAHLIAN KESEHATAN )",
  },
  {
    Code: "K.T&R",
    Name: "KIMIA ( KEAHLIAN TEKNOLOGI & REKAYASA)",
  },
  {
    Code: "PD",
    Name: "PEMROGRAMAN DASAR",
  },
  {
    Code: "PPAR",
    Name: "PENGANTAR PARIWISATA",
  },
  {
    Code: "KG",
    Name: "KEGEMARANKU",
  },
  {
    Code: "KL",
    Name: "KELUARGAKU",
  },
  {
    Code: "SBE",
    Name: "SELALU BERHEMAT ENERGI",
  },
  {
    Code: "BP",
    Name: "BERBAGI PEKERJAAN",
  },
  {
    Code: "HR",
    Name: "HIDUP RUKUN",
  },
  {
    Code: "TS",
    Name: "TUGASKU SEHARI-HARI",
  },
  {
    Code: "BLS",
    Name: "BENDA-BENDA DILINGKUNGAN SEKITAR",
  },
  {
    Code: "KDB",
    Name: "KERUKUNAN DALAM BERMASYARAKAT",
  },
  {
    Code: "BSBI",
    Name: "BANGGA SBG BANGSA INDONESIA",
  },
  {
    Code: "T1",
    Name: "TEMA 1",
  },
  {
    Code: "T3",
    Name: "TEMA 3",
  },
  {
    Code: "T5",
    Name: "TEMA 5",
  },
  {
    Code: "T7",
    Name: "TEMA 7",
  },
  {
    Code: "T9",
    Name: "TEMA 9",
  },
  {
    Code: "IT",
    Name: "IPA TERAPAN",
  },
  {
    Code: "AP.01",
    Name: "MEMAHAMI PRINSIP2 PENYELENGGARAAN ADM  PERKANTORAN",
  },
  {
    Code: "TSM08",
    Name: "M.PERBIKAN SISTEM HIDROLIK SPD MOTOR",
  },
  {
    Code: "AP.05",
    Name: "MENGOPERASIKAN APLIKASI PERANGKAT LUNAK",
  },
  {
    Code: "TSM09",
    Name: "MEMPERBAIKI SISTEM GAS BUANG",
  },
  {
    Code: "TSM10",
    Name: "MEMELIHARA BATERAI",
  },
  {
    Code: "AP.11",
    Name: "MENGELOLA SISTEM KEARSIPAN",
  },
  {
    Code: "AP.13",
    Name: "MEMPROSES PERJALANAN BISNIS",
  },
  {
    Code: "AP.15",
    Name: "MENGELOLA DANA KAS KECIL",
  },
  {
    Code: "AP.17",
    Name: "MENGELOLA DATA / INFORMASI DITEMPAT KERJA",
  },
  {
    Code: "AP.07",
    Name: "MENGELOLA PERALATAN KANTOR",
  },
  {
    Code: "P.01",
    Name: "FRONT OFFICE(AKOMODASI PERHOTELAN)",
  },
  {
    Code: "P.03",
    Name: "MELAKUKAN PENGELOLAAN USAHA JASA BOGA",
  },
  {
    Code: "P.05",
    Name: "MENGOLAH MAKANAN KONTINENTAL",
  },
  {
    Code: "AK.02",
    Name: "MELAKSANAKAN KOMUNIKASI BISNIS",
  },
  {
    Code: "AK.03",
    Name: "MENERAPKAN KESELAMATAN KESEHATAN KERJA  & LING HDP",
  },
  {
    Code: "AK.04",
    Name: "MENGELOLA DOKUMEN TRANSAKSI",
  },
  {
    Code: "AK.05",
    Name: "MEMPROSES DOKUMEN DANA KAS KECIL",
  },
  {
    Code: "AK.06",
    Name: "MEMPROSES DANA KAS DI BANK",
  },
  {
    Code: "AK.07",
    Name: "MEMPROSES ENTRI JURNAL",
  },
  {
    Code: "AK.08",
    Name: "MEMPROSES BUKU BESAR",
  },
  {
    Code: "AK.09",
    Name: "MENGELOLA KARTU PIUTANG",
  },
  {
    Code: "AK.10",
    Name: "MENGELOLA KARTU PERSEDIAAN",
  },
  {
    Code: "AK.11",
    Name: "MENGELOLA KARTU AKTIVA",
  },
  {
    Code: "AK.12",
    Name: "MENGELOLA KARTU UTANG",
  },
  {
    Code: "AK.13",
    Name: "MENYAJIKAN LAPORAN HARGA POKOK PRODUK",
  },
  {
    Code: "AK.14",
    Name: "MENYUSUN LAPORAN KEUANGAN",
  },
  {
    Code: "AK.15",
    Name: "MENYIAPKAN SURAT PEMBERITAHUAN PAJAK",
  },
  {
    Code: "AK.16",
    Name: "MENGOPERASIKAN PAKET PROGRAM PENGOLAH ANGKA",
  },
  {
    Code: "AK.17",
    Name: "MENGOPERASIKAN APLIKASI KOMPUTER AKUNTANSI",
  },
  {
    Code: "TKR01",
    Name: "MEMAHAMI DASAR2 MESIN",
  },
  {
    Code: "TKR02",
    Name: "MEMAHAMI PROSES2 DASAR PEMBENTUKAN LOGAM",
  },
  {
    Code: "TKR03",
    Name: "MENJELASKAN PROSES2 MESIN KONVERSI ENERGI",
  },
  {
    Code: "SP",
    Name: "SEJARAH PEMINATAN",
  },
  {
    Code: "SI",
    Name: "SEJARAH INDONESIA",
  },
  {
    Code: "MTKP",
    Name: "MATEMATIKA PEMINATAN",
  },
  {
    Code: "PEB",
    Name: "P. EKO & BISNIS",
  },
  {
    Code: "PAP",
    Name: "P. ADM & PERKANTORAN",
  },
  {
    Code: "LB",
    Name: "LINGKUNGAN BERSIH",
  },
  {
    Code: "PA",
    Name: "PERISTIWA ALAM",
  },
  {
    Code: "IN",
    Name: "INDAHNYA NEGERIKU",
  },
  {
    Code: "TT",
    Name: "TEMPAT TINGGALKU",
  },
  {
    Code: "B.SIG",
    Name: "BHS DAN SASTRA INGGRIS",
  },
  {
    Code: "FTIKO",
    Name: "FISIKA ( KEAHLIAN TEKNOLOGI & KOMUNIKASI )",
  },
  {
    Code: "K.KES",
    Name: "KIMIA ( KEAHLIAN KESEHATAN )",
  },
  {
    Code: "GT",
    Name: "GAMBAR TEKNIK",
  },
  {
    Code: "SK",
    Name: "SISTEM KOMPUTER",
  },
  {
    Code: "DR",
    Name: "DIRIKU",
  },
  {
    Code: "KGK",
    Name: "KEGIATANKU",
  },
  {
    Code: "IK",
    Name: "INDAHNYA KEBERSAMAAN",
  },
  {
    Code: "PTMH",
    Name: "PEDULI TERHADAP MAKHLUK HIDUP",
  },
  {
    Code: "BDL",
    Name: "BERMAIN DILINGKUNGAN",
  },
  {
    Code: "AS",
    Name: "AKU DAN SEKOLAHKU",
  },
  {
    Code: "PDK",
    Name: "PERISTIWA DALAM KEHIDUPAN",
  },
  {
    Code: "SIP",
    Name: "SEHAT ITU PENTING",
  },
  {
    Code: "T2",
    Name: "TEMA 2",
  },
  {
    Code: "T4",
    Name: "TEMA 4",
  },
  {
    Code: "T6",
    Name: "TEMA 6",
  },
  {
    Code: "T8",
    Name: "TEMA 8",
  },
  {
    Code: "AP.02",
    Name: "MENGAPLIKASIKAN KETRAMPILAN DASAR KOMUNIKASI",
  },
  {
    Code: "AP.04",
    Name: "M. PRINSIP2 KERJA SAMA DGN KOLEGA & PELANGGAN",
  },
  {
    Code: "AP.06",
    Name: "MENGOPERASIKAN APLIKASI PRESENTASI",
  },
  {
    Code: "AP.08",
    Name: "MELAKUKAN PROSEDUR ADM",
  },
  {
    Code: "AP.10",
    Name: "MENANGGANI SURAT/ DOKUMEN KANTOR",
  },
  {
    Code: "AP.12",
    Name: "MEMBUAT DOKUMEN",
  },
  {
    Code: "AP.14",
    Name: "MENGELOLA PERTEMUAN/ RAPAT",
  },
  {
    Code: "AP.16",
    Name: "MEMBERIKAN PELAYANAN KEPADA PELANGGAN",
  },
  {
    Code: "AP.18",
    Name: "MENGAPLIKASIKAN ADM PERKANTORAN DI TEMPAT KRJ",
  },
  {
    Code: "AP.09",
    Name: "MENANGGANI PENGGADAAN DOKUMEN",
  },
  {
    Code: "P.02",
    Name: "HOUSE KEEPING( AKOMODASI PERHOTELAN)",
  },
  {
    Code: "P.04",
    Name: "MENGOLAH MAKANAN INDONESIA",
  },
  {
    Code: "AK.01",
    Name: "MENERAPKAN PRINSIP PROFESIONAL BEKERJA",
  },
  {
    Code: "TKR04",
    Name: "MENGINTERPRETASIKAN GAMBAR TEKNIK",
  },
  {
    Code: "TSM11",
    Name: "M. OVERHAUL KEPALA SILENDER",
  },
  {
    Code: "TSM13",
    Name: "M.PERBAIKAN SISTEM BAHAN BAKAR SPD MOTOR",
  },
  {
    Code: "TSM15",
    Name: "M.PERBAIKAN UNIT KOPLING SPD MTR BRK KOMPONEN SIST",
  },
  {
    Code: "TSM17",
    Name: "M.PERBAIKAN SISTEM TRANSMISI MANUAL 2",
  },
  {
    Code: "TSM19",
    Name: "MELAKUKAN PERBAIKAN SISTEM SUSPENSI",
  },
  {
    Code: "TSM21",
    Name: "M.PERBAIKAN RINGAN PD RANGKAIAN SISTEM KELISTRIKAN",
  },
  {
    Code: "TSM23",
    Name: "MELAKUKAN PERBAIKAN SISTEM PENGISIAN",
  },
  {
    Code: "PM04",
    Name: "MEMAHAMI PRINSIP2 BISNIS",
  },
  {
    Code: "PM06",
    Name: "MELAKSANAKAN NEGOSIASI",
  },
  {
    Code: "PM08",
    Name: "MELAKSANAKAN PROSES ADM TRANSAKSI",
  },
  {
    Code: "PM10",
    Name: "MELAKSANAKAN PENAGIHAN PEMBAYARAN",
  },
  {
    Code: "PM12",
    Name: "MENEMUKAN PELUANG BARU DARI PELANGGAN",
  },
  {
    Code: "PM14",
    Name: "M. USAHA ECERAN / RITEL ( EXPANSION STORE OPENING",
  },
  {
    Code: "TKJ01",
    Name: "MERAKIT PERSONAL COMPUTER",
  },
  {
    Code: "TKJ03",
    Name: "MENERAPKAN KESELAMATAN KESEHATAN & LING HDP K3LH",
  },
  {
    Code: "TKJ05",
    Name: "MENERAPKAN FUNGSI PERIPHERAL & INSTALASI PC",
  },
  {
    Code: "TKJ07",
    Name: "MELAKUKAN PERBAIKAN & SETTING ULANG PC",
  },
  {
    Code: "TKJ09",
    Name: "MELAKUKAN PERAWATAN PC",
  },
  {
    Code: "TKJ11",
    Name: "MELAKUKAN INSTALASI SOFTWARE",
  },
  {
    Code: "TKJ13",
    Name: "M. PERMSL  PENGOPERASIAN PC YG TERSAMBUNG JARINGAN",
  },
  {
    Code: "TKJ15",
    Name: "M. INSTALASI SISTEM SISTEM OPERASI JAR BERBSS GUI",
  },
  {
    Code: "TKJ17",
    Name: "M. PERMSLH PERANGKAT YG TERSAMBUNG JAR BERBSS LUAS",
  },
  {
    Code: "TKJ19",
    Name: "M. P. & SETTING ULANG KONEKSI JARINGAN BSS LUAS",
  },
  {
    Code: "TKJ21",
    Name: "MERANCANG BANGUN & MENGANALISA WIDE AREA NETWORK",
  },
  {
    Code: "TKJ22",
    Name: "MERANCANG WEB DATA BASE UNTUK CONTENT SERVER",
  },
  {
    Code: "TBU",
    Name: "TATA BUSANA",
  },
  {
    Code: "DOBG",
    Name: "DIGITAL ONBOARDING BISNIS DIGITAL",
  },
  {
    Code: "SKI",
    Name: "SEJARAH KEBUDAYAAN ISLAM",
  },
  {
    Code: "TSM16",
    Name: "M.PERBAIKAN SISTEM TRANSMISI MANUAL",
  },
  {
    Code: "TSO17",
    Name: "M.PERBAIKAN SISTEM TRANSMISI OTOMATIS",
  },
  {
    Code: "GTO",
    Name: "GAMBAR TEKNIK OTOMOTIF",
  },
  {
    Code: "TDO",
    Name: "TEKNOLOGI DASAR OTOMOTIF",
  },
  {
    Code: "PDO",
    Name: "PEKERJAAN DASAR OTOMOTIF",
  },
  {
    Code: "PMKR",
    Name: "P. MESIN KENDARAAN RINGAN",
  },
  {
    Code: "PSPT",
    Name: "P.SASIS & PEMINDAH TENAGA KENDARAAN RINGAN",
  },
  {
    Code: "PKKR",
    Name: "P.KELISTRIKAN KENDARAAN RINGAN",
  },
  {
    Code: "PSSM",
    Name: "P. SASIS SEPEDA MOTOR",
  },
  {
    Code: "PKSM",
    Name: "P.KELISTRIKAN SEPEDA MOTOR",
  },
  {
    Code: "PMSM",
    Name: "P.MESIN SEPEDA MOTOR",
  },
  {
    Code: "KJD",
    Name: "KOMPUTER & JARINGAN DASAR",
  },
  {
    Code: "DG",
    Name: "DESAIN GRAFIS",
  },
  {
    Code: "AIJ",
    Name: "ADM INFRASTUKTUR JARINGAN",
  },
  {
    Code: "ASJ",
    Name: "ADM SISTEM JARINGAN",
  },
  {
    Code: "TLJ",
    Name: "TEKNOLOGI LAYANAN JARINGAN",
  },
  {
    Code: "PPL",
    Name: "PEMODELAN PERANGKAT LUNAK",
  },
  {
    Code: "BD",
    Name: "BASIS DATA RPL TI",
  },
  {
    Code: "PBO",
    Name: "PEMROGRAMAN BERORIENTASI OBJEK",
  },
  {
    Code: "DGP",
    Name: "DESAIN GRAFIS PERCETAKAN",
  },
  {
    Code: "TPAV",
    Name: "TEKNIK PENGOLAH AUDIO & VIDEO",
  },
  {
    Code: "TA3D",
    Name: "TEKNIK ANIMASI 2D & 3D",
  },
  {
    Code: "DMI",
    Name: "DESAIN MUDA INTERAKTIF",
  },
  {
    Code: "EP",
    Name: "ETIKA PROFESI",
  },
  {
    Code: "APAS",
    Name: "APLIKASI PENGOLAH ANGKA/ SPREADSHEET",
  },
  {
    Code: "KA",
    Name: "KOMPUTER AKUNTASI",
  },
  {
    Code: "ADP",
    Name: "ADMINISTRASI PAJAK",
  },
  {
    Code: "APJDM",
    Name: "P.AKT PERSH JASA, DAGANG & MANUFAKTUR",
  },
  {
    Code: "P.K",
    Name: "PENGELOLAAN KAS",
  },
  {
    Code: "APKM",
    Name: "AKT. PERBANKAN & KEUANGAN MIKRO LP",
  },
  {
    Code: "TP",
    Name: "TEKNOLOGI PERKANTORAN",
  },
  {
    Code: "KOR",
    Name: "KORESPONDENSI",
  },
  {
    Code: "KAP",
    Name: "KEARSIPAN",
  },
  {
    Code: "OTKP",
    Name: "OTOMISASI TATA KELOLA KEPEGAWAIAN",
  },
  {
    Code: "OTKK",
    Name: "OTOMISASI TATA KELOLA KEUANGAN",
  },
  {
    Code: "OTKSP",
    Name: "OTOMISASI TATA KELOLA SARANA & PRASARANA",
  },
  {
    Code: "OTKHK",
    Name: "OTOMISASI TATA KELOLA HUMAS & KEPROTOKOLAN",
  },
  {
    Code: "MKT",
    Name: "MARKETING PEMASARAN",
  },
  {
    Code: "PB",
    Name: "PERENCANAAN BISNIS",
  },
  {
    Code: "KB",
    Name: "KOMUNIKASI BISNIS",
  },
  {
    Code: "P.PRO",
    Name: "PENATAAN PRODUK",
  },
  {
    Code: "BO",
    Name: "BISNIS ONLINE",
  },
  {
    Code: "PBR",
    Name: "PENGELOLAAN BISNIS RITEL",
  },
  {
    Code: "TJBL",
    Name: "TEKNOLOGI JARINGAN BERBASIS LUAS",
  },
  {
    Code: "X.IND",
    Name: "BAHASA INDONESIA X",
  },
  {
    Code: "X.ING",
    Name: "BAHASA INGGRIS X",
  },
  {
    Code: "X.MMK",
    Name: "MATEMATIKA X",
  },
  {
    Code: "X.IPA",
    Name: "IPA X",
  },
  {
    Code: "X.IPS",
    Name: "IPS X",
  },
  {
    Code: "DOPBG",
    Name: "DIGITAL OPERATION BISNIS DIGITAL",
  },
  {
    Code: "PAUM",
    Name: "PENGELOLAAN ADM UMUM MANAGEMEN PERKANTORAN",
  },
  {
    Code: "PKU",
    Name: "PRAKARYA & KWU",
  },
  {
    Code: "AA",
    Name: "AQIDAH AKHLAK",
  },
  {
    Code: "BA",
    Name: "BAHASA ARAB",
  },
  {
    Code: "FQ",
    Name: "FIQIH",
  },
  {
    Code: "SKIS",
    Name: "SEJARAH KEBUDAYAAN ISLAM (SEJARAH KEAGAMAAN)",
  },
  {
    Code: "QHD",
    Name: "QURAN HADIST",
  },
  {
    Code: "AKT",
    Name: "P.  AKUNTANSI / AKT DASAR",
  },
  {
    Code: "TKR05",
    Name: "MENGGUNAKAN PERALATAN & PERLENGKAPAN DITEMPAT KRJ",
  },
  {
    Code: "TKR06",
    Name: "MENGGUNAKAN ALAT2 UKUR (MEASURING TOOLS)",
  },
  {
    Code: "TKR07",
    Name: "M. PROSEDUR KESELAMATAN, KSHT KRJ & LING T KRJ",
  },
  {
    Code: "TKR08",
    Name: "MEMPERBAIKI SISTEM HIDROLIK & KOMPRESOR UDARA",
  },
  {
    Code: "TKR09",
    Name: "M. PRO PENGELASAN, PEMATRIAN, PEMOTONGAN",
  },
  {
    Code: "TKR10",
    Name: "MELAKUKAN OVERHAUL SISTEM PENDINGIN & KOMPONENNYA",
  },
  {
    Code: "TKR11",
    Name: "M. SERVIS SISTEM BAHAN BAKAR BENSIN",
  },
  {
    Code: "TKR12",
    Name: "MEMPERBAIKI SISTEM INJEKSI BAHAN BAKAR DIESEL",
  },
  {
    Code: "TKR13",
    Name: "M. SERVIS ENGINE & KOMPONEN2NYA",
  },
  {
    Code: "TKR14",
    Name: "M. UNIT KOPLING & KOMPONEN2 SISTEM PENGOPRASIAN",
  },
  {
    Code: "TKR15",
    Name: "MEMELIHARA TRANSMISI",
  },
  {
    Code: "TKR16",
    Name: "MEMELIHARA UNIT FINAL DRIVE/GARDAN",
  },
  {
    Code: "TKR17",
    Name: "MEMPERBAIKI POROS PENGGERAK RODA",
  },
  {
    Code: "TKR18",
    Name: "MEMPERBAIKI RODA  & BAN",
  },
  {
    Code: "TKR19",
    Name: "MEMPERBAIKI SISTEM REM",
  },
  {
    Code: "TKR20",
    Name: "MEMPERBAIKI SISTEM KEMUDI",
  },
  {
    Code: "TKR21",
    Name: "MEMPERBAIKI SISTEM SUSPENSI",
  },
  {
    Code: "TKR22",
    Name: "MEMELIHARA BATERAI 2",
  },
  {
    Code: "TKR23",
    Name: "M.KERUSAKAN RINGAN  RANGKAIAN/ SISTEM KELISTRIKAN",
  },
  {
    Code: "TKR24",
    Name: "MEMPERBAIKI SISTEM PENGAPIAN",
  },
  {
    Code: "TKR25",
    Name: "MEMPERBAIKI SISTEM STARTER & PENGISIAN",
  },
  {
    Code: "TKR26",
    Name: "M. SERVIS SISTEM AC",
  },
  {
    Code: "MM.01",
    Name: "M.TEKNIK PENGAMBILAN GAMBAR PRODUKSI",
  },
  {
    Code: "MM.02",
    Name: "MENGGABUNGKAN AUDIO KE DALAM SAJIAN MULTIMEDIA",
  },
  {
    Code: "MM.03",
    Name: "MEMAHAMI CARA PENGGUNAAN PERALATAN TATA CAHAYA",
  },
  {
    Code: "TSM12",
    Name: "M.OVERHAUL SISTEM PENDINGIN BERIKUT KOMPONENNYA",
  },
  {
    Code: "TSM14",
    Name: "M.PERBAIKAN ENGINE SPD MOTOR BERIKUT KOMPONENNYA",
  },
  {
    Code: "SPSKE",
    Name: "SISTEM PENGMN DAN KONTROL ELEK KENDARAAM RINGAN ST",
  },
  {
    Code: "TSM18",
    Name: "M.PERBAIKAN SISTEM REM",
  },
  {
    Code: "TSM20",
    Name: "M.PEKERJAAN SERVIS PD RODA, BAN & RANTAI",
  },
  {
    Code: "TSM22",
    Name: "M. PERBAIKAN SISTEM STARTER",
  },
  {
    Code: "TSM24",
    Name: "M.PERBAIKAN SISTEM PENGAPIAN",
  },
  {
    Code: "PM05",
    Name: "MENATA PRODUK",
  },
  {
    Code: "PM07",
    Name: "M. KONFIRMASI KEPUTUSAN PELANGGAN",
  },
  {
    Code: "PM09",
    Name: "M.PENYERAHAN/PENGIRIMAN PRODUK",
  },
  {
    Code: "PM11",
    Name: "M.PERALATAN TRANSAKSI DILOKASI PENJUALAN",
  },
  {
    Code: "PM13",
    Name: "M.PELAYANAN PRIMA (SERVICE EXCELLENT )",
  },
  {
    Code: "PM15",
    Name: "MELAKUKAN PEMASARAN BARANG & JASA",
  },
  {
    Code: "TKJ02",
    Name: "MELAKUKAN INSTALASI OPERASI DASAR",
  },
  {
    Code: "TKJ04",
    Name: "MENERAPKAN TEKNIK ELEKTRONIKA ANALOG & DIGITAL DSR",
  },
  {
    Code: "TKJ06",
    Name: "MENDIAGNOSIS PERMSLAN PENGOPERASIAN PC & PERIFERAL",
  },
  {
    Code: "TKJ08",
    Name: "MELAKUKAN PERBAIKAN PERIFERAL",
  },
  {
    Code: "TKJ10",
    Name: "M. INSTALASI SISTEM OPERASI BERBASIS GUI",
  },
  {
    Code: "TKJ12",
    Name: "MELAKUKAN INSTALASI PERANGKAT JARINGAN LOKAL",
  },
  {
    Code: "TKJ14",
    Name: "M. PERBAIKAN & SETTING ULANG KONEKSI JARINGAN",
  },
  {
    Code: "TKJ16",
    Name: "M.INSTALASI PERANGKAT JARINGAN BERBASIS LUAS",
  },
  {
    Code: "TKJ18",
    Name: "MEMBUAT DESAIN SISTEM KEAMANAN JARINGAN",
  },
  {
    Code: "TKJ20",
    Name: "M. SERVER DALAM JARINGAN",
  },
  {
    Code: "TBO",
    Name: "TATA BOGA",
  },
  {
    Code: "PWRPL",
    Name: "PEMROGRAMAN WEB SETAHUN",
  },
  {
    Code: "PKPJ",
    Name: "PEMASANGAN DAN KONFIGURASI PERANGKAT JARINGAN",
  },
  {
    Code: "QH",
    Name: "QURAN HADIST (PEND AGAMA ISLAM)",
  },
  {
    Code: "AT",
    Name: "ADMINISTRASI TRANSAKSI",
  },
  {
    Code: "PALP",
    Name: "P. AKUNTANSI LEMBAGA/INST PEMERINTAH",
  },
  {
    Code: "K.TSA",
    Name: "KIMIA ( TKR)",
  },
  {
    Code: "F.TIK",
    Name: "FISIKA ( TIK )",
  },
  {
    Code: "KTKJ",
    Name: "KIMIA (TIK)",
  },
  {
    Code: "PPW",
    Name: "KEPARIWISATAAN",
  },
  {
    Code: "PBD",
    Name: "PERBANKAN DASAR",
  },
  {
    Code: "X.MTK",
    Name: "MATEMATIKA IPS X",
  },
  {
    Code: "MTK1",
    Name: "MATEMATIKA SMT 1",
  },
  {
    Code: "MTK2",
    Name: "MATEMATIKA SMT 2",
  },
  {
    Code: "IPA1",
    Name: "IPA SMT1",
  },
  {
    Code: "IPA2",
    Name: "IPA SMT 2",
  },
  {
    Code: "PJ",
    Name: "PJOK",
  },
  {
    Code: "PRA1",
    Name: "PRAKARYA SMT 1",
  },
  {
    Code: "PRA2",
    Name: "PRAKARYA SMT 2",
  },
  {
    Code: "SBK1",
    Name: "SENI BUDAYA SMT 1",
  },
  {
    Code: "SBK2",
    Name: "SENI BUDAYA SMT2",
  },
  {
    Code: "PKU1",
    Name: "PRAKARYA &KEWIRAUSAHAAN SMT 1",
  },
  {
    Code: "PKU2",
    Name: "PRAKARYA & KEWIRAUSAHAAN SMT2",
  },
  {
    Code: "SI1",
    Name: "SEJARAH INDONESIA SMT 1",
  },
  {
    Code: "SI2",
    Name: "SEJARAH INDONESIA SMT 2",
  },
  {
    Code: "SD",
    Name: "SIMULASI DIGITAL",
  },
  {
    Code: "SBUD",
    Name: "SENI BUDAYA",
  },
  {
    Code: "TJKN",
    Name: "TEKNOLOGI JARINGAN KABEL NIRKABEL",
  },
  {
    Code: "F.TEK",
    Name: "FISIKA (TKR )",
  },
  {
    Code: "PAI",
    Name: "PENDIDIKAN AGAMA ISLAM",
  },
  {
    Code: "SEJ",
    Name: "SEJARAH (BUKU 2)",
  },
  {
    Code: "KWU",
    Name: "KEWIRAUSAHAAN",
  },
  {
    Code: "EBAUP",
    Name: "EKONOMI  BISNIS & ADM UMUM",
  },
  {
    Code: "KBP",
    Name: "KOMUNIKASI BISNIS PEMASARAN",
  },
  {
    Code: "GTM",
    Name: "GAMBAR TEKNIK MANUFAKTUR",
  },
  {
    Code: "TM",
    Name: "GAMBAR TEKNIK MESIN",
  },
  {
    Code: "PDTM",
    Name: "PEKERJAAN DASAR TEKNIK MESIN",
  },
  {
    Code: "TPB",
    Name: "TEKNIK PEMESINAN BUBUT",
  },
  {
    Code: "TPF",
    Name: "TEKNIK PEMESINAN FRAIS",
  },
  {
    Code: "TPG",
    Name: "TEKNIK PEMESINAN GERINDA",
  },
  {
    Code: "TPNC",
    Name: "TEKNIK PEMESINAN NC/CNC & CAM",
  },
  {
    Code: "TI.K",
    Name: "INFORMATIKA",
  },
  {
    Code: "PKK4",
    Name: "P.KREATIF & KWU (TKR)",
  },
  {
    Code: "PKK",
    Name: "PRODUK KREATIF & KEWIRAUSAHAAN",
  },
  {
    Code: "PKK5",
    Name: "P.KREATIF & KWU (T. BISNIS SPD MOTOR )",
  },
  {
    Code: "PKK6",
    Name: "P. KREATIF & KWU (PEMESINAN )",
  },
  {
    Code: "PKK1",
    Name: "P.KREATIF & KWU ( T. KOM & JARINGAN )",
  },
  {
    Code: "PKK2",
    Name: "P. KRATIF & KWU ( R PERANGKAT LUNAK )",
  },
  {
    Code: "PKK3",
    Name: "P. KREATIF & KWU ( MULTMEDIA )",
  },
  {
    Code: "PKK7",
    Name: "P. KREATIF & KWU ( AKT & KEU LEMBAGA)",
  },
  {
    Code: "PKK8",
    Name: "P. KREATIF & KWU ( PERBANKAN & KEU MIKRO)",
  },
  {
    Code: "PKK9",
    Name: "P. KREATIF & KWU ( OTK PERKANTORAN)",
  },
  {
    Code: "PKK10",
    Name: "P.KREATIF & KWU ( B DARING & PEMASARAN )",
  },
  {
    Code: "LLK",
    Name: "LAYANAN LEMBAGA PERBANKAN & KEUANGAN MIKRO",
  },
  {
    Code: "PWM",
    Name: "PEMROGRAMAN WEB & PERANGKAT BERGERAK (LUNAK)",
  },
  {
    Code: "DPTK",
    Name: "DASAR PERANCANGAN TEKNIK MESIN",
  },
  {
    Code: "PADP",
    Name: "ADMINISTRASI UMUM",
  },
  {
    Code: "CD",
    Name: "CBT",
  },
  {
    Code: "PBSM",
    Name: "PENGELOLAAN BENGKEL SEPEDA MOTOR",
  },
  {
    Code: "SM",
    Name: "SENI MUSIK",
  },
  {
    Code: "SR",
    Name: "SENI RUPA",
  },
  {
    Code: "ST",
    Name: "SENI TARI",
  },
  {
    Code: "STT",
    Name: "SENI TEATER",
  },
  {
    Code: "IPAS",
    Name: "ILMU PENGETAHUAN ALAM & SOSIAL",
  },
  {
    Code: "PK.B",
    Name: "PRAKARYA - BUDIDAYA",
  },
  {
    Code: "PK.K",
    Name: "PRAKARYA - KERAJINAN",
  },
  {
    Code: "PK.P",
    Name: "PRAKARYA - PENGOLAHAN",
  },
  {
    Code: "PK.R",
    Name: "PRAKARYA - REKAYASA",
  },
  {
    Code: "SEJ.A",
    Name: "SEJARAH SMA/SMK",
  },
  {
    Code: "SEJ.K",
    Name: "SEJARAH (SMK )",
  },
  {
    Code: "MTKAT",
    Name: "MATEMATIKA TK LANJUT",
  },
  {
    Code: "B.SI",
    Name: "BHS & SAS IND / TINGKAT LANJUT",
  },
  {
    Code: "B.SIL",
    Name: "BHS & SAS INGG TINGKAT LANJUT",
  },
  {
    Code: "B.IGT",
    Name: "BHS.INGGRIS TEKNIK",
  },
  {
    Code: "IPB",
    Name: "IPAS BISNIS",
  },
  {
    Code: "IPKES",
    Name: "IPAS KESEHATAN",
  },
  {
    Code: "IPTEK",
    Name: "IPAS TEKNOLOGI",
  },
  {
    Code: "SEKRS",
    Name: "SISTEM ELEKTRIKAL KENDARAAN RINGAN SETAHUN",
  },
  {
    Code: "PBRPL",
    Name: "PEMROGRAMAN BASIS TEKS, GRAFIS & MULTIMEDIA",
  },
  {
    Code: "PPBRP",
    Name: "PEMROGRAMAN PERANGKAT BERGERAK SETAHUN",
  },
  {
    Code: "PPJ",
    Name: "PERENCANAAN DAN PENGALAMATAN JARINGAN",
  },
  {
    Code: "KJTKJ",
    Name: "KEAMANAN JARINGAN SETAHUN TKJ",
  },
  {
    Code: "PBBG",
    Name: "PERENCANAAN BISNIS PEMASARAN",
  },
  {
    Code: "DBBG",
    Name: "DIGITAL BRANDING BISNIS DIGITAL",
  },
  {
    Code: "DMBG",
    Name: "DIGITAL MARKETING BISNIS DIGITAL",
  },
  {
    Code: "EBMP",
    Name: "EKONOMI BISNIS MANAGEMEN PERKANTORAN",
  },
  {
    Code: "KKMP",
    Name: "KOMUNIKASI DI TEMPAT KERJA MP",
  },
  {
    Code: "TPMP",
    Name: "TEKNOLOGI PERKANTORAN MP",
  },
  {
    Code: "PKSMP",
    Name: "PENGELOLAAN KEUANGAN SEDERHANA MP",
  },
  {
    Code: "AKKEU",
    Name: "AKUTANSI DAN KEUANGAN LEMBAGA",
  },
  {
    Code: "MPL",
    Name: "MANAJEMEN PERKANTORAN DAN LAYANAN BISNIS",
  },
  {
    Code: "PPLG",
    Name: "PENGEMBANGAN PERANGKAT LUNAK DAM GIM",
  },
  {
    Code: "JKT",
    Name: "TEKNIK JARINGAN KOMPUTER DAN TELEKOMUNIKASI",
  },
  {
    Code: "TOI",
    Name: "TEKNIK OTOMOTIF_INTENS",
  },
  {
    Code: "KJS",
    Name: "KEAMANAN JARINGAN SETAHUN",
  },
  {
    Code: "DOB",
    Name: "DIGITAL ON BOARDING",
  },
  {
    Code: "DB",
    Name: "DIGITAL BRANDING",
  },
  {
    Code: "PPP",
    Name: "PENGEMASAN DAN PENDISTRIBUSIN PRODUK",
  },
  {
    Code: "PRP",
    Name: "PENGELOLAAN RAPAT / PERTEMUAN",
  },
  {
    Code: "PTKR",
    Name: "PEMINDAH TENAGA KENDARAAN RINGAN",
  },
  {
    Code: "STER",
    Name: "SENI TERPADU",
  },
  {
    Code: "PI",
    Name: "PERHOTELAN - INTENS",
  },
  {
    Code: "JUR",
    Name: "JURNAL 7 KEBIASAAN ANAK INDO",
  },
  {
    Code: "PBBD",
    Name: "PERENCANAAN BISNIS, BISNIS DIGITAL",
  },
  {
    Code: "CBR",
    Name: "CUSTOMER SERVICE BISNIS RETAIL",
  },
  {
    Code: "SBR",
    Name: "STRATEGI MARKETING VISUAL MERCHANDISING",
  },
  {
    Code: "PPSEM",
    Name: "PERAWATAN DAN PERBAIKAN ENGINE SEPEDA MOTOR",
  },
  {
    Code: "SPTSM",
    Name: "PRWT DAN PRBK SISTEM PEMINDAH TENAGA SEPEDA MOTOR",
  },
  {
    Code: "SMLH",
    Name: "PRWT DAN PRBK SEPEDA MOTOR LISTRIK DAN HYBRID",
  },
  {
    Code: "BSM",
    Name: "BENGKEL SEPEDA MOTOR",
  },
  {
    Code: "TKBK",
    Name: "TEKNIK KERJA BENGKEL DAN KELISTRIKAN",
  },
  {
    Code: "FTTX",
    Name: "FTTX",
  },
  {
    Code: "WA",
    Name: "WIRELESS ACCESS",
  },
  {
    Code: "PDDK",
    Name: "PRINSIP DASAR DESAIN DAN KOMUNIKASI",
  },
  {
    Code: "MDB",
    Name: "MENERAPKAN DESIGN BRIEF",
  },
  {
    Code: "PPD",
    Name: "PROSES PRODUKSI DESAIN",
  },
  {
    Code: "SETK",
    Name: "SEJARAH TK LANJU",
  },
  {
    Code: "KDG",
    Name: "KODING",
  },
  {
    Code: "BTA",
    Name: "BACA TULIS ALQURAN",
  },
  {
    Code: "FM",
    Name: "FISIK MOTORIK A",
  },
  {
    Code: "FM B",
    Name: "FISIK MOTORIK",
  },
  {
    Code: "SA",
    Name: "SEKOLAHKU A",
  },
  {
    Code: "AAIA",
    Name: "AKU ANAK INDONESIA A",
  },
  {
    Code: "BDSB",
    Name: "BELAJAR DI SEKOLAH B",
  },
  {
    Code: "ACIB",
    Name: "AKU CINTA INDONESIA B",
  },
  {
    Code: "KAB A",
    Name: "KASIH UNTUK AYAH IIBU",
  },
  {
    Code: "KTK B",
    Name: "KELUARGA DAN TEMAN KITA",
  },
  {
    Code: "ABJA",
    Name: "AKU BISA JAGA DIRI A",
  },
  {
    Code: "PSDM",
    Name: "PENGELOLAAN SDM MP",
  },
  {
    Code: "PSPMP",
    Name: "PENGELOLAAN SARANA PRASARANA MP",
  },
  {
    Code: "PHKMP",
    Name: "PENGELOLAAN HUMAS KEPROTOKOLAN  MP",
  },
  {
    Code: "EBALP",
    Name: "EKO BISNIS ADM UMUM LAYANAN PERBANKAN",
  },
  {
    Code: "PKLP",
    Name: "PENGELOLAAN KAS LAYANAN PERBANKAN",
  },
  {
    Code: "LPKLP",
    Name: "LAY LEMBAGA PERBANKANKEUNGAN MIKRO LP",
  },
  {
    Code: "KALP",
    Name: "KOMPUTER AKUTANSI LAY PERBANKAN",
  },
  {
    Code: "PLP",
    Name: "PERPAJAKAN LAYANAN PERBANKAN",
  },
  {
    Code: "EBAUA",
    Name: "EKONOMI BISNIS ADM UMUM AKUTANSI",
  },
  {
    Code: "APJA",
    Name: "AKUTANSI PERUSAHAAN JASA DAGANG MANUFAKTUR AK",
  },
  {
    Code: "AIPA",
    Name: "AKUTANSI  LEMBAGA/INSTANSI PEMERINTAH AKUNTANSI",
  },
  {
    Code: "AK",
    Name: "AKUTANSI KEUANGAN",
  },
  {
    Code: "PAK",
    Name: "PERPAJAKAN AKUTANSI",
  },
  {
    Code: "FNBP",
    Name: "FOOD AND BEVERAGE PERHOTELAN",
  },
  {
    Code: "MDBDK",
    Name: "MENERAPKAN DESIGN BRIEF DKV",
  },
  {
    Code: "KEKR",
    Name: "KONV ENERGI KENDARAAN RINGAN TKR",
  },
  {
    Code: "PMBKR",
    Name: "PROSES PELAYANAN MANAJEMEN BENGKEL KENDARAAN RINGA",
  },
  {
    Code: "PPKR",
    Name: "PROSEDUR PENGGUNAAN KENDARAAN RINGAN",
  },
  {
    Code: "PBKR",
    Name: "PERAWATAN BERKALA KENDARAAN RINGAN TKR",
  },
  {
    Code: "SEKR",
    Name: "SISTEM ENGINE KENDARAAN RINGAN",
  },
  {
    Code: "SSKR",
    Name: "SISTEM SASIS KENDARAAN RINGAN",
  },
  {
    Code: "SPTKR",
    Name: "SISTEM PEMINDAH TENAGA KENDARAAN RINGAN",
  },
  {
    Code: "DKV",
    Name: "DESAIN KOMUNIKASI VISUAL",
  },
  {
    Code: "PEM",
    Name: "PEMASARAN_INTENS",
  },
  {
    Code: "DSR",
    Name: "DASAR DASAR SENI RUPA",
  },
  {
    Code: "TMI",
    Name: "TEKNIK MESIN_INTENS",
  },
  {
    Code: "TJKT",
    Name: "TEKNIK JARINGAN KOMPUTER DAN TELEKOMUNIKASI 2",
  },
  {
    Code: "DM",
    Name: "DIGITAL MARKETING",
  },
  {
    Code: "DO",
    Name: "DIGITAL OPERATION",
  },
  {
    Code: "AP",
    Name: "AKUNTANSI PERUSAHAAN JASA, DAGANG, DAN MANUFAKTUR",
  },
  {
    Code: "PEK",
    Name: "PENGELOLAAN KEARSIPAN",
  },
  {
    Code: "PTER",
    Name: "PRAKARYA TERPADU",
  },
  {
    Code: "MBD",
    Name: "MARKETING BISNIS DIGITAL",
  },
  {
    Code: "MBR",
    Name: "MARKETING BISNIS RETAIL",
  },
  {
    Code: "KBR",
    Name: "KOMUNIKASI BISNIS, BISNIS RETAIL",
  },
  {
    Code: "ABR",
    Name: "ADMINISTRASI TRANSAKSI BISNIS RETAIL",
  },
  {
    Code: "PPSSM",
    Name: "PERAWATAN DAN PERBAIKAN SASIS SEPEDA MOTOR",
  },
  {
    Code: "SKSM",
    Name: "PRWT DAN PRBK SISTEM KELISTRIKAN SEPEDA MOTOR",
  },
  {
    Code: "EMSSM",
    Name: "PRWT DAN PRBK ENGINE MANAGEMENT SYSTEM SEPEDA MTR",
  },
  {
    Code: "SKEM",
    Name: "SISTEM KOMPUTER, ELEKTRONIKA, DAN MIKROPROSESOR",
  },
  {
    Code: "VSAT",
    Name: "VSAT",
  },
  {
    Code: "CPE",
    Name: "CUSTOMER PREMISE EQUIPMENT",
  },
  {
    Code: "PLD",
    Name: "PERANGKAT LUNAK DESAIN",
  },
  {
    Code: "KD",
    Name: "KARYA DESAIN",
  },
  {
    Code: "PGPRO",
    Name: "PG PRODUKTIF",
  },
  {
    Code: "MHIJ",
    Name: "MENGENAL HURUF HIJAIYAH",
  },
  {
    Code: "CA",
    Name: "CALISTUNG A",
  },
  {
    Code: "CB",
    Name: "CALISTUNG B",
  },
  {
    Code: "ACLA",
    Name: "AKU CINTA LINGKUNGAN A",
  },
  {
    Code: "LSA",
    Name: "LIBURANKU SERU A",
  },
  {
    Code: "ML B",
    Name: "MERAWAT LINGKUNGAN B",
  },
  {
    Code: "AL B",
    Name: "ASYIKNYA LIBURANKU B",
  },
  {
    Code: "BDS B",
    Name: "BELAJAR DI SEKOLAH B (LANJUTAN)",
  },
  {
    Code: "AMJA",
    Name: "AYO MENJAGA DIRI B",
  },
	}

	created := 0
	skipped := 0

	// delete all existing records before seeding
	db.Exec("DELETE FROM bidang_studi")

	//// Insert all records using FirstOrCreate
	for _, bidangStudi := range bidangStudiData {
		var result models.BidangStudi
		err := db.Where("code = ?", bidangStudi.Code).FirstOrCreate(&result, bidangStudi).Error
		if err != nil {
			return fmt.Errorf("failed to seed bidang_studi with code %s: %w", bidangStudi.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == bidangStudi.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("✓ Bidang Studi seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
