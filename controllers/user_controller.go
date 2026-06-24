package controllers

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"toktik/data"
)

type UserController struct{}

type PaymentMethod struct {
	Kode       string
	Nama       string
	Keterangan string
	Icon       string
	Active     bool
}

func (UserController) Beranda(c *fiber.Ctx) error {
	featured := []fiber.Map{}
	for _, m := range data.Films {
		if m.Featured && len(featured) < 3 {
			featured = append(featured, fiber.Map{
				"ID":          m.ID,
				"Judul":       m.Judul,
				"Genre":       m.Genre,
				"Durasi":      m.Durasi,
				"Rating":      m.Rating,
				"Synopsis":    m.Synopsis,
				"PosterColor": m.PosterColor,
				"Poster":      m.Poster,
				"Harga":       m.Harga,
				"Times":       m.Times,
			})
		}
	}

	now := data.Now()
	curMin := now.Hour()*60 + now.Minute()

	type soonItem struct {
		ID          string
		Judul       string
		Genre       string
		Durasi      string
		PosterColor string
		Poster      string
		Jam         string
		SoonIn      int
	}
	var soon []soonItem
	for _, m := range data.Films {
		var nextJam string
		nextDiff := 24 * 60
		for _, t := range m.Times {
			parts := strings.Split(t, ":")
			if len(parts) < 2 {
				continue
			}
			h, _ := strconv.Atoi(parts[0])
			min, _ := strconv.Atoi(parts[1])
			tMin := h*60 + min
			diff := tMin - curMin
			if diff < 0 {
				diff += 24 * 60
			}
			if diff < nextDiff {
				nextDiff = diff
				nextJam = t
			}
		}
		if nextJam != "" && nextDiff <= 6*60 {
			soon = append(soon, soonItem{
				ID: m.ID, Judul: m.Judul, Genre: m.Genre, Durasi: m.Durasi,
				PosterColor: m.PosterColor, Poster: m.Poster, Jam: nextJam, SoonIn: nextDiff,
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
			"ID":          s.ID,
			"Judul":       s.Judul,
			"Genre":       s.Genre,
			"Durasi":      s.Durasi,
			"PosterColor": s.PosterColor,
			"Poster":      s.Poster,
			"Jam":         s.Jam,
		})
	}

	return c.Render("user/beranda/index", fiber.Map{
		"Title":      "Beranda",
		"Active":     "home",
		"Initial":    "A",
		"Featured":   featured,
		"NowShowing": nowShowing,
	}, "layouts/user")
}

func (UserController) TiketIndex(c *fiber.Ctx) error {
	now := data.Now()
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

	type FilmSlot struct {
		FilmID     uint
		Film       string
		Studio     string
		Tanggal    string
		TanggalISO string
		Harga      string
		HargaNum   string
		Status     string
		FilmBg     string
		Poster     string
		Jams       []string
	}

	groupMap := map[string]*FilmSlot{}
	groupOrder := []string{}
	for _, j := range data.Jadwals {
		key := fmt.Sprintf("%d|%s", j.FilmID, j.TanggalISO)
		if existing, ok := groupMap[key]; ok {
			existing.Jams = append(existing.Jams, j.Jam)
		} else {
			film := data.FindFilmByID(fmt.Sprintf("%d", j.FilmID))
			poster := ""
			jams := []string{}
			if film != nil {
				poster = film.Poster
				jams = film.Times
			}
			slot := &FilmSlot{
				FilmID:     j.FilmID,
				Film:       j.Film,
				Studio:     j.Studio,
				Tanggal:    j.Tanggal,
				TanggalISO: j.TanggalISO,
				Harga:      j.Harga,
				HargaNum:   j.HargaNum,
				Status:     j.Status,
				FilmBg:     j.FilmBg,
				Poster:     poster,
				Jams:       jams,
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
			"FilmBg":     s.FilmBg,
			"Poster":     s.Poster,
			"Jams":       s.Jams,
		})
	}

	return c.Render("user/tiket/index", fiber.Map{
		"Title":   "Tiket",
		"Active":  "tiket",
		"Initial": "A",
		"Dates":   dates,
		"Jadwals": jadwals,
	}, "layouts/user")
}

func (UserController) TiketBeli(c *fiber.Ctx) error {
	id := c.Params("id")
	film := data.FindFilmByID(id)
	if film == nil {
		return c.Redirect("/tiket")
	}

	times := make([]fiber.Map, 0, len(film.Times))
	gridKeys := []string{}
	type timeCandidate struct {
		idx       int
		jam       string
		studio    string
		tipe      string
		soonIn    int
	}
	var candidates []timeCandidate
	now := data.Now()
	curMin := now.Hour()*60 + now.Minute()
	for i, t := range film.Times {
		studioIdx := (i % 4) + 1
		studio := data.Studios[studioIdx]
		tipe := studio.Tipe
		baris, _ := strconv.Atoi(studio.Baris)
		kursi, _ := strconv.Atoi(studio.KursiPerBaris)
		harga := 50000
		if tipe == "Reguler" {
			harga = 35000
		} else if tipe == "IMAX" {
			harga = 75000
		}
		key := studio.Nama + "|" + tipe
		if !contains(gridKeys, key) {
			gridKeys = append(gridKeys, key)
		}
		parts := strings.Split(t, ":")
		if len(parts) >= 2 {
			h, _ := strconv.Atoi(parts[0])
			min, _ := strconv.Atoi(parts[1])
			tMin := h*60 + min
			diff := tMin - curMin
			if diff < 0 {
				diff += 24 * 60
			}
			candidates = append(candidates, timeCandidate{
				idx: i, jam: t, studio: studio.Nama, tipe: tipe, soonIn: diff,
			})
		}
		times = append(times, fiber.Map{
			"Jam":           t,
			"Studio":        studio.Nama,
			"Tipe":          tipe,
			"Baris":         baris,
			"KursiPerBaris": kursi,
			"Harga":         harga,
			"Active":        false,
		})
	}

	soonestIdx := 0
	if len(candidates) > 0 {
		soonestIdx = candidates[0].idx
		best := candidates[0]
		for _, c := range candidates {
			if c.soonIn < best.soonIn {
				best = c
				soonestIdx = c.idx
			}
		}
	}

	soonestJam := film.Times[soonestIdx]
	soonestStudio := data.Studios[(soonestIdx%4)+1]
	soonestTipe := soonestStudio.Tipe
	for i := range times {
		if i == soonestIdx {
			times[i]["Active"] = true
		}
	}

	rowLetters := func(n int) []string {
		letters := []string{}
		for i := 0; i < n; i++ {
			letters = append(letters, string(rune('A'+i)))
		}
		return letters
	}

	buildGrid := func(rows, cols int, taken []string) []fiber.Map {
		takenMap := map[string]bool{}
		for _, t := range taken {
			takenMap[t] = true
		}
		var grid []fiber.Map
		aisleAt := cols/2 + cols%2
		for _, r := range rowLetters(rows) {
			var seats []fiber.Map
			for i := 1; i <= cols; i++ {
				code := fmt.Sprintf("%s%02d", r, i)
				status := "available"
				if takenMap[code] {
					status = "taken"
				}
				seats = append(seats, fiber.Map{"Code": code, "Status": status, "Aisle": i == aisleAt})
			}
			grid = append(grid, fiber.Map{"Row": r, "Seats": seats})
		}
		return grid
	}

	grids := map[string][]fiber.Map{}
	for _, key := range gridKeys {
		parts := strings.Split(key, "|")
		studioName := parts[0]
		for _, s := range data.Studios {
			if s.Nama != studioName {
				continue
			}
			b, _ := strconv.Atoi(s.Baris)
			k, _ := strconv.Atoi(s.KursiPerBaris)
			taken := []string{}
			if s.Nama == "Studio 1" {
				taken = []string{"A03", "A07", "B02", "B08", "C01", "C09", "D04", "D10", "E05", "F02", "F08", "G03", "G07", "H06"}
			} else {
				taken = []string{"A01", "A05", "A09", "B02", "B10", "C04", "C11", "D01", "D08", "E03", "E12", "F05", "F10", "G02", "G09", "H04", "H11", "I06", "I12", "J03", "J08"}
			}
			grids[key] = buildGrid(b, k, taken)
		}
	}

	activeKey := soonestStudio.Nama + "|" + soonestTipe
	if _, ok := grids[activeKey]; !ok {
		for _, k := range gridKeys {
			if strings.HasPrefix(k, soonestStudio.Nama) {
				activeKey = k
				break
			}
		}
	}
	if _, ok := grids[activeKey]; !ok && len(gridKeys) > 0 {
		activeKey = gridKeys[0]
	}

	return c.Render("user/tiket/beli_tiket", fiber.Map{
		"Title":      "Beli Tiket",
		"Active":     "tiket",
		"Initial":    "A",
		"ID":         id,
		"Film":       fiber.Map{"ID": film.ID, "Judul": film.Judul, "Genre": film.Genre, "Durasi": film.Durasi, "Synopsis": film.Synopsis, "Rating": film.Rating, "Poster": film.Poster, "PosterColor": film.PosterColor},
		"Times":      times,
		"Tipe":       soonestTipe,
		"Jam":        soonestJam,
		"Studio":     soonestStudio.Nama,
		"SeatGrids":  grids,
		"ActiveGrid": activeKey,
	}, "layouts/user")
}

func (UserController) TiketBeliSubmit(c *fiber.Ctx) error {
	return c.Redirect("/tiket/bayar/" + c.Params("id") + "?seats=" + c.FormValue("seats"))
}

func (UserController) TiketBayar(c *fiber.Ctx) error {
	id := c.Params("id")
	seatsRaw := c.Query("seats", "B04,C06,C07,E08")
	seats := strings.Split(seatsRaw, ",")
	sort.Strings(seats)

	payments := []PaymentMethod{
		{Kode: "kartu", Nama: "Kartu Kredit / Debit", Keterangan: "Visa, MasterCard, JCB", Icon: "credit-card", Active: true},
		{Kode: "qris", Nama: "QRIS", Keterangan: "Scan QR dari aplikasi pembayaran apa pun", Icon: "qr-code", Active: false},
		{Kode: "ewallet", Nama: "E-Wallet", Keterangan: "OVO, GoPay, DANA, ShopeePay", Icon: "wallet", Active: false},
		{Kode: "transfer", Nama: "Transfer Bank", Keterangan: "BCA, BNI, BRI, Mandiri", Icon: "building-2", Active: false},
		{Kode: "virtual", Nama: "Virtual Account", Keterangan: "Bayar via ATM atau mobile banking", Icon: "hash", Active: false},
	}

	pricePerSeat := 50000
	subtotal := len(seats) * pricePerSeat
	total := subtotal

	return c.Render("user/tiket/bayar", fiber.Map{
		"Title":      "Pembayaran",
		"Active":     "tiket",
		"Initial":    "A",
		"ID":         id,
		"Judul":      "Echoes of the Unknown",
		"Genre":      "Fiksi Ilmiah, 2j 28m",
		"Tanggal":    "23 Juni 2026",
		"Jam":        "20:00",
		"Studio":     "Studio 4 - Premiere",
		"Tipe":       "Premiere",
		"Seats":      seats,
		"SeatCount":  len(seats),
		"Subtotal":   data.FormatRupiah(subtotal),
		"Total":      data.FormatRupiah(total),
		"Payments":   payments,
		"Nama":       "Andi Pratama",
		"Email":      "andi@email.com",
	}, "layouts/user")
}

func (UserController) TiketBayarSubmit(c *fiber.Ctx) error {
	return c.Redirect("/tiket/berhasil/TT-20260623-0042")
}

func (UserController) TiketBerhasil(c *fiber.Ctx) error {
	return c.Render("user/tiket/berhasil", fiber.Map{
		"Title":       "Pembayaran Berhasil",
		"Active":      "tiket",
		"Initial":     "A",
		"KodeBooking": c.Params("id"),
		"Seats":       "B04, C06, C07, E08",
		"SeatCount":   4,
		"Subtotal":    "180.000",
		"Total":       "180.000",
		"MetodeBayar": "Kartu Kredit Visa",
	}, "layouts/user")
}

func (UserController) TiketSaya(c *fiber.Ctx) error {
	now := data.Now()
	tickets := make([]data.Ticket, len(data.Tickets))
	for i, t := range data.Tickets {
		showTime, _ := time.Parse("2006-01-02 15:04", t.TanggalISO+" "+t.JamISO)
		if t.Status == "Pending" && now.After(showTime) {
			t.Status = "Terlewat"
		}
		tickets[i] = t
	}
	return c.Render("user/tiket_saya/index", fiber.Map{
		"Title":   "Tiket Saya",
		"Active":  "tiket-saya",
		"Initial": "A",
		"Tickets": tickets,
	}, "layouts/user")
}

func (UserController) LihatTiket(c *fiber.Ctx) error {
	id := data.ParseID(c.Params("id"))
	ticket := data.FindTicketByID(id)
	if ticket == nil {
		return c.Redirect("/tiket-saya")
	}

	now := data.Now()
	showTime, _ := time.Parse("2006-01-02 15:04", ticket.TanggalISO+" "+ticket.JamISO)
	if ticket.Status == "Pending" && now.After(showTime) {
		ticket.Status = "Terlewat"
	}

	tickets := make([]fiber.Map, 0, len(ticket.Seats))
	for i, seat := range ticket.Seats {
		tickets = append(tickets, fiber.Map{
			"ID":          uint(i + 1),
			"Judul":       ticket.Judul,
			"Genre":       ticket.Genre,
			"Durasi":      ticket.Durasi,
			"Tanggal":     ticket.Tanggal,
			"Jam":         ticket.Jam,
			"Studio":      ticket.Studio,
			"Tipe":        ticket.Tipe,
			"SeatCode":    seat,
			"KodeBooking": fmt.Sprintf("%s-%02d", ticket.KodeBooking, i+1),
			"KodeQR":      fmt.Sprintf("%s-%02d-%s", ticket.KodeBooking, i+1, seat),
			"Status":      ticket.Status,
		})
	}

	tj, _ := json.Marshal(tickets)
	return c.Render("user/tiket_saya/lihat_tiket", fiber.Map{
		"Title":       "Lihat Tiket",
		"Active":      "tiket-saya",
		"Initial":     "A",
		"ID":          c.Params("id"),
		"Ticket":      ticket,
		"Tickets":     tickets,
		"TicketsJSON": string(tj),
	}, "layouts/user")
}

func (UserController) Profile(c *fiber.Ctx) error {
	u := data.Users[0]
	return c.Render("user/profile", fiber.Map{
		"Title":    "Profil Saya",
		"Active":   "profile",
		"Initial":  u.Initial,
		"Nama":     u.Nama,
		"Email":    u.Email,
		"AvatarBg": u.AvatarBg,
	}, "layouts/user")
}

func (UserController) ProfileUpdate(c *fiber.Ctx) error {
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
