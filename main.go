package main

import (
	"log"

	"github.com/dongri/emvco-qrcode/emvco"
)

func main() {
	emvqr := new(emvco.EMVQR)
	emvqr.PayloadFormatIndicator = "01"
	emvqr.PointOfInitiationMethod = "12" // 11 is static qrcode
	emvqr.MerchantAccountInformation = "ABCDEF1234567890"
	emvqr.MerchantCategoryCode = "5311"
	emvqr.TransactionCurrency = "392"
	emvqr.TransactionAmount = 999
	emvqr.CountryCode = "JP"
	emvqr.MerchantName = "DONGRI"
	emvqr.MerchantCity = "TOKYO"

	additionalTemplate := new(emvco.AdditionalDataFieldTemplate)
	additionalTemplate.BillNumber = "hoge"
	additionalTemplate.ReferenceLabel = "fuga"
	additionalTemplate.TerminalLabel = "piyo"

	emvqr.AdditionalDataFieldTemplate = *additionalTemplate

	qrcodeData, err := emvqr.GeneratePayload()
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(qrcodeData)
}
