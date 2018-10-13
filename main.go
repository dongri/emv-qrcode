package main

import (
	"fmt"

	"github.com/dongri/emvco-qrcode/emvco"
)

func main() {
	emvqr := new(emvco.EMVQR)
	emvqr.PayloadFormatIndicator = "01"
	emvqr.PointOfInitiationMethod = "12" // 11 is static qrcode
	emvqr.MerchantAccountInformation = "9090909090"
	emvqr.MerchantCategoryCode = "5311"
	emvqr.TransactionCurrency = "165"
	emvqr.TransactionAmount = 999
	emvqr.CountryCode = "JP"
	emvqr.MerchantName = "DongriShop"
	emvqr.MerchantCity = "Tokyo"
	emvqr.PostalCode = "1360076"
	s := emvqr.GeneratePayload()
	fmt.Println(s)
}
