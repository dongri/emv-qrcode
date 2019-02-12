package cpm

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// const ...
const (
	IDPayloadFormatIndicator                 = "85" // (M) Payload Format Indicator
	IDApplicationTemplate                    = "61" // (M) Application Template
	IDCommonDataTemplate                     = "62" // (O) Common Data Template
	IDApplicationSpecificTransparentTemplate = "63" // (O) Application Specific Transparent Template
	IDCommonDataTransparentTemplate          = "64" // (O) Common Data Transparent Template
)

// const ...
const (
	TagApplicationDefinitionFileName = "4F" // Application Definition File (ADF) Name
	TagApplicationLabel              = "50"
	TagTrack2EquivalentData          = "57"
	TagApplicationPAN                = "5A"
	TagCardholderName                = "5F20"
	TagLanguagePreference            = "5F2D"
	TagIssuerURL                     = "5F50"
	TagApplicationVersionNumber      = "9F08"
	TagTokenRequestorID              = "9F19"
	TagPaymentAccountReference       = "9F24"
	TagLast4DigitsOfPAN              = "9F25"
)

// EMVQR ...
type EMVQR struct {
	DataPayloadFormatIndicator string              // 85
	ApplicationTemplate        ApplicationTemplate // 61
	CommonDataTemplate         CommonDataTemplate  // 62
}

// ApplicationTemplate ...
type ApplicationTemplate struct {
	BERTLV
	// ApplicationSpecificTransparentTemplate ApplicationSpecificTransparentTemplate // 63
}

// CommonDataTemplate ...
type CommonDataTemplate struct {
	CommonDataTransparentTemplate CommonDataTransparentTemplate // 64
}

// CommonDataTransparentTemplate ...
type CommonDataTransparentTemplate struct {
	BERTLV
}

// ApplicationSpecificTransparentTemplate ...
type ApplicationSpecificTransparentTemplate struct {
	BERTLV
}

// BERTLV ...
type BERTLV struct {
	DataApplicationDefinitionFileName string // "4F"
	DataApplicationLabel              string // "50"
	DataTrack2EquivalentData          string // "57"
	DataApplicationPAN                string // "5A"
	DataCardholderName                string // "5F20"
	DataLanguagePreference            string // "5F2D"
	DataIssuerURL                     string // "5F50"
	DataApplicationVersionNumber      string // "9F08"
	DataTokenRequestorID              string // "9F19"
	DataPaymentAccountReference       string // "9F24"
	DataLast4DigitsOfPAN              string // "9F25"
}

// GeneratePayload ...
func (c *EMVQR) GeneratePayload() (string, error) {
	s := ""
	if c.DataPayloadFormatIndicator != "" {
		s += format(IDPayloadFormatIndicator, toHex(c.DataPayloadFormatIndicator))
	} else {
		return "", fmt.Errorf("DataPayloadFormatIndicator is mandatory")
	}

	if (ApplicationTemplate{}) != c.ApplicationTemplate {
		//t := c.ApplicationTemplate
		tt := c.ApplicationTemplate
		template := ""
		if tt.DataApplicationDefinitionFileName != "" {
			template += format(TagApplicationDefinitionFileName, toBinary(tt.DataApplicationDefinitionFileName))
		}
		if tt.DataApplicationLabel != "" {
			template += format(TagApplicationLabel, toHex(tt.DataApplicationLabel))
		}
		if tt.DataTrack2EquivalentData != "" {
			template += format(TagTrack2EquivalentData, toBinary(tt.DataTrack2EquivalentData))
		}
		if tt.DataApplicationPAN != "" {
			template += format(TagApplicationPAN, tt.DataApplicationPAN)
		}
		if tt.DataCardholderName != "" {
			template += format(TagCardholderName, toHex(tt.DataCardholderName))
		}
		if tt.DataLanguagePreference != "" {
			template += format(TagLanguagePreference, toHex(tt.DataLanguagePreference))
		}
		if tt.DataIssuerURL != "" {
			template += format(TagIssuerURL, toHex(tt.DataIssuerURL))
		}
		if tt.DataApplicationVersionNumber != "" {
			template += format(TagApplicationVersionNumber, toHex(tt.DataApplicationVersionNumber))
		}
		if tt.DataTokenRequestorID != "" {
			template += format(TagTokenRequestorID, toHex(tt.DataTokenRequestorID))
		}
		if tt.DataPaymentAccountReference != "" {
			template += format(TagPaymentAccountReference, toHex(tt.DataPaymentAccountReference))
		}
		if tt.DataLast4DigitsOfPAN != "" {
			template += format(TagLast4DigitsOfPAN, toHex(tt.DataLast4DigitsOfPAN))
		}
		//fmt.Println(template)
		//ts := format(IDApplicationSpecificTransparentTemplate, template)
		//fmt.Println(ts)
		s += format(IDApplicationTemplate, template)
		fmt.Println(s)
	}
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	s = base64.StdEncoding.EncodeToString([]byte(string(decoded)))
	return s, nil
}

func format(id, value string) string {
	length := utf8.RuneCountInString(value) / 2
	lengthStr := strconv.Itoa(length)
	fmt.Println(lengthStr)
	fmt.Printf("%X", lengthStr)
	lengthStr = "00" + fmt.Sprintf("%X", length)
	return id + lengthStr[len(lengthStr)-2:] + value
}

func formatStr(id, value string) string {
	length := utf8.RuneCountInString(value)
	lengthStr := strconv.Itoa(length)
	lengthStr = "00" + lengthStr
	return id + lengthStr[len(lengthStr)-2:] + value
}

func toBinary(s string) string {
	return s
}

func toHex(s string) string {
	src := []byte(s)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return string(dst)
	// fmt.Printf("%s\n", dst)
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
