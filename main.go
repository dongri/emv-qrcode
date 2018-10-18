package main

import (
	"fmt"

	"github.com/dongri/emvco-qrcode/emvco"
)

func main() {

	AcquirerIIN := "111"
	ForwardingIIN := "222"
	MerchantID := "123456789012345"
	mai := AcquirerIIN + ForwardingIIN + MerchantID

	emvqr := new(emvco.EMVQR)
	emvqr.PayloadFormatIndicator = "01"
	emvqr.PointOfInitiationMethod = "12" // 11 is static qrcode
	emvqr.MerchantAccountInformation = mai
	emvqr.MerchantCategoryCode = "5311"
	emvqr.TransactionCurrency = "392"
	emvqr.TransactionAmount = 999
	emvqr.CountryCode = "JP"
	emvqr.MerchantName = "DONGRI"
	emvqr.MerchantCity = "TOKYO"

	emvqr.BillNumber = "JPU20181018123456123456"
	emvqr.ReferenceLabel = "JPU20181018123456123456"
	emvqr.TerminalLabel = "123456"

	s := emvqr.GeneratePayload()
	fmt.Println(s)
}
