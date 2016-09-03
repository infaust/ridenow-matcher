-- Users: id, username, email, wave height range, allowed time range, name, surname, created
INSERT INTO user_profile VALUES (1, ' infaust', 'infaust.tg@gmail.com', '[0.1, 2.0]', '[9, 20]', 'Faust', 'Terrado');
-- Locations: id, user profile id, location id
INSERT INTO user_location VALUES (1, 1, 3536); -- Sitges
INSERT INTO user_location VALUES (2, 1, 3553); -- Son Bou (Menorca)
INSERT INTO user_location VALUES (3, 1, 3535); -- Barceloneta
