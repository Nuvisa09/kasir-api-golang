package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Produk
type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
	Category int `json:"category_id"`
}

type ProdukResponse struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
	Category string `json:"category"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40, Category: 2},
	{ID: 3, Nama: "Kecap", Harga: 12000, Stok: 28, Category: 3},
}

// Kategori

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Produk makanan siap saji"},
	{ID: 2, Name: "Minuman", Description: "Minuman botol dan kemasan"},
	{ID: 3, Name: "Bumbu", Description: "Bumbu dapur dan saus"},
}
var nextCategoryID = 4

// get Kategori keseluruhan

func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Validasi

func categoryExists(id int) bool {
	for _, c := range categories {
		if c.ID == id {
			return true
		}
	}
	return false
}

func getCategoryByIDLocal(id int) (Category, bool){
	for _, c := range categories {
		if c.ID == id {
			return c, true
		}
	}
	return Category{}, false
}

// Get Produk By ID
func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

// Get Kategori By ID
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
}

// Buat Kategori Baru
func createCategory(w http.ResponseWriter, r *http.Request) {

	var newCategory Category
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		http.Error(w, "Data tidak valid", http.StatusBadRequest)
		return
	}

	newCategory.ID = nextCategoryID
	nextCategoryID++
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}


// PUT localhost:8080/api/produk/{id}
func updateProduk(w http.ResponseWriter, r *http.Request) {
	// GET id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}
	// get data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// validasi
			if !categoryExists(updateProduk.Category){
				http.Error(w, "Category tidak ditemukan", http.StatusBadRequest)
				return
			}


	// loop produk, cari id, ganti sesuai request
	for i := range produk {
		if produk [i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}

	
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

// Update Kategori
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var input Category
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				http.Error(w, "Data tidak valid", http.StatusBadRequest)
				return
			}

	for i, c := range categories {
		if c.ID == id {
			
			categories[i].Name = input.Name
			categories[i].Description = input.Description

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories[i])
			return
		}
	}
	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
}

// Delete Produk
func deleteProduk(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	// ganti id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}
	// loop produk, cari ID, dapat index yang mau dihapus
	for i, p := range produk {
		if p.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			produk = append(produk[:i], produk[i+1:]...)
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
		})
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
	
}

// Delete Kategori
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Kategori dihapus",
			})
			return
		}
	}
	http.Error(w, "Kategorii tidak ditemukan", http.StatusNotFound)
}

func main() {

	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk (w, r)
		}

	})

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(produk)
			var result []ProdukResponse

			for _, p := range produk {
				cat, _ := getCategoryByIDLocal(p.Category)

				item := ProdukResponse{
					ID: p.ID,
					Nama: p.Nama,
					Harga: p.Harga,
					Stok: p.Stok,
					Category: cat.Name,
				}

				result = append(result, item)
			}
			json.NewEncoder(w).Encode(result)

		} else if r.Method == "POST" {
			// baca data dari request
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			// validasi
			if !categoryExists(produkBaru.Category){
				http.Error(w, "Category tidak ditemukan", http.StatusBadRequest)
				return
			}

			// masukkin data ke dalam variable produk
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) //201
			json.NewEncoder(w).Encode(produkBaru)
		}

	})

	// Kategori
	http.HandleFunc("/categories", func (w http.ResponseWriter, r *http.Request)  {
		if r.Method == "GET" {
			getCategories(w, r)
		} else if r.Method == "POST" {
			createCategory(w, r)
		}
	})

	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET"{
			getCategoryByID(w, r)
		}else if r.Method == "PUT"{
			updateCategory(w, r)
		}else if r.Method == "DELETE"{
			deleteCategory(w, r)
		}
	})

	//localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})
	fmt.Println("GServer running di localhost:8080")

	// err := http.ListenAndServe(":8080", nil)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Gserver running di port", port)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Gagal running server")
	}
}
