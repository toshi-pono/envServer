DROP DATABASE IF EXISTS sora;
CREATE DATABASE sora;
\c sora;
DROP TABLE IF EXISTS weather;

CREATE TABLE weather (
id SERIAL NOT NULL,
temperature FLOAT,
humidity FLOAT,
pressure FLOAT,
battery FLOAT,
created_at timestamp NOT NULL,
PRIMARY KEY(id)
);