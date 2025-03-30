CREATE TABLE author (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name CHARACTER VARYING(255),
    email CHARACTER VARYING(255) UNIQUE,
    date_joined timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX name_idx ON author (name);

CREATE TABLE phone (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name CHARACTER VARYING(255) NOT NULL UNIQUE,
    cols INTEGER NOT NULL,
    cols2 INTEGER NOT NULL DEFAULT -1
);

CREATE TABLE effect (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name CHARACTER VARYING(255) NOT NULL
);

CREATE TABLE ringtone (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name CHARACTER VARYING(255) NOT NULL,
    category int NOT NULL,
    downloads INTEGER DEFAULT 0,
    effect_id INTEGER REFERENCES effect (id),
    author_id INTEGER REFERENCES author (id),
    not_working INTEGER DEFAULT 0,
    time_added TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    glyphs TEXT,
    auto_generated BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE phone_and_ringtone (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    phone_id INTEGER REFERENCES phone (id),
    ringtone_id INTEGER REFERENCES ringtone (id) ON DELETE CASCADE
);

CREATE EXTENSION pg_trgm;

INSERT INTO
    phone (name, cols, cols2)
VALUES 
    ('(1)', 5, 15),
    ('(2)', 33, 5),
    ('(2a)', 26, -1),
    ('(3a)', 36, -1);

INSERT INTO
    effect (name)
VALUES 
    ('Dan'),
    ('Brrr'),
    ('606'),
    ('Weevil'),
    ('Modem'),
    ('Swedish House Mafia'),
    ('Sampha'),
    ('FM'),
    ('Fantasy'),
    ('Custom made');

/* ---------------------------------------------------------------------------------------------  */