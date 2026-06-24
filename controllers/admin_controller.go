package controllers

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"toktik/data"
)

type AdminController struct{}

func (AdminController) Dashboard(c *fiber.Ctx) error {
	stats := []fiber.Map{
		{"Label": "TIKET TERJUAL", "Value": "1.284", "Change": "+12.5%", "Sub": "vs minggu lalu", "Icon": "ticket", "BgColor": "bg-primary"},
		{"Label": "PENDAPATAN", "Value": "Rp 187jt", "Change": "+8.3%", "Sub": "vs minggu lalu", "Icon": "wallet", "BgColor": "bg-dark-soft"},
		{"Label": "TRANSAKSI", "Value": "342", "Change": "+5.1%", "Sub": "vs minggu lalu", "Icon": "receipt", "BgColor": "bg-text-secondary"},
		{"Label": "USER AKTIF", "Value": "1.847", "Change": "+3.2%", "Sub": "vs minggu lalu", "Icon": "users", "BgColor": "bg-dark-elevated"},
	}

	barData := []fiber.Map{
		{"Label": "Sen", "Value": 45},
		{"Label": "Sel", "Value": 62},
		{"Label": "Rab", "Value": 78},
		{"Label": "Kam", "Value": 52},
		{"Label": "Jum", "Value": 89},
		{"Label": "Sab", "Value": 95},
		{"Label": "Min", "Value": 71},
	}

	type FilmSold struct {
		ID         uint
		Judul      string
		TotalTiket int
		Persen     int
	}

	filmSoldMap := map[uint]*FilmSold{}
	for _, j := range data.Jadwals {
		parts := strings.Split(j.Tiket, "/")
		if len(parts) == 0 {
			continue
		}
		sold, _ := strconv.Atoi(parts[0])
		if existing, ok := filmSoldMap[j.FilmID]; ok {
			existing.TotalTiket += sold
		} else {
			filmSoldMap[j.FilmID] = &FilmSold{
				ID:         j.FilmID,
				Judul:      j.Film,
				TotalTiket: sold,
			}
		}
	}

	topFilms := make([]FilmSold, 0, len(filmSoldMap))
	for _, f := range filmSoldMap {
		maxSold := 0
		for _, x := range filmSoldMap {
			if x.TotalTiket > maxSold {
				maxSold = x.TotalTiket
			}
		}
		if maxSold > 0 {
			f.Persen = f.TotalTiket * 100 / maxSold
		}
		topFilms = append(topFilms, *f)
	}
	sort.Slice(topFilms, func(i, j int) bool { return topFilms[i].TotalTiket > topFilms[j].TotalTiket })
	if len(topFilms) > 5 {
		topFilms = topFilms[:5]
	}

	topFilmsList := make([]fiber.Map, 0, len(topFilms))
	for _, f := range topFilms {
		topFilmsList = append(topFilmsList, fiber.Map{
			"ID":         f.ID,
			"Judul":      f.Judul,
			"TotalTiket": f.TotalTiket,
			"Persen":     f.Persen,
		})
	}

	return c.Render("admin/dashboard/index", fiber.Map{
		"Title":    "Dashboard",
		"Active":   "dashboard",
		"Stats":    stats,
		"BarData":  barData,
		"TopFilms": topFilmsList,
	}, "layouts/admin")
}

func (AdminController) FilmIndex(c *fiber.Ctx) error {
	films := make([]fiber.Map, 0, len(data.Films))
	for _, f := range data.Films {
		films = append(films, fiber.Map{
			"ID":          f.ID,
			"Judul":       f.Judul,
			"Genre":       f.Genre,
			"Durasi":      f.Durasi,
			"Rating":      f.Rating,
			"Status":      f.Status,
			"PosterColor": f.PosterColor,
			"Poster":      f.Poster,
			"Tanggal":     f.Tanggal,
		})
	}
	return c.Render("admin/film/index", fiber.Map{
		"Title": "Manajemen Film", "Active": "film", "Films": films,
	}, "layouts/admin")
}

func (AdminController) FilmTambah(c *fiber.Ctx) error {
	return c.Render("admin/film/tambah", fiber.Map{
		"Title": "Tambah Film", "Active": "film",
	}, "layouts/admin")
}

func (AdminController) FilmTambahSubmit(c *fiber.Ctx) error {
	posterURL := strings.TrimSpace(c.FormValue("poster_url"))
	if posterURL == "" {
		posterURL = "https://picsum.photos/seed/" + strings.ToLower(strings.ReplaceAll(c.FormValue("judul"), " ", "-")) + "/400/600"
	}
	return c.Redirect("/admin/film?poster=" + posterURL)
}

func (AdminController) FilmEdit(c *fiber.Ctx) error {
	id := data.ParseID(c.Params("id"))
	var film *data.Film
	for i, f := range data.Films {
		if f.ID == strconv.FormatUint(uint64(id), 10) {
			film = &data.Films[i]
			break
		}
	}
	if film == nil {
		return c.Redirect("/admin/film")
	}
	return c.Render("admin/film/edit", fiber.Map{
		"Title": "Edit Film", "Active": "film", "Film": fiber.Map{
			"ID":          film.ID,
			"Judul":       film.Judul,
			"Genre":       film.Genre,
			"Durasi":      film.Durasi,
			"Rating":      film.Rating,
			"Synopsis":    film.Synopsis,
			"Status":      film.Status,
			"PosterColor": film.PosterColor,
			"Poster":      film.Poster,
			"Tanggal":     film.Tanggal,
		},
	}, "layouts/admin")
}

func (AdminController) FilmEditSubmit(c *fiber.Ctx) error {
	return c.Redirect("/admin/film")
}

func (AdminController) FilmHapus(c *fiber.Ctx) error {
	return c.Redirect("/admin/film")
}

func (AdminController) JadwalIndex(c *fiber.Ctx) error {
	jadwals := make([]fiber.Map, 0, len(data.Jadwals))
	for _, j := range data.Jadwals {
		jadwals = append(jadwals, fiber.Map{
			"ID":         j.ID,
			"FilmID":     j.FilmID,
			"StudioID":   j.StudioID,
			"Film":       j.Film,
			"Studio":     j.Studio,
			"Tanggal":    j.Tanggal,
			"Jam":        j.Jam,
			"Harga":      j.Harga,
			"HargaNum":   j.HargaNum,
			"Tiket":      j.Tiket,
			"Status":     j.Status,
			"FilmBg":     j.FilmBg,
			"StudioBg":   j.StudioBg,
			"TanggalISO": j.TanggalISO,
		})
	}
	return c.Render("admin/jadwal/index", fiber.Map{
		"Title": "Jadwal Tayang", "Active": "jadwal", "Jadwals": jadwals,
	}, "layouts/admin")
}

func (AdminController) JadwalTambah(c *fiber.Ctx) error {
	films := make([]fiber.Map, 0, len(data.Films))
	for _, f := range data.Films {
		films = append(films, fiber.Map{
			"ID":    f.ID,
			"Judul": f.Judul,
			"Genre": f.Genre,
		})
	}
	studios := make([]fiber.Map, 0, len(data.Studios))
	for _, s := range data.Studios {
		studios = append(studios, fiber.Map{
			"ID":            s.ID,
			"Nama":          s.Nama,
			"Tipe":          s.Tipe,
			"Baris":         s.Baris,
			"KursiPerBaris": s.KursiPerBaris,
		})
	}
	return c.Render("admin/jadwal/tambah", fiber.Map{
		"Title":   "Tambah Jadwal",
		"Active":  "jadwal",
		"Films":   films,
		"Studios": studios,
	}, "layouts/admin")
}

func (AdminController) JadwalTambahSubmit(c *fiber.Ctx) error {
	return c.Redirect("/admin/jadwal")
}

func (AdminController) JadwalEdit(c *fiber.Ctx) error {
	id := data.ParseID(c.Params("id"))
	jadwal := data.FindJadwalByID(id)
	if jadwal == nil {
		return c.Redirect("/admin/jadwal")
	}
	films := make([]fiber.Map, 0, len(data.Films))
	for _, f := range data.Films {
		films = append(films, fiber.Map{
			"ID":    f.ID,
			"Judul": f.Judul,
			"Genre": f.Genre,
		})
	}
	studios := make([]fiber.Map, 0, len(data.Studios))
	for _, s := range data.Studios {
		studios = append(studios, fiber.Map{
			"ID":            s.ID,
			"Nama":          s.Nama,
			"Tipe":          s.Tipe,
			"Baris":         s.Baris,
			"KursiPerBaris": s.KursiPerBaris,
		})
	}
	return c.Render("admin/jadwal/edit", fiber.Map{
		"Title":   "Edit Jadwal",
		"Active":  "jadwal",
		"ID":      c.Params("id"),
		"Jadwal":  jadwal,
		"Films":   films,
		"Studios": studios,
	}, "layouts/admin")
}

func (AdminController) JadwalEditSubmit(c *fiber.Ctx) error {
	return c.Redirect("/admin/jadwal")
}

func (AdminController) JadwalHapus(c *fiber.Ctx) error {
	return c.Redirect("/admin/jadwal")
}

func (AdminController) StudioIndex(c *fiber.Ctx) error {
	studios := make([]fiber.Map, 0, len(data.Studios))
	for _, s := range data.Studios {
		studios = append(studios, fiber.Map{
			"ID":            s.ID,
			"Nama":          s.Nama,
			"Tipe":          s.Tipe,
			"Baris":         s.Baris,
			"KursiPerBaris": s.KursiPerBaris,
			"Status":        s.Status,
			"Bg":            s.Bg,
		})
	}
	return c.Render("admin/studio/index", fiber.Map{
		"Title": "Manajemen Studio", "Active": "studio", "Studios": studios,
	}, "layouts/admin")
}

func (AdminController) StudioTambah(c *fiber.Ctx) error {
	return c.Render("admin/studio/tambah", fiber.Map{
		"Title": "Tambah Studio", "Active": "studio",
	}, "layouts/admin")
}

func (AdminController) StudioTambahSubmit(c *fiber.Ctx) error {
	return c.Redirect("/admin/studio")
}

func (AdminController) StudioEdit(c *fiber.Ctx) error {
	id := data.ParseID(c.Params("id"))
	studio := data.FindStudioByID(id)
	if studio == nil {
		return c.Redirect("/admin/studio")
	}
	return c.Render("admin/studio/edit", fiber.Map{
		"Title":  "Edit Studio",
		"Active": "studio",
		"Studio": fiber.Map{
			"ID":            studio.ID,
			"Nama":          studio.Nama,
			"Tipe":          studio.Tipe,
			"Baris":         studio.Baris,
			"KursiPerBaris": studio.KursiPerBaris,
			"Status":        studio.Status,
		},
	}, "layouts/admin")
}

func (AdminController) StudioEditSubmit(c *fiber.Ctx) error {
	return c.Redirect("/admin/studio")
}

func (AdminController) StudioHapus(c *fiber.Ctx) error {
	return c.Redirect("/admin/studio")
}

func (AdminController) TransaksiIndex(c *fiber.Ctx) error {
	transaksis := make([]fiber.Map, 0, len(data.Transaksis))
	for _, t := range data.Transaksis {
		transaksis = append(transaksis, fiber.Map{
			"ID":           t.ID,
			"Kode":         t.Kode,
			"Nama":         t.Nama,
			"Email":        t.Email,
			"Film":         t.Film,
			"Studio":       t.Studio,
			"Tanggal":      t.Tanggal,
			"Jam":          t.Jam,
			"Kursi":        t.Kursi,
			"Metode":       t.Metode,
			"Total":        t.Total,
			"Status":       t.Status,
			"StatusBayar":  t.StatusBayar,
			"Items":        t.Items,
			"FilmBg":       t.FilmBg,
			"StudioBg":     t.StudioBg,
			"NamaAvatar":   t.NamaAvatar,
			"NamaAvatarBg": t.NamaAvatarBg,
			"MetodeIcon":   t.MetodeIcon,
		})
	}
	return c.Render("admin/transaksi/index", fiber.Map{
		"Title": "Transaksi", "Active": "transaksi", "Transaksis": transaksis,
	}, "layouts/admin")
}

func (AdminController) TransaksiDetail(c *fiber.Ctx) error {
	id := data.ParseID(c.Params("id"))
	trx := data.FindTransaksiByID(id)
	if trx == nil {
		return c.Redirect("/admin/transaksi")
	}
	return c.Render("admin/transaksi/detail", fiber.Map{
		"Title":  "Detail Transaksi",
		"Active": "transaksi",
		"Trx":    trx,
	}, "layouts/admin")
}

func (AdminController) UserIndex(c *fiber.Ctx) error {
	users := make([]fiber.Map, 0, len(data.Users))
	for _, u := range data.Users {
		users = append(users, fiber.Map{
			"ID":       u.ID,
			"Nama":     u.Nama,
			"Email":    u.Email,
			"Role":     u.Role,
			"JoinDate": u.JoinDate,
			"Status":   u.Status,
			"AvatarBg": u.AvatarBg,
			"Initial":  u.Initial,
		})
	}
	return c.Render("admin/user/index", fiber.Map{
		"Title": "Manajemen User", "Active": "user", "Users": users,
	}, "layouts/admin")
}
