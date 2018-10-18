package emvco

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dongri/emvco-qrcode/crc16"
)

// const ....
const (
	IDPayloadFormatIndicator              = "00" // Payload Format Indicator
	IDPointOfInitiationMethod             = "01" // Point of Initiation Method
	IDMerchantAccountInformation          = "15" // 02-51 Merchant Account Information (At least one Merchant Account Information data object shall be present.)
	IDMerchantCategoryCode                = "52" // Merchant Category Code
	IDTransactionCurrency                 = "53" // Transaction Currency
	IDTransactionAmount                   = "54" // Transaction Amount
	IDTipOrConvenienceIndicator           = "55" // Tip or Convenience Indicator
	IDValueOfConvenienceFeeFixed          = "56" // Value of Convenience Fee Fixed
	IDValueOfConvenienceFeePercentage     = "57" // Value of Convenience Fee Percentage
	IDCountryCode                         = "58" // Country Code
	IDMerchantName                        = "59" // Merchant Name
	IDMerchantCity                        = "60" // Merchant City
	IDPostalCode                          = "61" // Postal Code
	IDAdditionalDataFieldTemplate         = "62" // Additional Data Field Template
	IDCRC                                 = "63" // CRC
	IDMerchantInformationLanguageTemplate = "64" // Merchant Information— Language Template
	IDRFUForEMVCo                         = "65" // 65-79 RFU for EMVCo
	IDUnreservedTemplates                 = "80" // 80-99 Unreserved Templates
)

// const ...
const (
	IDBillNumber                    = "01"
	IDMobileNumber                  = "02"
	IDStoreLabel                    = "03"
	IDLoyaltyNumber                 = "04"
	IDReferenceLabel                = "05"
	IDCustomerLabel                 = "06"
	IDTerminalLabel                 = "07"
	IDPurposeTransaction            = "08"
	IDAdditionalConsumerDataRequest = "09"

	IDLanguageReference             = "00"
	IDMerchantNameAlternateLanguage = "01"
)

// EMVQR ...
type EMVQR struct {
	PayloadFormatIndicator              string
	PointOfInitiationMethod             string
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

// GeneratePayload ...
func (c *EMVQR) GeneratePayload() string {
	s := format(IDPayloadFormatIndicator, c.PayloadFormatIndicator)
	s += format(IDPointOfInitiationMethod, c.PointOfInitiationMethod)
	s += format(IDMerchantAccountInformation, c.MerchantAccountInformation)
	s += format(IDMerchantCategoryCode, c.MerchantCategoryCode)
	s += format(IDTransactionCurrency, c.TransactionCurrency)
	s += format(IDTransactionAmount, formatAmount(c.TransactionAmount))
	s += format(IDCountryCode, c.CountryCode)
	s += format(IDMerchantName, c.MerchantName)
	s += format(IDMerchantCity, c.MerchantCity)
	s += format(IDPostalCode, c.PostalCode)
	s += format(IDAdditionalDataFieldTemplate,
		format(IDBillNumber, c.BillNumber)+
			format(IDReferenceLabel, c.ReferenceLabel)+
			format(IDTerminalLabel, c.TerminalLabel))

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
