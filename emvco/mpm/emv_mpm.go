package mpm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/dongri/emvco-qrcode/crc16"
)

// const ....
const (
	IDPayloadFormatIndicator              = "00" // (M) Payload Format Indicator
	IDPointOfInitiationMethod             = "01" // (O) Point of Initiation Method
	IDMerchantCategoryCode                = "52" // (M) Merchant Category Code
	IDTransactionCurrency                 = "53" // (M) Transaction Currency
	IDTransactionAmount                   = "54" // (C) Transaction Amount
	IDTipOrConvenienceIndicator           = "55" // (O) Tip or Convenience Indicator
	IDValueOfConvenienceFeeFixed          = "56" // (C) Value of Convenience Fee Fixed
	IDValueOfConvenienceFeePercentage     = "57" // (C) Value of Convenience Fee Percentage
	IDCountryCode                         = "58" // (M) Country Code
	IDMerchantName                        = "59" // (M) Merchant Name
	IDMerchantCity                        = "60" // (M) Merchant City
	IDPostalCode                          = "61" // (O) Postal Code
	IDAdditionalDataFieldTemplate         = "62" // (O) Additional Data Field Template
	IDMerchantInformationLanguageTemplate = "64" // (O) Merchant Information— Language Template
	IDRFUForEMVCo                         = "65" // (O) 65-79 RFU for EMVCo
	IDUnreservedTemplates                 = "80" // (O) 80-99 Unreserved Templates
	IDCRC                                 = "63" // (M) CRC
)

// Data Objects for Additional Data Field Template (ID "62")
const (
	AdditionalIDBillNumber                     = "01"
	AdditionalIDMobileNumber                   = "02"
	AdditionalIDStoreLabel                     = "03"
	AdditionalIDLoyaltyNumber                  = "04"
	AdditionalIDReferenceLabel                 = "05"
	AdditionalIDCustomerLabel                  = "06"
	AdditionalIDTerminalLabel                  = "07"
	AdditionalIDPurposeTransaction             = "08"
	AdditionalIDAdditionalConsumerDataRequest  = "09"
	AdditionalIDRFUforEMVCo                    = "10" // 10-49
	AdditionalIDPaymentSystemSpecificTemplates = "50" // 50-99
)

// Data Objects for Merchant Information—Language Template (ID "64")
const (
	MerchantInformationIDLanguagePreference = "00"
	MerchantInformationIDMerchantName       = "01"
	MerchantInformationIDMerchantCity       = "02"
	MerchantInformationIDRFUforEMVCo        = "03" // 03-99
)

// EMVQR ...
type EMVQR struct {
	PayloadFormatIndicator              string
	PointOfInitiationMethod             string
	MerchantAccountInformation          []MerchantAccountInformation // (M) Tag: 02-51 Merchant Account Information (At least one Merchant Account Information data object shall be present.)
	MerchantCategoryCode                string
	TransactionCurrency                 string
	TransactionAmount                   string
	TipOrConvenienceIndicator           string
	ValueOfConvenienceFeeFixed          string
	ValueOfConvenienceFeePercentage     string
	CountryCode                         string
	MerchantName                        string
	MerchantCity                        string
	PostalCode                          string
	AdditionalDataFieldTemplate         AdditionalDataFieldTemplate         // Tag: 62
	CRC                                 string                              // Tag: 63
	MerchantInformationLanguageTemplate MerchantInformationLanguageTemplate // Tag: 64
	RFUForEMVCo                         string                              // Tag: 65-79 RFU for EMVCo
	UnreservedTemplates                 string                              // Tag: 80-99 Unreserved Templates
}

// MerchantAccountInformation ...
type MerchantAccountInformation struct {
	Tag   string
	Value string
}

// AdditionalDataFieldTemplate ...
type AdditionalDataFieldTemplate struct {
	BillNumber                     string
	MobileNumber                   string
	StoreLabel                     string
	LoyaltyNumber                  string
	ReferenceLabel                 string
	CustomerLabel                  string
	TerminalLabel                  string
	PurposeTransaction             string
	AdditionalConsumerDataRequest  string
	RFUforEMVCo                    string // 10-49
	PaymentSystemSpecificTemplates string // 50-99
}

// MerchantInformationLanguageTemplate ...
type MerchantInformationLanguageTemplate struct {
	LanguagePreference string
	MerchantName       string
	MerchantCity       string
	RFUforEMVCo        string // 03-99
}

// GeneratePayload ...
func (c *EMVQR) GeneratePayload() (string, error) {
	s := ""
	if c.PayloadFormatIndicator != "" {
		s += format(IDPayloadFormatIndicator, c.PayloadFormatIndicator)
	} else {
		return "", fmt.Errorf("PayloadFormatIndicator is mandatory")
	}
	if c.PointOfInitiationMethod != "" {
		s += format(IDPointOfInitiationMethod, c.PointOfInitiationMethod)
	}
	if len(c.MerchantAccountInformation) > 0 {
		for _, t := range c.MerchantAccountInformation {
			s += format(t.Tag, t.Value)
		}
	} else {
		return "", fmt.Errorf("MerchantAccountInformation is mandatory")
	}
	if c.MerchantCategoryCode != "" && len(c.MerchantCategoryCode) == 4 {
		s += format(IDMerchantCategoryCode, c.MerchantCategoryCode)
	} else {
		return "", fmt.Errorf("MerchantCategoryCode is mandatory")
	}
	if c.TransactionCurrency != "" {
		s += format(IDTransactionCurrency, c.TransactionCurrency)
	} else {
		return "", fmt.Errorf("TransactionCurrency is mandatory")
	}
	if c.TransactionAmount != "" {
		s += format(IDTransactionAmount, c.TransactionAmount)
	}
	if c.TipOrConvenienceIndicator != "" {
		s += format(IDTipOrConvenienceIndicator, c.TipOrConvenienceIndicator)
	}
	if c.ValueOfConvenienceFeeFixed != "" {
		s += format(IDValueOfConvenienceFeeFixed, c.ValueOfConvenienceFeeFixed)
	}
	if c.ValueOfConvenienceFeePercentage != "" {
		s += format(IDValueOfConvenienceFeePercentage, c.ValueOfConvenienceFeePercentage)
	}
	if c.CountryCode != "" {
		s += format(IDCountryCode, c.CountryCode)
	} else {
		return "", fmt.Errorf("CountryCode is mandatory")
	}
	if c.MerchantName != "" {
		s += format(IDMerchantName, c.MerchantName)
	} else {
		return "", fmt.Errorf("MerchantName is mandatory")
	}
	if c.MerchantCity != "" {
		s += format(IDMerchantCity, c.MerchantCity)
	} else {
		return "", fmt.Errorf("MerchantCity is mandatory")
	}
	if c.PostalCode != "" {
		s += format(IDPostalCode, c.PostalCode)
	}
	if (AdditionalDataFieldTemplate{}) != c.AdditionalDataFieldTemplate {
		t := c.AdditionalDataFieldTemplate
		template := ""
		if t.BillNumber != "" {
			template += format(AdditionalIDBillNumber, t.BillNumber)
		}
		if t.MobileNumber != "" {
			template += format(AdditionalIDMobileNumber, t.MobileNumber)
		}
		if t.StoreLabel != "" {
			template += format(AdditionalIDStoreLabel, t.StoreLabel)
		}
		if t.LoyaltyNumber != "" {
			template += format(AdditionalIDLoyaltyNumber, t.LoyaltyNumber)
		}
		if t.ReferenceLabel != "" {
			template += format(AdditionalIDReferenceLabel, t.ReferenceLabel)
		}
		if t.CustomerLabel != "" {
			template += format(AdditionalIDCustomerLabel, t.CustomerLabel)
		}
		if t.TerminalLabel != "" {
			template += format(AdditionalIDTerminalLabel, t.TerminalLabel)
		}
		if t.PurposeTransaction != "" {
			template += format(AdditionalIDPurposeTransaction, t.PurposeTransaction)
		}
		if t.AdditionalConsumerDataRequest != "" {
			template += format(AdditionalIDAdditionalConsumerDataRequest, t.AdditionalConsumerDataRequest)
		}
		if t.RFUforEMVCo != "" {
			template += format(AdditionalIDRFUforEMVCo, t.RFUforEMVCo)
		} // 10-49
		if t.PaymentSystemSpecificTemplates != "" {
			template += format(AdditionalIDPaymentSystemSpecificTemplates, t.PaymentSystemSpecificTemplates)
		} // 50-99
		s += format(IDAdditionalDataFieldTemplate, template)
	}
	if (MerchantInformationLanguageTemplate{}) != c.MerchantInformationLanguageTemplate {
		t := c.MerchantInformationLanguageTemplate
		template := ""
		if t.LanguagePreference != "" {
			template += format(MerchantInformationIDLanguagePreference, t.LanguagePreference)
		}
		if t.MerchantName != "" {
			template += format(MerchantInformationIDMerchantName, t.MerchantName)
		}
		if t.MerchantCity != "" {
			template += format(MerchantInformationIDMerchantCity, t.MerchantCity)
		}
		if t.RFUforEMVCo != "" {
			template += format(MerchantInformationIDRFUforEMVCo, t.RFUforEMVCo)
		} // 03-99
		s += format(IDMerchantInformationLanguageTemplate, template)
	}
	if c.RFUForEMVCo != "" {
		s += format(IDRFUForEMVCo, c.RFUForEMVCo)
	}
	if c.UnreservedTemplates != "" {
		s += format(IDUnreservedTemplates, c.UnreservedTemplates)
	}
	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	crc := crc16.Checksum([]byte(s+IDCRC+"04"), table)
	crcStr := formatCrc(crc)
	s += format(IDCRC, crcStr)
	return s, nil
}

func format(id, value string) string {
	length := utf8.RuneCountInString(value)
	lengthStr := strconv.Itoa(length)
	lengthStr = "00" + lengthStr
	return id + lengthStr[len(lengthStr)-2:] + value
}

func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}

func formatCrc(crcValue uint16) string {
	crcValueString := strconv.FormatUint(uint64(crcValue), 16)
	s := "0000" + strings.ToUpper(crcValueString)
	return s[len(s)-4:]
}

// Decode ...
func Decode(payload string) *EMVQR {
	//payload
	return new(EMVQR)
}

// ParsePayload ...
func ParsePayload(emvString string) (*EMVQR, error) {
	var emvData = map[string]string{}
	var merchantAccountInformations []MerchantAccountInformation
	emvQR := new(EMVQR)
	inputText := emvString
	for len(inputText) > 0 {
		data, remainingText, err := readNext(inputText)
		if err != nil {
			return nil, err
		}
		emvData[data["id"]] = data["value"]
		id, err := strconv.Atoi(data["id"])
		if err != nil {
			return nil, err
		}
		if id >= 2 && id <= 51 {
			merchantAccountInformations = append(merchantAccountInformations, ParseMerchantAccountInformation(data))
		}
		if id == 62 {
			emvQR.AdditionalDataFieldTemplate, err = ParseAdditionalDataFieldTemplate(data)
			if err != nil {
				return nil, err
			}
		}
		if id == 64 {
			emvQR.MerchantInformationLanguageTemplate, err = ParseMerchantInformationLanguageTemplate(data)
			if err != nil {
				return nil, err
			}
		}
		inputText = remainingText
	}
	emvQR.MerchantAccountInformation = merchantAccountInformations

	emvQR.PayloadFormatIndicator = emvData[IDPayloadFormatIndicator]
	emvQR.PointOfInitiationMethod = emvData[IDPointOfInitiationMethod]
	emvQR.MerchantCategoryCode = emvData[IDMerchantCategoryCode]
	emvQR.TransactionCurrency = emvData[IDTransactionCurrency]
	if _, ok := emvData[IDTransactionAmount]; ok {
		emvQR.TransactionAmount = emvData[IDTransactionAmount]
	}
	emvQR.TipOrConvenienceIndicator = emvData[IDTipOrConvenienceIndicator]
	emvQR.ValueOfConvenienceFeeFixed = emvData[IDValueOfConvenienceFeeFixed]
	emvQR.ValueOfConvenienceFeePercentage = emvData[IDValueOfConvenienceFeePercentage]
	emvQR.CountryCode = emvData[IDCountryCode]
	emvQR.MerchantName = emvData[IDMerchantName]
	emvQR.MerchantCity = emvData[IDMerchantCity]
	emvQR.PostalCode = emvData[IDPostalCode]
	emvQR.CRC = emvData[IDCRC]
	emvQR.RFUForEMVCo = emvData[IDRFUForEMVCo]
	emvQR.UnreservedTemplates = emvData[IDUnreservedTemplates]
	return emvQR, nil
}

// ParseMerchantAccountInformation ...
func ParseMerchantAccountInformation(data map[string]string) MerchantAccountInformation {
	merchantAccountInformation := new(MerchantAccountInformation)
	merchantAccountInformation.Tag = data["id"]
	merchantAccountInformation.Value = data["value"]
	return *merchantAccountInformation
}

// ParseAdditionalDataFieldTemplate ...
func ParseAdditionalDataFieldTemplate(data map[string]string) (AdditionalDataFieldTemplate, error) {
	inputText := data["value"]
	additionalDataFieldTemplate := new(AdditionalDataFieldTemplate)
	for len(inputText) > 0 {
		data, remainingText, err := readNext(inputText)
		if err != nil {
			return *additionalDataFieldTemplate, err
		}
		value := data["value"]
		switch data["id"] {
		case AdditionalIDBillNumber:
			additionalDataFieldTemplate.BillNumber = value
		case AdditionalIDMobileNumber:
			additionalDataFieldTemplate.MobileNumber = value
		case AdditionalIDStoreLabel:
			additionalDataFieldTemplate.StoreLabel = value
		case AdditionalIDLoyaltyNumber:
			additionalDataFieldTemplate.LoyaltyNumber = value
		case AdditionalIDReferenceLabel:
			additionalDataFieldTemplate.ReferenceLabel = value
		case AdditionalIDCustomerLabel:
			additionalDataFieldTemplate.CustomerLabel = value
		case AdditionalIDTerminalLabel:
			additionalDataFieldTemplate.TerminalLabel = value
		case AdditionalIDPurposeTransaction:
			additionalDataFieldTemplate.PurposeTransaction = value
		case AdditionalIDAdditionalConsumerDataRequest:
			additionalDataFieldTemplate.AdditionalConsumerDataRequest = value
		case AdditionalIDRFUforEMVCo:
			additionalDataFieldTemplate.RFUforEMVCo = value
		case AdditionalIDPaymentSystemSpecificTemplates:
			additionalDataFieldTemplate.PaymentSystemSpecificTemplates = value
		}
		inputText = remainingText
	}
	return *additionalDataFieldTemplate, nil
}

// ParseMerchantInformationLanguageTemplate ...
func ParseMerchantInformationLanguageTemplate(data map[string]string) (MerchantInformationLanguageTemplate, error) {
	inputText := data["value"]
	merchantInformationLanguageTemplate := new(MerchantInformationLanguageTemplate)
	for len(inputText) > 0 {
		data, remainingText, err := readNext(inputText)
		if err != nil {
			return *merchantInformationLanguageTemplate, err
		}
		value := data["value"]
		switch data["id"] {
		case MerchantInformationIDLanguagePreference:
			merchantInformationLanguageTemplate.LanguagePreference = value
		case MerchantInformationIDMerchantName:
			merchantInformationLanguageTemplate.MerchantName = value
		case MerchantInformationIDMerchantCity:
			merchantInformationLanguageTemplate.MerchantCity = value
		case MerchantInformationIDRFUforEMVCo:
			merchantInformationLanguageTemplate.RFUforEMVCo = value
		}
		inputText = remainingText
	}
	return *merchantInformationLanguageTemplate, nil
}

func readNext(inputText string) (map[string]string, string, error) {
	id := substring(inputText, 0, 2)
	length, err := strconv.Atoi(substring(inputText, 2, 2))
	if err != nil {
		return nil, "", err
	}
	value := substring(inputText, 4, length)
	data := map[string]string{
		"id":    id,
		"value": value,
	}
	substringLength := len(inputText) - 4
	var remainingText string
	if substringLength > 0 {
		remainingText = substring(inputText, length+4, substringLength)
	}
	return data, remainingText, nil
}

func substring(str string, start, length int) string {
	if start < 0 || length <= 0 {
		return str
	}
	r := []rune(str)
	if start+length > len(r) {
		return string(r[start:])
	}
	return string(r[start : start+length])
}
