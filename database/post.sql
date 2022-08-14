CREATE TABLE IF NOT EXISTS posts
(
    id uuid not null primary key,
    name text ,
    description text,
    user_id uuid

);

CREATE TABLE IF NOT EXISTS medias(
    id uuid primary key not null,
    type text ,
    link text,
    post_id uuid,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
);