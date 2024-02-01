package entity_test

import (
	"blockchain/pkg/entity"
	"strings"
	"testing"
)

func TestTXParsing(t *testing.T) {
	txSt := "waddr1,waddr2,69.420,420.69,bHash,696969:digsign"
	m := entity.Miner{}
	m.AddTX(txSt)
	tx, ok := m.NextTx()
	if !ok {
		t.Fatal("TX havent been added")
	}
	if !strings.EqualFold(tx.WAddr, "waddr1") {
		t.Fatal("1 param")
	}
	if !strings.EqualFold(tx.RecWAddr, "waddr2") {
		t.Fatal("2 param")
	}
	if tx.Amount > 69.421 || tx.Amount < 69.419 {
		t.Fatal("3 param")
	}
	if tx.Change > 420.70 || tx.Change < 420.68 {
		t.Fatal("4 param")
	}
	if !strings.EqualFold(tx.BHash, "bHash") {
		t.Fatal("5 param")
	}
	if tx.Timestamp != 696969 {
		t.Fatal("6 param")
	}
	if !strings.EqualFold(tx.DS, "digsign") {
		t.Fatal("7 param")
	}
}
