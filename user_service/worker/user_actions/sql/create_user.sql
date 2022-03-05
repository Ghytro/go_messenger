CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(20) PRIMARY KEY NOT NULL,
    email VARCHAR(320) UNIQUE NOT NULL, -- max lenght of email address is defined by international standart
    password_md5_hash CHAR(32) NOT NULL,
    access_token CHAR(50) UNIQUE NOT NULL
);

INSERT INTO users VALUES ($1, $2, $3, $4);
