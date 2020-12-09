package commands

import (
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"strconv"
	"strings"
)

func GetDecryptCommand() components.Command {
	return components.Command{
		Name: "decrypt",
		Description: "Decrypt secret using Artifactory master key. Currently supports only encrypted messages of the form:" +
			" '<kid>.aesgcm256.<encrypted message>' or '<kid>.aesgcm128.<encrypted message>'",
		Aliases:   []string{"up"},
		Arguments: getDecryptArguments(),
		Flags:     getDecryptFlags(),
		EnvVars:   getDecryptEnvVar(),
		Action: func(c *components.Context) error {
			return decryptCmd(c)
		},
	}
}

func getDecryptArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "secret",
			Description: "The secret to decrypt",
		},
		{
			Name:        "key",
			Description: "Artifactory master key",
		},
	}
}

func getDecryptFlags() []components.Flag {
	return []components.Flag{}
}

func getDecryptEnvVar() []components.EnvVar {
	return []components.EnvVar{}
}

type decryptConfig struct {
	secret string
	key    string
}

func decryptCmd(c *components.Context) error {
	if len(c.Arguments) != 2 {
		return errors.New("Wrong number of arguments. Expected: 2, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}
	var conf = new(decryptConfig)
	conf.secret = c.Arguments[0]
	conf.key = c.Arguments[1]
	err := doDecrypt(conf)
	if err != nil {
		return err
	}
	return nil
}

func doDecrypt(c *decryptConfig) error {
	if c.secret == "" {
		return errors.New("secret can not be empty")
	}
	if c.key == "" {
		return errors.New("master key can not be empty")
	}
	key, _ := hex.DecodeString(c.key)
	split := strings.Split(c.secret, ".")
	if len(split) != 3 {
		return errors.New("encryption type is not supported")
	}
	nonce, _ := b64.URLEncoding.DecodeString(split[2][0:16])
	ciphertext, _ := b64.URLEncoding.DecodeString(split[2][16:])

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("could not decrypt")
	}
	log.Output(string(plaintext))
	return nil
}
