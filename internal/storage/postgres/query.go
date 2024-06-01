package postgres

const (
	getPassword = "SELECT password FROM auth WHERE username = $1"
	createUser  = "INSERT INTO auth(username, password) VALUES($1, $2)"
	isAdmin     = "SELECT isadmin FROM auth WHERE username = $1"
	deleteUser  = "DELETE FROM auth WHERE username = $1"

	createPost    = "INSERT INTO posts(author, name, description, content, commnets_allowed, createdAt) RETURNING id"
	createComment = "INSERT INTO comments(author, post_id, replies_to, text, createdAt) RETURNING id"

	getPosts    = "SELECT * FROM posts ORDER BY createdAt DESC LIMIT $1 OFFSET $2 "
	getComments = "SELECT * FROM comments WHERE post_id = $1 ORDER BY createdAt"
)
