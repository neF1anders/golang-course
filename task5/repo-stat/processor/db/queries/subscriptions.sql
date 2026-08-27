-- name: List :many
SELECT name, description, stars, forks, date FROM info
ORDER BY name;

-- name: Push :exec
INSERT INTO info (name, description, stars, forks, date)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (name) 
DO UPDATE SET
    description = EXCLUDED.description,
    stars       = EXCLUDED.stars,
    forks       = EXCLUDED.forks,
    date        = EXCLUDED.date;

-- name: Delete :exec
DELETE FROM info
WHERE name = $1;
