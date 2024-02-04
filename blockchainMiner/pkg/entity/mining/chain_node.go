package mining

type node struct {
	Header    string //header
	PHeader   string //parent header
	Nonce     int64  //nonce
	TX        tx     //tx
	Timestamp int64  //timestamp
}
