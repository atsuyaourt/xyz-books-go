-- Create authors table
CREATE TABLE authors (
    author_id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    middle_name TEXT DEFAULT '' NOT NULL,
    UNIQUE(first_name, middle_name, last_name)
);

-- Create books table
CREATE TABLE books (
    book_id INTEGER PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    isbn13 TEXT UNIQUE,
    isbn10 TEXT UNIQUE,
    price REAL NOT NULL,
    publication_year INTEGER NOT NULL,
    image_url TEXT,
    edition TEXT,
    publisher_id INTEGER NOT NULL,
    FOREIGN KEY (publisher_id) REFERENCES publishers(publisher_id)
);

-- Create publishers table
CREATE TABLE publishers (
    publisher_id INTEGER PRIMARY KEY,
    publisher_name TEXT NOT NULL
);

-- Create author_book junction table for many-to-many relationship
CREATE TABLE author_book (
    author_id INTEGER NOT NULL,
    book_id INTEGER NOT NULL,
    PRIMARY KEY (author_id, book_id),
    FOREIGN KEY (author_id) REFERENCES authors(author_id),
    FOREIGN KEY (book_id) REFERENCES books(book_id)
);