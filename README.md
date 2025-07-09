# 🍽️ POS Restaurant API (Go + PostgreSQL)

Ini adalah proyek backend sederhana untuk sistem **Point of Sale (POS) restoran**, dibuat sebagai latihan belajar **Golang**. Proyek ini menggunakan framework **Gin**, database **PostgreSQL**, dan arsitektur service-repository yang bersih dan terstruktur.

---

## ✨ Fitur Utama

- 📦 CRUD untuk entitas:
  - Customer, Staff, Outlet, Table
  - Orders, Order Items, Order Item Excluded Ingredients
  - Reservations, Customer Visits, Table Transfer
  - Bill & Split Bill
  - Bill Payments

- 🧾 Perhitungan otomatis:
  - Pajak & service charge dari outlet
  - Pembayaran split & pelacakan status pembayaran
  - Stok bahan berdasarkan item & ingredient yang dipesan

- 🔄 Soft delete (opsional) & validasi data yang konsisten

---

- 🧹 Pembelajaran yang belum sempat diterapkan:
  - Full menggunakan UUID
  - Updated_at menggunakan trigger
  - Beberapa pilihan ENUM dinamis dibuatkan table tersendiri
  - Full menggunakan Begin & Rollback pada Repository
  - Pembuatan model berbeda untuk tampilan yang berbeda
  - Menggunakan OmitEmpty pada model
  - Mencoba menggunakan PGX

---

## 🧰 Teknologi yang Digunakan

- [Golang](https://golang.org/)
- [Gin Gonic](https://github.com/gin-gonic/gin)
- [PostgreSQL](https://www.postgresql.org/)
- [github.com/google/uuid](https://pkg.go.dev/github.com/google/uuid)

---

## Swaggo - Swagger UI
```bash
# To Save/update
swag init --generalInfo cmd/pos-restaurant/main.go --output cmd/pos-restaurant/docs --parseDependency --parseInternal
```

## 🚀 Cara Menjalankan

### 1. Clone Repository

```bash
git clone https://github.com/username/pos-restaurant-go.git
cd pos-restaurant-go
```
### 3. Edit API/database/db-template.go
### 4. Rename to db.go
