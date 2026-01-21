package main

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

	// HEALTH CHECK
	if r.URL.Path == "/health" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
		return
	}

	// GET & POST /api/produk
	if r.URL.Path == "/api/produk" {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(produk)
			return
		}

		if r.Method == "POST" {
			var produkBaru Produk
			json.NewDecoder(r.Body).Decode(&produkBaru)
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(produkBaru)
			return
		}
	}

	// /api/produk/{id}
	if strings.HasPrefix(r.URL.Path, "/api/produk/") {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			for _, p := range produk {
				if p.ID == id {
					json.NewEncoder(w).Encode(p)
					return
				}
			}

		case "PUT":
			var update Produk
			json.NewDecoder(r.Body).Decode(&update)
			for i := range produk {
				if produk[i].ID == id {
					update.ID = id
					produk[i] = update
					json.NewEncoder(w).Encode(update)
					return
				}
			}

		case "DELETE":
			for i, p := range produk {
				if p.ID == id {
					produk = append(produk[:i], produk[i+1:]...)
					json.NewEncoder(w).Encode(map[string]string{
						"message": "sukses delete",
					})
					return
				}
			}
		}
	}

	http.NotFound(w, r)
}
