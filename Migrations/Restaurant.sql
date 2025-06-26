CREATE TABLE outlets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    service_charge_percentage DECIMAL(5,2) DEFAULT 10,
    tax_percentage DECIMAL(5,2) DEFAULT 10,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TYPE table_type AS ENUM ('Indoor', 'Outdoor');
CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    outlet_id INT REFERENCES outlets(id),
    table_number VARCHAR(50) NOT NULL,
    capacity INT NOT NULL,
    location_type table_type,
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'occupied', 'reserved', 'out_of_order'))
);

-- Contoh untuk disambungkan ke data dari HRD
CREATE TABLE staff (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('waiter', 'cashier', 'chef', 'manager', 'supervisor')),
    pin_code VARCHAR(10) NOT NULL, -- Untuk login POS
    is_active BOOLEAN DEFAULT TRUE
);

-- Menu
CREATE TABLE menu_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    qty INT DEFAULT 0,
    is_allergen BOOLEAN DEFAULT FALSE,  -- Bahan penyebab alergi umum
    is_active BOOLEAN DEFAULT TRUE,      -- Untuk toggle on/off
    description TEXT,                   -- Deskripsi alergi (e.g. "Kacang Almond")

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES menu_categories(id), -- Beverage, Main Dish, Dessert
    sku VARCHAR(50) UNIQUE,

    name VARCHAR(255) NOT NULL,
    description TEXT,

    price DECIMAL(10,2) NOT NULL,
    cost DECIMAL(10,2) NOT NULL,

    is_active BOOLEAN DEFAULT TRUE,
    preparation_time INT, -- Dalam menit

    tags JSONB, -- Untuk filtering (spicy, vegetarian, etc)

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE menu_ingredients (
    id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_items(id) ON DELETE CASCADE,   
    ingredient_id INT REFERENCES ingredients(id) ON DELETE CASCADE,
    quantity DECIMAL(6,2) NOT NULL,    -- Jumlah bahan (e.g. 100 gram)
    unit VARCHAR(20) NOT NULL,          -- Satuan (gram, ml, pcs)
    is_removable BOOLEAN DEFAULT TRUE,  -- Bisa dihapus saat pemesanan
    is_default BOOLEAN DEFAULT TRUE     -- Termasuk bahan pasti ada
);

CREATE TABLE sales_analysis_daily (
    id SERIAL PRIMARY KEY,
    outlet_id INT REFERENCES outlets(id),
    analysis_date DATE NOT NULL,
    total_sales DECIMAL(12,2) NOT NULL,
    total_covers INT NOT NULL, -- Jumlah tamu
    avg_spend_per_cover DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(12,2) NOT NULL,
    void_amount DECIMAL(12,2) NOT NULL, -- Total transaksi batal
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(outlet_id, analysis_date)
);
