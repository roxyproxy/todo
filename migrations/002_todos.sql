CREATE TABLE todos(
    id uuid DEFAULT uuid_generate_v4 (),
    name VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    userid uuid NOT NULL,
    FOREIGN KEY(userid)
        REFERENCES users (id),
    PRIMARY KEY (id)
);