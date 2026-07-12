CREATE TABLE authors (
    id bigserial PRIMARY KEY,
    name text
);

CREATE TABLE users (
    id bigserial PRIMARY KEY,
    name text,
    email text,
    password text
);

CREATE UNIQUE INDEX idx_users_email ON users (email);

CREATE TABLE books (
    id bigserial PRIMARY KEY,
    title text,
    author_id bigint
);

CREATE TABLE reviews (
    id bigserial PRIMARY KEY,
    book_id bigint NOT NULL,
    user_id bigint NOT NULL,
    rating bigint NOT NULL,
    comment text,
    CONSTRAINT chk_reviews_rating CHECK (rating >= 1 AND rating <= 5)
);

CREATE TABLE reading_lists (
    id bigserial PRIMARY KEY,
    user_id bigint,
    name text
);

CREATE TABLE reading_list_books (
    reading_list_id bigint NOT NULL,
    book_id bigint NOT NULL,
    PRIMARY KEY (reading_list_id, book_id),
    CONSTRAINT fk_reading_list_books_reading_list FOREIGN KEY (reading_list_id) REFERENCES reading_lists (id),
    CONSTRAINT fk_reading_list_books_book FOREIGN KEY (book_id) REFERENCES books (id)
);
