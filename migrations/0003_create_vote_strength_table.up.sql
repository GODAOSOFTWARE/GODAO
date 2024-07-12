-- migrations/0003_create_vote_strength_table.up.sql
CREATE TABLE vote_strength (
                               id INTEGER PRIMARY KEY AUTOINCREMENT,
                               wallet_address TEXT UNIQUE,
                               vote_power INTEGER
);
