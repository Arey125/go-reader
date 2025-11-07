CREATE TABLE user_words (
    user_id INTEGER NOT NULL,
    word_id INTEGER NOT NULL,
    status TEXT NOT NULL,
    UNIQUE (user_id, word_id),

    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(word_id) REFERENCES words(id)
);
