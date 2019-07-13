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

// const ...
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

// Data Objects for Merchant Information—Language Template (ID "64")
const (
	MerchantInformationIDLanguagePreference    ID = "00" // (M) Language Preference
	MerchantInformationIDMerchantName          ID = "01" // (M) Merchant Name
	MerchantInformationIDMerchantCity          ID = "02" // (O) Merchant City
	MerchantInformationIDRFUforEMVCoRangeStart ID = "03" // (O) 03-99 RFU for EMVCo
	MerchantInformationIDRFUforEMVCoRangeEnd   ID = "99" // (O) 03-99 RFU for EMVCo
)

// EMVQR ...
type EMVQR struct {
	PayloadFormatIndicator              TLV
	PointOfInitiationMethod             TLV
	MerchantAccountInformation          []TLV
	MerchantCategoryCode                TLV
	TransactionCurrency                 TLV
	TransactionAmount                   TLV
	TipOrConvenienceIndicator           TLV
	ValueOfConvenienceFeeFixed          TLV
	ValueOfConvenienceFeePercentage     TLV
	CountryCode                         TLV
	MerchantName                        TLV
	MerchantCity                        TLV
	PostalCode                          TLV
	AdditionalDataFieldTemplate         *AdditionalDataFieldTemplate
	CRC                                 TLV
	MerchantInformationLanguageTemplate MerchantInformationLanguageTemplate
	RFUforEMVCo                         []TLV
	UnreservedTemplates                 []TLV
}

// TLV ...
type TLV struct {
	Tag    ID
	Length string
	Value  string
}

func (tlv TLV) String() string {
	if tlv.Value == "" {
		return ""
	}
	return tlv.Tag.String() + tlv.Length + tlv.Value
}

// MerchantAccountInformation ...
type MerchantAccountInformation TLV

// AdditionalDataFieldTemplate ...
type AdditionalDataFieldTemplate struct {
	BillNumber                    TLV
	MobileNumber                  TLV
	StoreLabel                    TLV
	LoyaltyNumber                 TLV
	ReferenceLabel                TLV
	CustomerLabel                 TLV
	TerminalLabel                 TLV
	PurposeTransaction            TLV
	AdditionalConsumerDataRequest TLV
	RFUforEMVCo                   []TLV
	PaymentSystemSpecific         []TLV
}

// MerchantInformationLanguageTemplate ...
type MerchantInformationLanguageTemplate struct {
	LanguagePreference TLV
	MerchantName       TLV
	MerchantCity       TLV
	RFUforEMVCo        []TLV
}

// PointOfInitiationMethod ...
// type PointOfInitiationMethod TLV

// MerchantAccountInformation ...
// type MerchantAccountInformation TLV

// SetPayloadFormatIndicator ...
func (c *EMVQR) SetPayloadFormatIndicator(v string) EMVQR {
	tlv := TLV{
		Tag:    IDPayloadFormatIndicator,
		Length: l(v),
		Value:  v,
	}
	c.PayloadFormatIndicator = tlv
	return *c
}

// SetPointOfInitiationMethod ...
func (c *EMVQR) SetPointOfInitiationMethod(v string) EMVQR {
	tlv := TLV{
		Tag:    IDPointOfInitiationMethod,
		Length: l(v),
		Value:  v,
	}
	c.PointOfInitiationMethod = tlv
	return *c
}

// AddMerchantAccountInformation ...
func (c *EMVQR) AddMerchantAccountInformation(id ID, v string) EMVQR {
	tlv := TLV{
		Tag:    id,
		Length: l(v),
		Value:  v,
	}
	c.MerchantAccountInformation = append(c.MerchantAccountInformation, tlv)
	return *c
}

// SetMerchantCategoryCode ...
func (c *EMVQR) SetMerchantCategoryCode(v string) EMVQR {
	tlv := TLV{
		Tag:    IDMerchantCategoryCode,
		Length: l(v),
		Value:  v,
	}
	c.MerchantCategoryCode = tlv
	return *c
}

// SetTransactionCurrency ...
func (c *EMVQR) SetTransactionCurrency(v string) EMVQR {
	tlv := TLV{
		Tag:    IDTransactionCurrency,
		Length: l(v),
		Value:  v,
	}
	c.TransactionCurrency = tlv
	return *c
}

// SetTransactionAmount ...
func (c *EMVQR) SetTransactionAmount(v string) EMVQR {
	tlv := TLV{
		Tag:    IDTransactionAmount,
		Length: l(v),
		Value:  v,
	}
	c.TransactionAmount = tlv
	return *c
}

// SetTipOrConvenienceIndicator ...
func (c *EMVQR) SetTipOrConvenienceIndicator(v string) EMVQR {
	tlv := TLV{
		Tag:    IDTipOrConvenienceIndicator,
		Length: l(v),
		Value:  v,
	}
	c.TipOrConvenienceIndicator = tlv
	return *c
}

// SetValueOfConvenienceFeeFixed ...
func (c *EMVQR) SetValueOfConvenienceFeeFixed(v string) EMVQR {
	tlv := TLV{
		Tag:    IDValueOfConvenienceFeeFixed,
		Length: l(v),
		Value:  v,
	}
	c.ValueOfConvenienceFeeFixed = tlv
	return *c
}

// SetValueOfConvenienceFeePercentage ...
func (c *EMVQR) SetValueOfConvenienceFeePercentage(v string) EMVQR {
	tlv := TLV{
		Tag:    IDValueOfConvenienceFeePercentage,
		Length: l(v),
		Value:  v,
	}
	c.ValueOfConvenienceFeePercentage = tlv
	return *c
}

// SetCountryCode ...
func (c *EMVQR) SetCountryCode(v string) EMVQR {
	tlv := TLV{
		Tag:    IDCountryCode,
		Length: l(v),
		Value:  v,
	}
	c.CountryCode = tlv
	return *c
}

// SetMerchantName ...
func (c *EMVQR) SetMerchantName(v string) EMVQR {
	tlv := TLV{
		Tag:    IDMerchantName,
		Length: l(v),
		Value:  v,
	}
	c.MerchantName = tlv
	return *c
}

// SetMerchantCity ...
func (c *EMVQR) SetMerchantCity(v string) EMVQR {
	tlv := TLV{
		Tag:    IDMerchantCity,
		Length: l(v),
		Value:  v,
	}
	c.MerchantCity = tlv
	return *c
}

// SetPostalCode ...
func (c *EMVQR) SetPostalCode(v string) EMVQR {
	tlv := TLV{
		Tag:    IDPostalCode,
		Length: l(v),
		Value:  v,
	}
	c.PostalCode = tlv
	return *c
}

// SetAdditionalDataFieldTemplate ...
func (c *EMVQR) SetAdditionalDataFieldTemplate(v *AdditionalDataFieldTemplate) EMVQR {
	c.AdditionalDataFieldTemplate = v
	return *c
}

// SetCRC ...
func (c *EMVQR) SetCRC(v string) EMVQR {
	tlv := TLV{
		Tag:    IDCRC,
		Length: l(v),
		Value:  v,
	}
	c.CRC = tlv
	return *c
}

// SetMerchantInformationLanguageTemplate ...
func (c *EMVQR) SetMerchantInformationLanguageTemplate(v MerchantInformationLanguageTemplate) EMVQR {
	c.MerchantInformationLanguageTemplate = v
	return *c
}

// AddRFUforEMVCo ...
func (c *EMVQR) AddRFUforEMVCo(id ID, v string) EMVQR {
	tlv := TLV{
		Tag:    id,
		Length: l(v),
		Value:  v,
	}
	c.RFUforEMVCo = append(c.RFUforEMVCo, tlv)
	return *c
}

// AddUnreservedTemplates ...
func (c *EMVQR) AddUnreservedTemplates(id ID, v string) EMVQR {
	tlv := TLV{
		Tag:    id,
		Length: l(v),
		Value:  v,
	}
	c.UnreservedTemplates = append(c.UnreservedTemplates, tlv)
	return *c
}

// AdditionalDataFieldTemplate //

// SetBillNumber ...
func (s *AdditionalDataFieldTemplate) SetBillNumber(v string) {
	tlv := TLV{
		Tag:    AdditionalIDBillNumber,
		Length: l(v),
		Value:  v,
	}
	s.BillNumber = tlv
}

// SetMobileNumber ...
func (s *AdditionalDataFieldTemplate) SetMobileNumber(v string) {
	tlv := TLV{
		Tag:    AdditionalIDMobileNumber,
		Length: l(v),
		Value:  v,
	}
	s.MobileNumber = tlv
}

// SetStoreLabel ...
func (s *AdditionalDataFieldTemplate) SetStoreLabel(v string) {
	tlv := TLV{
		Tag:    AdditionalIDStoreLabel,
		Length: l(v),
		Value:  v,
	}
	s.StoreLabel = tlv
}

// SetLoyaltyNumber ...
func (s *AdditionalDataFieldTemplate) SetLoyaltyNumber(v string) {
	tlv := TLV{
		Tag:    AdditionalIDLoyaltyNumber,
		Length: l(v),
		Value:  v,
	}
	s.LoyaltyNumber = tlv
}

// SetReferenceLabel ...
func (s *AdditionalDataFieldTemplate) SetReferenceLabel(v string) {
	tlv := TLV{
		Tag:    AdditionalIDReferenceLabel,
		Length: l(v),
		Value:  v,
	}
	s.ReferenceLabel = tlv
}

// SetCustomerLabel ...
func (s *AdditionalDataFieldTemplate) SetCustomerLabel(v string) {
	tlv := TLV{
		Tag:    AdditionalIDCustomerLabel,
		Length: l(v),
		Value:  v,
	}
	s.CustomerLabel = tlv
}

// SetTerminalLabel ...
func (s *AdditionalDataFieldTemplate) SetTerminalLabel(v string) {
	tlv := TLV{
		Tag:    AdditionalIDTerminalLabel,
		Length: l(v),
		Value:  v,
	}
	s.TerminalLabel = tlv
}

// SetPurposeTransaction ...
func (s *AdditionalDataFieldTemplate) SetPurposeTransaction(v string) {
	tlv := TLV{
		Tag:    AdditionalIDPurposeTransaction,
		Length: l(v),
		Value:  v,
	}
	s.PurposeTransaction = tlv
}

// SetAdditionalConsumerDataRequest ...
func (s *AdditionalDataFieldTemplate) SetAdditionalConsumerDataRequest(v string) {
	tlv := TLV{
		Tag:    AdditionalIDAdditionalConsumerDataRequest,
		Length: l(v),
		Value:  v,
	}
	s.AdditionalConsumerDataRequest = tlv
}

// AddRFUforEMVCo ...
func (s *AdditionalDataFieldTemplate) AddRFUforEMVCo(id ID, v string) {
	tlv := TLV{
		Tag:    id,
		Length: l(v),
		Value:  v,
	}
	s.RFUforEMVCo = append(s.RFUforEMVCo, tlv)
}

// AddPaymentSystemSpecific ...
func (s *AdditionalDataFieldTemplate) AddPaymentSystemSpecific(id ID, v string) {
	tlv := TLV{
		Tag:    id,
		Length: l(v),
		Value:  v,
	}
	s.PaymentSystemSpecific = append(s.PaymentSystemSpecific, tlv)
}

// SetPaymentSystemSpecific ...
func (s AdditionalDataFieldTemplate) String() string {
	t := ""
	t += s.BillNumber.String()
	t += s.MobileNumber.String()
	t += s.StoreLabel.String()
	t += s.LoyaltyNumber.String()
	t += s.ReferenceLabel.String()
	t += s.CustomerLabel.String()
	t += s.TerminalLabel.String()
	t += s.PurposeTransaction.String()
	t += s.AdditionalConsumerDataRequest.String()
	for _, r := range s.RFUforEMVCo {
		t += r.String()
	}
	for _, p := range s.PaymentSystemSpecific {
		t += p.String()
	}
	tt := format(IDAdditionalDataFieldTemplate, t)
	return tt
}

// MerchantInformationLanguageTemplate //

// SetLanguagePreference ...
func (s MerchantInformationLanguageTemplate) SetLanguagePreference(v string) MerchantInformationLanguageTemplate {
	tlv := TLV{
		Tag:    MerchantInformationIDLanguagePreference,
		Length: l(v),
		Value:  v,
	}
	s.LanguagePreference = tlv
	return s
}

// SetMerchantName ..
func (s MerchantInformationLanguageTemplate) SetMerchantName(v string) MerchantInformationLanguageTemplate {
	tlv := TLV{
		Tag:    MerchantInformationIDMerchantName,
		Length: l(v),
		Value:  v,
	}
	s.MerchantName = tlv
	return s
}

// SetMerchantCity ...
func (s MerchantInformationLanguageTemplate) SetMerchantCity(v string) MerchantInformationLanguageTemplate {
	tlv := TLV{
		Tag:    MerchantInformationIDMerchantCity,
		Length: l(v),
		Value:  v,
	}
	s.MerchantCity = tlv
	return s
}

// AddRFUForEMVCo ...
func (s MerchantInformationLanguageTemplate) AddRFUForEMVCo(id ID, v string) MerchantInformationLanguageTemplate {
	tlv := TLV{
		Tag:    id,
		Length: l(v),
		Value:  v,
	}
	s.RFUforEMVCo = append(s.RFUforEMVCo, tlv)
	return s
}

// String() ...
func (s MerchantInformationLanguageTemplate) String() string {
	t := ""
	t += s.LanguagePreference.String()
	t += s.MerchantName.String()
	t += s.MerchantCity.String()
	for _, r := range s.RFUforEMVCo {
		t += r.String()
	}
	tt := format(IDMerchantInformationLanguageTemplate, t)
	return tt
}

//////////////////////////////////////////////////////////////////////////

// GeneratePayload ...
func (c *EMVQR) GeneratePayload() string {
	s := ""
	s += c.PayloadFormatIndicator.String()
	s += c.PointOfInitiationMethod.String()
	for _, m := range c.MerchantAccountInformation {
		s += m.String()
	}
	s += c.MerchantCategoryCode.String()
	s += c.TransactionCurrency.String()
	s += c.TransactionAmount.String()
	s += c.TipOrConvenienceIndicator.String()
	s += c.ValueOfConvenienceFeeFixed.String()
	s += c.ValueOfConvenienceFeePercentage.String()
	s += c.CountryCode.String()
	s += c.MerchantName.String()
	s += c.MerchantCity.String()
	s += c.PostalCode.String()
	s += c.AdditionalDataFieldTemplate.String()
	s += c.MerchantInformationLanguageTemplate.String()
	for _, r := range c.RFUforEMVCo {
		s += r.String()
	}
	for _, u := range c.UnreservedTemplates {
		s += u.String()
	}
	s += formatCrc(s)
	return s
}

// ParseEMVQR ...
func ParseEMVQR(payload string) (*EMVQR, error) {
	p := NewParser(payload)
	emvqr := &EMVQR{}
	for p.Next() {
		id := p.ID()
		// length := p.ValueLength()
		value := p.Value()
		switch id {
		case IDPayloadFormatIndicator:
			emvqr.SetPayloadFormatIndicator(value)
		case IDPointOfInitiationMethod:
			emvqr.SetPointOfInitiationMethod(value)
		case IDMerchantCategoryCode:
			emvqr.SetMerchantCategoryCode(value)
		case IDTransactionCurrency:
			emvqr.SetTransactionCurrency(value)
		case IDTransactionAmount:
			emvqr.SetTransactionAmount(value)
		case IDTipOrConvenienceIndicator:
			emvqr.SetTipOrConvenienceIndicator(value)
		case IDValueOfConvenienceFeeFixed:
			emvqr.SetValueOfConvenienceFeeFixed(value)
		case IDValueOfConvenienceFeePercentage:
			emvqr.SetValueOfConvenienceFeePercentage(value)
		case IDCountryCode:
			emvqr.SetCountryCode(value)
		case IDMerchantName:
			emvqr.SetMerchantName(value)
		case IDMerchantCity:
			emvqr.SetMerchantCity(value)
		case IDPostalCode:
			emvqr.SetPostalCode(value)
		case IDAdditionalDataFieldTemplate:
			adft, err := ParseAdditionalDataFieldTemplate(value)
			if err != nil {
				return nil, err
			}
			emvqr.AdditionalDataFieldTemplate = adft
		case IDCRC:
			emvqr.SetCRC(value)
		case IDMerchantInformationLanguageTemplate:
			t, err := ParseMerchantInformationLanguageTemplate(value)
			if err != nil {
				return nil, err
			}
			emvqr.SetMerchantInformationLanguageTemplate(*t)
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
				// if emvqr.MerchantAccountInformation == nil {
				// 	emvqr.MerchantAccountInformation = make(map[ID]*MerchantAccountInformation)
				// }
				emvqr.AddMerchantAccountInformation(id, value)
				continue
			}
			// RFUforEMVCo
			within, err = id.Between(IDRFUForEMVCoRangeStart, IDRFUForEMVCoRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				// if emvqr.RFUforEMVCo == nil {
				// 	emvqr.RFUforEMVCo = make(map[ID]*RFUforEMVCo)
				// }
				emvqr.AddRFUforEMVCo(id, value)
				continue
			}
			// Unreserved Tempaltes
			within, err = id.Between(IDUnreservedTemplatesRangeStart, IDUnreservedTemplatesRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				// if emvqr.UnreservedTemplates == nil {
				// 	emvqr.UnreservedTemplates = make(map[ID]*UnreservedTemplate)
				// }
				emvqr.AddUnreservedTemplates(id, value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return emvqr, nil
}

// // Validate ...
// func (c *EMVQR) Validate() error {
// 	// check mandatory
// 	if c.PayloadFormatIndicator == "" {
// 		return errors.New("PayloadFormatIndicator is mandatory")
// 	}
// 	if len(c.MerchantAccountInformation) <= 0 {
// 		return errors.New("MerchantAccountInformation is mandatory")
// 	}
// 	if c.MerchantCategoryCode == "" {
// 		return errors.New("MerchantCategoryCode is mandatory")
// 	}
// 	if c.TransactionCurrency == "" {
// 		return errors.New("TransactionCurrency is mandatory")
// 	}
// 	if c.CountryCode == "" {
// 		return errors.New("CountryCode is mandatory")
// 	}
// 	if c.MerchantName == "" {
// 		return errors.New("MerchantName is mandatory")
// 	}
// 	if c.MerchantCity == "" {
// 		return errors.New("MerchantCity is mandatory")
// 	}
// 	// check validate
// 	if c.PointOfInitiationMethod != "" {
// 		if err := c.PointOfInitiationMethod.Validate(); err != nil {
// 			return err
// 		}
// 	}
// 	for _, t := range c.MerchantAccountInformation {
// 		if err := t.Validate(); err != nil {
// 			return err
// 		}
// 	}
// 	if c.AdditionalDataFieldTemplate != nil {
// 		if err := c.AdditionalDataFieldTemplate.Validate(); err != nil {
// 			return err
// 		}
// 	}
// 	if c.MerchantInformationLanguageTemplate != nil {
// 		if err := c.MerchantInformationLanguageTemplate.Validate(); err != nil {
// 			return err
// 		}
// 	}
// 	for _, t := range c.RFUforEMVCo {
// 		if err := t.Validate(); err != nil {
// 			return err
// 		}
// 	}
// 	for _, t := range c.UnreservedTemplates {
// 		if err := t.Validate(); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// const (
// 	// PointOfInitiationMethodStatic ...
// 	PointOfInitiationMethodStatic PointOfInitiationMethod = "11"
// 	// PointOfInitiationMethodDynamic ...
// 	PointOfInitiationMethodDynamic PointOfInitiationMethod = "12"
// )

// // ParsePointOfInitiationMethod ...
// func ParsePointOfInitiationMethod(payload string) PointOfInitiationMethod {
// 	return PointOfInitiationMethod(payload)
// }

// // GeneratePayload ...
// func (c PointOfInitiationMethod) GeneratePayload() string {
// 	return format(IDPointOfInitiationMethod, string(c))
// }

// // IsStaticMethod ...
// func (c PointOfInitiationMethod) IsStaticMethod() bool {
// 	return c == PointOfInitiationMethodStatic
// }

// // IsDynamicMethod ...
// func (c PointOfInitiationMethod) IsDynamicMethod() bool {
// 	return c == PointOfInitiationMethodDynamic
// }

// // Validate ...
// func (c PointOfInitiationMethod) Validate() error {
// 	if !c.IsStaticMethod() && !c.IsDynamicMethod() {
// 		return fmt.Errorf("PointOfInitiationMethod should be \"11\" or \"12\", PointOfInitiationMethod: %s", c)
// 	}
// 	return nil
// }

// // RFUforEMVCo ...
// type RFUforEMVCo struct {
// 	Value string
// }

// // ParseRFUforEMVCo ...
// func ParseRFUforEMVCo(payload string) *RFUforEMVCo {
// 	return &RFUforEMVCo{
// 		Value: payload,
// 	}
// }

// // GeneratePayload ...
// func (c *RFUforEMVCo) GeneratePayload() string {
// 	return c.Value
// }

// // Validate ...
// func (c *RFUforEMVCo) Validate() error {
// 	return nil
// }

// // UnreservedTemplate ...
// type UnreservedTemplate struct {
// 	Value string
// }

// // ParseUnreservedTemplate ...
// func ParseUnreservedTemplate(payload string) *UnreservedTemplate {
// 	return &UnreservedTemplate{
// 		Value: payload,
// 	}
// }

// // GeneratePayload ...
// func (c *UnreservedTemplate) GeneratePayload() string {
// 	return c.Value
// }

// // Validate ...
// func (c *UnreservedTemplate) Validate() error {
// 	return nil
// }

// ParseMerchantAccountInformation ...
// func ParseMerchantAccountInformation(id ID, length string, payload string) *MerchantAccountInformation {
// 	return &MerchantAccountInformation{
// 		Tag:    id,
// 		Length: length,
// 		Value:  payload,
// 	}
// }

// // GeneratePayload ...
// func (c *MerchantAccountInformation) GeneratePayload() string {
// 	return c.Value
// }

// // Validate ...
// func (c *MerchantAccountInformation) Validate() error {
// 	return nil
// }

// ParseAdditionalDataFieldTemplate ...
func ParseAdditionalDataFieldTemplate(payload string) (*AdditionalDataFieldTemplate, error) {
	p := NewParser(payload)
	additionalDataFieldTemplate := &AdditionalDataFieldTemplate{}
	for p.Next() {
		id := p.ID()
		value := p.Value()
		switch id {
		case AdditionalIDBillNumber:
			additionalDataFieldTemplate.SetBillNumber(value)
		case AdditionalIDMobileNumber:
			additionalDataFieldTemplate.SetMobileNumber(value)
		case AdditionalIDStoreLabel:
			additionalDataFieldTemplate.SetStoreLabel(value)
		case AdditionalIDLoyaltyNumber:
			additionalDataFieldTemplate.SetLoyaltyNumber(value)
		case AdditionalIDReferenceLabel:
			additionalDataFieldTemplate.SetReferenceLabel(value)
		case AdditionalIDCustomerLabel:
			additionalDataFieldTemplate.SetCustomerLabel(value)
		case AdditionalIDTerminalLabel:
			additionalDataFieldTemplate.SetTerminalLabel(value)
		case AdditionalIDPurposeTransaction:
			additionalDataFieldTemplate.SetPurposeTransaction(value)
		case AdditionalIDAdditionalConsumerDataRequest:
			additionalDataFieldTemplate.SetAdditionalConsumerDataRequest(value)
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
				// if additionalDataFieldTemplate.PaymentSystemSpecific == nil {
				// 	additionalDataFieldTemplate.PaymentSystemSpecific = make(map[ID]*PaymentSystemSpecific)
				// }
				additionalDataFieldTemplate.AddPaymentSystemSpecific(id, value)
				continue
			}
			// RFU for EMVCo
			within, err = id.Between(AdditionalIDRFUforEMVCoRangeStart, AdditionalIDRFUforEMVCoRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				// if additionalDataFieldTemplate.RFUforEMVCo == nil {
				// 	additionalDataFieldTemplate.RFUforEMVCo = make(map[ID]*RFUforEMVCo)
				// }
				additionalDataFieldTemplate.AddRFUforEMVCo(id, value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return additionalDataFieldTemplate, nil
}

// // GeneratePayload ...
// func (c *AdditionalDataFieldTemplate) GeneratePayload() string {
// 	s := ""
// 	if c.BillNumber != "" {
// 		s += format(AdditionalIDBillNumber, c.BillNumber)
// 	}
// 	if c.MobileNumber != "" {
// 		s += format(AdditionalIDMobileNumber, c.MobileNumber)
// 	}
// 	if c.StoreLabel != "" {
// 		s += format(AdditionalIDStoreLabel, c.StoreLabel)
// 	}
// 	if c.LoyaltyNumber != "" {
// 		s += format(AdditionalIDLoyaltyNumber, c.LoyaltyNumber)
// 	}
// 	if c.ReferenceLabel != "" {
// 		s += format(AdditionalIDReferenceLabel, c.ReferenceLabel)
// 	}
// 	if c.CustomerLabel != "" {
// 		s += format(AdditionalIDCustomerLabel, c.CustomerLabel)
// 	}
// 	if c.TerminalLabel != "" {
// 		s += format(AdditionalIDTerminalLabel, c.TerminalLabel)
// 	}
// 	if c.PurposeTransaction != "" {
// 		s += format(AdditionalIDPurposeTransaction, c.PurposeTransaction)
// 	}
// 	if c.AdditionalConsumerDataRequest != "" {
// 		s += format(AdditionalIDAdditionalConsumerDataRequest, c.AdditionalConsumerDataRequest)
// 	}
// 	if len(c.RFUforEMVCo) > 0 {
// 		for k, t := range c.RFUforEMVCo {
// 			s += format(ID(k), t.GeneratePayload())
// 		}
// 	}
// 	if len(c.PaymentSystemSpecific) > 0 {
// 		for k, t := range c.PaymentSystemSpecific {
// 			s += format(ID(k), t.GeneratePayload())
// 		}
// 	}
// 	return s
// }

// // Validate ...
// func (c *AdditionalDataFieldTemplate) Validate() error {
// 	return nil
// }

// // PaymentSystemSpecific ...
// type PaymentSystemSpecific struct {
// 	Value string
// }

// // ParsePaymentSystemSpecific ...
// func ParsePaymentSystemSpecific(value string) *PaymentSystemSpecific {
// 	return &PaymentSystemSpecific{Value: value}
// }

// // GeneratePayload ...
// func (v *PaymentSystemSpecific) GeneratePayload() string {
// 	return v.Value
// }

// ParseMerchantInformationLanguageTemplate ...
func ParseMerchantInformationLanguageTemplate(value string) (*MerchantInformationLanguageTemplate, error) {
	p := NewParser(value)
	merchantInformationLanguageTemplate := &MerchantInformationLanguageTemplate{}
	for p.Next() {
		id := p.ID()
		value := p.Value()
		switch id {
		case MerchantInformationIDLanguagePreference:
			merchantInformationLanguageTemplate.SetLanguagePreference(value)
		case MerchantInformationIDMerchantName:
			merchantInformationLanguageTemplate.SetMerchantName(value)
		case MerchantInformationIDMerchantCity:
			merchantInformationLanguageTemplate.SetMerchantCity(value)
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
				// if merchantInformationLanguageTemplate.RFUForEMVCo == nil {
				// 	merchantInformationLanguageTemplate.RFUForEMVCo = make(map[ID]*RFUforEMVCo)
				// }
				merchantInformationLanguageTemplate.AddRFUForEMVCo(id, value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return merchantInformationLanguageTemplate, nil
}

// // GeneratePayload ...
// func (c *MerchantInformationLanguageTemplate) GeneratePayload() string {
// 	s := ""
// 	if c.LanguagePreference != "" {
// 		s += format(MerchantInformationIDLanguagePreference, c.LanguagePreference)
// 	}
// 	if c.MerchantName != "" {
// 		s += format(MerchantInformationIDMerchantName, c.MerchantName)
// 	}
// 	if c.MerchantCity != "" {
// 		s += format(MerchantInformationIDMerchantCity, c.MerchantCity)
// 	}
// 	if len(c.RFUForEMVCo) > 0 {
// 		for id, t := range c.RFUForEMVCo {
// 			s += format(id, t.GeneratePayload())
// 		}
// 	}
// 	return s
// }

// // Validate ...
// func (c *MerchantInformationLanguageTemplate) Validate() error {
// 	// check mandatory
// 	if c.LanguagePreference == "" {
// 		return errors.New("LanguagePreference is mandatory")
// 	}
// 	if c.MerchantName == "" {
// 		return errors.New("MerchantName is mandatory")
// 	}
// 	// check validate
// 	for _, t := range c.RFUForEMVCo {
// 		if err := t.Validate(); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

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

func l(v string) string {
	return fmt.Sprintf("%02d", utf8.RuneCountInString(v))
}
