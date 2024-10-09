INSERT INTO comments (parent_id, article_id, author, text, censored, pub_time) VALUES (currval('comments_id_seq'), 1, 'Bob', 'Hey there!', false, 1728308797);
INSERT INTO comments (parent_id, article_id, author, text, censored, pub_time) VALUES (currval('comments_id_seq'), 1, 'Alice', 'Booring, next', false, 1728318797);
INSERT INTO comments (parent_id, article_id, author, text, censored, pub_time) VALUES (2, 1, 'Pajeet', 'Hey pretty lady', false, 1728408797);
