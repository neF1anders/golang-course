CREATE TABLE IF NOT EXISTS info (
    name VARCHAR(255) NOT NULL,
    description  VARCHAR(255),
    stars  INT,
    forks  INT,
    date TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (name)
);