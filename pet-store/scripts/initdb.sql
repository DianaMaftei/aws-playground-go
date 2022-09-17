CREATE TABLE pet (
                     id BIGINT NOT NULL AUTO_INCREMENT,
                     name VARCHAR(20) NOT NULL,
                     status VARCHAR(20) NOT NULL,
                     PRIMARY KEY (id)
);

CREATE TABLE category (
                     id BIGINT NOT NULL AUTO_INCREMENT,
                     name VARCHAR(20) NOT NULL,
                     PRIMARY KEY (id)
);

INSERT INTO category (name) VALUES ('Dogs');
INSERT INTO category (name) VALUES ('Cats');
INSERT INTO pet (name, status) VALUES ('Fluffy', 'available');