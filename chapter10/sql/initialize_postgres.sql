CREATE TABLE pluginsrc
(
    a_key smallserial NOT NULL,
    a_string text NOT NULL,
    a_number integer,
    a_dtg timestamp without time zone,
    a_decimal numeric,
    PRIMARY KEY (a_key)
);

ALTER TABLE IF EXISTS pluginsrc
    OWNER to "postgresUser";
