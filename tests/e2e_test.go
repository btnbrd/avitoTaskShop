package main

import "testing"

func main() {
	testing.RunTests(func(pat, str string) (bool, error) {
		return true, nil
	},
		[]testing.InternalTest{
			{"TestCreate", TestSignIN},
			{"TestBuyMerch", TestBuyMerch},
			{"TestBankruptBecome", TestBecameBankrupt},
			{"BankruptSend", TestBankruptSend},
			{"Bankrupt buy", TestBankruptBuy},
			{"Send", TestSend},
			{"Send to yourself", TestSendToYouself},
			{"Send to absent", TestSendToAbsent}})

}
