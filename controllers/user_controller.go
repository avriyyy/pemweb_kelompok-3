package controllers

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"

	"toktik/database"
	"toktik/models"
)

type UserController struct{}

func (UserController) Beranda(c *fiber.Ctx) error {
	var films []models.Film
	database.DB.Where("status = ?", "active").Find(&films)

	featured := []fiber.Map{}
	for _, m := range films {
		if len(featured) >= 3 {
			break
		}
		featured = append(featured, fiber.Map{
			"ID":       fmt.Sprintf("%d", m.ID),
			"Judul":    m.Judul,
			"Genre":    m.Genre,
			"Durasi":   fmt.Sprintf("%dj %dm", m.Durasi/60, m.Durasi%60),
			"Rating":   m.Rating,
			"Synopsis": m.Sinopsis,
			"Poster":   m.Poster,
			"Harga":    fmt.Sprintf("%.0f", m.Harga),
		})
	}

	now := time.Now()
	curMin := now.Hour()*60 + now.Minute()

	type soonItem struct {
		ID     string
		Judul  string
		Genre  string
		Durasi string
		Poster string
		Jam    string
		SoonIn int
	}
	var soon []soonItem
	for _, m := range films {
		var schedules []models.Schedule
		database.DB.Where("film_id = ? AND tanggal_tayang >= ?", m.ID, now.Format("2006-01-02")).Find(&schedules)

		var nextJam string
		nextDiff := 24 * 60
		for _, s := range schedules {
			parts := strings.Split(s.JamTayang, ":")
			if len(parts) < 2 {
				continue
			}
			h, _ := strconv.Atoi(parts[0])
			min, _ := strconv.Atoi(parts[1])
			tMin := h*60 + min
			sDate := s.TanggalTayang.Format("2006-01-02")
			today := now.Format("2006-01-02")
			var diff int
			if sDate == today {
				diff = tMin - curMin
			} else {
				diff = 0
			}
			if diff < 0 {
				diff += 24 * 60
			}
			if diff < nextDiff {
				nextDiff = diff
				nextJam = s.JamTayang
			}
		}
		if nextJam != "" && nextDiff <= 6*60 {
			soon = append(soon, soonItem{
				ID: fmt.Sprintf("%d", m.ID), Judul: m.Judul, Genre: m.Genre,
				Durasi: fmt.Sprintf("%dj %dm", m.Durasi/60, m.Durasi%60),
				Poster: m.Poster, Jam: nextJam, SoonIn: nextDiff,
			})
		}
	}
	sort.Slice(soon, func(i, j int) bool { return soon[i].SoonIn < soon[j].SoonIn })
	if len(soon) > 5 {
		soon = soon[:5]
	}

	nowShowing := make([]fiber.Map, 0, len(soon))
	for _, s := range soon {
		nowShowing = append(nowShowing, fiber.Map{
			"ID":     s.ID,
			"Judul":  s.Judul,
			"Genre":  s.Genre,
			"Durasi": s.Durasi,
			"Poster": s.Poster,
			"Jam":    s.Jam,
		})
	}

	userID := c.Cookies("user_id")
	var user models.User
	database.DB.First(&user, userID)

	initial := ""
	if len(user.Nama) > 0 {
		initial = string(user.Nama[0])
	}

	return c.Render("user/beranda/index", fiber.Map{
		"Title":      "Beranda",
		"Active":     "home",
		"Initial":    initial,
		"Nama":       user.Nama,
		"Email":      user.Email,
		"Featured":   featured,
		"NowShowing": nowShowing,
	}, "layouts/user")
}

func (UserController) TiketIndex(c *fiber.Ctx) error {
	now := time.Now()
	dayNames := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	monthNames := []string{"", "JAN", "FEB", "MAR", "APR", "MEI", "JUN", "JUL", "AGU", "SEP", "OKT", "NOV", "DES"}

	dates := make([]fiber.Map, 0, 7)
	for i := 0; i < 7; i++ {
		d := now.AddDate(0, 0, i)
		label := dayNames[d.Weekday()]
		if i == 0 {
			label = "Hari Ini"
		} else if i == 1 {
			label = "Besok"
		}
		dates = append(dates, fiber.Map{
			"Label":   label,
			"Day":     fmt.Sprintf("%d", d.Day()),
			"Month":   monthNames[int(d.Month())],
			"ISO":     d.Format("2006-01-02"),
			"IsToday": i == 0,
		})
	}

	var schedules []models.Schedule
	database.DB.
		Preload("Film").
		Preload("Studio").
		Where("status = ? AND tanggal_tayang >= ?", "Aktif", now.Format("2006-01-02")).
		Order("tanggal_tayang ASC, jam_tayang ASC").
		Find(&schedules)

	type FilmSlot struct {
		FilmID     uint
		Film       string
		Studio     string
		Tanggal    string
		TanggalISO string
		Harga      string
		HargaNum   string
		Status     string
		Poster     string
		Jams       []string
	}

	groupMap := map[string]*FilmSlot{}
	groupOrder := []string{}
	for _, j := range schedules {
		key := fmt.Sprintf("%d|%s", j.FilmID, j.TanggalTayang.Format("2006-01-02"))
		tanggalStr := j.TanggalTayang.Format("2 Jan 2006")
		hargaNumStr := fmt.Sprintf("%.0f", j.Harga)

		if existing, ok := groupMap[key]; ok {
			existing.Jams = append(existing.Jams, j.JamTayang)
		} else {
			slot := &FilmSlot{
				FilmID:     j.FilmID,
				Film:       j.Film.Judul,
				Studio:     j.Studio.NamaStudio,
				Tanggal:    tanggalStr,
				TanggalISO: j.TanggalTayang.Format("2006-01-02"),
				Harga:      fmt.Sprintf("Rp %s", models.FormatRupiah(int(j.Harga))),
				HargaNum:   hargaNumStr,
				Status:     j.Status,
				Poster:     j.Film.Poster,
				Jams:       []string{j.JamTayang},
			}
			groupMap[key] = slot
			groupOrder = append(groupOrder, key)
		}
	}

	jadwals := make([]fiber.Map, 0, len(groupOrder))
	for _, key := range groupOrder {
		s := groupMap[key]
		jadwals = append(jadwals, fiber.Map{
			"FilmID":     s.FilmID,
			"Film":       s.Film,
			"Studio":     s.Studio,
			"Tanggal":    s.Tanggal,
			"TanggalISO": s.TanggalISO,
			"Harga":      s.Harga,
			"HargaNum":   s.HargaNum,
			"Status":     s.Status,
			"Poster":     s.Poster,
			"Jams":       s.Jams,
		})
	}

	userID := c.Cookies("user_id")
	var user models.User
	database.DB.First(&user, userID)

	initial := ""
	if len(user.Nama) > 0 {
		initial = string(user.Nama[0])
	}

	return c.Render("user/tiket/index", fiber.Map{
		"Title":   "Tiket",
		"Active":  "tiket",
		"Initial": initial,
		"Nama":    user.Nama,
		"Dates":   dates,
		"Jadwals": jadwals,
	}, "layouts/user")
}

func (UserController) TiketBeli(c *fiber.Ctx) error {
	id := c.Params("id")
	var film models.Film
	if err := database.DB.First(&film, id).Error; err != nil {
		return c.Redirect("/tiket")
	}

	var schedules []models.Schedule
	database.DB.
		Preload("Studio").
		Where("film_id = ? AND status = ?", id, "Aktif").
		Order("jam_tayang ASC").
		Find(&schedules)

	if len(schedules) == 0 {
		return c.Redirect("/tiket")
	}

	var scheduleIDs []uint
	for _, s := range schedules {
		scheduleIDs = append(scheduleIDs, s.ID)
	}

	type takenRow struct {
		JadwalID  uint
		NomorKursi string
	}
	var takenRows []takenRow
	database.DB.Model(&models.Tiket{}).
		Select("jadwal_id, nomor_kursi").
		Where("jadwal_id IN (?)", scheduleIDs).
		Find(&takenRows)
	takenMap := map[uint]map[string]bool{}
	for _, row := range takenRows {
		if takenMap[row.JadwalID] == nil {
			takenMap[row.JadwalID] = map[string]bool{}
		}
		takenMap[row.JadwalID][row.NomorKursi] = true
	}

	now := time.Now()
	curMin := now.Hour()*60 + now.Minute()

	times := make([]fiber.Map, 0, len(schedules))
	var soonestIdx int
	soonestDiff := 24 * 60

	for i, s := range schedules {
		parts := strings.Split(s.JamTayang, ":")
		tMin := 0
		if len(parts) >= 2 {
			h, _ := strconv.Atoi(parts[0])
			min, _ := strconv.Atoi(parts[1])
			tMin = h*60 + min
		}
		diff := tMin - curMin
		if diff < 0 {
			diff += 24 * 60
		}
		if diff < soonestDiff {
			soonestDiff = diff
			soonestIdx = i
		}

		times = append(times, fiber.Map{
			"ScheduleID":    s.ID,
			"Jam":           s.JamTayang,
			"Studio":        s.Studio.NamaStudio,
			"Tipe":          s.Studio.Tipe,
			"Baris":         s.Studio.JumlahBaris,
			"KursiPerBaris": s.Studio.JumlahKolom,
			"Harga":         int(s.Harga),
			"HargaFmt":      models.FormatRupiah(int(s.Harga)),
			"Active":        false,
		})
	}

	if len(times) > 0 && soonestIdx < len(times) {
		times[soonestIdx]["Active"] = true
	}

	rowLetters := func(n int) []string {
		letters := []string{}
		for i := 0; i < n; i++ {
			letters = append(letters, string(rune('A'+i)))
		}
		return letters
	}

	seatGrids := map[string][]fiber.Map{}
	gridKeys := []string{}
	for _, s := range schedules {
		gridKey := fmt.Sprintf("sched-%d", s.ID)
		gridKeys = append(gridKeys, gridKey)
		taken := takenMap[s.ID]
		baris := s.Studio.JumlahBaris
		kursi := s.Studio.JumlahKolom
		aisleAt := kursi/2 + kursi%2

		var grid []fiber.Map
		for _, r := range rowLetters(baris) {
			var seats []fiber.Map
			for i := 1; i <= kursi; i++ {
				code := fmt.Sprintf("%s%02d", r, i)
				status := "available"
				if taken[code] {
					status = "taken"
				}
				seats = append(seats, fiber.Map{"Code": code, "Status": status, "Aisle": i == aisleAt})
			}
			grid = append(grid, fiber.Map{"Row": r, "Seats": seats})
		}
		seatGrids[gridKey] = grid
	}

	activeKey := fmt.Sprintf("sched-%d", schedules[soonestIdx].ID)
	selectedJam := schedules[soonestIdx].JamTayang
	selectedStudio := schedules[soonestIdx].Studio.NamaStudio
	selectedTipe := schedules[soonestIdx].Studio.Tipe

	userID := c.Cookies("user_id")
	var user models.User
	database.DB.First(&user, userID)

	initial := ""
	if len(user.Nama) > 0 {
		initial = string(user.Nama[0])
	}

	tanggalStr := schedules[soonestIdx].TanggalTayang.Format("Monday, 2 January 2006")

	return c.Render("user/tiket/beli_tiket", fiber.Map{
		"Title":      "Beli Tiket",
		"Active":     "tiket",
		"Initial":    initial,
		"Nama":       user.Nama,
		"ID":         id,
		"Tanggal":    tanggalStr,
		"Film": fiber.Map{
			"ID":          film.ID,
			"Judul":       film.Judul,
			"Genre":       film.Genre,
			"Durasi":      fmt.Sprintf("%dj %dm", film.Durasi/60, film.Durasi%60),
			"Synopsis":    film.Sinopsis,
			"Rating":      film.Rating,
			"Poster":      film.Poster,
		},
		"Times":      times,
		"Jam":        selectedJam,
		"Studio":     selectedStudio,
		"Tipe":       selectedTipe,
		"SeatGrids":  seatGrids,
		"ActiveGrid": activeKey,
	}, "layouts/user")
}

func (UserController) TiketBeliSubmit(c *fiber.Ctx) error {
	scheduleID := c.FormValue("schedule_id")
	seatsRaw := c.FormValue("seats")
	if scheduleID == "" || seatsRaw == "" {
		return c.Redirect("/tiket")
	}

	var schedule models.Schedule
	if err := database.DB.First(&schedule, scheduleID).Error; err != nil {
		return c.Redirect("/tiket")
	}

	userID, _ := strconv.Atoi(c.Cookies("user_id"))
	seats := strings.Split(seatsRaw, ",")
	hargaPerKursi := int(schedule.Harga)
	totalHarga := len(seats) * hargaPerKursi
	kodeBooking := fmt.Sprintf("TT-%s-%04d", time.Now().Format("20060102"), uint(time.Now().UnixNano()%10000))

	transaksi := models.Transaksi{
		UserID:      uint(userID),
		JadwalID:    schedule.ID,
		TotalHarga:  totalHarga,
		Status:      "pending",
		KodeBooking: kodeBooking,
	}
	if err := database.DB.Create(&transaksi).Error; err != nil {
		return c.Redirect("/tiket?error=gagal")
	}

	for _, seat := range seats {
		seat = strings.TrimSpace(seat)
		if seat == "" {
			continue
		}
		tiket := models.Tiket{
			TransaksiID: transaksi.ID,
			JadwalID:    schedule.ID,
			NomorKursi:  seat,
			Harga:       hargaPerKursi,
			Status:      "aktif",
		}
		database.DB.Create(&tiket)
	}

	return c.Redirect(fmt.Sprintf("/tiket/bayar/%d", transaksi.ID))
}

func (UserController) TiketBayar(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaksi models.Transaksi
	if err := database.DB.
		Preload("Schedule.Film").
		Preload("Schedule.Studio").
		Preload("Tiket").
		First(&transaksi, id).Error; err != nil {
		return c.Redirect("/tiket-saya")
	}

	userID := c.Cookies("user_id")
	if fmt.Sprintf("%d", transaksi.UserID) != userID {
		return c.Redirect("/tiket-saya")
	}

	if transaksi.Status != "pending" {
		return c.Redirect("/tiket-saya")
	}

	schedule := transaksi.Schedule

	seats := []string{}
	for _, t := range transaksi.Tiket {
		seats = append(seats, t.NomorKursi)
	}
	seatsStr := strings.Join(seats, ",")

	var user models.User
	database.DB.First(&user, userID)

	initial := ""
	if len(user.Nama) > 0 {
		initial = string(user.Nama[0])
	}

	tanggalStr := schedule.TanggalTayang.Format("2 Jan 2006")
	studioStr := fmt.Sprintf("%s - %s", schedule.Studio.NamaStudio, schedule.Studio.Tipe)

	return c.Render("user/tiket/bayar", fiber.Map{
		"Title":     "Pembayaran",
		"Active":    "tiket",
		"Initial":   initial,
		"Nama":      user.Nama,
		"Email":     user.Email,
		"ID":        id,
		"Judul":     schedule.Film.Judul,
		"Genre":     fmt.Sprintf("%s, %dj %dm", schedule.Film.Genre, schedule.Film.Durasi/60, schedule.Film.Durasi%60),
		"Tanggal":   tanggalStr,
		"Jam":       schedule.JamTayang,
		"Studio":    studioStr,
		"Tipe":      schedule.Studio.Tipe,
		"Seats":     seats,
		"SeatsStr":  seatsStr,
		"SeatCount": len(seats),
		"Subtotal":  models.FormatRupiah(transaksi.TotalHarga),
		"Total":     models.FormatRupiah(transaksi.TotalHarga),
	}, "layouts/user")
}

func (UserController) TiketBayarSubmit(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaksi models.Transaksi
	if err := database.DB.First(&transaksi, id).Error; err != nil {
		return c.Redirect("/tiket-saya")
	}

	userID := c.Cookies("user_id")
	if fmt.Sprintf("%d", transaksi.UserID) != userID {
		return c.Redirect("/tiket-saya")
	}

	if transaksi.Status != "pending" {
		return c.Redirect("/tiket-saya")
	}

	transaksi.MetodeBayar = "qris"
	transaksi.Status = "lunas"
	database.DB.Save(&transaksi)

	return c.Redirect("/tiket/berhasil/" + transaksi.KodeBooking)
}

func (UserController) TiketBayarUlang(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaksi models.Transaksi
	if err := database.DB.First(&transaksi, id).Error; err != nil {
		return c.Redirect("/tiket-saya")
	}

	userID := c.Cookies("user_id")
	if fmt.Sprintf("%d", transaksi.UserID) != userID {
		return c.Redirect("/tiket-saya")
	}

	if transaksi.Status != "pending" {
		return c.Redirect("/tiket-saya")
	}

	return c.Redirect(fmt.Sprintf("/tiket/bayar/%d", transaksi.ID))
}

func (UserController) TiketBerhasil(c *fiber.Ctx) error {
	kodeBooking := c.Params("id")

	var transaksi models.Transaksi
	if err := database.DB.
		Preload("Schedule.Film").
		Preload("Schedule.Studio").
		Preload("Tiket").
		Where("kode_booking = ?", kodeBooking).
		First(&transaksi).Error; err != nil {
		return c.Redirect("/tiket-saya")
	}

	userID := c.Cookies("user_id")
	var user models.User
	database.DB.First(&user, userID)

	initial := ""
	if len(user.Nama) > 0 {
		initial = string(user.Nama[0])
	}

	seats := []string{}
	for _, t := range transaksi.Tiket {
		seats = append(seats, t.NomorKursi)
	}
	seatStr := strings.Join(seats, ", ")
	tanggalStr := transaksi.Schedule.TanggalTayang.Format("Monday, 2 January 2006")
	jamStr := transaksi.Schedule.JamTayang + " WIB"
	studioStr := fmt.Sprintf("%s · %s", transaksi.Schedule.Studio.NamaStudio, transaksi.Schedule.Studio.Tipe)

	durasiStr := fmt.Sprintf("%dj %dm", transaksi.Schedule.Film.Durasi/60, transaksi.Schedule.Film.Durasi%60)

	return c.Render("user/tiket/berhasil", fiber.Map{
		"Title":       "Pembayaran Berhasil",
		"Active":      "tiket",
		"Initial":     initial,
		"Nama":        user.Nama,
		"KodeBooking": transaksi.KodeBooking,
		"Seats":       seatStr,
		"SeatCount":   len(seats),
		"Subtotal":    models.FormatRupiah(transaksi.TotalHarga),
		"Total":       models.FormatRupiah(transaksi.TotalHarga),
		"MetodeBayar": transaksi.MetodeBayar,
		"Judul":       transaksi.Schedule.Film.Judul,
		"Genre":       transaksi.Schedule.Film.Genre,
		"Durasi":      durasiStr,
		"Tanggal":     tanggalStr,
		"Jam":         jamStr,
		"Studio":      studioStr,
		"Tipe":        transaksi.Schedule.Studio.Tipe,
	}, "layouts/user")
}

func (UserController) TiketSaya(c *fiber.Ctx) error {
	userID := c.Cookies("user_id")

	var transaksis []models.Transaksi
	database.DB.
		Preload("Tiket").
		Preload("Schedule.Film").
		Preload("Schedule.Studio").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&transaksis)

	now := time.Now()

	type ticketItem struct {
		ID          uint
		Judul       string
		Genre       string
		Durasi      string
		Tanggal     string
		TanggalISO  string
		Jam         string
		JamISO      string
		Studio      string
		Tipe        string
		Seats       []string
		KodeBooking string
		TotalNum    string
		Status      string
		Poster      string
		FilmID      string
	}

	tickets := make([]ticketItem, 0)
	for _, trx := range transaksis {
		jamISO := trx.Schedule.JamTayang
		showTime, _ := time.Parse("2006-01-02 15:04", trx.Schedule.TanggalTayang.Format("2006-01-02")+" "+jamISO)
		status := trx.Status
		if status == "pending" && now.After(showTime) {
			status = "Terlewat"
		} else if status == "pending" {
			status = "Pending"
		} else if status == "selesai" || status == "Selesai" || status == "lunas" {
			status = "Lunas"
		}

		seats := []string{}
		for _, t := range trx.Tiket {
			seats = append(seats, t.NomorKursi)
		}
		durasi := fmt.Sprintf("%dj %dm", trx.Schedule.Film.Durasi/60, trx.Schedule.Film.Durasi%60)
		tanggalStr := trx.Schedule.TanggalTayang.Format("2 Jan 2006")

		tickets = append(tickets, ticketItem{
			ID:          trx.ID,
			Judul:       trx.Schedule.Film.Judul,
			Genre:       trx.Schedule.Film.Genre,
			Durasi:      durasi,
			Tanggal:     tanggalStr,
			TanggalISO:  trx.Schedule.TanggalTayang.Format("2006-01-02"),
			Jam:         jamISO + " WIB",
			JamISO:      jamISO,
			Studio:      trx.Schedule.Studio.NamaStudio,
			Tipe:        trx.Schedule.Studio.Tipe,
			Seats:       seats,
			KodeBooking: trx.KodeBooking,
			TotalNum:    models.FormatRupiah(trx.TotalHarga),
			Status:      status,
			Poster:      trx.Schedule.Film.Poster,
			FilmID:      fmt.Sprintf("%d", trx.Schedule.FilmID),
		})
	}

	var user models.User
	database.DB.First(&user, userID)

	initial := ""
	if len(user.Nama) > 0 {
		initial = string(user.Nama[0])
	}

	return c.Render("user/tiket_saya/index", fiber.Map{
		"Title":   "Tiket Saya",
		"Active":  "tiket-saya",
		"Initial": initial,
		"Nama":    user.Nama,
		"Tickets": tickets,
	}, "layouts/user")
}

func (UserController) LihatTiket(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaksi models.Transaksi
	if err := database.DB.
		Preload("Tiket").
		Preload("Schedule.Film").
		Preload("Schedule.Studio").
		First(&transaksi, id).Error; err != nil {
		return c.Redirect("/tiket-saya")
	}

	userID := c.Cookies("user_id")
	if fmt.Sprintf("%d", transaksi.UserID) != userID {
		return c.Redirect("/tiket-saya")
	}

	now := time.Now()
	jamISO := transaksi.Schedule.JamTayang
	showTime, _ := time.Parse("2006-01-02 15:04", transaksi.Schedule.TanggalTayang.Format("2006-01-02")+" "+jamISO)

	status := transaksi.Status
	if status == "pending" && now.After(showTime) {
		status = "Terlewat"
	} else if status == "pending" {
		status = "Pending"
	} else if status == "selesai" || status == "Selesai" || status == "lunas" {
		status = "Lunas"
	}

	durasi := fmt.Sprintf("%dj %dm", transaksi.Schedule.Film.Durasi/60, transaksi.Schedule.Film.Durasi%60)
	genre := fmt.Sprintf("%s · %s · %s", transaksi.Schedule.Film.Genre, durasi, transaksi.Schedule.Studio.Tipe)
	tanggalStr := transaksi.Schedule.TanggalTayang.Format("Monday, 2 January 2006")
	jamStr := jamISO + " WIB"

	tickets := make([]fiber.Map, 0, len(transaksi.Tiket))
	for i, t := range transaksi.Tiket {
		tickets = append(tickets, fiber.Map{
			"ID":          t.ID,
			"Judul":       transaksi.Schedule.Film.Judul,
			"Genre":       genre,
			"Durasi":      durasi,
			"Tanggal":     tanggalStr,
			"Jam":         jamStr,
			"Studio":      transaksi.Schedule.Studio.NamaStudio,
			"Tipe":        transaksi.Schedule.Studio.Tipe,
			"SeatCode":    t.NomorKursi,
			"KodeBooking": fmt.Sprintf("%s-%02d", transaksi.KodeBooking, i+1),
			"KodeQR":      fmt.Sprintf("%s-%02d-%s", transaksi.KodeBooking, i+1, t.NomorKursi),
			"Status":      status,
		})
	}

	var user models.User
	database.DB.First(&user, userID)

	initial := ""
	if len(user.Nama) > 0 {
		initial = string(user.Nama[0])
	}

	tj, _ := json.Marshal(tickets)
	return c.Render("user/tiket_saya/lihat_tiket", fiber.Map{
		"Title":       "Lihat Tiket",
		"Active":      "tiket-saya",
		"Initial":     initial,
		"Nama":        user.Nama,
		"ID":          id,
		"Ticket":      transaksi,
		"Tickets":     tickets,
		"TicketsJSON": string(tj),
	}, "layouts/user")
}

func (UserController) Profile(c *fiber.Ctx) error {
	userID := c.Cookies("user_id")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Redirect("/login")
	}

	initial := ""
	if len(user.Nama) > 0 {
		initial = string(user.Nama[0])
	}

	return c.Render("user/profile", fiber.Map{
		"Title":   "Profil Saya",
		"Active":  "profile",
		"Initial": initial,
		"Nama":    user.Nama,
		"Email":   user.Email,
		"Role":    user.Role,
	}, "layouts/user")
}

func (UserController) ProfileUpdate(c *fiber.Ctx) error {
	userID := c.Cookies("user_id")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Redirect("/login")
	}

	user.Nama = c.FormValue("nama")
	user.Email = c.FormValue("email")

	currentPassword := c.FormValue("current_password")
	newPassword := c.FormValue("new_password")
	confirmPassword := c.FormValue("confirm_password")

	if newPassword != "" {
		err := bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(currentPassword),
		)
		if err != nil {
			return c.Redirect("/profile?error=Password+lama+salah")
		}
		if newPassword != confirmPassword {
			return c.Redirect("/profile?error=Konfirmasi+password+tidak+cocok")
		}
		hash, _ := bcrypt.GenerateFromPassword(
			[]byte(newPassword),
			bcrypt.DefaultCost,
		)
		user.Password = string(hash)
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Redirect("/profile?error=Gagal+menyimpan")
	}

	return c.Redirect("/profile?success=1")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
