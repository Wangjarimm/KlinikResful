package belajarHandler

import (
	"fmt"
	"sistem-informasi-klinik/database"
	"sistem-informasi-klinik/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	{
		var users []model.Pasien
		// Find all users in database
		result := database.DB.Preload("JadwalDokter.Dokter").Preload("JadwalDokter.Hari").Preload("JadwalDokter.Jam").Preload("JadwalDokter.Ruangan").Find(&users)
		// Check for errors during query execution
		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": result.Error.Error(),
			})
		}
		// Return list of users
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Data User Berhasil Ditampilkan!",
			"data":    users,
		})
	}
}

func CreateUser(c *fiber.Ctx) error {
	// Parse request body
	var user model.Pasien
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	if user.TglReservasi.IsZero() {
		user.TglReservasi = time.Now()
	}

	// Peroleh jadwal dokter berdasarkan ID jadwal
	jadwalDokter := model.JadwalDokter{}
	result := database.DB.Preload("Hari").Preload("Jam").First(&jadwalDokter, user.IdJadwal)
	if result.Error != nil {
		return result.Error
	}

	// Tentukan tanggal reservasi berdasarkan hari dan jam jadwal dokter
	user.Reservasi = calculateReservationDate(jadwalDokter.Hari.Hari, jadwalDokter.Jam.Jam)

	// Cek kuota reservasi dokter
	// Cek kuota reservasi dokter
	quotaExceeded, err := isQuotaExceeded(jadwalDokter.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error checking reservation quota",
			"error":   err.Error(),
		})
	}

	if quotaExceeded {
		// Jika kuota terpenuhi, pindahkan reservasi ke minggu depan
		user.Reservasi = calculateNextWeekReservation(jadwalDokter.Hari.Hari, jadwalDokter.Jam.Jam)
	} else {
		// Jika kuota belum terpenuhi, lanjutkan dengan tanggal reservasi yang dihitung sebelumnya
		user.Reservasi = calculateReservationDate(jadwalDokter.Hari.Hari, jadwalDokter.Jam.Jam)
	}

	// Insert new user into database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Return success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Data Berhasil Ditambahkan!",
		"data":    user.Id,
	})
}

func isQuotaExceeded(jadwalID int) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Pasien{}).
		Where("id_jadwal = ?", jadwalID).
		Where("DATE(reservasi) >= CURRENT_DATE").
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count >= 2, nil
}

// Fungsi untuk menghitung tanggal reservasi minggu depan
func calculateNextWeekReservation(hari string, jam string) string {
	now := time.Now()

	// Temukan hari berdasarkan nama hari dari jadwal dokter
	targetDay := getDayOfWeekFromString(hari)
	daysUntilTarget := (targetDay - int(now.Weekday()) + 7) % 7

	// Pindahkan ke minggu depan
	daysUntilTarget += 7

	// Hitung tanggal reservasi
	targetDate := now.AddDate(0, 0, daysUntilTarget)
	reservationDate := fmt.Sprintf("%s %s:00", targetDate.Format("2006-01-02"), jam)

	return reservationDate
}

// Fungsi untuk menghitung tanggal reservasi berdasarkan hari dan jam jadwal dokter
func calculateReservationDate(hari string, jam string) string {
	now := time.Now()

	// Temukan hari berdasarkan nama hari dari jadwal dokter
	targetDay := getDayOfWeekFromString(hari)
	daysUntilTarget := (targetDay - int(now.Weekday()) + 7) % 7

	// Hitung tanggal reservasi
	targetDate := now.AddDate(0, 0, daysUntilTarget)
	reservationDate := fmt.Sprintf("%s %s:00", targetDate.Format("2006-01-02"), jam)

	return reservationDate
}

// Fungsi untuk mendapatkan hari dalam bentuk angka dari nama hari
func getDayOfWeekFromString(day string) int {
	days := map[string]int{
		"Senin":  1,
		"Selasa": 2,
		"Rabu":   3,
		"Kamis":  4,
		"Jum'at": 5, // Jika diperlukan, karena Minggu adalah hari terakhir dalam seminggu
	}

	return days[day]
}

func GetUser(c *fiber.Ctx) error {
	// Get id_user parameter from query string
	Id := c.Query("id")
	// Find user by id in the database
	var user model.Pasien
	result := database.DB.Preload("JadwalDokter.Dokter").Preload("JadwalDokter.Hari").Preload("JadwalDokter.Jam").Preload("JadwalDokter.Ruangan").Where("id = ?", Id).First(&user)
	// Check if user exists
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Tidak Ditemukan!",
		})
	}
	// Check for errors during query
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	// Set default tanggal reservasi jika nil
	if user.TglReservasi.IsZero() {
		user.TglReservasi = time.Now()
	}
	// Return user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data": fiber.Map{
			"id":            user.Id,
			"nama_lengkap":  user.Namalengkap,
			"nik":           user.Nik,
			"jenis_kelamin": user.Jeniskelamin,
			"tempat_lahir":  user.Tempatlahir,
			"tanggal_lahir": user.Tanggallahir,
			"alamat":        user.Alamat,
			"no_hp":         user.Nohp,
			"id_jadwal":     user.IdJadwal,
			"tgl_reservasi": user.TglReservasi,
			"reservasi":     user.Reservasi,
			"jadwal_dokter": user.JadwalDokter,
		},
	})
}

func UpdateUser(c *fiber.Ctx) error {
	// Parse request body
	var newUser model.Pasien
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Get id_user from the parsed body
	id := newUser.Id

	// Find user by id_user in database
	var user model.Pasien
	result := database.DB.First(&user, id)
	// Check if user exists
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Tidak Ditemukan",
		})
	}

	// Update user in database
	result = database.DB.Model(&user).Updates(newUser)
	// Check for errors during update
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User Berhasil Diperbarui!",
		"data":    user,
	})
}
func DeleteUser(c *fiber.Ctx) error {
	// Get id_user parameter from query string
	id := c.Query("id")
	// Find user by id_user in database
	var user model.Pasien
	result := database.DB.First(&user, id)
	// Check if user exists
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Tidak Ditemukan",
		})
	}
	// Delete user from database
	result = database.DB.Delete(&user)
	// Check for errors during deletion
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User Berhasil Dihapus!",
		"data":    user,
	})
}

func GetJadwalById(c *fiber.Ctx) error {
	// Get id_user parameter from request url
	id := c.Params("id")
	// Find user by id_user in database
	var user model.Pasien
	result := database.DB.Preload("Pasien.Id").Preload("JadwalDokter.Dokter").Preload("JadwalDokter.Hari").Preload("JadwalDokter.Jam").Preload("JadwalDokter.Ruangan").Preload("Reservasi.Hari").Preload("Reservasi.Jam").Where("id = ?", id).First(&user, id)
	// Check if user exists
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Tidak Ditemukan!",
		})
	}
	// Check for errors during query
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	// Return user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    user,
	})
}
