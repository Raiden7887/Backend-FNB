# BE Kafe - Backend API

Aplikasi backend untuk sistem kafe yang menyediakan API untuk autentikasi pengguna dan manajemen menu.

## Setup Database

### 1. Buat Database dan Tabel

Jalankan script SQL berikut di MySQL:

```sql
-- Create database if not exists
CREATE DATABASE IF NOT EXISTS db_kafe;

-- Use the database
USE db_kafe;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create menus table
CREATE TABLE IF NOT EXISTS menus (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    foto TEXT NOT NULL,
    harga INT NOT NULL,
    deskripsi TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create index
CREATE INDEX idx_username ON users(username);
```

### 2. Konfigurasi Environment Variables

Buat file `.env` di root directory dengan konfigurasi berikut:

```env
DB_USER=root
DB_PW=26j13a05
DB_HOST=localhost
DB_PORT=3306
DB_NAME=db_kafe
JWT_SECRET=your-secret-key-here
```

## Menjalankan Aplikasi

1. Pastikan MySQL server berjalan
2. Pastikan database dan tabel sudah dibuat
3. Jalankan aplikasi:

```bash
go run main.go
```

Aplikasi akan berjalan di port 8000.

## API Endpoints

### Public Endpoints (Tidak memerlukan autentikasi)

#### 1. Register User
- **URL**: `POST /register`
- **Body**:
```json
{
    "username": "user123",
    "password": "password123"
}
```

#### 2. Login User
- **URL**: `POST /login`
- **Body**:
```json
{
    "username": "user123",
    "password": "password123"
}
```
- **Response**:
```json
{
    "message": "Login successful",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "username": "user123"
    }
}
```

### Protected Endpoints (Memerlukan Bearer Token)

**Catatan**: Semua endpoint di bawah ini memerlukan header `Authorization: Bearer <token>`

#### 3. Get User Info
- **URL**: `GET /user?username=user123`
- **Headers**: `Authorization: Bearer <token>`

#### 4. Create Menu
- **URL**: `POST /menu`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
    "nama": "Nasi Goreng",
    "foto": "https://example.com/nasi-goreng.jpg",
    "harga": 20000,
    "deskripsi": "Nasi goreng spesial dengan telur dan ayam"
}
```

#### 5. Edit Menu
- **URL**: `PUT /menu/{id}`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
```json
{
    "nama": "Nasi Goreng Spesial",
    "foto": "https://example.com/nasi-goreng-spesial.jpg",
    "harga": 25000,
    "deskripsi": "Nasi goreng spesial dengan telur, ayam, dan sayuran"
}
```

#### 6. Delete Menu
- **URL**: `DELETE /menu/{id}`
- **Headers**: `Authorization: Bearer <token>`

## Cara Test dengan cURL

### 1. Register User
```bash
curl -X POST http://localhost:8000/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 2. Login dan Dapatkan Token
```bash
curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 3. Gunakan Token untuk Akses Protected Endpoint
```bash
# Ganti <TOKEN> dengan token yang didapat dari login
curl -X GET "http://localhost:8000/user?username=testuser" \
  -H "Authorization: Bearer <TOKEN>"
```

### 4. Create Menu dengan Token
```bash
curl -X POST http://localhost:8000/menu \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "nama": "Mie Goreng",
    "foto": "https://example.com/mie-goreng.jpg",
    "harga": 18000,
    "deskripsi": "Mie goreng dengan bumbu special"
  }'
```

## JWT Token

- **Expiration**: Token berlaku selama 24 jam
- **Algorithm**: HS256
- **Secret Key**: Dapat diubah di `utils/jwt.go` (dalam produksi gunakan environment variable)

## Dependencies

Pastikan semua dependencies terinstall:

```bash
go mod tidy
```

Dependencies yang digunakan:
- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/gorilla/mux` - HTTP router
- `github.com/golang-jwt/jwt/v5` - JWT token
- `golang.org/x/crypto/bcrypt` - Password hashing 