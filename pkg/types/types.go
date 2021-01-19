package types

// Data represents data from data_table
type Data struct {
	ID     int    `json:"id"`
	String string `json:"string"`
}

// User is a system user
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Fullname     string `json:"fullname"`
	PasswordSalt string `json:"passwordSalt"`
	PasswordHash string `json:"passwordHash"`
	Email        string `json:"email"`
	IsDisabled   bool   `json:"isDisabled"`
}
