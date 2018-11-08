# emvco-qrcode

```
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

qrcodeData := emvqr.GeneratePayload()
fmt.Println(qrcodeData)
```
