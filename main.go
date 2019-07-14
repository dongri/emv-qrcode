package main

import (
	"log"

	"github.com/dongri/emvco-qrcode/emvco/cpm"
	"github.com/dongri/emvco-qrcode/emvco/mpm"
)

func main() {

	// MPM Generate
	emvqr := new(mpm.EMVQR)
	emvqr.SetPayloadFormatIndicator("01")
	emvqr.SetPointOfInitiationMethod("12") // 11 is static qrcode
	merchantAccountInformationJCB := new(mpm.MerchantAccountInformation)
	merchantAccountInformationJCB.SetGloballyUniqueIdentifier("D123456")
	merchantAccountInformationJCB.SetPaymentNetworkSpecific("13", "JCB1234567890")
	emvqr.AddMerchantAccountInformation(mpm.ID("29"), merchantAccountInformationJCB)

	merchantAccountInformationMaster := new(mpm.MerchantAccountInformation)
	merchantAccountInformationMaster.SetGloballyUniqueIdentifier("M123456")
	merchantAccountInformationMaster.SetPaymentNetworkSpecific("04", "MASTER1234567890")
	emvqr.AddMerchantAccountInformation(mpm.ID("31"), merchantAccountInformationMaster)

	emvqr.SetMerchantCategoryCode("5311")
	emvqr.SetTransactionCurrency("392")
	emvqr.SetTransactionAmount("999.123")
	emvqr.SetCountryCode("JP")
	emvqr.SetMerchantName("DONGRI")
	emvqr.SetMerchantCity("TOKYO")
	additionalTemplate := new(mpm.AdditionalDataFieldTemplate)
	additionalTemplate.SetBillNumber("hoge")
	additionalTemplate.SetReferenceLabel("fuga")
	additionalTemplate.SetTerminalLabel("piyo")
	emvqr.SetAdditionalDataFieldTemplate(additionalTemplate)
	code, err := mpm.Encode(emvqr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(code)
	// 0002010102121313JCB12345678900416MASTER12345678905204531153033925407999.1235802JP5906DONGRI6005TOKYO62240104hoge0504fuga0704piyo6304C343

	// MPM Parse
	emvqr, err = mpm.Decode("00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING64200002ZH0104最佳运输0202北京540523.7253031565502016233030412340603***0708A60086670902ME91320016A0112233449988770708123456786304A13A")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(emvqr)

	raw := emvqr.RawData()
	log.Println("\n" + raw)

	binary := emvqr.BinaryData()
	log.Println("\n" + binary)

	json := emvqr.JSON()
	log.Println(json)

	// CPM
	qr := new(cpm.EMVQR)
	qr.DataPayloadFormatIndicator = "CPV01"

	appTemplate1 := new(cpm.ApplicationTemplate)
	appTemplate1.DataApplicationDefinitionFileName = "A0000000555555"
	appTemplate1.DataApplicationLabel = "Product1"
	qr.ApplicationTemplates = append(qr.ApplicationTemplates, *appTemplate1)

	appTemplate2 := new(cpm.ApplicationTemplate)
	appTemplate2.DataApplicationDefinitionFileName = "A0000000666666"
	appTemplate2.DataApplicationLabel = "Product2"
	qr.ApplicationTemplates = append(qr.ApplicationTemplates, *appTemplate2)

	cdt := new(cpm.CommonDataTemplate)
	cdt.DataApplicationPAN = "1234567890123458"
	cdt.DataCardholderName = "CARDHOLDER/EMV"
	cdt.DataLanguagePreference = "ruesdeen"

	cdtt := new(cpm.CommonDataTransparentTemplate)
	cdtt.DataIssuerApplicationData = "06010A03000000"
	cdtt.DataApplicationCryptogram = "584FD385FA234BCC"
	cdtt.DataApplicationTransactionCounter = "0001"
	cdtt.DataUnpredictableNumber = "6D58EF13"
	cdt.CommonDataTransparentTemplates = append(cdt.CommonDataTransparentTemplates, *cdtt)

	qr.CommonDataTemplates = append(qr.CommonDataTemplates, *cdt)

	comQRCode, err := qr.GeneratePayload()
	if err != nil {
		log.Println(err)
	}
	log.Println(comQRCode)
	// hQVDUFYwMWETTwegAAAAVVVVUAhQcm9kdWN0MWETTwegAAAAZmZmUAhQcm9kdWN0MmJJWggSNFZ4kBI0WF8gDkNBUkRIT0xERVIvRU1WXy0IcnVlc2RlZW5kIZ8QBwYBCgMAAACfJghYT9OF+iNLzJ82AgABnzcEbVjvEw==

}
