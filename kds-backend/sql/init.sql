DROP TABLE IF EXISTS likes CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS tables CASCADE;

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
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
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

CREATE TABLE IF NOT EXISTS tables (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    data TEXT NOT NULL, -- Хранит данные, разделенные пробелами, например: "cross_0673 5 femboygayfunny 5"
    created_by INT NOT NULL,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_users_nickname ON users (nickname);
CREATE INDEX idx_users_joined_at ON users (joined_at);
CREATE INDEX idx_users_likes ON users (likes);
CREATE INDEX idx_users_roles ON users USING GIN (roles);
CREATE INDEX idx_users_achievements ON users USING GIN (achievements);
CREATE INDEX idx_users_links ON users USING GIN (links);
CREATE INDEX idx_users_medals ON users USING GIN (medals);

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