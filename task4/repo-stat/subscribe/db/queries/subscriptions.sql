-- name: List :many
SELECT owner, repo FROM subscriptions
ORDER BY owner, repo;

-- name: Push :exec
INSERT INTO subscriptions (owner, repo)
VALUES ($1, $2)
ON CONFLICT (owner, repo) DO NOTHING;

-- name: Delete :exec
DELETE FROM subscriptions
WHERE owner = $1 AND repo = $2;
