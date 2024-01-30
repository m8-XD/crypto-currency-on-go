package ui

import (
	"fmt"
	"sync"

	"github.com/gen2brain/iup-go/iup"
)

func Start(wg *sync.WaitGroup) {

	defer wg.Done()
	iup.Open()
	defer iup.Close()

	privateKey := iup.Text().SetAttributes(`SIZE=80x`)

	genKeyPairBtn := iup.Button("generate private key")
	walletsFrame := iup.Frame(
		iup.Hbox(
			privateKey,
			genKeyPairBtn,
		),
	).SetAttribute("TITLE", "Private key")

	balance := iup.Text().SetHandle("balance")
	balance.SetAttribute("READONLY", "YES")
	balance.SetAttribute("VALUE", "0")

	updateBalanceBtn := iup.Button("update")

	balanceFrame := iup.Frame(
		iup.Vbox(
			balance,
			updateBalanceBtn,
		),
	).SetAttribute("TITLE", "Balance:")

	hbox1 := iup.Hbox(
		walletsFrame,
		balanceFrame,
	)

	pubKey1 := iup.Text().SetAttributes(`SIZE=80x`)
	pubKey2 := iup.Text().SetAttributes(`SIZE=80x`)
	amount1 := iup.Text().SetAttributes(`SIZE=80x`)
	amount2 := iup.Text().SetAttributes(`SIZE=80x`)
	change := iup.Text().SetAttributes(`SIZE=80x`)

	key1Frame := iup.Frame(
		iup.Hbox(
			iup.Vbox(
				iup.Label("Enter public key 1"),
				pubKey1,
			),
			iup.Vbox(
				iup.Label("Enter amount to send"),
				amount1,
			),
		),
	).SetAttribute("TITLE", "reciever 1")

	key2Frame := iup.Frame(
		iup.Hbox(
			iup.Vbox(
				iup.Label("Enter public key 2"),
				pubKey2,
			),
			iup.Vbox(
				iup.Label("Enter amount to send"),
				amount2,
			),
		),
	).SetAttribute("TITLE", "reciever 2")

	doneBtn := iup.Button("Send")

	doneBtnCb := func(ih iup.Ihandle, button, pressed, x, y int, status string) int {
		if pressed == 1 {
			privateKeyText := iup.GetAttribute(privateKey, "VALUE")
			fmt.Println("privateKeyText: " + privateKeyText)
			pubKey1Text := iup.GetAttribute(pubKey1, "VALUE")
			fmt.Println("pubKey1Text " + pubKey1Text)
			pubKey2Text := iup.GetAttribute(pubKey2, "VALUE")
			fmt.Println("pubKey2Text" + pubKey2Text)
			amount1Text := iup.GetAttribute(amount1, "VALUE")
			fmt.Println("amount1Text" + amount1Text)
			amount2Text := iup.GetAttribute(amount2, "VALUE")
			fmt.Println("amount2Text " + amount2Text)
			changeText := iup.GetAttribute(change, "VALUE")
			fmt.Println("changeText " + changeText)
		}
		return iup.DEFAULT
	}

	updateBalanceBtnCb := func(ih iup.Ihandle, button, pressed, x, y int, status string) int {
		if pressed == 1 {
			iup.GetHandle("balance").SetAttribute("VALUE", "123")

		}
		return iup.DEFAULT
	}

	iup.SetCallback(doneBtn, "BUTTON_CB", iup.ButtonFunc(doneBtnCb))
	iup.SetCallback(updateBalanceBtn, "BUTTON_CB", iup.ButtonFunc(updateBalanceBtnCb))

	vbox1 := iup.Vbox(
		hbox1,
		key1Frame,
		key2Frame,
		iup.Label("Enter change amount"),
		change,
		doneBtn,
	).SetAttributes("MARGIN=5x5, GAP=5")

	dlg := iup.Dialog(vbox1).SetAttributes(`TITLE="Sample", MENU=menu, ICON=img1`)
	dlg.SetHandle("dlg")

	iup.Map(dlg)

	iup.Show(dlg)
	iup.MainLoop()
}

func Error(msg string) {

}
