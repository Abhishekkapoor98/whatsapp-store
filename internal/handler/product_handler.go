package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"whatsapp-store/internal/domain"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	productUsecase domain.ProductUsecase
}

func NewProductHandler(u domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: u}
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stockCount, _ := strconv.ParseInt(r.FormValue("stock_count"), 10, 32)

	// 1. Get image from form
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. Initialize Cloudinary
	cldEnv := os.Getenv("CLOUDINARY_URL")
	if cldEnv == "" {
		http.Error(w, "CLOUDINARY_URL is not set on the server", http.StatusInternalServerError)
		return
	}

	cld, err := cloudinary.NewFromURL(cldEnv)
	if err != nil {
		http.Error(w, "Failed to connect to Cloudinary", http.StatusInternalServerError)
		return
	}

	// 3. Upload image to Cloudinary directly from memory
	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: "whatsapp_store", // Yeh Cloudinary mein ek folder bana dega
	})
	if err != nil {
		http.Error(w, "Failed to upload image to Cloudinary", http.StatusInternalServerError)
		return
	}

	// 4. Save Cloudinary's secure URL to our database
	product := &domain.Product{
		Name:        name,
		Description: description,
		ImageURL:    uploadResult.SecureURL,
		Price:       price,
		StockCount:  int32(stockCount),
	}

	createdProduct, err := h.productUsecase.AddProduct(r.Context(), product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

func (h *ProductHandler) GetAvailableProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productUsecase.FetchAvailableProducts(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productUsecase.FetchAllProducts(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 32)

	var req struct {
		StockCount int32 `json:"stock_count"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	updatedProduct, err := h.productUsecase.ModifyStock(r.Context(), int32(id), req.StockCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

// DeleteProduct handles the HTTP DELETE request to remove a product
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.productUsecase.RemoveProduct(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Product deleted successfully"}`))
}
