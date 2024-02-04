package ui

import (
	"blockchain/pkg/cryptography"
	"blockchain/pkg/entity"
	"blockchain/pkg/utils"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gen2brain/iup-go/iup"
)

func Start(c *entity.Client, wg *sync.WaitGroup) {

	defer wg.Done()
	iup.Open()
	defer iup.Close()

	privateKey := iup.Text().SetHandle("pk")
	privateKey.SetAttributes(`SIZE=80x`)

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

	pubKey := iup.Text().SetAttributes(`SIZE=80x`)
	amount := iup.Text().SetAttributes(`SIZE=80x`)
	change := iup.Text().SetAttributes(`SIZE=80x`)

	key1Frame := iup.Frame(
		iup.Hbox(
			iup.Vbox(
				iup.Label("Enter public key 1"),
				pubKey,
			),
			iup.Vbox(
				iup.Label("Enter amount to send"),
				amount,
			),
		),
	).SetAttribute("TITLE", "reciever 1")

	doneBtn := iup.Button("Send")

	doneBtnCb := func(ih iup.Ihandle, button, pressed, x, y int, status string) int {
		if pressed == 1 {
			privateKeyText := iup.GetAttribute(privateKey, "VALUE")
			pubKeyText := iup.GetAttribute(pubKey, "VALUE")
			amountText := iup.GetAttribute(amount, "VALUE")
			changeText := iup.GetAttribute(change, "VALUE")

			createTX(c, privateKeyText, pubKeyText, amountText, changeText)
		}
		return iup.DEFAULT
	}

	updateBalanceBtnCb := func(ih iup.Ihandle, button, pressed, x, y int, status string) int {
		if pressed == 1 {
			iup.GetHandle("balance").SetAttribute("VALUE", "123")
		}
		return iup.DEFAULT
	}

	genKeyBtnCb := func(ih iup.Ihandle, button, pressed, x, y int, status string) int {
		if pressed == 1 {
			kPair, err := cryptography.GenerateKeyPair()
			if err != nil {
				return iup.DEFAULT
			}
			privKey := kPair.PrivateHex()

			iup.GetHandle("pk").SetAttribute("VALUE", privKey)
		}
		return iup.DEFAULT
	}

	iup.SetCallback(doneBtn, "BUTTON_CB", iup.ButtonFunc(doneBtnCb))
	iup.SetCallback(updateBalanceBtn, "BUTTON_CB", iup.ButtonFunc(updateBalanceBtnCb))
	iup.SetCallback(genKeyPairBtn, "BUTTON_CB", iup.ButtonFunc(genKeyBtnCb))

	vbox1 := iup.Vbox(
		hbox1,
		key1Frame,
		iup.Label("Enter change:"),
		change,
		doneBtn,
	).SetAttributes("MARGIN=5x5, GAP=5")

	dlg := iup.Dialog(vbox1).SetAttributes(`TITLE="Sample", MENU=menu, ICON=img1`)
	dlg.SetHandle("dlg")

	iup.Map(dlg)

	iup.Show(dlg)
	iup.MainLoop()
}

func createTX(c *entity.Client, senderPrivKey string, recieverPubKey string, amount string, change string) {
	kPair, err := cryptography.GenKeysFromPrivate(senderPrivKey)
	if err != nil {
		fmt.Println("invalid private key")
		return
	}
	pubKey := kPair.Public()
	if err != nil {
		fmt.Println("couldn't return public key, err: " + err.Error())
		return
	}
	wAddr := cryptography.WaletAddr(pubKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	amount = strings.Replace(amount, ",", ".", -1)
	change = strings.Replace(change, ",", ".", -1)

	if !utils.IsNumber(amount) || !utils.IsNumber(change) {
		fmt.Println("passed amount or change is not a number")
		return
	}
	fAmount, _ := strconv.ParseFloat(amount, 64)
	fChange, _ := strconv.ParseFloat(change, 64)
	bHash := utils.ChooseBlock(fAmount + fChange)

	timestamp := fmt.Sprint(time.Now().Unix())

	payload := strings.Join([]string{wAddr, recieverPubKey, amount, change, bHash, timestamp}, ",")
	utils.SendTX(c, kPair, payload)
}
