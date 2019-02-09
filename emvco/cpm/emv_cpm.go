package cpm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// const ....
const (
	IDPayloadFormatIndicator = "00" // (M) Payload Format Indicator
)

// Data Objects for Additional Data Field Template (ID "62")
const (
	AdditionalIDBillNumber = "01"
)

// Data Objects for Merchant Information—Language Template (ID "64")
const (
	MerchantInformationIDLanguagePreference = "00"
)

// EMVQR ...
type EMVQR struct {
	PayloadFormatIndicator string
}

// GeneratePayload ...
func (c *EMVQR) GeneratePayload() (string, error) {
	s := ""
	if c.PayloadFormatIndicator != "" {
		s += format(IDPayloadFormatIndicator, c.PayloadFormatIndicator)
	} else {
		return "", fmt.Errorf("PayloadFormatIndicator is mandatory")
	}
	return s, nil
}

func format(id, value string) string {
	length := utf8.RuneCountInString(value)
	lengthStr := strconv.Itoa(length)
	lengthStr = "00" + lengthStr
	return id + lengthStr[len(lengthStr)-2:] + value
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.0f", amount)
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
