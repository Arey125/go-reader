-- name: AddWord :exec
insert or ignore into words (word, pos) values (?, ?);

-- name: UpdateWordDefinition :exec
update words set definitions = ? where word = ? and pos = ?;

-- name: GetWordDefinition :one
select definitions from words where word = ? and pos = ?;

-- name: GetUserWords :many
select sqlc.embed(words) from user_words left join words on words.id = user_words.word_id where user_id = ?;

-- name: AddUserWord :exec
insert or ignore into user_words(user_id, word_id, status) select ? as user_id, id as word_id, ? as status from words where word = ? and pos = ?;
