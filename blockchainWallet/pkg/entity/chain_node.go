package entity

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	Header    string //header
	PHeader   string //parent header
	Nonce     int64  //nonce
	TX        tx     //tx
	Timestamp int64  //timestamp
}

// transaction
type tx struct {
	WAddr     string
	RecWAddr  string //reciever wallet address
	Amount    float64
	Change    float64
	BHash     string //block hash
	Timestamp int64
	DS        string //digital signature
	Payload   string
}

type Base64 = func([]byte) string

var (
	//alias for base64 encoding
	encB64 func([]byte) string = base64.StdEncoding.EncodeToString
	//alias for base64 decoding
	decB64 func(string) ([]byte, error) = base64.StdEncoding.DecodeString
)

func (n Node) Pack() string {
	nonce := fmt.Sprint(n.Nonce)
	timestamp := fmt.Sprint(n.Timestamp)
	return fmt.Sprintf("{%v,%v,%v,%v,%v}",
		encB64([]byte(n.Header)),
		encB64([]byte(n.PHeader)),
		encB64([]byte(nonce)),
		encB64([]byte(n.TX.Payload)),
		encB64([]byte(timestamp)))
}

func Unpack(raw string) (*Node, error) {
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

	return &Node{
		Header:    string(header),
		PHeader:   string(pHeader),
		Nonce:     int64(nonce),
		TX:        tx,
		Timestamp: timestamp}, nil
}

func parseTX(txRaw string) (txn tx, err error) {
	parsedTX := strings.Split(txRaw, ":")
	if len(parsedTX) != 2 {
		err = errors.New("invalid TX")
		return
	}
	payloadRaw, ds := parsedTX[0], parsedTX[1]
	dsBytes, err := base64.StdEncoding.DecodeString(ds)
	ds = string(dsBytes)
	if err != nil {
		return
	}
	payload := strings.Split(payloadRaw, ",")
	amount, err := strconv.ParseFloat(payload[2], 64)
	if err != nil {
		return
	}
	change, err := strconv.ParseFloat(payload[3], 64)
	if err != nil {
		return
	}
	timestamp, err := strconv.ParseInt(payload[5], 10, 64)
	if err != nil {
		return
	}
	txn = tx{
		WAddr:     payload[0],
		RecWAddr:  payload[1],
		Amount:    amount,
		Change:    change,
		BHash:     payload[4],
		Timestamp: timestamp,
		DS:        ds,
		Payload:   payloadRaw}
	return
}
