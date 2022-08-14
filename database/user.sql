CREATE TABLE IF NOT EXISTS users
(
    id uuid not null primary key,
    first_name text,
    last_name text,
    bio text,
    email text,
    status text,
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone ,
    deleted_at timestamp without time zone,
    phone_number text[],
    username text,
    password text not null
);

CREATE TABLE IF NOT EXISTS adress
(
    user_id uuid ,
    city text,
    country text,
    district text,
    postal_code integer,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);