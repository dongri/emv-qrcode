package main

import (
	"log"

	"github.com/dongri/emvco-qrcode/emvco/cpm"
	"github.com/dongri/emvco-qrcode/emvco/mpm"
)

func main() {
	// MPM
	emvqr := new(mpm.EMVQR)
	emvqr.PayloadFormatIndicator = "01"
	emvqr.PointOfInitiationMethod = "12" // 11 is static qrcode
	emvqr.MerchantAccountInformation = "ABCDEF1234567890"
	emvqr.MerchantCategoryCode = "5311"
	emvqr.TransactionCurrency = "392"
	emvqr.TransactionAmount = 999
	emvqr.CountryCode = "JP"
	emvqr.MerchantName = "DONGRI"
	emvqr.MerchantCity = "TOKYO"

	additionalTemplate := new(mpm.AdditionalDataFieldTemplate)
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

	// CPM
	qr := new(cpm.EMVQR)
	qr.DataPayloadFormatIndicator = "CPV01"

	template := new(cpm.ApplicationTemplate)
	template.DataApplicationDefinitionFileName = "A0000000555555"
	template.DataApplicationLabel = "Product1"
	qr.ApplicationTemplates = append(qr.ApplicationTemplates, *template)

	template0 := new(cpm.ApplicationTemplate)
	template0.DataApplicationDefinitionFileName = "A0000000666666"
	template0.DataApplicationLabel = "Product2"
	qr.ApplicationTemplates = append(qr.ApplicationTemplates, *template0)

	template1 := new(cpm.CommonDataTemplate)
	template1.DataApplicationPAN = "1234567890123458"
	template1.DataCardholderName = "CARDHOLDER/EMV"
	template1.DataLanguagePreference = "ruesdeen"

	template2 := new(cpm.CommonDataTransparentTemplate)
	template2.DataIssuerApplicationData = "06010A03000000"
	template2.DataApplicationCryptogram = "584FD385FA234BCC"
	template2.DataApplicationTransactionCounter = "0001"
	template2.DataUnpredictableNumber = "6D58EF13"
	template1.CommonDataTransparentTemplates = append(template1.CommonDataTransparentTemplates, *template2)

	qr.CommonDataTemplates = append(qr.CommonDataTemplates, *template1)

	qrcode, err := qr.GeneratePayload()
	if err != nil {
		log.Println(err)
	}
	log.Println(qrcode)
	if qrcode != "hQVDUFYwMWETTwegAAAAVVVVUAhQcm9kdWN0MWETTwegAAAAZmZmUAhQcm9kdWN0MmJJWggSNFZ4kBI0WF8gDkNBUkRIT0xERVIvRU1WXy0IcnVlc2RlZW5kIZ8QBwYBCgMAAACfJghYT9OF+iNLzJ82AgABnzcEbVjvEw==" {
		log.Println("Diff")
	}
}
