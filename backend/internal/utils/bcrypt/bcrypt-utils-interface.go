package bcrypt

type IBcryptUtils interface {
	HashPassword(password string) (string, error)
	CheckPassword(password string, hashedPassword string) error
}
