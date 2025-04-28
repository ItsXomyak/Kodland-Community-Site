DROP TABLE IF EXISTS likes CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS tables CASCADE;
DROP TABLE IF EXISTS login_attempts CASCADE;
DROP TABLE IF EXISTS password_reset_tokens CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;

-- таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE,
    avatar_url TEXT NOT NULL,
    links JSONB NOT NULL DEFAULT '[]',
    achievements JSONB NOT NULL DEFAULT '[]',
    roles JSONB NOT NULL DEFAULT '[]',
    medals JSONB NOT NULL DEFAULT '{"gold": 0, "silver": 0, "bronze": 0}',
    likes INTEGER NOT NULL DEFAULT 0,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMP,
    is_locked BOOLEAN DEFAULT FALSE,
    lock_until TIMESTAMP,
    failed_login_attempts INTEGER DEFAULT 0
);

-- таблица аккаунтов
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email TEXT UNIQUE,
    discord_id TEXT UNIQUE,
    telegram_id TEXT UNIQUE,
    google_id TEXT UNIQUE,
    password_hash TEXT,
    password_salt TEXT NOT NULL,
    password_changed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    two_factor_secret TEXT,
    recovery_codes TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- таблица лайков
CREATE TABLE likes (
    id SERIAL PRIMARY KEY,
    from_account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    to_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (from_account_id, to_user_id)
);

-- таблица таблиц
CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    data TEXT NOT NULL,
    created_by INT NOT NULL,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- таблица попыток входа
CREATE TABLE login_attempts (
    id SERIAL PRIMARY KEY,
    ip_address TEXT NOT NULL,
    email TEXT NOT NULL,
    success BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- таблица токенов сброса пароля
CREATE TABLE password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- таблица сессий
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    ip_address TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_users_nickname ON users (nickname);
CREATE INDEX idx_users_joined_at ON users (joined_at);
CREATE INDEX idx_users_likes ON users (likes);
CREATE INDEX idx_users_roles ON users USING GIN (roles);
CREATE INDEX idx_users_achievements ON users USING GIN (achievements);
CREATE INDEX idx_users_links ON users USING GIN (links);
CREATE INDEX idx_users_medals ON users USING GIN (medals);
CREATE INDEX idx_login_attempts_ip ON login_attempts (ip_address);
CREATE INDEX idx_login_attempts_email ON login_attempts (email);
CREATE INDEX idx_sessions_token ON sessions (token);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens (token);
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens (user_id);

-- Триггеры
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tables_updated_at
    BEFORE UPDATE ON tables
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Функции безопасности
CREATE OR REPLACE FUNCTION check_password_complexity()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.password_hash IS NOT NULL AND
       (LENGTH(NEW.password_hash) < 8 OR
        NEW.password_hash !~ '[A-Z]' OR
        NEW.password_hash !~ '[a-z]' OR
        NEW.password_hash !~ '[0-9]') THEN
        RAISE EXCEPTION 'Пароль должен содержать минимум 8 символов, включая заглавные и строчные буквы, а также цифры';
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER check_password_complexity_trigger
    BEFORE INSERT OR UPDATE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION check_password_complexity();

-- Вставка тестовых данных
INSERT INTO users (nickname, avatar_url, links, achievements, roles, medals, likes, joined_at)
VALUES 
(
    'scrameee',
    'https://robohash.org/scrameee',
    '["https://github.com/ItsXomyak", "https://itch.io/profile/scrameee", "https://github.com/ItsXomyak"]',
    '["Написал сайт для KDS"]',
    '["Мемолог", "16+ лет", "Германия", "Web"]',
    '{"gold": 2, "silver": 1, "bronze": 0}',
    1984,
    NOW() - INTERVAL '300 days'
),
(
    'pixelqueen',
    'https://robohash.org/pixelqueen',
    '["https://github.com/pixelqueen", "https://kodland.com/hub/pixelqueen"]',
    '["Организовала 5 ивентов"]',
    '["Организатор", "18+ лет", "Франция", "Event Manager"]',
    '{"gold": 1, "silver": 3, "bronze": 2}',
    15,
    NOW() - INTERVAL '450 days'
),
(
    'codewizard',
    'https://robohash.org/codewizard',
    '["https://github.com/codewizard", "https://itch.io/profile/codewizard"]',
    '["Написал бота для Discord сервера"]',
    '["Технарь", "17 лет", "Россия", "Backend"]',
    '{"gold": 0, "silver": 0, "bronze": 1}',
    10,
    NOW() - INTERVAL '200 days'
);