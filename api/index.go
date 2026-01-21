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

	// /health  -> akses: /api/health
	if path == "/health" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
		})
		return
	}

	// /produk -> akses: /api/produk
	if path == "/produk" {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
			return
		}
	}

	// /produk/{id} -> akses: /api/produk/{id}
	if strings.HasPrefix(path, "/produk/") {
		idStr := strings.TrimPrefix(path, "/produk/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		for _, p := range produk {
			if p.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(p)
				return
			}
		}

		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}

	http.NotFound(w, r)
}
