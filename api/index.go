package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "Kecap", Harga: 12000, Stok: 28},
}

func Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// health check
	if path == "/health" {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
		})
		return
	}

	// /api/produk
	if path == "/api/produk" {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(produk)
			return
		}
	}

	// /api/produk/{id}
	if strings.HasPrefix(path, "/api/produk/") {
		idStr := strings.TrimPrefix(path, "/api/produk/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		for _, p := range produk {
			if p.ID == id {
				json.NewEncoder(w).Encode(p)
				return
			}
		}
		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}

	http.NotFound(w, r)
}
