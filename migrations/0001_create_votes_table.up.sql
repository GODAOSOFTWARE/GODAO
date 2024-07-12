-- migrations/0001_create_votes_table.up.sql
CREATE TABLE votes (
                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                       title TEXT,
                       subtitle TEXT,
                       description TEXT,
                       voter TEXT,
                       choice TEXT,
                       vote_power INT,
                       wallet_address TEXT
);