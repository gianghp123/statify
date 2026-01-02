package bcrypt

import "github.com/stretchr/testify/mock"

type BcryptUtilsMock struct {
	mock.Mock
}

func (b *BcryptUtilsMock) HashPassword(password string) (string, error) {
	args := b.Called(password)
	return args.String(0), args.Error(1)
}

func (b *BcryptUtilsMock) CheckPassword(password string, hashedPassword string) error {
	args := b.Called(password, hashedPassword)
	return args.Error(0)
}
