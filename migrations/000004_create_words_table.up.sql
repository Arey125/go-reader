CREATE TABLE words (
    id INTEGER PRIMARY KEY,
	word TEXT NOT NULL,
	pos TEXT NOT NULL,
	definition TEXT,
    unique (word, pos)
);
