package emvco

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dongri/emvco-qrcode/crc16"
)

// BasePayload ...
type BasePayload interface {
	PayloadFormatIndicator() string
	MerchantAccountInformation() string
	MerchantCategoryCode() string
	TransactionCurrency() string
	CountryCode() string
	MerchantName() string
	MerchantCity() string
	CRC() string
}

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
	AdditionalDataFieldTemplate         AdditionalDataFieldTemplate
	CRC                                 string
	MerchantInformationLanguageTemplate MerchantInformationLanguageTemplate
	RFUForEMVCo                         string
	UnreservedTemplates                 string
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
func (c *EMVQR) GeneratePayload() string {
	s := ""
	if c.PayloadFormatIndicator != "" {
		s += format(IDPayloadFormatIndicator, c.PayloadFormatIndicator)
	}
	if c.PointOfInitiationMethod != "" {
		s += format(IDPointOfInitiationMethod, c.PointOfInitiationMethod)
	}
	if c.MerchantAccountInformation != "" {
		s += format(IDMerchantAccountInformation, c.MerchantAccountInformation)
	}
	if c.MerchantCategoryCode != "" {
		s += format(IDMerchantCategoryCode, c.MerchantCategoryCode)
	}
	if c.TransactionCurrency != "" {
		s += format(IDTransactionCurrency, c.TransactionCurrency)
	}
	if c.TransactionAmount > 0 {
		s += format(IDTransactionAmount, formatAmount(c.TransactionAmount))
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
	}
	if c.MerchantName != "" {
		s += format(IDMerchantName, c.MerchantName)
	}
	if c.MerchantCity != "" {
		s += format(IDMerchantCity, c.MerchantCity)
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
	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	crc := crc16.Checksum([]byte(s+IDCRC+"04"), table)
	crcStr := formatCrc(crc)
	s += format(IDCRC, crcStr)
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
