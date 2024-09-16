package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Producto struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion,omitempty"`
	Precio      float64 `json:"precio"`
	CategoriaID int     `json:"categoria_id"`
}

func (p *Producto) UnmarshalJSON(data []byte) error {
	type Alias Producto
	aux := &struct {
		Descripcion *string `json:"descripcion"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Descripcion != nil {
		p.Descripcion = *aux.Descripcion
	}
	return nil
}

func (p Producto) MarshalJSON() ([]byte, error) {
	type Alias Producto
	return json.Marshal(&struct {
		Descripcion string `json:"descripcion,omitempty"`
		Alias
	}{
		Descripcion: p.Descripcion,
		Alias:       (Alias)(p),
	})
}

var db *sql.DB

func main() {
	// Initialize database connection
	initDB()

	// Configure Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Configure trusted proxies
	trustedProxies := []string{
		"127.0.0.1",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	err := r.SetTrustedProxies(trustedProxies)
	if err != nil {
		log.Fatal("Error configuring trusted proxies: ", err)
	}

	// Define routes
	r.GET("/productos", getProductos)
	r.GET("/productos/:id", getProducto)
	r.POST("/productos", createProducto)
	r.PUT("/productos/:id", updateProducto)
	r.DELETE("/productos/:id", deleteProducto)

	// Start the server
	log.Println("Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func initDB() {
	var err error
	db, err = sql.Open("mysql", "creperia_user:creperia_password@tcp(mariadb:3306)/creperia_db")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func getProductos(c *gin.Context) {
	rows, err := db.Query("SELECT id, nombre, descripcion, precio, categoria_id FROM productos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var productos []Producto
	for rows.Next() {
		var p Producto
		var descripcion sql.NullString
		if err := rows.Scan(&p.ID, &p.Nombre, &descripcion, &p.Precio, &p.CategoriaID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		p.Descripcion = descripcion.String
		productos = append(productos, p)
	}

	c.JSON(http.StatusOK, productos)
}

func getProducto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var p Producto
	var descripcion sql.NullString
	err = db.QueryRow("SELECT id, nombre, descripcion, precio, categoria_id FROM productos WHERE id = ?", id).
		Scan(&p.ID, &p.Nombre, &descripcion, &p.Precio, &p.CategoriaID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	p.Descripcion = descripcion.String
	c.JSON(http.StatusOK, p)
}

func createProducto(c *gin.Context) {
	var p Producto
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO productos (nombre, descripcion, precio, categoria_id) VALUES (?, ?, ?, ?)",
		p.Nombre, sql.NullString{String: p.Descripcion, Valid: p.Descripcion != ""}, p.Precio, p.CategoriaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	c.JSON(http.StatusCreated, p)
}

func updateProducto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var p Producto
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec("UPDATE productos SET nombre = ?, descripcion = ?, precio = ?, categoria_id = ? WHERE id = ?",
		p.Nombre, sql.NullString{String: p.Descripcion, Valid: p.Descripcion != ""}, p.Precio, p.CategoriaID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p.ID = id
	c.JSON(http.StatusOK, p)
}

func deleteProducto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := db.Exec("DELETE FROM productos WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}