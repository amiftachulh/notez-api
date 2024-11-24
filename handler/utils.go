package handler

import "github.com/alexedwards/argon2id"

func hashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, &argon2id.Params{
		Memory:      19456,
		Iterations:  2,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})
	if err != nil {
		return "", err
	}
	return hash, nil
}
