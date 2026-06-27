package routes

import (
	"github.com/gofiber/fiber/v2"

	"toktik/controllers"
	"toktik/middleware"
)

func Web(app *fiber.App) {
	auth := controllers.AuthController{}
	usr := controllers.UserController{}
	adm := controllers.AdminController{}

	app.Get("/",
    middleware.AuthOnly(),
    usr.Beranda,
)

	app.Get("/login", auth.LoginPage)
	app.Post("/login", auth.LoginSubmit)
	app.Get("/register", auth.RegisterPage)
	app.Post("/register", auth.RegisterSubmit)
	app.Get("/logout", auth.Logout)

	app.Get("/tiket", middleware.AuthOnly(), usr.TiketIndex)

	app.Get("/tiket/beli/:id",
		middleware.AuthOnly(),
		usr.TiketBeli)

	app.Post("/tiket/beli/:id",
		middleware.AuthOnly(),
		usr.TiketBeliSubmit)

	app.Get("/tiket/bayar/:id",
		middleware.AuthOnly(),
		usr.TiketBayar)

	app.Post("/tiket/bayar/:id",
		middleware.AuthOnly(),
		usr.TiketBayarSubmit)

	app.Get("/tiket/bayar-ulang/:id",
		middleware.AuthOnly(),
		usr.TiketBayarUlang)

	app.Get("/tiket/berhasil/:id",
		middleware.AuthOnly(),
		usr.TiketBerhasil)

	app.Get("/tiket-saya",
		middleware.AuthOnly(),
		usr.TiketSaya)

	app.Get("/tiket-saya/:id",
		middleware.AuthOnly(),
		usr.LihatTiket)

	app.Get("/profile",
		middleware.AuthOnly(),
		usr.Profile)

	app.Post("/profile",
		middleware.AuthOnly(),
		usr.ProfileUpdate)

	admin := app.Group("/admin", middleware.AdminOnly())
	admin.Get("/", adm.Dashboard)

	admin.Get("/film", adm.FilmIndex)
	admin.Get("/film/tambah", adm.FilmTambah)
	admin.Post("/film/tambah", adm.FilmTambahSubmit)
	admin.Get("/film/edit/:id", adm.FilmEdit)
	admin.Post("/film/edit/:id", adm.FilmEditSubmit)
	admin.Post("/film/hapus/:id", adm.FilmHapus)

	admin.Get("/jadwal", adm.JadwalIndex)
	admin.Get("/jadwal/tambah", adm.JadwalTambah)
	admin.Post("/jadwal/tambah", adm.JadwalTambahSubmit)
	admin.Get("/jadwal/edit/:id", adm.JadwalEdit)
	admin.Post("/jadwal/edit/:id", adm.JadwalEditSubmit)
	admin.Post("/jadwal/hapus/:id", adm.JadwalHapus)

	admin.Get("/studio", adm.StudioIndex)
	admin.Get("/studio/tambah", adm.StudioTambah)
	admin.Post("/studio/tambah", adm.StudioTambahSubmit)
	admin.Get("/studio/edit/:id", adm.StudioEdit)
	admin.Post("/studio/edit/:id", adm.StudioEditSubmit)
	admin.Post("/studio/hapus/:id", adm.StudioHapus)

	admin.Get("/transaksi", adm.TransaksiIndex)
	admin.Get("/transaksi/:id", adm.TransaksiDetail)

	admin.Get("/user", adm.UserIndex)

	admin.Post("/user/admin/:id", adm.JadikanAdmin)
	admin.Post("/user/user/:id", adm.JadikanUser)
	admin.Post("/user/hapus/:id", adm.UserHapus)
}
