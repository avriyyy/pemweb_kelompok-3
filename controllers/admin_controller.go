package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"toktik/database"
	"toktik/models"
)

type AdminController struct{}

func (AdminController) Dashboard(c *fiber.Ctx) error {
	var totalTiket int64
	database.DB.Model(&models.Tiket{}).Count(&totalTiket)

	var totalPendapatan int64
	database.DB.Model(&models.Transaksi{}).
		Where("status IN ?", []string{"Selesai", "selesai", "lunas"}).
		Select("COALESCE(SUM(total_harga), 0)").
		Scan(&totalPendapatan)

	var totalTransaksi int64
	database.DB.Model(&models.Transaksi{}).Count(&totalTransaksi)

	var totalUser int64
	database.DB.Model(&models.User{}).Count(&totalUser)

	stats := []fiber.Map{
		{"Label": "TIKET TERJUAL", "Value": fmt.Sprintf("%d", totalTiket), "Change": "-", "Sub": "total semua tiket", "Icon": "ticket", "BgColor": "bg-primary"},
		{"Label": "PENDAPATAN", "Value": fmt.Sprintf("Rp %s", models.FormatRupiah(int(totalPendapatan))), "Change": "-", "Sub": "total pendapatan", "Icon": "wallet", "BgColor": "bg-dark-soft"},
		{"Label": "TRANSAKSI", "Value": fmt.Sprintf("%d", totalTransaksi), "Change": "-", "Sub": "total transaksi", "Icon": "receipt", "BgColor": "bg-text-secondary"},
		{"Label": "USER AKTIF", "Value": fmt.Sprintf("%d", totalUser), "Change": "-", "Sub": "total user terdaftar", "Icon": "users", "BgColor": "bg-dark-elevated"},
	}

	dayNames := []string{"Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"}
	barData := []fiber.Map{}
	var barMax int64
	for i := 6; i >= 0; i-- {
		day := time.Now().AddDate(0, 0, -i)
		dayStart := day.Format("2006-01-02")
		dayEnd := day.AddDate(0, 0, 1).Format("2006-01-02")
		var count int64
		database.DB.Model(&models.Transaksi{}).
			Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
			Count(&count)
		if count > barMax {
			barMax = count
		}
		barData = append(barData, fiber.Map{
			"Label": dayNames[day.Weekday()],
			"Value": count,
		})
	}
	// Hitung persentase untuk setiap bar
	for i := range barData {
		v := barData[i]["Value"].(int64)
		pct := 0
		if barMax > 0 {
			pct = int(v * 100 / barMax)
		}
		barData[i]["Pct"] = pct
	}

	type FilmSold struct {
		ID         uint
		Judul      string
		TotalTiket int64
	}
	var topFilms []FilmSold
	database.DB.Raw(`
		SELECT f.id, f.judul, COUNT(t.id) as total_tiket
		FROM tikets t
		JOIN schedules s ON s.id = t.jadwal_id
		JOIN films f ON f.id = s.film_id
		GROUP BY f.id, f.judul
		ORDER BY total_tiket DESC
		LIMIT 5
	`).Scan(&topFilms)

	topFilmsList := make([]fiber.Map, 0, len(topFilms))
	for _, f := range topFilms {
		topFilmsList = append(topFilmsList, fiber.Map{
			"ID":         f.ID,
			"Judul":      f.Judul,
			"TotalTiket": f.TotalTiket,
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
	tanggalRilis, _ := time.Parse("2006-01-02", c.FormValue("tanggal_rilis"))

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
		"Title": "Edit Film", "Active": "film", "Film": film,
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
	tanggalRilis, _ := time.Parse("2006-01-02", c.FormValue("tanggal_rilis"))

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

func (AdminController) JadwalIndex(c *fiber.Ctx) error {
	var jadwals []models.Schedule
	database.DB.
		Preload("Film").
		Preload("Studio").
		Find(&jadwals)
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
		"Title": "Tambah Jadwal", "Active": "jadwal",
		"Films": films, "Studios": studios,
	}, "layouts/admin")
}

func (AdminController) JadwalTambahSubmit(c *fiber.Ctx) error {
	filmID, _ := strconv.Atoi(c.FormValue("film_id"))
	studioID, _ := strconv.Atoi(c.FormValue("studio_id"))
	harga, _ := strconv.ParseFloat(c.FormValue("harga"), 64)
	tanggal, _ := time.Parse("2006-01-02", c.FormValue("tanggal"))

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
		"Title": "Edit Jadwal", "Active": "jadwal",
		"Jadwal": jadwal, "Films": films, "Studios": studios,
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
	tanggal, _ := time.Parse("2006-01-02", c.FormValue("tanggal"))

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
	var studios []models.Studio
	database.DB.Find(&studios)
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
	baris, _ := strconv.Atoi(c.FormValue("baris"))
	kolom, _ := strconv.Atoi(c.FormValue("kursi_per_baris"))

	studio := models.Studio{
		NamaStudio:  c.FormValue("nama"),
		Tipe:        c.FormValue("tipe"),
		JumlahBaris: baris,
		JumlahKolom: kolom,
	}
	if err := database.DB.Create(&studio).Error; err != nil {
		return c.SendString(err.Error())
	}
	return c.Redirect("/admin/studio")
}

func (AdminController) StudioEdit(c *fiber.Ctx) error {
	id := c.Params("id")
	var studio models.Studio
	if err := database.DB.First(&studio, id).Error; err != nil {
		return c.Redirect("/admin/studio")
	}
	return c.Render("admin/studio/edit", fiber.Map{
		"Title": "Edit Studio", "Active": "studio", "Studio": studio,
	}, "layouts/admin")
}

func (AdminController) StudioEditSubmit(c *fiber.Ctx) error {
	id := c.Params("id")
	var studio models.Studio
	if err := database.DB.First(&studio, id).Error; err != nil {
		return c.Redirect("/admin/studio")
	}
	baris, _ := strconv.Atoi(c.FormValue("baris"))
	kolom, _ := strconv.Atoi(c.FormValue("kursi_per_baris"))

	studio.NamaStudio = c.FormValue("nama")
	studio.Tipe = c.FormValue("tipe")
	studio.JumlahBaris = baris
	studio.JumlahKolom = kolom
	database.DB.Save(&studio)
	return c.Redirect("/admin/studio")
}

func (AdminController) StudioHapus(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.Studio{}, id)
	return c.Redirect("/admin/studio")
}

func (AdminController) TransaksiIndex(c *fiber.Ctx) error {
	var transaksis []models.Transaksi
	database.DB.
		Preload("User").
		Preload("Schedule.Film").
		Preload("Schedule.Studio").
		Order("created_at DESC").
		Find(&transaksis)

	items := make([]fiber.Map, 0, len(transaksis))
	for _, t := range transaksis {
		status := t.Status
		if status == "pending" {
			status = "Menunggu"
		} else if status == "selesai" || status == "Selesai" || status == "lunas" {
			status = "Selesai"
		} else if status == "batal" || status == "Batal" {
			status = "Batal"
		}
		tanggalStr := t.Schedule.TanggalTayang.Format("2 Jan 2006")

		items = append(items, fiber.Map{
			"ID":      t.ID,
			"Kode":    t.KodeBooking,
			"Nama":    t.User.Nama,
			"Film":    t.Schedule.Film.Judul,
			"Tanggal": tanggalStr,
			"Metode":  t.MetodeBayar,
			"Total":   models.FormatRupiah(t.TotalHarga),
			"Status":  status,
		})
	}

	return c.Render("admin/transaksi/index", fiber.Map{
		"Title": "Transaksi", "Active": "transaksi", "Transaksis": items,
	}, "layouts/admin")
}

func (AdminController) TransaksiDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var transaksi models.Transaksi
	if err := database.DB.
		Preload("User").
		Preload("Schedule.Film").
		Preload("Schedule.Studio").
		Preload("Tiket").
		First(&transaksi, id).Error; err != nil {
		return c.Redirect("/admin/transaksi")
	}

	status := transaksi.Status
	if status == "pending" {
		status = "Menunggu"
	} else if status == "selesai" || status == "Selesai" || status == "lunas" {
		status = "Selesai"
	} else {
		status = "Batal"
	}

	tanggalStr := transaksi.Schedule.TanggalTayang.Format("2 Jan 2006")
	items := fmt.Sprintf("%d Tiket %s", len(transaksi.Tiket), transaksi.Schedule.Studio.Tipe)

	trxMap := fiber.Map{
		"ID":      transaksi.ID,
		"Kode":    transaksi.KodeBooking,
		"Nama":    transaksi.User.Nama,
		"Film":    transaksi.Schedule.Film.Judul,
		"Items":   items,
		"Tanggal": tanggalStr,
		"Metode":  transaksi.MetodeBayar,
		"Total":   models.FormatRupiah(transaksi.TotalHarga),
		"Status":  status,
		"FilmBg":  "bg-primary",
	}

	return c.Render("admin/transaksi/detail", fiber.Map{
		"Title": "Detail Transaksi", "Active": "transaksi", "Trx": trxMap,
	}, "layouts/admin")
}

func (AdminController) UserIndex(c *fiber.Ctx) error {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return c.SendString("Gagal mengambil data user")
	}
	return c.Render("admin/user/index", fiber.Map{
		"Title": "Manajemen User", "Active": "user", "Users": users,
	}, "layouts/admin")
}

func (AdminController) JadikanAdmin(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	database.DB.Model(&models.User{}).Where("id = ?", id).Update("role", "admin")
	return c.Redirect("/admin/user")
}

func (AdminController) JadikanUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	database.DB.Model(&models.User{}).Where("id = ?", id).Update("role", "user")
	return c.Redirect("/admin/user")
}

func (AdminController) UserHapus(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	database.DB.Delete(&models.User{}, id)
	return c.Redirect("/admin/user")
}
