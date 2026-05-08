CREATE TABLE IF NOT EXISTS subscriptions (
    owner VARCHAR(255) NOT NULL,
    repo  VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (owner, repo)
);