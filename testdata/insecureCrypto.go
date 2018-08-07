package main

import(
	"crypto/des"
	"crypto/md5"
	"crypto/sha1"
)

func badCrypto() int {
	var key []byte;
	h := md5.New();
	h = sha1.New();
	c, err := des.NewTripleDESCipher(key);
	if err != nil {
		return 1;
	}
	if h != nil && c != nil {
		return 0;
	}
	return 0;
}
