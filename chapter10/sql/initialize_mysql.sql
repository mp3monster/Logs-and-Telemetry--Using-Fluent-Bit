create database demo;
use demo;

CREATE TABLE pluginsrc
(
    a_key varchar(100) NOT NULL,
    a_string text NOT NULL,
    a_number integer,
    a_dtg timestamp,
    a_decimal decimal(20,10),
    PRIMARY KEY (a_key)
);
CREATE TABLE plugindest
(
    a_key varchar(100) NOT NULL,
    a_string text NOT NULL,
    a_number integer,
    a_dtg timestamp,
    a_decimal decimal(20,10),
    primkey int NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (primkey)
);