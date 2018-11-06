package emvco

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dongri/emvco-qrcode/crc16"
)

// const ....
const (
	IDPayloadFormatIndicator              = "00" // (M) Payload Format Indicator
	IDPointOfInitiationMethod             = "01" // (O) Point of Initiation Method
	IDMerchantAccountInformation          = "15" // (M) 02-51 Merchant Account Information (At least one Merchant Account Information data object shall be present.)
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
	IDCRC                                 = "63" // (M) CRC
	IDMerchantInformationLanguageTemplate = "64" // (O) Merchant Information— Language Template
	IDRFUForEMVCo                         = "65" // (O) 65-79 RFU for EMVCo
	IDUnreservedTemplates                 = "80" // (O) 80-99 Unreserved Templates
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
	PayloadFormatIndicator              PayloadFormatIndicator
	PointOfInitiationMethod             PointOfInitiationMethod
	MerchantAccountInformation          string
	MerchantCategoryCode                string
	TransactionCurrency                 string
	TransactionAmount                   float64
	TipOrConvenienceIndicator           string
	ValueOfConvenienceFeeFixed          string
	ValueOfConvenienceFeePercentage     string
	CountryCode                         string
	MerchantName                        string
	MerchantCity                        string
	PostalCode                          string
	AdditionalDataFieldTemplate         string
	CRC                                 string
	MerchantInformationLanguageTemplate string
	RFUForEMVCo                         string
	UnreservedTemplates                 string
	BillNumber                          string
	ReferenceLabel                      string
	TerminalLabel                       string
	LanguageReference                   string
	MerchantNameAlternateLanguage       string
}

// PayloadFormatIndicator ...
type PayloadFormatIndicator struct {
	ID    string
	Value string
}

// PointOfInitiationMethod ...
type PointOfInitiationMethod struct {
	ID    string
	Value string
}

// SetPayloadFormatIndicator ...
func (c *EMVQR) SetPayloadFormatIndicator(value string) {
	payloadFormatIndicator := new(PayloadFormatIndicator)
	payloadFormatIndicator.ID = "00"
	payloadFormatIndicator.Value = value
	c.PayloadFormatIndicator = *payloadFormatIndicator
}

// GeneratePayload ...
func (c *EMVQR) GeneratePayload() string {
	s := format(c.PayloadFormatIndicator.ID, c.PayloadFormatIndicator.Value)
	s += format(c.PointOfInitiationMethod.ID, c.PointOfInitiationMethod.Value)
	s += format(IDMerchantAccountInformation, c.MerchantAccountInformation)
	s += format(IDMerchantCategoryCode, c.MerchantCategoryCode)
	s += format(IDTransactionCurrency, c.TransactionCurrency)
	s += format(IDTransactionAmount, formatAmount(c.TransactionAmount))
	s += format(IDCountryCode, c.CountryCode)
	s += format(IDMerchantName, c.MerchantName)
	s += format(IDMerchantCity, c.MerchantCity)
	s += format(IDPostalCode, c.PostalCode)
	s += format(IDAdditionalDataFieldTemplate,
		format(AdditionalIDBillNumber, c.BillNumber)+
			format(AdditionalIDReferenceLabel, c.ReferenceLabel)+
			format(AdditionalIDTerminalLabel, c.TerminalLabel))

	//s = "00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING64200002ZH0104最佳运输0202北京540523.7253031565502016233030412340603***0708A60086670902ME91320016A011223344998877070812345678"

	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	crc := crc16.Checksum([]byte(s+IDCRC+"04"), table)

	crcStr := formatCrc(crc)
	s += format(IDCRC, crcStr)
	return s
}

func format(id, value string) string {
	length := strconv.Itoa(len(value))
	length = "00" + length
	return id + length[len(length)-2:] + value
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.0f", amount)
}

func formatCrc(crcValue uint16) string {
	crcValueString := strconv.FormatUint(uint64(crcValue), 16)
	s := "0000" + strings.ToUpper(crcValueString)
	return s[len(s)-4:]
}
