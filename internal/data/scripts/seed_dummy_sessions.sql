-- this script is used to verify background service worker can cleanup old sessions.
INSERT INTO sessions (id, token, user_id, expires_at)
VALUES
    ('a6cd7fc1-d7f0-11f0-8d6c-3c7c3f215201', 'empty', '5002ffb1-d7bd-11f0-b31f-3c7c3f215201', CURRENT_TIMESTAMP),
    ('a49897f8-d7f0-11f0-a092-3c7c3f215201', 'empty', '5002ffb1-d7bd-11f0-b31f-3c7c3f215201', CURRENT_TIMESTAMP),
    ('a736c374-d7f0-11f0-a692-3c7c3f215201', 'empty', '5002ffb1-d7bd-11f0-b31f-3c7c3f215201', CURRENT_TIMESTAMP),
    ('a77bf8d8-d7f0-11f0-a6f3-3c7c3f215201', 'empty', '5002ffb1-d7bd-11f0-b31f-3c7c3f215201', CURRENT_TIMESTAMP),
    ('edeeb4d0-d7f0-11f0-9952-3c7c3f215201', 'empty', '5002ffb1-d7bd-11f0-b31f-3c7c3f215201', CURRENT_TIMESTAMP),
    ('eed4f2f3-d7f0-11f0-bf22-3c7c3f215201', 'empty', '5002ffb1-d7bd-11f0-b31f-3c7c3f215201', CURRENT_TIMESTAMP);