-- migrations/0002_create_user_votes_table.up.sql
CREATE TABLE user_votes (
                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                            vote_id INTEGER,
                            voter TEXT,
                            choice TEXT,
                            vote_power INTEGER
);
