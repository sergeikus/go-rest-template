CREATE TABLE users (
    id int GENERATED ALWAYS AS IDENTITY,
    username VARCHAR(50) NOT NULL UNIQUE,
    fullname VARCHAR(255),
    password_salt VARCHAR(50) NOT NULL,
    password_hash VARCHAR(128) NOT NULL,
    email VARCHAR(100) NOT NULL,
    is_disabled BOOLEAN,
    PRIMARY KEY(id)
);

CREATE TABLE user_sessions
(
    session_key TEXT,
    user_id INT NOT NULL,
    session_created TIMESTAMPTZ NOT NULL,
    last_seen TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (session_key),
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
);

CREATE TABLE data_table
(
    id SERIAL PRIMARY KEY,
    string varchar(255)
);

INSERT INTO data_table (string) VALUES ('data1');
INSERT INTO data_table (string) VALUES ('data2');
INSERT INTO data_table (string) VALUES ('data3');
INSERT INTO data_table (string) VALUES ('data4');
INSERT INTO data_table (string) VALUES ('data5');

/* username=test password=password */
INSERT INTO users (username, fullname, password_salt, password_hash, email, is_disabled) 
VALUES ('test', 'Test User', 'TESTSALT', '7fc909bc1888eb3b5c717dfcca83a2b1b031ecb40ed6ad4278399d78d29ea0212805336b86df2c7254e9d206eee53b5300edeaeee6f35bb96b8f0890c693d24f', 'test@email.com', false);
