/*
* @Author: detailyang
* @Date:   2015-08-16 14:07:18
* @Last Modified by:   detailyang
* @Last Modified time: 2015-08-20 01:09:52
 */

package hikarian

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rc4"
	"io"
)

type Cipher struct {
	enc cipher.Stream
	dec cipher.Stream
}

func truncateSecretToSize(secret string, size int) []byte {
	h := md5.New()
	if size > md5.Size {
		buf := make([]byte, size)
		for i := 0; i < size/md5.Size; i++ {
			copy(buf[i*md5.Size:(i+1)*md5.Size-1], h.Sum(nil))
		}
		return buf[:size]
	}

	io.WriteString(h, secret)

	return h.Sum(nil)[:size]
}

func NewChiper(algo, secret string) (*Cipher, error) {
	if algo == "rc4" {
		c, err := rc4.NewCipher(truncateSecretToSize(secret, 32))
		if err != nil {
			return nil, err
		}
		return &Cipher{
			enc: c,
			dec: c,
		}, nil
	} else if algo == "aes" {
		key := truncateSecretToSize(secret, 32)
		c, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		return &Cipher{
			enc: cipher.NewCFBEncrypter(c, key[:c.BlockSize()]),
			dec: cipher.NewCFBDecrypter(c, key[:c.BlockSize()]),
		}, nil
	}

	cipher, err := rc4.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &Cipher{
		enc: cipher,
		dec: cipher,
	}, nil
}

func (self *Cipher) Encrypt(cleartext, ciphertext []byte) {
	self.enc.XORKeyStream(ciphertext, cleartext)
}

func (self *Cipher) Decrypt(ciphertext, cleartext []byte) {
	self.dec.XORKeyStream(cleartext, ciphertext)
}
