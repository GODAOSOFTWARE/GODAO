-- Функция для создания таблицы user_votes
CREATE TABLE IF NOT EXISTS user_votes (
                                          id INTEGER PRIMARY KEY AUTOINCREMENT,
                                          vote_id INTEGER,
                                          voter TEXT,
                                          choice TEXT,
                                          vote_power INTEGER
);
