-- name: AllTextsWithoutContent :many
select id, title, user_id, created_at from texts; 

-- name: GetText :one
select * from texts where id = ?; 

-- name: AddText :exec
insert into texts (title, content, user_id, created_at) values (?, ?, ?, ?);
