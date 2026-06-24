package data

import "time"

type Film struct {
	ID          string
	Judul       string
	Genre       string
	Durasi      string
	Rating      string
	Synopsis    string
	PosterColor string
	Poster      string
	Harga       string
	Times       []string
	Featured    bool
	Sinop       string
	PosterBg    string
	Status      string
	Tanggal     string
}

type Studio struct {
	ID            uint
	Nama          string
	Tipe          string
	Baris         string
	KursiPerBaris string
	Status        string
	Bg            string
}

type Jadwal struct {
	ID         uint
	FilmID     uint
	StudioID   uint
	Film       string
	Studio     string
	Tanggal    string
	TanggalISO string
	Jam        string
	Harga      string
	HargaNum   string
	Tiket      string
	Status     string
	FilmBg     string
	StudioBg   string
}

type Transaksi struct {
	ID             uint
	Kode           string
	Nama           string
	Email          string
	Film           string
	Studio         string
	Tanggal        string
	Jam            string
	Kursi          string
	Metode         string
	Total          string
	TotalNum       string
	Status         string
	StatusBayar    string
	TanggalISO     string
	JamISO         string
	FilmBg         string
	StudioBg       string
	NamaAvatar     string
	NamaAvatarBg   string
	MetodeIcon     string
	Items          string
}

type User struct {
	ID        uint
	Nama      string
	Email     string
	Phone     string
	Role      string
	JoinDate  string
	Status    string
	AvatarBg  string
	Initial   string
}

type Ticket struct {
	ID          uint
	FilmID      string
	Judul       string
	Genre       string
	Durasi      string
	PosterColor string
	Tanggal     string
	TanggalISO  string
	Jam         string
	JamISO      string
	Studio      string
	Tipe        string
	Seats       []string
	SeatCode    string
	KodeBooking string
	KodeQR      string
	Total       string
	TotalNum    string
	Harga       string
	Status      string
	MetodeBayar string
	Nama        string
	Email       string
}

var Films = []Film{
	{ID: "1", Judul: "Echoes of the Unknown", Genre: "Fiksi Ilmiah · Petualangan", Durasi: "2j 28m", Rating: "13+", Synopsis: "Perjalanan fiksi ilmiah yang mendebarkan melalui dunia yang terlupakan.", PosterColor: "bg-primary", Poster: "https://picsum.photos/seed/echoes-unknown/400/600", Harga: "45.000", Times: []string{"13:00", "15:30", "18:00", "20:30", "22:00"}, Featured: true, Status: "Aktif", Tanggal: "23 Jun 2026"},
	{ID: "2", Judul: "Midnight Tide", Genre: "Drama · Romansa", Durasi: "2j 05m", Rating: "13+", Synopsis: "Kisah cinta yang terjalin di pesisir pantai saat malam.", PosterColor: "bg-primary", Poster: "https://picsum.photos/seed/midnight-tide/400/600", Harga: "45.000", Times: []string{"14:00", "16:45", "19:15", "21:30"}, Featured: true, Status: "Aktif", Tanggal: "23 Jun 2026"},
	{ID: "3", Judul: "Crimson Peak", Genre: "Horor · Misteri", Durasi: "1j 55m", Rating: "17+", Synopsis: "Rumah tua berdarah menyimpan rahasia kelam.", PosterColor: "bg-dark-soft", Poster: "https://picsum.photos/seed/crimson-peak/400/600", Harga: "50.000", Times: []string{"13:30", "16:00", "18:30", "21:00"}, Featured: true, Status: "Aktif", Tanggal: "23 Jun 2026"},
	{ID: "4", Judul: "Neon Skyline", Genre: "Aksi · Fiksi Ilmiah", Durasi: "2j 12m", Rating: "13+", Synopsis: "Petualangan di kota metropolitan masa depan.", PosterColor: "bg-text-secondary", Poster: "https://picsum.photos/seed/neon-skyline/400/600", Harga: "55.000", Times: []string{"12:45", "15:15", "17:45", "20:15", "22:30"}, Featured: false, Status: "Aktif", Tanggal: "24 Jun 2026"},
	{ID: "5", Judul: "Velvet Hours", Genre: "Romantis", Durasi: "1j 48m", Rating: "13+", Synopsis: "Sebuah pertemuan tak terduga di sore yang hujan.", PosterColor: "bg-primary", Poster: "https://picsum.photos/seed/velvet-hours/400/600", Harga: "40.000", Times: []string{"13:00", "15:30", "18:00", "20:30"}, Featured: false, Status: "Aktif", Tanggal: "25 Jun 2026"},
	{ID: "6", Judul: "The Last Echoing Letter", Genre: "Drama", Durasi: "1j 50m", Rating: "13+", Synopsis: "Surat terakhir yang mengubah segalanya.", PosterColor: "bg-dark-soft", Poster: "https://picsum.photos/seed/last-letter/400/600", Harga: "45.000", Times: []string{"14:00", "16:30", "19:00", "21:30"}, Featured: false, Status: "Tidak Aktif", Tanggal: "26 Jun 2026"},
	{ID: "7", Judul: "Iron Wolves", Genre: "Aksi", Durasi: "1j 55m", Rating: "17+", Synopsis: "Sekelompok pejuang legendaris.", PosterColor: "bg-text-secondary", Poster: "https://picsum.photos/seed/iron-wolves/400/600", Harga: "50.000", Times: []string{"13:00", "16:00", "19:00", "22:00"}, Featured: false, Status: "Tidak Aktif", Tanggal: "27 Jun 2026"},
	{ID: "8", Judul: "Winter Echo", Genre: "Drama", Durasi: "2j 05m", Rating: "13+", Synopsis: "Kisah musim dingin yang menghangatkan.", PosterColor: "bg-primary", Poster: "https://picsum.photos/seed/winter-echo/400/600", Harga: "40.000", Times: []string{"12:30", "15:00", "17:30", "20:00"}, Featured: false, Status: "Aktif", Tanggal: "28 Jun 2026"},
}

var Studios = []Studio{
	{ID: 1, Nama: "Studio 1", Tipe: "Reguler", Baris: "8", KursiPerBaris: "10", Status: "Aktif", Bg: "bg-primary"},
	{ID: 2, Nama: "Studio 2", Tipe: "Premiere", Baris: "10", KursiPerBaris: "12", Status: "Aktif", Bg: "bg-primary"},
	{ID: 3, Nama: "Studio 3", Tipe: "Premiere", Baris: "8", KursiPerBaris: "12", Status: "Aktif", Bg: "bg-primary"},
	{ID: 4, Nama: "Studio 4", Tipe: "Premiere", Baris: "10", KursiPerBaris: "12", Status: "Aktif", Bg: "bg-primary"},
	{ID: 5, Nama: "Studio 5", Tipe: "Premiere", Baris: "8", KursiPerBaris: "12", Status: "Aktif", Bg: "bg-primary"},
	{ID: 6, Nama: "Studio 6", Tipe: "IMAX", Baris: "12", KursiPerBaris: "14", Status: "Tidak Aktif", Bg: "bg-primary"},
}

var Jadwals = []Jadwal{
	{ID: 1, FilmID: 1, StudioID: 2, Film: "Echoes of the Unknown", Studio: "Studio 2", Tanggal: "23 Jun 2026", TanggalISO: "2026-06-23", Jam: "13:00", Harga: "Rp 45.000", HargaNum: "45000", Tiket: "98/120", Status: "Aktif", FilmBg: "bg-primary", StudioBg: "bg-primary"},
	{ID: 2, FilmID: 2, StudioID: 4, Film: "Midnight Tide", Studio: "Studio 4", Tanggal: "23 Jun 2026", TanggalISO: "2026-06-23", Jam: "15:30", Harga: "Rp 45.000", HargaNum: "45000", Tiket: "105/120", Status: "Aktif", FilmBg: "bg-dark-soft", StudioBg: "bg-primary"},
	{ID: 3, FilmID: 3, StudioID: 2, Film: "Crimson Peak", Studio: "Studio 2", Tanggal: "23 Jun 2026", TanggalISO: "2026-06-23", Jam: "18:00", Harga: "Rp 50.000", HargaNum: "50000", Tiket: "115/120", Status: "Aktif", FilmBg: "bg-text-secondary", StudioBg: "bg-primary"},
	{ID: 4, FilmID: 5, StudioID: 1, Film: "Velvet Hours", Studio: "Studio 1", Tanggal: "23 Jun 2026", TanggalISO: "2026-06-23", Jam: "20:30", Harga: "Rp 40.000", HargaNum: "40000", Tiket: "75/80", Status: "Aktif", FilmBg: "bg-primary", StudioBg: "bg-primary"},
	{ID: 5, FilmID: 6, StudioID: 3, Film: "The Last Echoing Letter", Studio: "Studio 3", Tanggal: "24 Jun 2026", TanggalISO: "2026-06-24", Jam: "13:00", Harga: "Rp 45.000", HargaNum: "45000", Tiket: "85/96", Status: "Selesai", FilmBg: "bg-dark-soft", StudioBg: "bg-primary"},
	{ID: 6, FilmID: 7, StudioID: 5, Film: "Iron Wolves", Studio: "Studio 5", Tanggal: "25 Jun 2026", TanggalISO: "2026-06-25", Jam: "16:00", Harga: "Rp 50.000", HargaNum: "50000", Tiket: "90/96", Status: "Aktif", FilmBg: "bg-text-secondary", StudioBg: "bg-primary"},
	{ID: 7, FilmID: 8, StudioID: 6, Film: "Winter Echo", Studio: "Studio 6", Tanggal: "26 Jun 2026", TanggalISO: "2026-06-26", Jam: "19:30", Harga: "Rp 40.000", HargaNum: "40000", Tiket: "120/168", Status: "Aktif", FilmBg: "bg-primary", StudioBg: "bg-primary"},
	{ID: 8, FilmID: 4, StudioID: 2, Film: "Neon Skyline", Studio: "Studio 2", Tanggal: "27 Jun 2026", TanggalISO: "2026-06-27", Jam: "21:00", Harga: "Rp 55.000", HargaNum: "55000", Tiket: "110/120", Status: "Aktif", FilmBg: "bg-dark-soft", StudioBg: "bg-primary"},
}

var Transaksis = []Transaksi{
	{ID: 1, Kode: "TT-20260623-0042", Nama: "Andi Pratama", Email: "andi@email.com", Film: "Echoes of the Unknown", Studio: "Studio 4", Tanggal: "23 Jun 2026", TanggalISO: "2026-06-23", Jam: "20:00 WIB", JamISO: "20:00", Kursi: "B04, C06, C07, E08", Metode: "Kartu Kredit", Total: "Rp 200.000", TotalNum: "200000", Status: "Selesai", StatusBayar: "Lunas", Items: "4 Kursi Premiere", FilmBg: "bg-primary", StudioBg: "bg-primary", NamaAvatar: "A", NamaAvatarBg: "bg-primary", MetodeIcon: "credit-card"},
	{ID: 2, Kode: "TT-20260624-0118", Nama: "Siti Nurhaliza", Email: "siti@email.com", Film: "Midnight Tide", Studio: "Studio 2", Tanggal: "24 Jun 2026", TanggalISO: "2026-06-24", Jam: "19:15 WIB", JamISO: "19:15", Kursi: "F05, F06", Metode: "QRIS", Total: "Rp 100.000", TotalNum: "100000", Status: "Selesai", StatusBayar: "Lunas", Items: "2 Kursi Premiere", FilmBg: "bg-dark-soft", StudioBg: "bg-primary", NamaAvatar: "S", NamaAvatarBg: "bg-primary", MetodeIcon: "qr-code"},
	{ID: 3, Kode: "TT-20260625-0156", Nama: "Budi Santoso", Email: "budi@email.com", Film: "Crimson Peak", Studio: "Studio 2", Tanggal: "25 Jun 2026", TanggalISO: "2026-06-25", Jam: "18:00 WIB", JamISO: "18:00", Kursi: "D08, D09", Metode: "E-Wallet", Total: "Rp 100.000", TotalNum: "100000", Status: "Selesai", StatusBayar: "Lunas", Items: "2 Kursi Premiere", FilmBg: "bg-text-secondary", StudioBg: "bg-primary", NamaAvatar: "B", NamaAvatarBg: "bg-primary", MetodeIcon: "wallet"},
	{ID: 4, Kode: "TT-20260626-0203", Nama: "Dewi Lestari", Email: "dewi@email.com", Film: "Velvet Hours", Studio: "Studio 1", Tanggal: "26 Jun 2026", TanggalISO: "2026-06-26", Jam: "20:30 WIB", JamISO: "20:30", Kursi: "E05, E06, E07", Metode: "Transfer Bank", Total: "Rp 120.000", TotalNum: "120000", Status: "Selesai", StatusBayar: "Lunas", Items: "3 Kursi Reguler", FilmBg: "bg-primary", StudioBg: "bg-primary", NamaAvatar: "D", NamaAvatarBg: "bg-primary", MetodeIcon: "building-2"},
	{ID: 5, Kode: "TT-20260627-0289", Nama: "Rian Hidayat", Email: "rian@email.com", Film: "Echoes of the Unknown", Studio: "Studio 4", Tanggal: "27 Jun 2026", TanggalISO: "2026-06-27", Jam: "18:00 WIB", JamISO: "18:00", Kursi: "A05, A06", Metode: "Virtual Account", Total: "Rp 90.000", TotalNum: "90000", Status: "Selesai", StatusBayar: "Lunas", Items: "2 Kursi Premiere", FilmBg: "bg-primary", StudioBg: "bg-primary", NamaAvatar: "R", NamaAvatarBg: "bg-dark-soft", MetodeIcon: "hash"},
	{ID: 6, Kode: "TT-20260628-0301", Nama: "Maya Sari", Email: "maya@email.com", Film: "The Last Echoing Letter", Studio: "Studio 3", Tanggal: "28 Jun 2026", TanggalISO: "2026-06-28", Jam: "14:00 WIB", JamISO: "14:00", Kursi: "G07, G08", Metode: "QRIS", Total: "Rp 90.000", TotalNum: "90000", Status: "Menunggu", StatusBayar: "Pending", Items: "2 Kursi Premiere", FilmBg: "bg-dark-soft", StudioBg: "bg-primary", NamaAvatar: "M", NamaAvatarBg: "bg-primary", MetodeIcon: "qr-code"},
	{ID: 7, Kode: "TT-20260629-0315", Nama: "Andi Pratama", Email: "andi@email.com", Film: "Iron Wolves", Studio: "Studio 5", Tanggal: "29 Jun 2026", TanggalISO: "2026-06-29", Jam: "19:00 WIB", JamISO: "19:00", Kursi: "B10", Metode: "Kartu Kredit", Total: "Rp 50.000", TotalNum: "50000", Status: "Batal", StatusBayar: "Gagal", Items: "1 Kursi Premiere", FilmBg: "bg-text-secondary", StudioBg: "bg-primary", NamaAvatar: "A", NamaAvatarBg: "bg-primary", MetodeIcon: "credit-card"},
}

var Users = []User{
	{ID: 1, Nama: "Andi Pratama", Email: "andi@email.com", Role: "User", JoinDate: "12 Jan 2026", Status: "Aktif", AvatarBg: "bg-primary", Initial: "A"},
	{ID: 2, Nama: "Siti Nurhaliza", Email: "siti@email.com", Role: "User", JoinDate: "18 Feb 2026", Status: "Aktif", AvatarBg: "bg-primary", Initial: "S"},
	{ID: 3, Nama: "Budi Santoso", Email: "budi@email.com", Role: "User", JoinDate: "5 Mar 2026", Status: "Aktif", AvatarBg: "bg-primary", Initial: "B"},
	{ID: 4, Nama: "Dewi Lestari", Email: "dewi@email.com", Role: "User", JoinDate: "22 Mar 2026", Status: "Aktif", AvatarBg: "bg-primary", Initial: "D"},
	{ID: 5, Nama: "Rian Hidayat", Email: "rian@email.com", Role: "Admin", JoinDate: "1 Apr 2026", Status: "Aktif", AvatarBg: "bg-dark-soft", Initial: "R"},
	{ID: 6, Nama: "Maya Sari", Email: "maya@email.com", Role: "User", JoinDate: "15 Apr 2026", Status: "Nonaktif", AvatarBg: "bg-primary", Initial: "M"},
	{ID: 7, Nama: "Admin TokTik", Email: "admin@toktik.com", Role: "Admin", JoinDate: "1 Jan 2026", Status: "Aktif", AvatarBg: "bg-dark-soft", Initial: "A"},
}

var Tickets = []Ticket{
	{ID: 1, FilmID: "1", Judul: "Echoes of the Unknown", Genre: "Fiksi Ilmiah · Petualangan", Durasi: "2j 28m", PosterColor: "bg-primary", Tanggal: "Senin, 23 Juni 2026", TanggalISO: "2026-06-23", Jam: "20:00 WIB", JamISO: "20:00", Studio: "Studio 4", Tipe: "Premiere", Seats: []string{"B04", "C06", "C07", "E08"}, SeatCode: "4 Tiket", KodeBooking: "TT-20260623-0042", Total: "Rp 200.000", TotalNum: "200.000", Harga: "50000", Status: "Lunas", MetodeBayar: "Kartu Kredit Visa", Nama: "Andi Pratama", Email: "andi@email.com"},
	{ID: 2, FilmID: "2", Judul: "Midnight Tide", Genre: "Drama · Romansa", Durasi: "2j 05m", PosterColor: "bg-dark-soft", Tanggal: "Selasa, 30 Juni 2026", TanggalISO: "2026-06-30", Jam: "19:15 WIB", JamISO: "19:15", Studio: "Studio 2", Tipe: "Premiere", Seats: []string{"F05", "F06"}, SeatCode: "2 Tiket", KodeBooking: "TT-20260621-0102", Total: "Rp 100.000", TotalNum: "100.000", Harga: "50000", Status: "Pending", MetodeBayar: "QRIS", Nama: "Andi Pratama", Email: "andi@email.com"},
	{ID: 3, FilmID: "3", Judul: "Crimson Peak", Genre: "Horor · Misteri", Durasi: "1j 55m", PosterColor: "bg-text-secondary", Tanggal: "Senin, 15 Juni 2026", TanggalISO: "2026-06-15", Jam: "20:30 WIB", JamISO: "20:30", Studio: "Studio 1", Tipe: "Reguler", Seats: []string{"D08"}, SeatCode: "1 Tiket", KodeBooking: "TT-20260615-0008", Total: "Rp 35.000", TotalNum: "35.000", Harga: "35000", Status: "Lunas", MetodeBayar: "E-Wallet (OVO)", Nama: "Andi Pratama", Email: "andi@email.com"},
	{ID: 4, FilmID: "5", Judul: "Velvet Hours", Genre: "Drama", Durasi: "1j 48m", PosterColor: "bg-primary", Tanggal: "Jumat, 26 Juni 2026", TanggalISO: "2026-06-26", Jam: "21:00 WIB", JamISO: "21:00", Studio: "Studio 5", Tipe: "Premiere", Seats: []string{"G03", "G04"}, SeatCode: "2 Tiket", KodeBooking: "TT-20260618-0050", Total: "Rp 80.000", TotalNum: "80.000", Harga: "40000", Status: "Lunas", MetodeBayar: "Transfer Bank (BCA)", Nama: "Andi Pratama", Email: "andi@email.com"},
}

func FormatRupiah(n int) string {
	s := ""
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if s == "" {
		s = "0"
	}
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}
	if neg {
		result = "-" + result
	}
	return result
}

func ParseID(s string) uint {
	var id uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		id = id*10 + uint(c-'0')
	}
	return id
}

func Now() time.Time {
	return time.Now()
}

func FindFilmByID(id string) *Film {
	for i, f := range Films {
		if f.ID == id {
			return &Films[i]
		}
	}
	return nil
}

func FindStudioByID(id uint) *Studio {
	for i, s := range Studios {
		if s.ID == id {
			return &Studios[i]
		}
	}
	return nil
}

func FindJadwalByID(id uint) *Jadwal {
	for i, j := range Jadwals {
		if j.ID == id {
			return &Jadwals[i]
		}
	}
	return nil
}

func FindTransaksiByID(id uint) *Transaksi {
	for i, t := range Transaksis {
		if t.ID == id {
			return &Transaksis[i]
		}
	}
	return nil
}

func FindUserByID(id uint) *User {
	for i, u := range Users {
		if u.ID == id {
			return &Users[i]
		}
	}
	return nil
}

func FindTicketByID(id uint) *Ticket {
	for i, t := range Tickets {
		if t.ID == id {
			return &Tickets[i]
		}
	}
	return nil
}
