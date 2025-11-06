-- name: AddWord :exec
insert or ignore into words (word, pos) values (?, ?);

-- name: UpdateWordDefinition :exec
update words set definitions = ? where word = ? and pos = ?;

-- name: GetWordDefinition :one
select definitions from words where word = ? and pos = ?;
