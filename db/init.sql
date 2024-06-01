CREATE TABLE public.posts(
    id BIGSERIAL PRIMARY KEY,
    author VARCHAR(64) NOT NULL,
    name VARCHAR(256) NOT NULL,
    description VARCHAR(1024) NOT NULL,
    content text NOT NULL,
    comments_allowed boolean DEFAULT TRUE,
    createdAt timestampz,
    updatedAt timestampz
)

CREATE TABLE public.comments(
    id BIGSERIAL PRIMARY KEY,
    author VARCHAR(64) NOT NULL,
    post_id BIGSERIAL CHECK(post_id > 0),
    replies_to BIGSERIAL CHECK(replies_to > 0),
    text VARCHAR(2048) NOT NULL,
    createdAt timestampz,
    updatedAt timestampz,
    FOREIGN KEY post_id REFERENCES posts(id)
)

CREATE TABLE public.auth(
    username VARCHAR(64) PRIMARY KEY,
    password VARCHAR(64) NOT NULL,
    isadmin boolean DEFAULT FALSE
)