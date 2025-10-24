package cron

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/md5"
    "encoding/base64"
    "encoding/hex"
    "errors"
    "fmt"
    "hash"
    "io"
)

const (
    pkcs5SaltLen = 8
    aes256KeyLen = 32
)

// DecryptCryptoJsAesMsg 解密 CryptoJS.AES.encrypt(msg, password) 的密文
func DecryptCryptoJsAesMsg(password string, ciphertext string) ([]byte, error) {
    const blocklen = aes.BlockSize

    rawEncrypted, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return nil, fmt.Errorf("base64 decode error: %v", err)
    }

    if len(rawEncrypted) < 16+blocklen || string(rawEncrypted[:8]) != "Salted__" {
        return nil, errors.New("invalid ciphertext format")
    }

    salt := rawEncrypted[8:16]
    encrypted := rawEncrypted[16:]

    key, iv := BytesToKey(salt, []byte(password), md5.New(), aes256KeyLen, blocklen)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("aes cipher init error: %v", err)
    }

    mode := cipher.NewCBCDecrypter(block, iv)
    decrypted := make([]byte, len(encrypted))
    mode.CryptBlocks(decrypted, encrypted)

    return pkcs7Strip(decrypted, blocklen)
}

func BytesToKey(salt, data []byte, h hash.Hash, keyLen, blockLen int) (key, iv []byte) {
    if len(salt) > 0 && len(salt) != pkcs5SaltLen {
        panic(fmt.Sprintf("salt length %d != %d", len(salt), pkcs5SaltLen))
    }

    var (
        concat   []byte
        lastHash []byte
        totalLen = keyLen + blockLen
    )

    for len(concat) < totalLen {
        h.Reset()
        h.Write(append(lastHash, append(data, salt...)...))
        lastHash = h.Sum(nil)
        concat = append(concat, lastHash...)
    }

    return concat[:keyLen], concat[keyLen:totalLen]
}

func Md5String(inputs ...string) string {
    keyHash := md5.New()
    for _, str := range inputs {
        io.WriteString(keyHash, str)
    }
    return hex.EncodeToString(keyHash.Sum(nil))
}

func pkcs7Strip(data []byte, blockSize int) ([]byte, error) {
    length := len(data)
    if length == 0 {
        return nil, errors.New("pkcs7: empty data")
    }
    if length%blockSize != 0 {
        return nil, errors.New("pkcs7: not block aligned")
    }
    padLen := int(data[length-1])
    if padLen == 0 || padLen > blockSize {
        return nil, errors.New("pkcs7: invalid pad length")
    }
    if !bytes.HasSuffix(data, bytes.Repeat([]byte{byte(padLen)}, padLen)) {
        return nil, errors.New("pkcs7: invalid padding content")
    }
    return data[:length-padLen], nil
}
