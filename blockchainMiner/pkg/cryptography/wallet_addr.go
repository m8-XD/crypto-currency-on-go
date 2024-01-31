package cryptography

func WaletAddr(publicKey string) (wAddr string, err error) {
	pkSHA256, err := HashSHA256(publicKey)
	if err != nil {
		return
	}

	pkSHA256MD5, err := HashMD5(string(pkSHA256))
	if err != nil {
		return
	}

	wAddr = EncodeBase32(string(pkSHA256MD5))
	return
}
