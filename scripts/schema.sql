DROP TABLE IF EXISTS comments;

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    parent_id INTEGER DEFAULT 0,
    FOREIGN KEY (parent_id) REFERENCES comments(id),
    article_id INTEGER DEFAULT 0,
    author TEXT  NOT NULL,
    text TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0
);

CREATE INDEX idx_comments_pub_time ON comments (pub_time);
