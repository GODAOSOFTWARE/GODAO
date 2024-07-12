-- Функция для создания таблицы vote_strength
CREATE TABLE IF NOT EXISTS vote_strength (
                                             id INTEGER PRIMARY KEY AUTOINCREMENT,
                                             wallet_address TEXT UNIQUE,
                                             vote_power INTEGER
);
