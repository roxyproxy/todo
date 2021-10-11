CREATE TABLE users(
    id uuid DEFAULT uuid_generate_v4 (),
    username VARCHAR(50) NOT NULL,
    firstname VARCHAR(255) NULL,
    lastname VARCHAR(255) NULL,
    password VARCHAR(255) NOT NULL,
    location VARCHAR(50) NOT NULL,
    PRIMARY KEY (id)
);


