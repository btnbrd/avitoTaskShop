package models

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Coins    int    `db:"coins"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
