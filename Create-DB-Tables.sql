CREATE TABLE IF NOT EXISTS Users (
	user_uuid VARCHAR(36) PRIMARY KEY NOT NULL,
	username TEXT UNIQUE,
	email TEXT UNIQUE,
	password TEXT,
	first_name TEXT,
	last_name TEXT,
	birth_date TIME,
	gender TEXT,
	role TEXT,
	profile_picture MEDIUMTEXT,
	created_at TIME
);

CREATE TABLE IF NOT EXISTS Posts (
	post_uuid VARCHAR(36) PRIMARY KEY NOT NULL,
	user_uuid VARCHAR(36) NOT NULL,
	title TEXT,
	content TEXT,
	username TEXT,
	profile_picture MEDIUMTEXT,
	likes INT,
	dislikes INT,
	created_at TIME
);

CREATE TABLE IF NOT EXISTS Comments (
	comment_id VARCHAR(36) PRIMARY KEY NOT NULL,
	post_uuid  VARCHAR(36) NOT NULL,
	user_uuid  VARCHAR(36) NOT NULL,
	content TEXT,
	username TEXT,
	profile_picture MEDIUMTEXT,
	created_at TIME
);