package mpm

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/dongri/emv-qrcode/crc16"
)

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

// Data Object ID Allocation in Merchant Account Information Template ...
const (
	MerchantAccountInformationIDGloballyUniqueIdentifier    ID = "00"
	MerchantAccountInformationIDPaymentNetworkSpecificStart ID = "01" // (O) 03-99 RFU for EMVCo
	MerchantAccountInformationIDPaymentNetworkSpecificEnd   ID = "99" // (O) 03-99 RFU for EMVCo
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

// Data Object ID Allocation in Merchant Account Information Template ...
const (
	UnreservedTemplateIDGloballyUniqueIdentifier ID = "00"
	UnreservedTemplateIDContextSpecificDataStart ID = "01" // (O) 03-99 RFU for EMVCo
	UnreservedTemplateIDContextSpecificDataEnd   ID = "99" // (O) 03-99 RFU for EMVCo
)

const (
	// PointOfInitiationMethodStatic ...
	PointOfInitiationMethodStatic = "11"
	// PointOfInitiationMethodDynamic ...
	PointOfInitiationMethodDynamic = "12"
)

// EMVQR ...
type EMVQR struct {
	PayloadFormatIndicator              TLV                                  `json:"Payload Format Indicator"`
	PointOfInitiationMethod             TLV                                  `json:"Point of Initiation Method"`
	MerchantAccountInformation          map[ID]MerchantAccountInformationTLV `json:"Merchant Account Information"`
	MerchantCategoryCode                TLV                                  `json:"Merchant Category Code"`
	TransactionCurrency                 TLV                                  `json:"Transaction Currency"`
	TransactionAmount                   TLV                                  `json:"Transaction Amount"`
	TipOrConvenienceIndicator           TLV                                  `json:"Tip or Convenience Indicator"`
	ValueOfConvenienceFeeFixed          TLV                                  `json:"Value of Convenience Fee Fixed"`
	ValueOfConvenienceFeePercentage     TLV                                  `json:"Value of Convenience Fee Percentage"`
	CountryCode                         TLV                                  `json:"Country Code"`
	MerchantName                        TLV                                  `json:"Merchant Name"`
	MerchantCity                        TLV                                  `json:"Merchant City"`
	PostalCode                          TLV                                  `json:"Postal Code"`
	AdditionalDataFieldTemplate         *AdditionalDataFieldTemplate         `json:"Additional Data Field Template"`
	CRC                                 TLV                                  `json:"CRC"`
	MerchantInformationLanguageTemplate *MerchantInformationLanguageTemplate `json:"Merchant Information - Language Template"`
	RFUforEMVCo                         []TLV                                `json:"RFU for EMVCo"`
	UnreservedTemplates                 map[ID]UnreservedTemplateTLV         `json:"Unreserved Templates"`
}

// MerchantAccountInformationTLV ...
type MerchantAccountInformationTLV struct {
	Tag    ID
	Length string
	Value  *MerchantAccountInformation
}

// MerchantAccountInformation ...
type MerchantAccountInformation struct {
	GloballyUniqueIdentifier TLV   `json:"Globally Unique Identifier"`
	PaymentNetworkSpecific   []TLV `json:"Payment network specific"`
}

// AdditionalDataFieldTemplate ...
type AdditionalDataFieldTemplate struct {
	BillNumber                    TLV   `json:"Bill Number"`
	MobileNumber                  TLV   `json:"Country Code"`
	StoreLabel                    TLV   `json:"Store Label"`
	LoyaltyNumber                 TLV   `json:"Loyalty Number"`
	ReferenceLabel                TLV   `json:"Reference Label"`
	CustomerLabel                 TLV   `json:"Customer Label"`
	TerminalLabel                 TLV   `json:"Terminal Label"`
	PurposeTransaction            TLV   `json:"Purpose of Transaction"`
	AdditionalConsumerDataRequest TLV   `json:"Additional Consumer Data Request"`
	RFUforEMVCo                   []TLV `json:"RFU for EMVCo"`
	PaymentSystemSpecific         []TLV `json:"Payment System specific templates"`
}

// MerchantInformationLanguageTemplate ...
type MerchantInformationLanguageTemplate struct {
	LanguagePreference TLV   `json:"Language Preference"`
	MerchantName       TLV   `json:"Merchant Name"`
	MerchantCity       TLV   `json:"Merchant City"`
	RFUforEMVCo        []TLV `json:"RFU for EMVCo"`
}

// UnreservedTemplateTLV ...
type UnreservedTemplateTLV struct {
	Tag    ID
	Length string
	Value  *UnreservedTemplate
}

// UnreservedTemplate ...
type UnreservedTemplate struct {
	GloballyUniqueIdentifier TLV   `json:"Globally Unique Identifier"`
	ContextSpecificData      []TLV `json:"Context Specific Data"`
}

// DataType ...
type DataType string

// const ...
const (
	DataTypeBinary DataType = "binary"
	DataTypeRaw    DataType = "raw"
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

// DataWithType ...
func (tlv TLV) DataWithType(dataType DataType, indent string) string {
	if tlv.Value == "" {
		return ""
	}
	if dataType == DataTypeBinary {
		rep := regexp.MustCompile("(.{2})")
		hexStr := hex.EncodeToString([]byte(tlv.Value))
		hexArray := rep.FindAllString(hexStr, -1)
		return indent + tlv.Tag.String() + " " + tlv.Length + " " + strings.Join(hexArray, " ") + "\n"
	}
	if dataType == DataTypeRaw {
		return indent + tlv.Tag.String() + " " + tlv.Length + " " + tlv.Value + "\n"
	}
	return ""
}

// BinaryData ...
func (c *EMVQR) BinaryData() string {
	return c.dataWithType(DataTypeBinary)
}

// RawData ...
func (c *EMVQR) RawData() string {
	return c.dataWithType(DataTypeRaw)
}

// JSON ...
func (c *EMVQR) JSON() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (c *EMVQR) dataWithType(dataType DataType) string {
	indent := ""
	s := ""
	s += c.PayloadFormatIndicator.DataWithType(dataType, indent)
	s += c.PointOfInitiationMethod.DataWithType(dataType, indent)
	for _, m := range c.MerchantAccountInformation {
		s += m.DataWithType(dataType, " ")
	}
	s += c.MerchantCategoryCode.DataWithType(dataType, indent)
	s += c.TransactionCurrency.DataWithType(dataType, indent)
	s += c.TransactionAmount.DataWithType(dataType, indent)
	s += c.TipOrConvenienceIndicator.DataWithType(dataType, indent)
	s += c.ValueOfConvenienceFeeFixed.DataWithType(dataType, indent)
	s += c.ValueOfConvenienceFeePercentage.DataWithType(dataType, indent)
	s += c.CountryCode.DataWithType(dataType, indent)
	s += c.MerchantName.DataWithType(dataType, indent)
	s += c.MerchantCity.DataWithType(dataType, indent)
	s += c.PostalCode.DataWithType(dataType, indent)
	s += c.AdditionalDataFieldTemplate.DataWithType(dataType, " ")
	s += c.MerchantInformationLanguageTemplate.DataWithType(dataType, " ")
	for _, r := range c.RFUforEMVCo {
		s += r.DataWithType(dataType, " ")
	}
	for _, u := range c.UnreservedTemplates {
		s += u.DataWithType(dataType, " ")
	}
	s += c.CRC.DataWithType(dataType, indent)
	return s
}

// SetPayloadFormatIndicator ...
func (c *EMVQR) SetPayloadFormatIndicator(v string) {
	tlv := TLV{
		Tag:    IDPayloadFormatIndicator,
		Length: l(v),
		Value:  v,
	}
	c.PayloadFormatIndicator = tlv
}

// SetPointOfInitiationMethod ...
func (c *EMVQR) SetPointOfInitiationMethod(v string) {
	tlv := TLV{
		Tag:    IDPointOfInitiationMethod,
		Length: l(v),
		Value:  v,
	}
	c.PointOfInitiationMethod = tlv
}

// AddMerchantAccountInformation ...
func (c *EMVQR) AddMerchantAccountInformation(id ID, v *MerchantAccountInformation) {
	tlv := MerchantAccountInformationTLV{
		Tag:    id,
		Length: l(v.String()),
		Value:  v,
	}
	if c.MerchantAccountInformation == nil {
		c.MerchantAccountInformation = make(map[ID]MerchantAccountInformationTLV)
	}
	c.MerchantAccountInformation[id] = tlv
}

// SetMerchantCategoryCode ...
func (c *EMVQR) SetMerchantCategoryCode(v string) {
	tlv := TLV{
		Tag:    IDMerchantCategoryCode,
		Length: l(v),
		Value:  v,
	}
	c.MerchantCategoryCode = tlv
}

// SetTransactionCurrency ...
func (c *EMVQR) SetTransactionCurrency(v string) {
	tlv := TLV{
		Tag:    IDTransactionCurrency,
		Length: l(v),
		Value:  v,
	}
	c.TransactionCurrency = tlv
}

// SetTransactionAmount ...
func (c *EMVQR) SetTransactionAmount(v string) {
	tlv := TLV{
		Tag:    IDTransactionAmount,
		Length: l(v),
		Value:  v,
	}
	c.TransactionAmount = tlv
}

// SetTipOrConvenienceIndicator ...
func (c *EMVQR) SetTipOrConvenienceIndicator(v string) {
	tlv := TLV{
		Tag:    IDTipOrConvenienceIndicator,
		Length: l(v),
		Value:  v,
	}
	c.TipOrConvenienceIndicator = tlv
}

// SetValueOfConvenienceFeeFixed ...
func (c *EMVQR) SetValueOfConvenienceFeeFixed(v string) {
	tlv := TLV{
		Tag:    IDValueOfConvenienceFeeFixed,
		Length: l(v),
		Value:  v,
	}
	c.ValueOfConvenienceFeeFixed = tlv
}

// SetValueOfConvenienceFeePercentage ...
func (c *EMVQR) SetValueOfConvenienceFeePercentage(v string) {
	tlv := TLV{
		Tag:    IDValueOfConvenienceFeePercentage,
		Length: l(v),
		Value:  v,
	}
	c.ValueOfConvenienceFeePercentage = tlv
}

// SetCountryCode ...
func (c *EMVQR) SetCountryCode(v string) {
	tlv := TLV{
		Tag:    IDCountryCode,
		Length: l(v),
		Value:  v,
	}
	c.CountryCode = tlv
}

// SetMerchantName ...
func (c *EMVQR) SetMerchantName(v string) {
	tlv := TLV{
		Tag:    IDMerchantName,
		Length: l(v),
		Value:  v,
	}
	c.MerchantName = tlv
}

// SetMerchantCity ...
func (c *EMVQR) SetMerchantCity(v string) {
	tlv := TLV{
		Tag:    IDMerchantCity,
		Length: l(v),
		Value:  v,
	}
	c.MerchantCity = tlv
}

// SetPostalCode ...
func (c *EMVQR) SetPostalCode(v string) {
	tlv := TLV{
		Tag:    IDPostalCode,
		Length: l(v),
		Value:  v,
	}
	c.PostalCode = tlv
}

// SetAdditionalDataFieldTemplate ...
func (c *EMVQR) SetAdditionalDataFieldTemplate(v *AdditionalDataFieldTemplate) {
	c.AdditionalDataFieldTemplate = v
}

// SetCRC ...
func (c *EMVQR) SetCRC(v string) {
	tlv := TLV{
		Tag:    IDCRC,
		Length: l(v),
		Value:  v,
	}
	c.CRC = tlv
}

// SetMerchantInformationLanguageTemplate ...
func (c *EMVQR) SetMerchantInformationLanguageTemplate(v *MerchantInformationLanguageTemplate) {
	c.MerchantInformationLanguageTemplate = v
}

// AddRFUforEMVCo ...
func (c *EMVQR) AddRFUforEMVCo(id ID, v string) {
	tlv := TLV{
		Tag:    id,
		Length: l(v),
		Value:  v,
	}
	c.RFUforEMVCo = append(c.RFUforEMVCo, tlv)
}

// AddUnreservedTemplates ...
func (c *EMVQR) AddUnreservedTemplates(id ID, v *UnreservedTemplate) {
	tlv := UnreservedTemplateTLV{
		Tag:    id,
		Length: l(v.String()),
		Value:  v,
	}
	if c.UnreservedTemplates == nil {
		c.UnreservedTemplates = make(map[ID]UnreservedTemplateTLV)
	}
	c.UnreservedTemplates[id] = tlv
}

// MerchantAccountInformation //

func (s *MerchantAccountInformationTLV) String() string {
	if s == nil {
		return ""
	}
	t := ""
	t += s.Tag.String() + s.Length + s.Value.String()
	return t
}

// DataWithType ..
func (s *MerchantAccountInformationTLV) DataWithType(dataType DataType, indent string) string {
	if s == nil {
		return ""
	}
	return s.Tag.String() + " " + s.Length + "\n" + s.Value.DataWithType(dataType, indent)
}

// SetGloballyUniqueIdentifier ...
func (s *MerchantAccountInformation) SetGloballyUniqueIdentifier(v string) {
	tlv := TLV{
		Tag:    MerchantAccountInformationIDGloballyUniqueIdentifier,
		Length: l(v),
		Value:  v,
	}
	s.GloballyUniqueIdentifier = tlv
}

// AddPaymentNetworkSpecific ...
func (s *MerchantAccountInformation) AddPaymentNetworkSpecific(id ID, v string) {
	tlv := TLV{
		Tag:    id,
		Length: l(v),
		Value:  v,
	}
	s.PaymentNetworkSpecific = append(s.PaymentNetworkSpecific, tlv)
}

func (s *MerchantAccountInformation) String() string {
	if s == nil {
		return ""
	}
	t := ""
	t += s.GloballyUniqueIdentifier.String()
	sort.Slice(s.PaymentNetworkSpecific, func(i, j int) bool {
		return s.PaymentNetworkSpecific[i].Tag < s.PaymentNetworkSpecific[j].Tag
	})
	for _, pns := range s.PaymentNetworkSpecific {
		t += pns.String()
	}
	return t
}

// DataWithType ...
func (s *MerchantAccountInformation) DataWithType(dataType DataType, indent string) string {
	if s == nil {
		return ""
	}
	var pnsData string
	for _, pns := range s.PaymentNetworkSpecific {
		pnsData += indent + pns.DataWithType(dataType, indent)
	}
	return indent + s.GloballyUniqueIdentifier.DataWithType(dataType, indent) + pnsData
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

// Stirng ...
func (s *AdditionalDataFieldTemplate) String() string {
	if s == nil {
		return ""
	}
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

// DataWithType ...
func (s *AdditionalDataFieldTemplate) DataWithType(dataType DataType, indent string) string {
	if s == nil {
		return ""
	}
	t := ""
	t += s.BillNumber.DataWithType(dataType, indent)
	t += s.MobileNumber.DataWithType(dataType, indent)
	t += s.StoreLabel.DataWithType(dataType, indent)
	t += s.LoyaltyNumber.DataWithType(dataType, indent)
	t += s.ReferenceLabel.DataWithType(dataType, indent)
	t += s.CustomerLabel.DataWithType(dataType, indent)
	t += s.TerminalLabel.DataWithType(dataType, indent)
	t += s.PurposeTransaction.DataWithType(dataType, indent)
	t += s.AdditionalConsumerDataRequest.DataWithType(dataType, indent)
	for _, r := range s.RFUforEMVCo {
		t += r.DataWithType(dataType, indent)
	}
	for _, p := range s.PaymentSystemSpecific {
		t += p.DataWithType(dataType, indent)
	}
	tt := IDAdditionalDataFieldTemplate.String() + " " + ll(s.String()) + "\n" + t
	return tt
}

// MerchantInformationLanguageTemplate //

// SetLanguagePreference ...
func (s *MerchantInformationLanguageTemplate) SetLanguagePreference(v string) {
	tlv := TLV{
		Tag:    MerchantInformationIDLanguagePreference,
		Length: l(v),
		Value:  v,
	}
	s.LanguagePreference = tlv
}

// SetMerchantName ..
func (s *MerchantInformationLanguageTemplate) SetMerchantName(v string) {
	tlv := TLV{
		Tag:    MerchantInformationIDMerchantName,
		Length: l(v),
		Value:  v,
	}
	s.MerchantName = tlv
}

// SetMerchantCity ...
func (s *MerchantInformationLanguageTemplate) SetMerchantCity(v string) {
	tlv := TLV{
		Tag:    MerchantInformationIDMerchantCity,
		Length: l(v),
		Value:  v,
	}
	s.MerchantCity = tlv
}

// AddRFUForEMVCo ...
func (s *MerchantInformationLanguageTemplate) AddRFUForEMVCo(id ID, v string) {
	tlv := TLV{
		Tag:    id,
		Length: l(v),
		Value:  v,
	}
	s.RFUforEMVCo = append(s.RFUforEMVCo, tlv)
}

// String() ...
func (s *MerchantInformationLanguageTemplate) String() string {
	if s == nil {
		return ""
	}
	t := ""
	t += s.LanguagePreference.String()
	t += s.MerchantName.String()
	t += s.MerchantCity.String()
	for _, r := range s.RFUforEMVCo {
		t += r.String()
	}
	t = format(IDMerchantInformationLanguageTemplate, t)
	return t
}

// DataWithType ...
func (s *MerchantInformationLanguageTemplate) DataWithType(dataType DataType, indent string) string {
	if s == nil {
		return ""
	}
	t := ""
	t += s.LanguagePreference.DataWithType(dataType, indent)
	t += s.MerchantName.DataWithType(dataType, indent)
	t += s.MerchantCity.DataWithType(dataType, indent)
	for _, r := range s.RFUforEMVCo {
		t += r.DataWithType(dataType, indent)
	}
	t = IDMerchantInformationLanguageTemplate.String() + " " + ll(s.String()) + "\n" + t
	return t
}

// MerchantAccountInformation //

func (s *UnreservedTemplateTLV) String() string {
	if s == nil {
		return ""
	}
	t := ""
	t += s.Tag.String() + s.Length + s.Value.String()
	return t
}

// DataWithType ..
func (s *UnreservedTemplateTLV) DataWithType(dataType DataType, indent string) string {
	if s == nil {
		return ""
	}
	return s.Tag.String() + " " + s.Length + "\n" + s.Value.DataWithType(dataType, indent)
}

// SetGloballyUniqueIdentifier ...
func (s *UnreservedTemplate) SetGloballyUniqueIdentifier(v string) {
	tlv := TLV{
		Tag:    UnreservedTemplateIDGloballyUniqueIdentifier,
		Length: l(v),
		Value:  v,
	}
	s.GloballyUniqueIdentifier = tlv
}

// AddContextSpecificData ...
func (s *UnreservedTemplate) AddContextSpecificData(id ID, v string) {
	tlv := TLV{
		Tag:    id,
		Length: l(v),
		Value:  v,
	}
	s.ContextSpecificData = append(s.ContextSpecificData, tlv)

}

func (s *UnreservedTemplate) String() string {
	if s == nil {
		return ""
	}
	t := ""
	t += s.GloballyUniqueIdentifier.String()
	for _, c := range s.ContextSpecificData {
		t += c.String()
	}
	return t
}

// DataWithType ...
func (s *UnreservedTemplate) DataWithType(dataType DataType, indent string) string {
	if s == nil {
		return ""
	}
	var csData string
	for _, cs := range s.ContextSpecificData {
		csData += indent + cs.DataWithType(dataType, indent)
	}
	return indent + s.GloballyUniqueIdentifier.DataWithType(dataType, indent) + csData
}

//////////////////////////////////////////////////////////////////////////

// GeneratePayload ...
func (c *EMVQR) GeneratePayload() string {
	s := ""
	s += c.PayloadFormatIndicator.String()
	s += c.PointOfInitiationMethod.String()
	var keys []string
	for k := range c.MerchantAccountInformation {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := c.MerchantAccountInformation[ID(k)]
		s += v.String()
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
			emvqr.SetMerchantInformationLanguageTemplate(t)
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
				t, err := ParseMerchantAccountInformation(value)
				if err != nil {
					return nil, err
				}
				emvqr.AddMerchantAccountInformation(id, t)
				continue
			}
			// RFUforEMVCo
			within, err = id.Between(IDRFUForEMVCoRangeStart, IDRFUForEMVCoRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				emvqr.AddRFUforEMVCo(id, value)
				continue
			}
			// Unreserved Tempaltes
			within, err = id.Between(IDUnreservedTemplatesRangeStart, IDUnreservedTemplatesRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
				t, err := ParseUnreservedTemplate(value)
				if err != nil {
					return nil, err
				}
				emvqr.AddUnreservedTemplates(id, t)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return emvqr, nil
}

// Validate ...
func (c *EMVQR) Validate() error {
	// check mandatory
	if c.PayloadFormatIndicator.Value == "" {
		return errors.New("PayloadFormatIndicator is mandatory")
	}
	if len(c.MerchantAccountInformation) <= 0 {
		return errors.New("MerchantAccountInformation is mandatory")
	}
	if c.TransactionCurrency.Value == "" {
		return errors.New("TransactionCurrency is mandatory")
	}
	if c.CountryCode.Value == "" {
		return errors.New("CountryCode is mandatory")
	}
	if c.MerchantCity.Value == "" {
		return errors.New("MerchantCity is mandatory")
	}

	// if c.MerchantCategoryCode.Value == "" {
	// 	return errors.New("MerchantCategoryCode is mandatory")
	// }
	// if c.MerchantName.Value == "" {
	// 	return errors.New("MerchantName is mandatory")
	// }

	// check validate
	if c.PointOfInitiationMethod.Value != "" {
		if c.PointOfInitiationMethod.Value != PointOfInitiationMethodStatic && c.PointOfInitiationMethod.Value != PointOfInitiationMethodDynamic {
			return fmt.Errorf("PointOfInitiationMethod should be \"11\" or \"12\", PointOfInitiationMethod: %s", c)
		}
	}
	if c.MerchantInformationLanguageTemplate != nil {
		if err := c.MerchantInformationLanguageTemplate.Validate(); err != nil {
			return err
		}
	}
	return nil
}

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
				additionalDataFieldTemplate.AddPaymentSystemSpecific(id, value)
				continue
			}
			// RFU for EMVCo
			within, err = id.Between(AdditionalIDRFUforEMVCoRangeStart, AdditionalIDRFUforEMVCoRangeEnd)
			if err != nil {
				return nil, err
			}
			if within {
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

// ParseMerchantAccountInformation ...
func ParseMerchantAccountInformation(value string) (*MerchantAccountInformation, error) {
	p := NewParser(value)
	merchantAccountInformation := &MerchantAccountInformation{}
	for p.Next() {
		id := p.ID()
		value := p.Value()
		switch id {
		case MerchantAccountInformationIDGloballyUniqueIdentifier:
			merchantAccountInformation.SetGloballyUniqueIdentifier(value)
		default:
			var (
				within bool
				err    error
			)
			within, err = id.Between(MerchantAccountInformationIDPaymentNetworkSpecificStart, MerchantAccountInformationIDPaymentNetworkSpecificEnd)
			if err != nil {
				return nil, err
			}
			if within {
				merchantAccountInformation.AddPaymentNetworkSpecific(id, value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return merchantAccountInformation, nil
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

// ParseUnreservedTemplate ...
func ParseUnreservedTemplate(value string) (*UnreservedTemplate, error) {
	p := NewParser(value)
	unreservedTemplate := &UnreservedTemplate{}
	for p.Next() {
		id := p.ID()
		value := p.Value()
		switch id {
		case UnreservedTemplateIDGloballyUniqueIdentifier:
			unreservedTemplate.SetGloballyUniqueIdentifier(value)
		default:
			var (
				within bool
				err    error
			)
			within, err = id.Between(UnreservedTemplateIDContextSpecificDataStart, UnreservedTemplateIDContextSpecificDataEnd)
			if err != nil {
				return nil, err
			}
			if within {
				unreservedTemplate.AddContextSpecificData(id, value)
				continue
			}
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return unreservedTemplate, nil
}

// Validate ...
func (s *MerchantInformationLanguageTemplate) Validate() error {
	// check mandatory
	if s.LanguagePreference.Value == "" {
		return errors.New("LanguagePreference is mandatory")
	}
	if s.MerchantName.Value == "" {
		return errors.New("MerchantName is mandatory")
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

func l(v string) string {
	return fmt.Sprintf("%02d", utf8.RuneCountInString(v))
}

func ll(v string) string {
	return fmt.Sprintf("%02d", utf8.RuneCountInString(v)-4)
}
