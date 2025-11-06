CREATE TABLE texts (
    id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	user_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL
);
