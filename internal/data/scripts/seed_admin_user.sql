-- Add a password hash by running the hasher command line tool to generate a new admin password
-- Add an id by running the id_generator command line tool
INSERT INTO users (id, username, email, first_name, last_name, password_hash)
VALUES (null, 'admin', 'admin@localhost', 'Local', 'Administrator', null);

-- replace null user id with id from above
INSERT INTO user_roles (user_id, role_id)
SELECT null AS user_id, r.id AS role_id
FROM roles r
WHERE r.Name = 'Administrator';