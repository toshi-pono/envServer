DROP DATABASE IF EXISTS sora;
CREATE DATABASE sora;
\c sora;
DROP TABLE IF EXISTS weather;

CREATE TABLE weather
(
    id SERIAL NOT NULL,
    temperature FLOAT,
    humidity FLOAT,
    pressure FLOAT,
    battery FLOAT,
    created_at timestamp NOT NULL,
    PRIMARY KEY(id)
);

DROP TABLE IF EXISTS users;

CREATE TABLE users
(
    id SERIAL NOT NULL,
    userId text NOT NULL,
    userType SMALLINT NOT NULL,
    reply SMALLINT,
    PRIMARY KEY(id)
);