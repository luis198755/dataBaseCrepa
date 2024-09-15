-- Crear la base de datos
CREATE DATABASE IF NOT EXISTS creperia_db;
USE creperia_db;

-- Tabla de categorías
CREATE TABLE categorias (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL
);

-- Tabla de productos
CREATE TABLE productos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    precio DECIMAL(6,2) NOT NULL,
    categoria_id INT,
    FOREIGN KEY (categoria_id) REFERENCES categorias(id)
);

-- Tabla de ingredientes
CREATE TABLE ingredientes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL,
    precio_extra DECIMAL(6,2)
);

-- Tabla de relación productos-ingredientes
CREATE TABLE producto_ingrediente (
    producto_id INT,
    ingrediente_id INT,
    PRIMARY KEY (producto_id, ingrediente_id),
    FOREIGN KEY (producto_id) REFERENCES productos(id),
    FOREIGN KEY (ingrediente_id) REFERENCES ingredientes(id)
);

-- Tabla de bebidas
CREATE TABLE bebidas (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL,
    tipo ENUM('fria', 'caliente') NOT NULL,
    precio DECIMAL(6,2) NOT NULL
);

-- Insertar categorías
INSERT INTO categorias (nombre) VALUES
('Crepas Saladas'),
('Especialidades Dulces'),
('Fruta Natural'),
('Crepas Sencillas');

-- Insertar algunos productos (ejemplos)
INSERT INTO productos (nombre, precio, categoria_id) VALUES
('Queso manchego', 63.00, 1),
('Jamón', 69.00, 1),
('Philadelphia, Galletas oreo', 55.00, 2),
('Nutella, Galletas oreo', 55.00, 2),
('Fresa con Nutella', 61.00, 3),
('Manzana con Nutella', 57.00, 3),
('Mermelada de fresa', 39.00, 4);

-- Insertar ingredientes extra
INSERT INTO ingredientes (nombre, precio_extra) VALUES
('Gansito', 15.00),
('Kinder delice', 15.00),
('Nuez', 10.00),
('Almendras', 10.00),
('Arándanos', 10.00),
('Chispas de chocolate Hershey\'s', 10.00),
('Lechera', 10.00),
('Chocolate Hershey\'s', 10.00),
('Bombones', 10.00);

-- Insertar algunas bebidas
INSERT INTO bebidas (nombre, tipo, precio) VALUES
('Soda Italiana', 'fria', 67.00),
('Frappe Nutella', 'fria', 71.00),
('Malteada', 'fria', 67.00),
('Americano', 'caliente', 39.00),
('Capuchino', 'caliente', 45.00),
('Chocolate', 'caliente', 51.00);