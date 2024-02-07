package mining

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type node struct {
	Header    string //header
	PHeader   string //parent header
	Nonce     int64  //nonce
	TX        tx     //tx
	Timestamp int64  //timestamp
}

type Base64 = func([]byte) string

var (
	//alias for base64 encoding
	encB64 func([]byte) string = base64.StdEncoding.EncodeToString
	//alias for base64 decoding
	decB64 func(string) ([]byte, error) = base64.StdEncoding.DecodeString
)

func (n node) Pack() string {
	nonce := fmt.Sprint(n.Nonce)
	timestamp := fmt.Sprint(n.Timestamp)
	return fmt.Sprintf("{%v,%v,%v,%v,%v}",
		encB64([]byte(n.Header)),
		encB64([]byte(n.PHeader)),
		encB64([]byte(nonce)),
		encB64([]byte(n.TX.Payload)),
		encB64([]byte(timestamp)))
}

func Unpack(raw string) (*node, error) {
	raw = strings.TrimPrefix(raw, "{")
	raw = strings.TrimSuffix(raw, "}")

	fields := strings.Split(raw, ",")

	errs := make([]error, 0, 5)

	header, err := decB64(fields[0])
	errs = append(errs, err)
	pHeader, err := decB64(fields[1])
	errs = append(errs, err)
	nonceRaw, err := decB64(fields[2])
	errs = append(errs, err)
	txRaw, err := decB64(fields[3])
	errs = append(errs, err)
	timestampRaw, err := decB64(fields[4])
	errs = append(errs, err)

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	nonce, err := strconv.ParseFloat(string(nonceRaw), 64)
	if err != nil {
		return nil, err
	}

	timestamp, err := strconv.ParseInt(string(timestampRaw), 10, 64)
	if err != nil {
		return nil, err
	}

	tx, err := parseTX(string(txRaw))
	if err != nil {
		return nil, err
	}

	return &node{
		Header:    string(header),
		PHeader:   string(pHeader),
		Nonce:     int64(nonce),
		TX:        tx,
		Timestamp: timestamp}, nil
}
