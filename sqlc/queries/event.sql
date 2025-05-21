-- name: NewEvent :exec
INSERT INTO event (id, event_type, event_data, status)
VALUES (@id, @event_type, @event_data, 'ACTIVE');

-- name: GetActiveEvent :many
SELECT *
FROM event
WHERE status = 'ACTIVE';

-- name: GetEventByType :many
SELECT *
FROM event
WHERE status = 'ACTIVE'
AND event_type = @event_type;