CREATE TABLE permissions (
    id INTEGER PRIMARY KEY,
    slug TEXT,
    email TEXT DEFAULT ''
);

INSERT INTO permissions (slug)
VALUES
    ('app'),
    ('users.permissions');

CREATE TABLE user_permissions (
    user_id INTEGER,
    permission_id INTEGER,

    PRIMARY KEY (user_id, permission_id)

    FOREIGN KEY(user_id) REFERENCES users(id)
    FOREIGN KEY(permission_id) REFERENCES permissions(id)
);
