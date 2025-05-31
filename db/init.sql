CREATE TABLE posts (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    user_session_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    user_avatar_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE comments (
    id TEXT PRIMARY KEY,
    post_id TEXT REFERENCES posts(id),
    parent_id TEXT, -- references either posts.id or comments.id
    content TEXT NOT NULL,
    user_session_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    user_avatar_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);