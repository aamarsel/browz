CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица клиентов
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    phone TEXT UNIQUE NOT NULL,
    telegram_id INT UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица доступных временных слотов
CREATE TABLE IF NOT EXISTS available_slots (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    time TIME NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,  
    UNIQUE(date, time)
);

-- Таблица записей клиентов
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    slot_id INT NOT NULL REFERENCES available_slots(id) ON DELETE CASCADE,
    status TEXT CHECK (status IN ('pending', 'confirmed', 'cancelled')) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(client_id, slot_id)
);

-- Таблица услуг
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price INT NOT NULL CHECK (price >= 0), -- Цена в рублях
    duration INTERVAL NOT NULL -- Длительность услуги
);

CREATE INDEX IF NOT EXISTS idx_available_slots_date ON available_slots(date);
CREATE INDEX IF NOT EXISTS idx_bookings_client ON bookings(client_id);
