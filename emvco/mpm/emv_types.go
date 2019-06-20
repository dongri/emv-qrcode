package mpm

import (
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
	MerchantAccountInformation          map[ID]*MerchantAccountInformationTemplate
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

func parseEMVQR(payload string) (*EMVQR, error) {
	p := NewParser(payload)
	emvqr := &EMVQR{}
	for p.Next() {
		id := p.ID()
		value := p.Value()
		switch id {
		case IDPayloadFormatIndicator:
			emvqr.PayloadFormatIndicator = value
		case IDPointOfInitiationMethod:
			emvqr.PointOfInitiationMethod = PointOfInitiationMethod(value)
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
			adft, err := parseAdditionalDataFieldTemplate(value)
			if err != nil {
				return nil, err
			}
			emvqr.AdditionalDataFieldTemplate = adft
		case IDCRC:
			emvqr.CRC = value
		case IDMerchantInformationLanguageTemplate:
			t, err := parseMerchantInformationLanguageTemplate(value)
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
					emvqr.MerchantAccountInformation = make(map[ID]*MerchantAccountInformationTemplate)
				}
				t, err := parseMerchantAccountInformationTemplate(value)
				if err != nil {
					return nil, err
				}
				emvqr.MerchantAccountInformation[id] = t
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
				emvqr.RFUforEMVCo[id] = parseRFUforEMVCo(value)
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
				emvqr.UnreservedTemplates[id] = parseUnreservedTemplate(value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return emvqr, nil
}

// Stringify ...
func (c *EMVQR) Stringify() (string, error) {
	s := ""
	s += format(IDPayloadFormatIndicator, c.PayloadFormatIndicator)
	if c.PointOfInitiationMethod != "" {
		s += c.PointOfInitiationMethod.Stringify()
	}
	if len(c.MerchantAccountInformation) > 0 {
		for id, t := range c.MerchantAccountInformation {
			s += format(id, t.Stringify())
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
		s += format(IDAdditionalDataFieldTemplate, c.AdditionalDataFieldTemplate.Stringify())
	}
	if c.MerchantInformationLanguageTemplate != nil {
		s += format(IDMerchantInformationLanguageTemplate, c.MerchantInformationLanguageTemplate.Stringify())
	}
	if len(c.RFUforEMVCo) > 0 {
		for id, t := range c.RFUforEMVCo {
			s += format(id, t.Stringify())
		}
	}
	if len(c.UnreservedTemplates) > 0 {
		for id, t := range c.UnreservedTemplates {
			s += format(id, t.Stringify())
		}
	}
	s += formatCrc(s)
	return s, nil
}

func (c *EMVQR) Validate() error {
	return nil
}

type PointOfInitiationMethod string

const (
	PointOfInitiationMethodStatic  PointOfInitiationMethod = "11"
	PointOfInitiationMethodDynamic PointOfInitiationMethod = "12"
)

// IsStaticMethod ...
func (m PointOfInitiationMethod) IsStaticMethod() bool {
	return m == PointOfInitiationMethodStatic
}

// IsDynamicMethod ...
func (m PointOfInitiationMethod) IsDynamicMethod() bool {
	return m == PointOfInitiationMethodDynamic
}

// Stringify ...
func (v PointOfInitiationMethod) Stringify() string {
	return format(IDPointOfInitiationMethod, string(v))
}

// RFUforEMVCo ...
type RFUforEMVCo struct {
	Value string
}

func parseRFUforEMVCo(value string) *RFUforEMVCo {
	return &RFUforEMVCo{
		Value: value,
	}
}

// Stringify ...
func (v RFUforEMVCo) Stringify() string {
	return v.Value
}

// UnreservedTemplate ...
type UnreservedTemplate struct {
	Value string
}

func parseUnreservedTemplate(value string) *UnreservedTemplate {
	return &UnreservedTemplate{
		Value: value,
	}
}

// Stringify ...
func (v UnreservedTemplate) Stringify() string {
	return v.Value
}

// Merchant Account Information (IDs "02" to "51")
const (
	MerchantAccountIDGloballyUniqueIdentifier         = "00" // (M) Globally Unique Identifier
	MerchantAccountIDPaymentNetworkSpecificRangeStart = "01" // (O) 01-99 Payment network specific
	MerchantAccountIDPaymentNetworkSpecificRangeEnd   = "99" // (O) 01-99 Payment network specific
)

// MerchantAccountInformation ...
type MerchantAccountInformationTemplate struct {
	GloballyUniqueIdentifier string
	PaymentNetworkSpecific   map[ID]*PaymentNetworkSpecific
}

func parseMerchantAccountInformationTemplate(payload string) (*MerchantAccountInformationTemplate, error) {
	p := NewParser(payload)
	merchantAccountInformation := &MerchantAccountInformationTemplate{}
	for p.Next() {
		id := p.ID()
		value := p.Value()
		switch id {
		case MerchantAccountIDGloballyUniqueIdentifier:
			merchantAccountInformation.GloballyUniqueIdentifier = value
		default:
			var (
				within bool
				err    error
			)
			// Payment Network Specific
			within, err = id.Between(MerchantAccountIDPaymentNetworkSpecificRangeStart, MerchantAccountIDPaymentNetworkSpecificRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				if merchantAccountInformation.PaymentNetworkSpecific == nil {
					merchantAccountInformation.PaymentNetworkSpecific = make(map[ID]*PaymentNetworkSpecific)
				}
				merchantAccountInformation.PaymentNetworkSpecific[id] = parsePaymentNetworkSpecific(value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return merchantAccountInformation, nil
}

func (c MerchantAccountInformationTemplate) Stringify() string {
	s := ""
	s += format(MerchantAccountIDGloballyUniqueIdentifier, c.GloballyUniqueIdentifier)
	if len(c.PaymentNetworkSpecific) > 0 {
		for id, t := range c.PaymentNetworkSpecific {
			s += format(id, t.Stringify())
		}
	}
	return s
}

// PaymentNetworkSpecific ...
type PaymentNetworkSpecific struct {
	Value string
}

func parsePaymentNetworkSpecific(value string) *PaymentNetworkSpecific {
	return &PaymentNetworkSpecific{
		Value: value,
	}
}

func (c PaymentNetworkSpecific) Stringify() string {
	return c.Value
}

// Data Objects for Additional Data Field Template (ID "62")
const (
	AdditionalIDBillNumber                               ID = "01"
	AdditionalIDMobileNumber                             ID = "02"
	AdditionalIDStoreLabel                               ID = "03"
	AdditionalIDLoyaltyNumber                            ID = "04"
	AdditionalIDReferenceLabel                           ID = "05"
	AdditionalIDCustomerLabel                            ID = "06"
	AdditionalIDTerminalLabel                            ID = "07"
	AdditionalIDPurposeTransaction                       ID = "08"
	AdditionalIDAdditionalConsumerDataRequest            ID = "09"
	AdditionalIDRFUforEMVCoRangeStart                    ID = "10"
	AdditionalIDRFUforEMVCoRangeEnd                      ID = "49"
	AdditionalIDPaymentSystemSpecificTemplatesRangeStart ID = "50"
	AdditionalIDPaymentSystemSpecificTemplatesRangeEnd   ID = "99"
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
	PaymentSystemSpecific         map[ID]*PaymentSystemSpecificTemplate
}

func parseAdditionalDataFieldTemplate(payload string) (*AdditionalDataFieldTemplate, error) {
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
					additionalDataFieldTemplate.PaymentSystemSpecific = make(map[ID]*PaymentSystemSpecificTemplate)
				}
				paymentSystemSpecific, err := parsePaymentSystemSpecificTemplate(value)
				if err != nil {
					return nil, err
				}
				additionalDataFieldTemplate.PaymentSystemSpecific[id] = paymentSystemSpecific
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
				additionalDataFieldTemplate.RFUforEMVCo[id] = parseRFUforEMVCo(value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return additionalDataFieldTemplate, nil
}

// Stringify ...
func (c *AdditionalDataFieldTemplate) Stringify() string {
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
			s += format(ID(k), t.Stringify())
		}
	}
	if len(c.PaymentSystemSpecific) > 0 {
		for k, t := range c.PaymentSystemSpecific {
			s += format(ID(k), t.Stringify())
		}
	}
	return s
}

// Additional Payment System Specific
const (
	PaymentSystemIDGloballyUniqueIdentifier        = "00" // (M) Globally Unique Identifier
	PaymentSystemIDPaymentSystemSpecificRangeStart = "01" // (O) 01-99 Payment System specific
	PaymentSystemIDPaymentSystemSpecificRangeEnd   = "99" // (O) 01-99 Payment System specific
)

// PaymentSystemSpecificTemplate ...
type PaymentSystemSpecificTemplate struct {
	GloballyUniqueIdentifier string
	PaymentSystemSpecific    map[ID]*PaymentSystemSpecific
}

func parsePaymentSystemSpecificTemplate(payload string) (*PaymentSystemSpecificTemplate, error) {
	p := NewParser(payload)
	paymentSystemSpecificTemplate := &PaymentSystemSpecificTemplate{}
	for p.Next() {
		id := p.ID()
		value := p.Value()
		switch id {
		case PaymentSystemIDGloballyUniqueIdentifier:
			paymentSystemSpecificTemplate.GloballyUniqueIdentifier = value
		default:
			var (
				within bool
				err    error
			)
			// Payment System Specific
			within, err = id.Between(PaymentSystemIDPaymentSystemSpecificRangeStart, PaymentSystemIDPaymentSystemSpecificRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				if paymentSystemSpecificTemplate.PaymentSystemSpecific == nil {
					paymentSystemSpecificTemplate.PaymentSystemSpecific = make(map[ID]*PaymentSystemSpecific)
				}
				paymentSystemSpecificTemplate.PaymentSystemSpecific[id] = parsePaymentSystemSpecific(value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return paymentSystemSpecificTemplate, nil
}

// Stringify ...
func (c *PaymentSystemSpecificTemplate) Stringify() string {
	s := ""
	if c.GloballyUniqueIdentifier != "" {
		s += format(PaymentSystemIDGloballyUniqueIdentifier, c.GloballyUniqueIdentifier)
	}
	if len(c.PaymentSystemSpecific) > 0 {
		for id, t := range c.PaymentSystemSpecific {
			s += format(id, t.Stringify())
		}
	}
	return s
}

// PaymentSystemSpecific ...
type PaymentSystemSpecific struct {
	Value string
}

func parsePaymentSystemSpecific(value string) *PaymentSystemSpecific {
	return &PaymentSystemSpecific{Value: value}
}

// Stringify ...
func (v *PaymentSystemSpecific) Stringify() string {
	return v.Value
}

// Data Objects for Merchant Information—Language Template (ID "64")
const (
	MerchantInformationIDLanguagePreference    = "00"
	MerchantInformationIDMerchantName          = "01"
	MerchantInformationIDMerchantCity          = "02"
	MerchantInformationIDRFUforEMVCoRangeStart = "03"
	MerchantInformationIDRFUforEMVCoRangeEnd   = "99"
)

// MerchantInformationLanguageTemplate ...
type MerchantInformationLanguageTemplate struct {
	LanguagePreference string
	MerchantName       string
	MerchantCity       string
	RFUForEMVCo        map[ID]*RFUforEMVCo
}

func parseMerchantInformationLanguageTemplate(value string) (*MerchantInformationLanguageTemplate, error) {
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
				merchantInformationLanguageTemplate.RFUForEMVCo[id] = parseRFUforEMVCo(value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return merchantInformationLanguageTemplate, nil
}

func (c *MerchantInformationLanguageTemplate) Stringify() string {
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
			s += format(id, t.Stringify())
		}
	}
	return s
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
