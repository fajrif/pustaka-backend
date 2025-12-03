package helpers

import (
	// "fmt"
	"errors"
	"pustaka-backend/models"
)

// ValidateProjectData melakukan validasi bisnis lintas bidang pada struct Project
func ValidateProjectData(p *models.Project) error {

    // Jika TanggalMulai adalah zero value (kasus edge jika validasi frontend gagal)
    if p.TanggalMulai.IsZero() {
        return errors.New("Tanggal mulai tidak valid atau kosong")
    }

    // 1. Validasi tanggal_selesai > tanggal_mulai
    if p.TanggalSelesai != nil {
        // Kita ingin tanggal selesai STRICTLY AFTER tanggal mulai
        if !p.TanggalSelesai.After(p.TanggalMulai) {
            return errors.New("Tanggal selesai harus setelah tanggal mulai")
        }
    }

    // 2. Validasi tanggal_perjanjian >= tanggal_mulai
		// fmt.Printf("Tanggal perjanjian: %+v\n", p.TanggalPerjanjian)
		// fmt.Printf("Tanggal mulai: %+v\n", p.TanggalMulai)
		if p.TanggalPerjanjian != nil {
        // Kita cek apakah TanggalPerjanjian persis sama atau sebelumnya
        isBeforeOrEqual := p.TanggalPerjanjian.Before(p.TanggalMulai) ||
                           (p.TanggalPerjanjian.Year() == p.TanggalMulai.Year() &&
                            p.TanggalPerjanjian.Month() == p.TanggalMulai.Month() &&
                            p.TanggalPerjanjian.Day() == p.TanggalMulai.Day())

        // Jika kondisinya TIDAK BeforeOrEqual, berarti tanggal perjanjian lebih besar (Setelah)
        if !isBeforeOrEqual {
            return errors.New("Tanggal perjanjian tidak boleh setelah tanggal mulai")
        }
    }

    // 3. Validasi management_fee < nilai_pekerjaan
    if p.ManagementFee != nil {
        // Pastikan Anda melakukan dereferensi pointer dengan '*'
        if *p.ManagementFee >= p.NilaiPekerjaan {
            return errors.New("Management fee tidak boleh melebihi nilai pekerjaan")
        }
    }

    return nil
}

// ValidateTransactionData melakukan validasi pada struct Transaction
func ValidateTransactionData(t *models.Transaction) error {

    // Validasi JumlahRealisasi: wajib diisi dan > 100000
    if t.JumlahRealisasi <= 100000 {
        return errors.New("Jumlah realisasi wajib diisi dan harus lebih besar dari 100.000")
    }

    return nil
}
