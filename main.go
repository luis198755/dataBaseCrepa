package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Producto struct {
	ID          int            `json:"id"`
	Nombre      string         `json:"nombre"`
	Descripcion sql.NullString `json:"descripcion"`
	Precio      float64        `json:"precio"`
	CategoriaID int            `json:"categoria_id"`
}

var db *sql.DB

func main() {
	// Inicializar la conexión a la base de datos
	initDB()

	// Configurar el router de Gin
	gin.SetMode(gin.ReleaseMode)  // Opcional: Establece el modo de producción
	r := gin.Default()

	// Configurar múltiples proxies confiables
	trustedProxies := []string{
		"127.0.0.1",        // localhost
		"10.0.0.0/8",       // Red privada clase A
		"172.16.0.0/12",    // Red privada clase B
		"192.168.0.0/16",   // Red privada clase C
		// Agrega aquí las direcciones IP o rangos CIDR de tus proxies adicionales
	}

	err := r.SetTrustedProxies(trustedProxies)
	if err != nil {
		log.Fatal("Error al configurar proxies confiables: ", err)
	}

	// Rutas para los productos
	r.GET("/productos", getProductos)
	r.GET("/productos/:id", getProducto)
	r.POST("/productos", createProducto)
	r.PUT("/productos/:id", updateProducto)
	r.DELETE("/productos/:id", deleteProducto)

	// Iniciar el servidor
	log.Println("Servidor iniciando en http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
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
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.CategoriaID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		productos = append(productos, p)
	}

	c.JSON(http.StatusOK, productos)
}

func getProducto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var p Producto
	err = db.QueryRow("SELECT id, nombre, descripcion, precio, categoria_id FROM productos WHERE id = ?", id).
		Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.CategoriaID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func createProducto(c *gin.Context) {
	var p Producto
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO productos (nombre, descripcion, precio, categoria_id) VALUES (?, ?, ?, ?)",
		p.Nombre, p.Descripcion, p.Precio, p.CategoriaID)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var p Producto
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec("UPDATE productos SET nombre = ?, descripcion = ?, precio = ?, categoria_id = ? WHERE id = ?",
		p.Nombre, p.Descripcion, p.Precio, p.CategoriaID, id)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	result, err := db.Exec("DELETE FROM productos WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado"})
}
