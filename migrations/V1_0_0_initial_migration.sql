CREATE TABLE ingredient
(
    id               SERIAL PRIMARY KEY,
    english_category varchar(30),
    english_name     varchar(150) NOT NULL UNIQUE,
    native_name      varchar(180),
    translated_name  varchar(150),
    shopping_link    varchar
);
CREATE INDEX idx_english_category ON ingredient (english_category);

CREATE TABLE source
(
    id          SERIAL PRIMARY KEY,
    description varchar(150),
    location    varchar,
    category    smallint
);

CREATE TABLE recipe
(
    id                    SERIAL PRIMARY KEY,
    attribution_source_id integer,
    cuisine_source_id     integer,
    english_name          varchar(180) NOT NULL,
    native_name           varchar(250),
    note                  varchar[],
    instruction           jsonb[] NOT NULL,
    FOREIGN KEY (attribution_source_id) REFERENCES source (id),
    FOREIGN KEY (cuisine_source_id) REFERENCES source (id)
);

CREATE TABLE recipe_source
(
    id              SERIAL PRIMARY KEY,
    recipe_id       integer NOT NULL,
    emoji_source_id integer NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipe (id),
    FOREIGN KEY (emoji_source_id) REFERENCES source (id)
);

CREATE TABLE ingredient_specification
(
    id                   SERIAL PRIMARY KEY,
    component            varchar(50) NOT NULL,
    recipe_id            integer NOT NULL,
    ingredient_id        integer NOT NULL,
    note                 varchar(25),
    amount_quantity      float(3)[2],
    amount_mass          smallint,
    preparation_quantity float(3),
    preparation_type     varchar(20),
    preparation_length   smallint,
    FOREIGN KEY (recipe_id) references recipe (id),
    FOREIGN KEY (ingredient_id) REFERENCES ingredient (id)
);