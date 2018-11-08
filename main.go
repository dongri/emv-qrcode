package main

import (
	"fmt"

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
	additionalTemplate.BillNumber = "JPU20181018123456123456"
	additionalTemplate.ReferenceLabel = "JPU20181018123456123456"
	additionalTemplate.TerminalLabel = "123456"

	emvqr.AdditionalDataFieldTemplate = *additionalTemplate

	s := emvqr.GeneratePayload()
	fmt.Println(s)
}
