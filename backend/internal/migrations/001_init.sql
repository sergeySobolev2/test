-- Миграция для системы расчета звукоизоляции перегородок

-- 1) USERS
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_moderator BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 2) PARTITIONS (типы перегородок для звукоизоляции)
CREATE TABLE IF NOT EXISTS partitions (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,                    -- Название: "Гипсокартон 12мм"
    category VARCHAR(100),                          -- Категория: "Легкие", "Тяжелые"
    description TEXT,                               -- Описание конструкции
    noise_reduction VARCHAR(50),                    -- Снижение шума: "30-35 дБ"
    thickness VARCHAR(50),                          -- Толщина: "10-15 см"
    material VARCHAR(100),                          -- Материал: "ГКЛ + минвата"
    price_per_sqm VARCHAR(50),                      -- Цена: "500-800 руб/м²"
    image_url VARCHAR(200),                         -- URL изображения
    is_active BOOLEAN NOT NULL DEFAULT TRUE,        -- Активность
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_partitions_active ON partitions(is_active);
CREATE INDEX IF NOT EXISTS idx_partitions_category ON partitions(category);

-- 3) CALCULATIONS (расчеты звукоизоляции)
CREATE TABLE IF NOT EXISTS calculations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL,                    -- черновик, сформирован, завершен, отклонен, удалён
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    formed_at TIMESTAMP,
    completed_at TIMESTAMP,
    moderator_id INTEGER,
    room_area DECIMAL(6,2),                         -- Площадь помещения (м²)
    noise_reduction_db DECIMAL(5,2),                -- Требуемое снижение шума (дБ)
    required_thickness DECIMAL(5,2),                -- Рекомендуемая толщина (см)
    expert_comment TEXT,                            -- Комментарий эксперта
    CONSTRAINT calculations_status_check
        CHECK (status IN ('черновик', 'удален', 'удалён', 'сформирован', 'завершен', 'отклонен'))
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_calculations_user_id') THEN
        ALTER TABLE calculations ADD CONSTRAINT fk_calculations_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT;
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_calculations_moderator_id') THEN
        ALTER TABLE calculations ADD CONSTRAINT fk_calculations_moderator_id
        FOREIGN KEY (moderator_id) REFERENCES users(id) ON DELETE RESTRICT;
    END IF;
END$$;

CREATE UNIQUE INDEX IF NOT EXISTS one_draft_per_user ON calculations (user_id) WHERE status = 'черновик';
CREATE INDEX IF NOT EXISTS idx_calculations_user_id ON calculations(user_id);
CREATE INDEX IF NOT EXISTS idx_calculations_status ON calculations(status);

-- 4) CALCULATION_ITEMS (связь M:N между расчетами и перегородками)
CREATE TABLE IF NOT EXISTS calculation_items (
    calculation_id INTEGER NOT NULL,
    partition_id INTEGER NOT NULL,
    quantity INTEGER,                               -- Количество элементов
    is_main BOOLEAN DEFAULT FALSE,                  -- Основной тип перегородки
    comment TEXT,                                   -- Примечание
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (calculation_id, partition_id),
    FOREIGN KEY (calculation_id) REFERENCES calculations(id) ON DELETE CASCADE,
    FOREIGN KEY (partition_id) REFERENCES partitions(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_calculation_items_calculation ON calculation_items(calculation_id);
CREATE INDEX IF NOT EXISTS idx_calculation_items_partition ON calculation_items(partition_id);
