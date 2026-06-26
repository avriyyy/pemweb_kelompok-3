package controllers

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"toktik/data"
	"toktik/database"
	"toktik/models"
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

	var films []models.Film

	database.DB.Find(&films)

	return c.Render("admin/film/index", fiber.Map{
		"Title":  "Manajemen Film",
		"Active": "film",
		"Films":  films,
	}, "layouts/admin")
}

func (AdminController) FilmTambah(c *fiber.Ctx) error {
	return c.Render("admin/film/tambah", fiber.Map{
		"Title": "Tambah Film", "Active": "film",
	}, "layouts/admin")
}

func (AdminController) FilmTambahSubmit(c *fiber.Ctx) error {

	durasi, _ := strconv.Atoi(c.FormValue("durasi"))
	harga, _ := strconv.ParseFloat(c.FormValue("harga"), 64)

	tanggalRilis, _ := time.Parse(
		"2006-01-02",
		c.FormValue("tanggal_rilis"),
	)

	film := models.Film{
		Judul:        c.FormValue("judul"),
		Genre:        c.FormValue("genre"),
		Durasi:       durasi,
		Sinopsis:     c.FormValue("sinopsis"),
		Poster:       c.FormValue("poster_url"),
		Rating:       c.FormValue("rating"),
		TanggalRilis: tanggalRilis,
		Status:       c.FormValue("status"),
		Harga:        harga,
	}

	if err := database.DB.Create(&film).Error; err != nil {
		return c.SendString(err.Error())
	}

	return c.Redirect("/admin/film")
}

func (AdminController) FilmEdit(c *fiber.Ctx) error {

	id := c.Params("id")

	var film models.Film

	if err := database.DB.First(&film, id).Error; err != nil {
		return c.Redirect("/admin/film")
	}

	return c.Render("admin/film/edit", fiber.Map{
		"Title":  "Edit Film",
		"Active": "film",
		"Film":   film,
	}, "layouts/admin")
}

func (AdminController) FilmEditSubmit(c *fiber.Ctx) error {

	id := c.Params("id")

	var film models.Film

	if err := database.DB.First(&film, id).Error; err != nil {
		return c.Redirect("/admin/film")
	}

	durasi, _ := strconv.Atoi(c.FormValue("durasi"))
	harga, _ := strconv.ParseFloat(c.FormValue("harga"), 64)

	tanggalRilis, _ := time.Parse(
		"2006-01-02",
		c.FormValue("tanggal_rilis"),
	)

	film.Judul = c.FormValue("judul")
	film.Genre = c.FormValue("genre")
	film.Durasi = durasi
	film.Sinopsis = c.FormValue("sinopsis")
	film.Poster = c.FormValue("poster_url")
	film.Rating = c.FormValue("rating")
	film.TanggalRilis = tanggalRilis
	film.Status = c.FormValue("status")
	film.Harga = harga

	database.DB.Save(&film)

	return c.Redirect("/admin/film")
}

func (AdminController) FilmHapus(c *fiber.Ctx) error {

	id := c.Params("id")

	database.DB.Delete(&models.Film{}, id)

	return c.Redirect("/admin/film")
}

// func (AdminController) JadwalIndex(c *fiber.Ctx) error {

// 	var jadwals []models.Schedule

// 	database.DB.
// 		Preload("Film").
// 		Preload("Studio").
// 		Find(&jadwals)

//		return c.Render("admin/jadwal/index", fiber.Map{
//			"Title":   "Jadwal Tayang",
//			"Active":  "jadwal",
//			"Jadwals": jadwals,
//		}, "layouts/admin")
//	}
func (AdminController) JadwalIndex(c *fiber.Ctx) error {

	var jadwals []models.Schedule

	database.DB.
		Preload("Film").
		Preload("Studio").
		Find(&jadwals)

	for _, j := range jadwals {
		println(
			"Film =", j.Film.Judul,
			"Studio =", j.Studio.NamaStudio,
		)
	}

	return c.Render("admin/jadwal/index", fiber.Map{
		"Title":   "Jadwal Tayang",
		"Active":  "jadwal",
		"Jadwals": jadwals,
	}, "layouts/admin")
}

func (AdminController) JadwalTambah(c *fiber.Ctx) error {

	var films []models.Film
	var studios []models.Studio

	database.DB.Find(&films)
	database.DB.Find(&studios)

	return c.Render("admin/jadwal/tambah", fiber.Map{
		"Title":   "Tambah Jadwal",
		"Active":  "jadwal",
		"Films":   films,
		"Studios": studios,
	}, "layouts/admin")
}

func (AdminController) JadwalTambahSubmit(c *fiber.Ctx) error {

	filmID, _ := strconv.Atoi(c.FormValue("film_id"))
	studioID, _ := strconv.Atoi(c.FormValue("studio_id"))

	harga, _ := strconv.ParseFloat(c.FormValue("harga"), 64)

	tanggal, _ := time.Parse(
		"2006-01-02",
		c.FormValue("tanggal"),
	)

	jadwal := models.Schedule{
		FilmID:        uint(filmID),
		StudioID:      uint(studioID),
		TanggalTayang: tanggal,
		JamTayang:     c.FormValue("jam"),
		Harga:         harga,
		Status:        c.FormValue("status"),
	}

	if err := database.DB.Create(&jadwal).Error; err != nil {
		return c.SendString(err.Error())
	}

	return c.Redirect("/admin/jadwal")
}

func (AdminController) JadwalEdit(c *fiber.Ctx) error {

	id := c.Params("id")

	var jadwal models.Schedule

	if err := database.DB.
		Preload("Film").
		Preload("Studio").
		First(&jadwal, id).Error; err != nil {

		return c.Redirect("/admin/jadwal")
	}

	tanggalISO := jadwal.TanggalTayang.Format("2006-01-02")

	var films []models.Film
	var studios []models.Studio

	database.DB.Find(&films)
	database.DB.Find(&studios)

	return c.Render("admin/jadwal/edit", fiber.Map{
		"Title":      "Edit Jadwal",
		"Active":     "jadwal",
		"Jadwal":     jadwal,
		"Films":      films,
		"Studios":    studios,
		"TanggalISO": tanggalISO,
	}, "layouts/admin")
}

func (AdminController) JadwalEditSubmit(c *fiber.Ctx) error {

	id := c.Params("id")

	var jadwal models.Schedule

	if err := database.DB.First(&jadwal, id).Error; err != nil {
		return c.Redirect("/admin/jadwal")
	}

	filmID, _ := strconv.Atoi(c.FormValue("film_id"))
	studioID, _ := strconv.Atoi(c.FormValue("studio_id"))
	harga, _ := strconv.ParseFloat(c.FormValue("harga"), 64)

	tanggal, _ := time.Parse(
		"2006-01-02",
		c.FormValue("tanggal"),
	)

	jadwal.FilmID = uint(filmID)
	jadwal.StudioID = uint(studioID)
	jadwal.TanggalTayang = tanggal
	jadwal.JamTayang = c.FormValue("jam")
	jadwal.Harga = harga
	jadwal.Status = c.FormValue("status")

	database.DB.Save(&jadwal)

	return c.Redirect("/admin/jadwal")
}

func (AdminController) JadwalHapus(c *fiber.Ctx) error {

	id := c.Params("id")

	database.DB.Delete(&models.Schedule{}, id)

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
