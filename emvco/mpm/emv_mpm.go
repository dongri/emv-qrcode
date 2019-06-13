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
	IDPayloadFormatIndicator               = 0  // (M) Payload Format Indicator
	IDPointOfInitiationMethod              = 1  // (O) Point of Initiation Method
	IDMerchantAccountInformationRangeStart = 2  // (M) Merchant Account Information
	IDMerchantAccountInformationRangeEnd   = 51 // (M) Merchant Account Information
	IDMerchantCategoryCode                 = 52 // (M) Merchant Category Code
	IDTransactionCurrency                  = 53 // (M) Transaction Currency
	IDTransactionAmount                    = 54 // (C) Transaction Amount
	IDTipOrConvenienceIndicator            = 55 // (O) Tip or Convenience Indicator
	IDValueOfConvenienceFeeFixed           = 56 // (C) Value of Convenience Fee Fixed
	IDValueOfConvenienceFeePercentage      = 57 // (C) Value of Convenience Fee Percentage
	IDCountryCode                          = 58 // (M) Country Code
	IDMerchantName                         = 59 // (M) Merchant Name
	IDMerchantCity                         = 60 // (M) Merchant City
	IDPostalCode                           = 61 // (O) Postal Code
	IDAdditionalDataFieldTemplate          = 62 // (O) Additional Data Field Template
	IDMerchantInformationLanguageTemplate  = 64 // (O) Merchant Information— Language Template
	IDRFUForEMVCoRangeStart                = 65 // (O) RFU for EMVCo
	IDRFUForEMVCoRangeEnd                  = 79 // (O) RFU for EMVCo
	IDUnreservedTemplatesRangeStart        = 80 // (O) Unreserved Templates
	IDUnreservedTemplatesRangeEnd          = 99 // (O) Unreserved Templates
	IDCRC                                  = 63 // (M) CRC
)

// Data Objects for Additional Data Field Template (ID "62")
const (
	AdditionalIDBillNumber                               = 1
	AdditionalIDMobileNumber                             = 2
	AdditionalIDStoreLabel                               = 3
	AdditionalIDLoyaltyNumber                            = 4
	AdditionalIDReferenceLabel                           = 5
	AdditionalIDCustomerLabel                            = 6
	AdditionalIDTerminalLabel                            = 7
	AdditionalIDPurposeTransaction                       = 8
	AdditionalIDAdditionalConsumerDataRequest            = 9
	AdditionalIDRFUforEMVCoRangeStart                    = 10
	AdditionalIDRFUforEMVCoRangeEnd                      = 49
	AdditionalIDPaymentSystemSpecificTemplatesRangeStart = 50
	AdditionalIDPaymentSystemSpecificTemplatesRangeEnd   = 99
)

// Data Objects for Merchant Information—Language Template (ID "64")
const (
	MerchantInformationIDLanguagePreference    = 0
	MerchantInformationIDMerchantName          = 1
	MerchantInformationIDMerchantCity          = 2
	MerchantInformationIDRFUforEMVCoRangeStart = 3
	MerchantInformationIDRFUforEMVCoRangeEnd   = 99
)

// EMVQR ...
type EMVQR struct {
	PayloadFormatIndicator              string
	PointOfInitiationMethod             string
	MerchantAccountInformationTemplates []*MerchantAccountInformationTemplate // (M) Merchant Account Information
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
	AdditionalDataFieldTemplate         *AdditionalDataFieldTemplate         // Tag: 62
	CRC                                 string                               // Tag: 63
	MerchantInformationLanguageTemplate *MerchantInformationLanguageTemplate // Tag: 64
	RFUForEMVCoTemplates                []*RFUForEMVCoTemplate               // Tag: 65-79 RFU for EMVCo
	UnreservedTemplates                 []*UnreservedTemplate                // Tag: 80-99 Unreserved Templates
}

// MerchantAccountInformationTemplate ...
type MerchantAccountInformationTemplate struct {
	ID    int64
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
	RFUForEMVCoTemplates           []*AdditionalRFUForEMVCoTemplate
	PaymentSystemSpecificTemplates []*AdditionalPaymentSystemSpecificTemplate
}

// AdditionalRFUForEMVCoTemplate ...
type AdditionalRFUForEMVCoTemplate struct {
	ID    int64
	Value string
}

// AdditionalPaymentSystemSpecificTemplate ...
type AdditionalPaymentSystemSpecificTemplate struct {
	ID    int64
	Value string
}

// RFUforEMVCo ...
type RFUForEMVCoTemplate struct {
	ID    int64
	Value string
}

// UnreservedTemplate ...
type UnreservedTemplate struct {
	ID    int64
	Value string
}

// MerchantInformationLanguageTemplate ...
type MerchantInformationLanguageTemplate struct {
	LanguagePreference   string
	MerchantName         string
	MerchantCity         string
	RFUForEMVCoTemplates []*MerchantInformationRFUForEMVCoTemplate
}

// RFUforEMVCo ...
type MerchantInformationRFUForEMVCoTemplate struct {
	ID    int64
	Value string
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
	if len(c.MerchantAccountInformationTemplates) > 0 {
		for _, t := range c.MerchantAccountInformationTemplates {
			s += format(t.ID, t.Value)
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
	if c.AdditionalDataFieldTemplate != nil {
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
		if len(t.RFUForEMVCoTemplates) > 0 {
			for _, t := range t.RFUForEMVCoTemplates {
				template += format(t.ID, t.Value)
			}
		} // 10-49
		if len(t.PaymentSystemSpecificTemplates) > 0 {
			for _, t := range t.PaymentSystemSpecificTemplates {
				template += format(t.ID, t.Value)
			}
		} // 50-99
		s += format(IDAdditionalDataFieldTemplate, template)
	}
	if c.MerchantInformationLanguageTemplate != nil {
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
		if len(t.RFUForEMVCoTemplates) > 0 {
			for _, t := range t.RFUForEMVCoTemplates {
				template += format(t.ID, t.Value)
			}
		}
		s += format(IDMerchantInformationLanguageTemplate, template)
	}
	if len(c.RFUForEMVCoTemplates) > 0 {
		for _, t := range c.RFUForEMVCoTemplates {
			s += format(t.ID, t.Value)
		}
	}
	if len(c.UnreservedTemplates) > 0 {
		for _, t := range c.UnreservedTemplates {
			s += format(t.ID, t.Value)
		}
	}
	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	crc := crc16.Checksum([]byte(fmt.Sprintf("%s%02d04", s, IDCRC)), table)
	crcStr := formatCrc(crc)
	s += format(IDCRC, crcStr)
	return s, nil
}

func format(id int64, value string) string {
	valueLength := utf8.RuneCountInString(value)
	return fmt.Sprintf("%02d%02d%s", id, valueLength, value)
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
func ParsePayload(payload string) (*EMVQR, error) {
	p, err := NewParser(payload)
	if err != nil {
		return nil, err
	}
	emvQR := new(EMVQR)
	for p.Next() {
		id := p.ID()
		switch id {
		case IDPayloadFormatIndicator:
			emvQR.PayloadFormatIndicator = p.Value()
		case IDPointOfInitiationMethod:
			emvQR.PointOfInitiationMethod = p.Value()
		case IDMerchantCategoryCode:
			emvQR.MerchantCategoryCode = p.Value()
		case IDTransactionCurrency:
			emvQR.TransactionCurrency = p.Value()
		case IDTransactionAmount:
			emvQR.TransactionAmount = p.Value()
		case IDTipOrConvenienceIndicator:
			emvQR.TipOrConvenienceIndicator = p.Value()
		case IDValueOfConvenienceFeeFixed:
			emvQR.ValueOfConvenienceFeeFixed = p.Value()
		case IDValueOfConvenienceFeePercentage:
			emvQR.ValueOfConvenienceFeePercentage = p.Value()
		case IDCountryCode:
			emvQR.CountryCode = p.Value()
		case IDMerchantName:
			emvQR.MerchantName = p.Value()
		case IDMerchantCity:
			emvQR.MerchantCity = p.Value()
		case IDPostalCode:
			emvQR.PostalCode = p.Value()
		case IDCRC:
			emvQR.CRC = p.Value()
		}
		if id >= IDMerchantAccountInformationRangeStart && id <= IDMerchantAccountInformationRangeEnd {
			merchantAccountInformationTemplate, err := parseMerchantAccountInformationTemplate(id, p.Value())
			if err != nil {
				return nil, err
			}
			emvQR.MerchantAccountInformationTemplates = append(emvQR.MerchantAccountInformationTemplates, merchantAccountInformationTemplate)
		}
		if id == IDAdditionalDataFieldTemplate {
			additionalDataFieldTemplate, err := parseAdditionalDataFieldTemplate(p.Value())
			if err != nil {
				return nil, err
			}
			emvQR.AdditionalDataFieldTemplate = additionalDataFieldTemplate
		}
		if id == IDMerchantInformationLanguageTemplate {
			merchantInformationLanguageTemplate, err := parseMerchantInformationLanguageTemplate(p.Value())
			if err != nil {
				return nil, err
			}
			emvQR.MerchantInformationLanguageTemplate = merchantInformationLanguageTemplate
		}
		if id >= IDRFUForEMVCoRangeStart && id <= IDRFUForEMVCoRangeEnd {
			rfuForEMVCoTemplate, err := parseRFUForEMVCoTemplate(id, p.Value())
			if err != nil {
				return nil, err
			}
			emvQR.RFUForEMVCoTemplates = append(emvQR.RFUForEMVCoTemplates, rfuForEMVCoTemplate)
		}
		if id >= IDUnreservedTemplatesRangeStart && id <= IDUnreservedTemplatesRangeEnd {
			unreservedTemplate, err := parseUnreservedTemplate(id, p.Value())
			if err != nil {
				return nil, err
			}
			emvQR.UnreservedTemplates = append(emvQR.UnreservedTemplates, unreservedTemplate)
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return emvQR, nil
}

func parseMerchantAccountInformationTemplate(id int64, value string) (*MerchantAccountInformationTemplate, error) {
	return &MerchantAccountInformationTemplate{
		ID:    id,
		Value: value,
	}, nil
}

func parseAdditionalDataFieldTemplate(value string) (*AdditionalDataFieldTemplate, error) {
	p, err := NewParser(value)
	if err != nil {
		return nil, err
	}
	additionalDataFieldTemplate := new(AdditionalDataFieldTemplate)
	for p.Next() {
		switch p.ID() {
		case AdditionalIDBillNumber:
			additionalDataFieldTemplate.BillNumber = p.Value()
		case AdditionalIDMobileNumber:
			additionalDataFieldTemplate.MobileNumber = p.Value()
		case AdditionalIDStoreLabel:
			additionalDataFieldTemplate.StoreLabel = p.Value()
		case AdditionalIDLoyaltyNumber:
			additionalDataFieldTemplate.LoyaltyNumber = p.Value()
		case AdditionalIDReferenceLabel:
			additionalDataFieldTemplate.ReferenceLabel = p.Value()
		case AdditionalIDCustomerLabel:
			additionalDataFieldTemplate.CustomerLabel = p.Value()
		case AdditionalIDTerminalLabel:
			additionalDataFieldTemplate.TerminalLabel = p.Value()
		case AdditionalIDPurposeTransaction:
			additionalDataFieldTemplate.PurposeTransaction = p.Value()
		case AdditionalIDAdditionalConsumerDataRequest:
			additionalDataFieldTemplate.AdditionalConsumerDataRequest = p.Value()
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return additionalDataFieldTemplate, nil
}

func parseMerchantInformationLanguageTemplate(value string) (*MerchantInformationLanguageTemplate, error) {
	p, err := NewParser(value)
	if err != nil {
		return nil, err
	}
	merchantInformationLanguageTemplate := new(MerchantInformationLanguageTemplate)
	for p.Next() {
		switch p.ID() {
		case MerchantInformationIDLanguagePreference:
			merchantInformationLanguageTemplate.LanguagePreference = p.Value()
		case MerchantInformationIDMerchantName:
			merchantInformationLanguageTemplate.MerchantName = p.Value()
		case MerchantInformationIDMerchantCity:
			merchantInformationLanguageTemplate.MerchantCity = p.Value()
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return merchantInformationLanguageTemplate, nil
}

func parseRFUForEMVCoTemplate(id int64, value string) (*RFUForEMVCoTemplate, error) {
	return &RFUForEMVCoTemplate{
		ID:    id,
		Value: value,
	}, nil
}

func parseUnreservedTemplate(id int64, value string) (*UnreservedTemplate, error) {
	return &UnreservedTemplate{
		ID:    id,
		Value: value,
	}, nil
}
