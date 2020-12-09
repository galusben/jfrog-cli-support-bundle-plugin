package commands

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"io"
	"strconv"
)

func GetEncryptCommand() components.Command {
	return components.Command{
		Name: "encrypt",
		Description: "Encrypt secret using Artifactory master key. Output will be in the form" +
			" '<kid>.aesgcm256.<encrypted message>' or '<kid>.aesgcm128.<encrypted message>' depends on the key length",
		Aliases:   []string{"up"},
		Arguments: getEncryptArguments(),
		Flags:     getEncryptFlags(),
		EnvVars:   getEncryptEnvVar(),
		Action: func(c *components.Context) error {
			return encryptCmd(c)
		},
	}
}

func getEncryptArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "plaintext",
			Description: "Plain text to encrypt",
		},
		{
			Name:        "key",
			Description: "Artifactory master key",
		},
	}
}

func getEncryptFlags() []components.Flag {
	return []components.Flag{}
}

func getEncryptEnvVar() []components.EnvVar {
	return []components.EnvVar{}
}

type encryptConfig struct {
	plaintext string
	key       string
}

func encryptCmd(c *components.Context) error {
	if len(c.Arguments) != 2 {
		return errors.New("Wrong number of arguments. Expected: 2, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}
	var conf = new(encryptConfig)
	conf.plaintext = c.Arguments[0]
	conf.key = c.Arguments[1]
	err := doEncrypt(conf)
	if err != nil {
		return err
	}
	return nil
}

func doEncrypt(c *encryptConfig) error {
	if c.plaintext == "" {
		return errors.New("secret can not be empty")
	}
	if c.key == "" {
		return errors.New("master key can not be empty")
	}
	key, _ := hex.DecodeString(c.key)
	plaintext := []byte(c.plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	secret := append(nonce, ciphertext...)
	keysha256 := sha256.Sum256(key)
	fmt.Printf("%v.%v%v.%v", hex.EncodeToString(keysha256[:])[0:6], "aesgcm", 8*len(key),
		b64.URLEncoding.EncodeToString(secret))
	return nil
}
