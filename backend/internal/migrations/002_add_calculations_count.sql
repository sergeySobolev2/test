-- Миграция: добавление счетчика результатов расчета

-- Добавляем поле calculations_count в таблицу calculations
ALTER TABLE calculations 
ADD COLUMN IF NOT EXISTS calculations_count INTEGER NOT NULL DEFAULT 0;

-- Комментарий к полю
COMMENT ON COLUMN calculations.calculations_count IS 'Количество выполненных расчетов (результатов от async сервиса)';
