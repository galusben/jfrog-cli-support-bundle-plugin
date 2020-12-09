package crypto

type encryptor interface {
	encrypt(key string, plaintext string) (string, error)
}
