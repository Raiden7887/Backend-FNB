package models

// Menu represents the menu data structure for JSON request/response
type Menu struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama"`
	Foto      string `json:"foto"`
	Harga     int    `json:"harga"`
	Deskripsi string `json:"deskripsi"`
}
