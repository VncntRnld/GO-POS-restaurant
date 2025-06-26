-- Customers
CREATE TYPE customer_type AS ENUM ('hotel_guest', 'non-guest');
CREATE TABLE customers (
    cust_id SERIAL PRIMARY KEY, -- Perlu ditambahin HYBRID sama UUID kah?
    hotel_guest_id VARCHAR(50) NULL, -- ID akan dikoneksikan dengan sistem hotel
    tipe customer_type NOT NULL,

    nama VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    visit_count INT DEFAULT 0,
    last_visit TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


-- Restaurant
CREATE TYPE visit_type AS ENUM ('breakfast', 'lunch', 'dinner', 'event');
CREATE TABLE customer_visits (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(cust_id),
    visit_type visit_type NOT NULL,
    visit_date TIMESTAMP NOT NULL,
    
    room_number VARCHAR(20), -- Jika tamu hotel
    reservation_id INT NULL, -- FK ke tabel reservations
    outlet_id INT REFERENCES outlets(id), -- Untuk tracking outlet restoran

    total_spent DECIMAL(12,2),
    pax INT, -- Jumlah orang
    -- pax_children INT, -- jumlah anak-anak

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TYPE status_reservation AS ENUM ('confirmed', 'waiting', 'canceled', 'no_show', 'seated');
CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(cust_id),

    reservation_time TIMESTAMP NOT NULL,
    pax INT NOT NULL,
    table_id INT REFERENCES tables(id),

    status status_reservation NOT NULL,
    special_request TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


-- Order & Billing
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL, -- UUID
    table_id INT REFERENCES tables(id),
    customer_id INT REFERENCES customers(cust_id),
    hotel_room VARCHAR(20) NULL, -- Untuk charge ke kamar
    waiter_id INT REFERENCES staff(id),
    outlet_id INT REFERENCES outlets(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('open', 'settled', 'void', 'transferred')),
    order_type VARCHAR(20) NOT NULL CHECK (order_type IN ('dine_in', 'takeaway', 'delivery', 'room_service')),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TYPE status_bill AS ENUM ('open', 'paid', 'partial', 'void');
CREATE TABLE bills (
    id SERIAL PRIMARY KEY,
    bill_number VARCHAR(50) UNIQUE NOT NULL, --UUID
    order_id INT REFERENCES orders(id),
    
    original_bill_id INT NULL REFERENCES bills(id), -- Untuk split bill
    status status_bill NOT NULL,

    subtotal DECIMAL(12,2) NOT NULL,
    tax_amount DECIMAL(12,2) NOT NULL,
    service_charge DECIMAL(12,2) NOT NULL,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) NOT NULL,
    paid_amount DECIMAL(12,2) DEFAULT 0,
    balance_due DECIMAL(12,2) GENERATED ALWAYS AS (total_amount - paid_amount) STORED,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE bill_payments (
    id SERIAL PRIMARY KEY,
    bill_id INT REFERENCES bills(id),
    payment_method VARCHAR(50) NOT NULL CHECK (payment_method IN ('cash', 'credit_card', 'debit_card', 'room_charge', 'voucher', 'split')),
    amount DECIMAL(12,2) NOT NULL,
    reference_number VARCHAR(100), -- untuk pembayaran room cth: ROOM-401
    room_charge_approved_by INT REFERENCES staff(id),
    payment_time TIMESTAMP DEFAULT NOW()
);

CREATE TABLE table_transfers (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id),
    from_table_id INT REFERENCES tables(id),
    to_table_id INT REFERENCES tables(id),
    transferred_by INT REFERENCES staff(id),
    transferred_at TIMESTAMP DEFAULT NOW(),
    reason TEXT
);