-- Создание таблицы статусов активов
CREATE TABLE IF NOT EXISTS asset_statuses (
                                              id SERIAL PRIMARY KEY,
                                              name VARCHAR(50) NOT NULL UNIQUE
    );

-- Создание таблицы местоположений
CREATE TABLE IF NOT EXISTS locations (
                                         id SERIAL PRIMARY KEY,
                                         address VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL
    );

-- Создание таблицы отделов
CREATE TABLE IF NOT EXISTS departments (
                                           id SERIAL PRIMARY KEY,
                                           name VARCHAR(100) NOT NULL,
    location VARCHAR(255) NOT NULL,
    head_id INTEGER
    );

-- Создание таблицы сотрудников
CREATE TABLE IF NOT EXISTS employees (
                                         id SERIAL PRIMARY KEY,
                                         full_name VARCHAR(100) NOT NULL,
    position VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'employee',
    department_id INTEGER REFERENCES departments(id) ON DELETE SET NULL
    );

-- Обновление внешнего ключа для глав отделов
ALTER TABLE departments
    ADD CONSTRAINT fk_departments_head
        FOREIGN KEY (head_id) REFERENCES employees(id) ON DELETE SET NULL;

-- Создание таблицы активов
CREATE TABLE IF NOT EXISTS assets (
                                      id SERIAL PRIMARY KEY,
                                      name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    acquisition_date DATE NOT NULL,
    cost DECIMAL(10, 2) NOT NULL,
    status_id INTEGER NOT NULL REFERENCES asset_statuses(id),
    current_location_id INTEGER NOT NULL REFERENCES locations(id),
    department_id INTEGER REFERENCES departments(id) ON DELETE SET NULL
    );

-- Создание таблицы перемещений активов
CREATE TABLE IF NOT EXISTS asset_transfers (
                                               id SERIAL PRIMARY KEY,
                                               asset_id INTEGER NOT NULL REFERENCES assets(id),
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    from_location_id INTEGER NOT NULL REFERENCES locations(id),
    to_location_id INTEGER NOT NULL REFERENCES locations(id),
    transfer_date TIMESTAMP NOT NULL DEFAULT NOW(),
    notes TEXT
    );

-- Создание таблицы сессий
CREATE TABLE IF NOT EXISTS sessions (
                                        id SERIAL PRIMARY KEY,
                                        employee_id INTEGER NOT NULL REFERENCES employees(id),
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Начальные данные
INSERT INTO asset_statuses (name) VALUES
                                      ('В эксплуатации'),
                                      ('На складе'),
                                      ('В ремонте'),
                                      ('Списан');

INSERT INTO locations (address, type) VALUES
                                          ('Главный офис, ул. Центральная, 1', 'Офис'),
                                          ('Склад №1, ул. Ленина, 10', 'Склад'),
                                          ('Филиал, ул. Парковая, 5', 'Офис');

INSERT INTO departments (name, location) VALUES
                                             ('IT', 'Главный офис'),
                                             ('Бухгалтерия', 'Главный офис'),
                                             ('Логистика', 'Склад №1');

-- Создание индексов для улучшения производительности
CREATE INDEX idx_assets_status ON assets(status_id);
CREATE INDEX idx_assets_location ON assets(current_location_id);
CREATE INDEX idx_assets_department ON assets(department_id);
CREATE INDEX idx_transfers_asset ON asset_transfers(asset_id);
CREATE INDEX idx_transfers_employee ON asset_transfers(employee_id);
CREATE INDEX idx_transfers_date ON asset_transfers(transfer_date);
CREATE INDEX idx_employees_department ON employees(department_id);