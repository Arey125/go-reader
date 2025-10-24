CREATE TABLE words (
    id INTEGER PRIMARY KEY,
	word TEXT NOT NULL,
	pos TEXT NOT NULL,
	definitions TEXT,
    unique (word, pos)
);
