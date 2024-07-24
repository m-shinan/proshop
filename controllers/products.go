package controllers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/database"
	"github.com/m-shinan/project-shop/models"
)

func AdminProducts(c *gin.Context) {
	products, shouldReturn := GetProducts(c)
	if shouldReturn {
		return
	}

	categories, err := GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.HTML(http.StatusOK, "admin_product.html", gin.H{
		"Products":           products,
		"Categories":         categories,
		"ProductSearchTerm":  c.Query("search"),
		"CategorySearchTerm": c.Query("search"),
	})
}

func GetProducts(c *gin.Context) ([]models.Products, bool) {
	db := database.DB
	rows, err := db.Query(`
        SELECT p.id, p.product_name, p.description, p.stock, p.price, p.category_id, p.image, c.category_name
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return nil, true
	}
	defer rows.Close()

	var products []models.Products
	for rows.Next() {
		var product models.Products
		var categoryName string
		if err := rows.Scan(&product.ID, &product.ProductName, &product.Description, &product.Stock, &product.Price, &product.CategoryId, &product.Image, &categoryName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product"})
			return nil, true
		}
		product.Category.CategoryName = categoryName
		products = append(products, product)
	}
	return products, false
}

type Product struct {
	ProductName string  `form:"product_name" binding:"required"`
	Description string  `form:"description" binding:"required"`
	Stock       int     `form:"stock" binding:"required"`
	Price       float64 `form:"price" binding:"required"`
	CategoryId  uint    `form:"category_id" binding:"required"`
	Image       string
}

func AdminAddProduct(c *gin.Context) {
	// Extract form values manually
	productName := c.PostForm("product_name")
	description := c.PostForm("description")
	stockStr := c.PostForm("stock")
	priceStr := c.PostForm("price")
	categoryIdStr := c.PostForm("category_id")

	// Convert string values to appropriate types
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock value"})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
		return
	}

	categoryId, err := strconv.ParseUint(categoryIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Handle file upload
	file, err := c.FormFile("image")
	var imagePath string
	if err != nil {
		fmt.Println("Error retrieving file:", err)
		imagePath = "default_image.jpg"
	} else {
		// Save the file
		imagePath = filepath.Join("uploads", file.Filename)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
			return
		}
		imagePath = "/" + imagePath
	}

	// Create product struct
	product := Product{
		ProductName: productName,
		Description: description,
		Stock:       stock,
		Price:       price,
		CategoryId:  uint(categoryId),
		Image:       filepath.ToSlash(imagePath),
	}
	product.Image = filepath.ToSlash(imagePath)
	// Prepare the SQL query
	query := `INSERT INTO products (product_name, description, stock, price, category_id, image) VALUES ($1, $2, $3, $4, $5, $6)`

	db := database.DB
	// Execute the query
	_, err = db.Exec(query, product.ProductName, product.Description, product.Stock, product.Price, product.CategoryId, product.Image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Redirect to the product management page
	c.Redirect(http.StatusFound, "/admin/products")
}

func AdminAddProductPage(c *gin.Context) {
	categories, err := GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "admin_add_product.html", gin.H{"Categories": categories})
}

func AdminEditProductPage(c *gin.Context) {
	productID := c.Param("id")
	var product models.Products
	var categories []models.Categories

	db := database.DB

	// Fetch product details
	err := db.QueryRow("SELECT id, product_name, description, stock, price, category_id, image FROM products WHERE id = $1", productID).Scan(
		&product.ID, &product.ProductName, &product.Description, &product.Stock, &product.Price, &product.CategoryId, &product.Image,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product details"})
		return
	}

	// Fetch categories
	rows, err := db.Query("SELECT id, category_name FROM categories")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Categories
		if err := rows.Scan(&category.ID, &category.CategoryName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan category"})
			return
		}
		categories = append(categories, category)
	}

	c.HTML(http.StatusOK, "admin_edit_product.html", gin.H{
		"Product":    product,
		"Categories": categories,
	})
}

func AdminEditProduct(c *gin.Context) {
	var product models.Products

	// Bind form data
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Handle file upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload image"})
		return
	}

	// Save the uploaded file
	imagePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}
	product.Image = filepath.ToSlash(imagePath)

	db := database.DB

	// Update the product in the database
	_, err = db.Exec(
		"UPDATE products SET product_name=$1, description=$2, stock=$3, price=$4, category_id=$5, image=$6 WHERE id=$7",
		product.ProductName, product.Description, product.Stock, product.Price, product.CategoryId, product.Image, product.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/products")
}

func AdminDeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	db := database.DB

	// Delete the product from the database
	_, err := db.Exec("DELETE FROM products WHERE id = $1", productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	log.Printf("Deleted product with ID: %s", productID)

	c.Redirect(http.StatusSeeOther, "/admin/products")
}

///////////////////////////////////   CATEGORY     ////////////////////////////////////////////////

func GetCategories() ([]models.Categories, error) {
	db := database.DB
	rows, err := db.Query("SELECT id, category_name FROM categories")
	if err != nil {
		log.Println("Error fetching categories:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Categories
	for rows.Next() {
		var category models.Categories
		if err := rows.Scan(&category.ID, &category.CategoryName); err != nil {
			log.Println("Error scanning category row:", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	log.Println("Fetched categories:", categories)
	return categories, nil
}

type Category struct {
	ID           int    `json:"id"`
	CategoryName string `form:"category_name" binding:"required"`
}

func AdminAddCategory(c *gin.Context) {
	var category Category

	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prepare the SQL query
	query := `INSERT INTO categories (category_name) VALUES ($1)` // Use $1 for parameterized query for PostgreSQL

	// Execute the query
	db := database.DB
	_, err := db.Exec(query, category.CategoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Redirect to the product management page
	c.Redirect(http.StatusFound, "/admin/products")
}

func AdminEditcatPage(c *gin.Context) {
	categoryID := c.Param("id")
	var category models.Categories

	db := database.DB

	// Fetch category details
	err := db.QueryRow("SELECT id, category_name FROM categories WHERE id = $1", categoryID).Scan(
		&category.ID, &category.CategoryName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category details"})
		return
	}

	c.HTML(http.StatusOK, "admin_edit_cat.html", gin.H{
		"Category": category,
	})
}

func AdminEditCat(c *gin.Context) {
	var category models.Categories

	// Bind form data
	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	log.Printf("Received category data: ID=%d, Name=%s", category.ID, category.CategoryName)

	db := database.DB

	// Update the category in the database
	_, err := db.Exec(
		"UPDATE categories SET category_name=$1 WHERE id=$2",
		category.CategoryName, category.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/products")
}

func AdminDeleteCat(c *gin.Context) {
	categoryID := c.Param("id")

	db := database.DB

	// Delete the category from the database
	_, err := db.Exec("DELETE FROM categories WHERE id = $1", categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	log.Printf("Deleted category with ID: %s", categoryID)

	c.Redirect(http.StatusSeeOther, "/admin/products")
}
