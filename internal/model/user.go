package model

// User adalah representasi langsung dari tabel users di database.
// Struct ini dipakai oleh repository — jangan expose langsung ke client.
type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // field sensitif — tidak boleh pernah dikirim ke client
}
