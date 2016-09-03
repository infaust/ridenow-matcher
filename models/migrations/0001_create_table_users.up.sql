CREATE TABLE user_profile (
	id serial PRIMARY KEY,
	username varchar(100) NOT NULL,
	email varchar(256) UNIQUE,
	wave_height_range numrange, -- meters
	allowed_time_range int4range, -- hour within a day
	-- password hash, oauth token, etc.
	name varchar(100),
	surname varchar(200),
	created timestamp NOT NULL default current_timestamp
);

CREATE TABLE user_location (
    id serial PRIMARY KEY,
    user_profile_id serial references user_profile(id),
    location_id serial
)