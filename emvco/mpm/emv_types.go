package mpm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/dongri/emvco-qrcode/crc16"
)

// ID ...
type ID string

// String ...
func (id ID) String() string {
	return string(id)
}

// ParseInt ...
func (id ID) ParseInt() (int64, error) {
	return strconv.ParseInt(id.String(), 10, 64)
}

// Equal ...
func (id ID) Equal(val ID) bool {
	return id == val
}

// Between ...
func (id ID) Between(start ID, end ID) (bool, error) {
	idNum, err := id.ParseInt()
	if err != nil {
		return false, err
	}
	startNum, err := start.ParseInt()
	if err != nil {
		return false, err
	}
	endNum, err := end.ParseInt()
	if err != nil {
		return false, err
	}
	return idNum >= startNum && idNum <= endNum, nil
}

const (
	IDPayloadFormatIndicator               ID = "00" // (M) Payload Format Indicator
	IDPointOfInitiationMethod              ID = "01" // (O) Point of Initiation Method
	IDMerchantAccountInformationRangeStart ID = "02" // (M) 2-51 Merchant Account Information
	IDMerchantAccountInformationRangeEnd   ID = "51" // (M) 2-51 Merchant Account Information
	IDMerchantCategoryCode                 ID = "52" // (M) Merchant Category Code
	IDTransactionCurrency                  ID = "53" // (M) Transaction Currency
	IDTransactionAmount                    ID = "54" // (C) Transaction Amount
	IDTipOrConvenienceIndicator            ID = "55" // (O) Tip or Convenience Indicator
	IDValueOfConvenienceFeeFixed           ID = "56" // (C) Value of Convenience Fee Fixed
	IDValueOfConvenienceFeePercentage      ID = "57" // (C) Value of Convenience Fee Percentage
	IDCountryCode                          ID = "58" // (M) Country Code
	IDMerchantName                         ID = "59" // (M) Merchant Name
	IDMerchantCity                         ID = "60" // (M) Merchant City
	IDPostalCode                           ID = "61" // (O) Postal Code
	IDAdditionalDataFieldTemplate          ID = "62" // (O) Additional Data Field Template
	IDCRC                                  ID = "63" // (M) CRC
	IDMerchantInformationLanguageTemplate  ID = "64" // (O) Merchant Information— Language Template
	IDRFUForEMVCoRangeStart                ID = "65" // (O) 65-79 RFU for EMVCo
	IDRFUForEMVCoRangeEnd                  ID = "79" // (O) 65-79 RFU for EMVCo
	IDUnreservedTemplatesRangeStart        ID = "80" // (O) 80-99 Unreserved Templates
	IDUnreservedTemplatesRangeEnd          ID = "99" // (O) 80-99 Unreserved Templates
)

// EMVQR ...
type EMVQR struct {
	PayloadFormatIndicator              string
	PointOfInitiationMethod             PointOfInitiationMethod
	MerchantAccountInformation          map[ID]*MerchantAccountInformation
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
	AdditionalDataFieldTemplate         *AdditionalDataFieldTemplate
	CRC                                 string
	MerchantInformationLanguageTemplate *MerchantInformationLanguageTemplate
	RFUforEMVCo                         map[ID]*RFUforEMVCo
	UnreservedTemplates                 map[ID]*UnreservedTemplate
}

// ParseEMVQR ...
func ParseEMVQR(payload string) (*EMVQR, error) {
	p := NewParser(payload)
	emvqr := &EMVQR{}
	for p.Next() {
		id := p.ID()
		value := p.Value()
		switch id {
		case IDPayloadFormatIndicator:
			emvqr.PayloadFormatIndicator = value
		case IDPointOfInitiationMethod:
			emvqr.PointOfInitiationMethod = ParsePointOfInitiationMethod(value)
		case IDMerchantCategoryCode:
			emvqr.MerchantCategoryCode = value
		case IDTransactionCurrency:
			emvqr.TransactionCurrency = value
		case IDTransactionAmount:
			emvqr.TransactionAmount = value
		case IDTipOrConvenienceIndicator:
			emvqr.TipOrConvenienceIndicator = value
		case IDValueOfConvenienceFeeFixed:
			emvqr.ValueOfConvenienceFeeFixed = value
		case IDValueOfConvenienceFeePercentage:
			emvqr.ValueOfConvenienceFeePercentage = value
		case IDCountryCode:
			emvqr.CountryCode = value
		case IDMerchantName:
			emvqr.MerchantName = value
		case IDMerchantCity:
			emvqr.MerchantCity = value
		case IDPostalCode:
			emvqr.PostalCode = value
		case IDAdditionalDataFieldTemplate:
			adft, err := ParseAdditionalDataFieldTemplate(value)
			if err != nil {
				return nil, err
			}
			emvqr.AdditionalDataFieldTemplate = adft
		case IDCRC:
			emvqr.CRC = value
		case IDMerchantInformationLanguageTemplate:
			t, err := ParseMerchantInformationLanguageTemplate(value)
			if err != nil {
				return nil, err
			}
			emvqr.MerchantInformationLanguageTemplate = t
		default:
			var (
				within bool
				err    error
			)
			// Merchant Account Information
			within, err = id.Between(IDMerchantAccountInformationRangeStart, IDMerchantAccountInformationRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				if emvqr.MerchantAccountInformation == nil {
					emvqr.MerchantAccountInformation = make(map[ID]*MerchantAccountInformation)
				}
				emvqr.MerchantAccountInformation[id] = ParseMerchantAccountInformation(value)
				continue
			}
			// RFUforEMVCo
			within, err = id.Between(IDRFUForEMVCoRangeStart, IDRFUForEMVCoRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				if emvqr.RFUforEMVCo == nil {
					emvqr.RFUforEMVCo = make(map[ID]*RFUforEMVCo)
				}
				emvqr.RFUforEMVCo[id] = ParseRFUforEMVCo(value)
				continue
			}
			// Unreserved Tempaltes
			within, err = id.Between(IDUnreservedTemplatesRangeStart, IDUnreservedTemplatesRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				if emvqr.UnreservedTemplates == nil {
					emvqr.UnreservedTemplates = make(map[ID]*UnreservedTemplate)
				}
				emvqr.UnreservedTemplates[id] = ParseUnreservedTemplate(value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return emvqr, nil
}

// GeneratePayload
func (c *EMVQR) GeneratePayload() (string, error) {
	s := ""
	s += format(IDPayloadFormatIndicator, c.PayloadFormatIndicator)
	if c.PointOfInitiationMethod != "" {
		s += c.PointOfInitiationMethod.GeneratePayload()
	}
	if len(c.MerchantAccountInformation) > 0 {
		for id, t := range c.MerchantAccountInformation {
			s += format(id, t.GeneratePayload())
		}
	}
	if c.MerchantCategoryCode != "" {
		s += format(IDMerchantCategoryCode, c.MerchantCategoryCode)
	}
	if c.TransactionCurrency != "" {
		s += format(IDTransactionCurrency, c.TransactionCurrency)
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
	if c.AdditionalDataFieldTemplate != nil {
		s += format(IDAdditionalDataFieldTemplate, c.AdditionalDataFieldTemplate.GeneratePayload())
	}
	if c.MerchantInformationLanguageTemplate != nil {
		s += format(IDMerchantInformationLanguageTemplate, c.MerchantInformationLanguageTemplate.GeneratePayload())
	}
	if len(c.RFUforEMVCo) > 0 {
		for id, t := range c.RFUforEMVCo {
			s += format(id, t.GeneratePayload())
		}
	}
	if len(c.UnreservedTemplates) > 0 {
		for id, t := range c.UnreservedTemplates {
			s += format(id, t.GeneratePayload())
		}
	}
	s += formatCrc(s)
	return s, nil
}

func (c *EMVQR) Validate() error {
	// check mandatory
	if c.PayloadFormatIndicator == "" {
		return errors.New("PayloadFormatIndicator is mandatory")
	}
	if len(c.MerchantAccountInformation) <= 0 {
		return errors.New("MerchantAccountInformation is mandatory")
	}
	if c.MerchantCategoryCode == "" {
		return errors.New("MerchantCategoryCode is mandatory")
	}
	if c.TransactionCurrency == "" {
		return errors.New("TransactionCurrency is mandatory")
	}
	if c.CountryCode == "" {
		return errors.New("CountryCode is mandatory")
	}
	if c.MerchantName == "" {
		return errors.New("MerchantName is mandatory")
	}
	if c.MerchantCity == "" {
		return errors.New("MerchantCity is mandatory")
	}
	// check validate
	if c.PointOfInitiationMethod != "" {
		if err := c.PointOfInitiationMethod.Validate(); err != nil {
			return err
		}
	}
	for _, t := range c.MerchantAccountInformation {
		if err := t.Validate(); err != nil {
			return err
		}
	}
	if c.AdditionalDataFieldTemplate != nil {
		if err := c.AdditionalDataFieldTemplate.Validate(); err != nil {
			return err
		}
	}
	if c.MerchantInformationLanguageTemplate != nil {
		if err := c.MerchantInformationLanguageTemplate.Validate(); err != nil {
			return err
		}
	}
	for _, t := range c.RFUforEMVCo {
		if err := t.Validate(); err != nil {
			return err
		}
	}
	for _, t := range c.UnreservedTemplates {
		if err := t.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// PointOfInitiationMethod ...
type PointOfInitiationMethod string

const (
	// PointOfInitiationMethodStatic ...
	PointOfInitiationMethodStatic PointOfInitiationMethod = "11"
	// PointOfInitiationMethodDynamic ...
	PointOfInitiationMethodDynamic PointOfInitiationMethod = "12"
)

// ParsePointOfInitiationMethod ...
func ParsePointOfInitiationMethod(payload string) PointOfInitiationMethod {
	return PointOfInitiationMethod(payload)
}

// GeneratePayload ...
func (c PointOfInitiationMethod) GeneratePayload() string {
	return format(IDPointOfInitiationMethod, string(c))
}

// IsStaticMethod ...
func (c PointOfInitiationMethod) IsStaticMethod() bool {
	return c == PointOfInitiationMethodStatic
}

// IsDynamicMethod ...
func (c PointOfInitiationMethod) IsDynamicMethod() bool {
	return c == PointOfInitiationMethodDynamic
}

// Validate
func (c PointOfInitiationMethod) Validate() error {
	if !c.IsStaticMethod() && !c.IsDynamicMethod() {
		return fmt.Errorf("PointOfInitiationMethod should be \"11\" or \"12\", PointOfInitiationMethod: %s", c)
	}
	return nil
}

// RFUforEMVCo ...
type RFUforEMVCo struct {
	Value string
}

// ParseRFUforEMVCo ...
func ParseRFUforEMVCo(payload string) *RFUforEMVCo {
	return &RFUforEMVCo{
		Value: payload,
	}
}

// GeneratePayload ...
func (c *RFUforEMVCo) GeneratePayload() string {
	return c.Value
}

// Validate ...
func (c *RFUforEMVCo) Validate() error {
	return nil
}

// UnreservedTemplate ...
type UnreservedTemplate struct {
	Value string
}

// ParseUnreservedTemplate ...
func ParseUnreservedTemplate(payload string) *UnreservedTemplate {
	return &UnreservedTemplate{
		Value: payload,
	}
}

// GeneratePayload ...
func (c *UnreservedTemplate) GeneratePayload() string {
	return c.Value
}

// Validate ...
func (c *UnreservedTemplate) Validate() error {
	return nil
}

// MerchantAccountInformation ...
type MerchantAccountInformation struct {
	Value string
}

// ParseMerchantAccountInformation ...
func ParseMerchantAccountInformation(payload string) *MerchantAccountInformation {
	return &MerchantAccountInformation{
		Value: payload,
	}
}

// GeneratePayload ...
func (c *MerchantAccountInformation) GeneratePayload() string {
	return c.Value
}

// GeneratePayload ...
func (c *MerchantAccountInformation) Validate() error {
	return nil
}

// Data Objects for Additional Data Field Template (ID "62")
const (
	AdditionalIDBillNumber                               ID = "01" // (O) Bill Number
	AdditionalIDMobileNumber                             ID = "02" // (O) Mobile Number
	AdditionalIDStoreLabel                               ID = "03" // (O) Store Label
	AdditionalIDLoyaltyNumber                            ID = "04" // (O) Loyalty Number
	AdditionalIDReferenceLabel                           ID = "05" // (O) Reference Label
	AdditionalIDCustomerLabel                            ID = "06" // (O) Customer Label
	AdditionalIDTerminalLabel                            ID = "07" // (O) Terminal Label
	AdditionalIDPurposeTransaction                       ID = "08" // (O) Purpose Transaction
	AdditionalIDAdditionalConsumerDataRequest            ID = "09" // (O) Additional Consumer Data Request
	AdditionalIDRFUforEMVCoRangeStart                    ID = "10" // (O) RFU for EMVCo
	AdditionalIDRFUforEMVCoRangeEnd                      ID = "49" // (O) RFU for EMVCo
	AdditionalIDPaymentSystemSpecificTemplatesRangeStart ID = "50" // (O) Payment System Specific Templates
	AdditionalIDPaymentSystemSpecificTemplatesRangeEnd   ID = "99" // (O) Payment System Specific Templates
)

// AdditionalDataFieldTemplate ...
type AdditionalDataFieldTemplate struct {
	BillNumber                    string
	MobileNumber                  string
	StoreLabel                    string
	LoyaltyNumber                 string
	ReferenceLabel                string
	CustomerLabel                 string
	TerminalLabel                 string
	PurposeTransaction            string
	AdditionalConsumerDataRequest string
	RFUforEMVCo                   map[ID]*RFUforEMVCo
	PaymentSystemSpecific         map[ID]*PaymentSystemSpecific
}

func ParseAdditionalDataFieldTemplate(payload string) (*AdditionalDataFieldTemplate, error) {
	p := NewParser(payload)
	additionalDataFieldTemplate := &AdditionalDataFieldTemplate{}
	for p.Next() {
		id := p.ID()
		value := p.Value()
		switch id {
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
		default:
			var (
				within bool
				err    error
			)
			// Payment System Specific
			within, err = id.Between(AdditionalIDPaymentSystemSpecificTemplatesRangeStart, AdditionalIDPaymentSystemSpecificTemplatesRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				if additionalDataFieldTemplate.PaymentSystemSpecific == nil {
					additionalDataFieldTemplate.PaymentSystemSpecific = make(map[ID]*PaymentSystemSpecific)
				}
				additionalDataFieldTemplate.PaymentSystemSpecific[id] = ParsePaymentSystemSpecific(value)
				continue
			}
			// RFU for EMVCo
			within, err = id.Between(AdditionalIDRFUforEMVCoRangeStart, AdditionalIDRFUforEMVCoRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				if additionalDataFieldTemplate.RFUforEMVCo == nil {
					additionalDataFieldTemplate.RFUforEMVCo = make(map[ID]*RFUforEMVCo)
				}
				additionalDataFieldTemplate.RFUforEMVCo[id] = ParseRFUforEMVCo(value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return additionalDataFieldTemplate, nil
}

// GeneratePayload ...
func (c *AdditionalDataFieldTemplate) GeneratePayload() string {
	s := ""
	if c.BillNumber != "" {
		s += format(AdditionalIDBillNumber, c.BillNumber)
	}
	if c.MobileNumber != "" {
		s += format(AdditionalIDMobileNumber, c.MobileNumber)
	}
	if c.StoreLabel != "" {
		s += format(AdditionalIDStoreLabel, c.StoreLabel)
	}
	if c.LoyaltyNumber != "" {
		s += format(AdditionalIDLoyaltyNumber, c.LoyaltyNumber)
	}
	if c.ReferenceLabel != "" {
		s += format(AdditionalIDReferenceLabel, c.ReferenceLabel)
	}
	if c.CustomerLabel != "" {
		s += format(AdditionalIDCustomerLabel, c.CustomerLabel)
	}
	if c.TerminalLabel != "" {
		s += format(AdditionalIDTerminalLabel, c.TerminalLabel)
	}
	if c.PurposeTransaction != "" {
		s += format(AdditionalIDPurposeTransaction, c.PurposeTransaction)
	}
	if c.AdditionalConsumerDataRequest != "" {
		s += format(AdditionalIDAdditionalConsumerDataRequest, c.AdditionalConsumerDataRequest)
	}
	if len(c.RFUforEMVCo) > 0 {
		for k, t := range c.RFUforEMVCo {
			s += format(ID(k), t.GeneratePayload())
		}
	}
	if len(c.PaymentSystemSpecific) > 0 {
		for k, t := range c.PaymentSystemSpecific {
			s += format(ID(k), t.GeneratePayload())
		}
	}
	return s
}

// Validate ...
func (c *AdditionalDataFieldTemplate) Validate() error {
	return nil
}

// PaymentSystemSpecific ...
type PaymentSystemSpecific struct {
	Value string
}

// ParsePaymentSystemSpecific ...
func ParsePaymentSystemSpecific(value string) *PaymentSystemSpecific {
	return &PaymentSystemSpecific{Value: value}
}

// GeneratePayload ...
func (v *PaymentSystemSpecific) GeneratePayload() string {
	return v.Value
}

// Data Objects for Merchant Information—Language Template (ID "64")
const (
	MerchantInformationIDLanguagePreference    = "00" // (M) Language Preference
	MerchantInformationIDMerchantName          = "01" // (M) Merchant Name
	MerchantInformationIDMerchantCity          = "02" // (O) Merchant City
	MerchantInformationIDRFUforEMVCoRangeStart = "03" // (O) 03-99 RFU for EMVCo
	MerchantInformationIDRFUforEMVCoRangeEnd   = "99" // (O) 03-99 RFU for EMVCo
)

// MerchantInformationLanguageTemplate ...
type MerchantInformationLanguageTemplate struct {
	LanguagePreference string
	MerchantName       string
	MerchantCity       string
	RFUForEMVCo        map[ID]*RFUforEMVCo
}

// ParseMerchantInformationLanguageTemplate ...
func ParseMerchantInformationLanguageTemplate(value string) (*MerchantInformationLanguageTemplate, error) {
	p := NewParser(value)
	merchantInformationLanguageTemplate := &MerchantInformationLanguageTemplate{}
	for p.Next() {
		id := p.ID()
		value := p.Value()
		switch id {
		case MerchantInformationIDLanguagePreference:
			merchantInformationLanguageTemplate.LanguagePreference = value
		case MerchantInformationIDMerchantName:
			merchantInformationLanguageTemplate.MerchantName = value
		case MerchantInformationIDMerchantCity:
			merchantInformationLanguageTemplate.MerchantCity = value
		default:
			var (
				within bool
				err    error
			)
			// RFU for EMVCo
			within, err = id.Between(MerchantInformationIDRFUforEMVCoRangeStart, MerchantInformationIDRFUforEMVCoRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				if merchantInformationLanguageTemplate.RFUForEMVCo == nil {
					merchantInformationLanguageTemplate.RFUForEMVCo = make(map[ID]*RFUforEMVCo)
				}
				merchantInformationLanguageTemplate.RFUForEMVCo[id] = ParseRFUforEMVCo(value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return merchantInformationLanguageTemplate, nil
}

// GeneratePayload ...
func (c *MerchantInformationLanguageTemplate) GeneratePayload() string {
	s := ""
	if c.LanguagePreference != "" {
		s += format(MerchantInformationIDLanguagePreference, c.LanguagePreference)
	}
	if c.MerchantName != "" {
		s += format(MerchantInformationIDMerchantName, c.MerchantName)
	}
	if c.MerchantCity != "" {
		s += format(MerchantInformationIDMerchantCity, c.MerchantCity)
	}
	if len(c.RFUForEMVCo) > 0 {
		for id, t := range c.RFUForEMVCo {
			s += format(id, t.GeneratePayload())
		}
	}
	return s
}

// Validate ...
func (c *MerchantInformationLanguageTemplate) Validate() error {
	// check mandatory
	if c.LanguagePreference == "" {
		return errors.New("LanguagePreference is mandatory")
	}
	if c.MerchantName == "" {
		return errors.New("MerchantName is mandatory")
	}
	// check validate
	for _, t := range c.RFUForEMVCo {
		if err := t.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func format(id ID, value string) string {
	valueLength := utf8.RuneCountInString(value)
	return fmt.Sprintf("%s%02d%s", id.String(), valueLength, value)
}

func formatCrc(value string) string {
	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	crcValue := crc16.Checksum([]byte(value+IDCRC.String()+"04"), table)
	crcValueString := strconv.FormatUint(uint64(crcValue), 16)
	s := "0000" + strings.ToUpper(crcValueString)
	return format(IDCRC, s[len(s)-4:])
}
