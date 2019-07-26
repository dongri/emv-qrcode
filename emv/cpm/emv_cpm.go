package cpm

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
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
	TagApplicationDefinitionFileName = "4F"
	TagApplicationLabel              = "50"
	TagTrack2EquivalentData          = "57"
	TagApplicationPAN                = "5A"
	TagCardholderName                = "5F20"
	TagLanguagePreference            = "5F2D"
	TagIssuerURL                     = "5F50"
	TagApplicationVersionNumber      = "9F08"
	TagIssuerApplicationData         = "9F10"
	TagTokenRequestorID              = "9F19"
	TagPaymentAccountReference       = "9F24"
	TagLast4DigitsOfPAN              = "9F25"
	TagApplicationCryptogram         = "9F26"
	TagApplicationTransactionCounter = "9F36"
	TagUnpredictableNumber           = "9F37"
)

// EMVQR ...
type EMVQR struct {
	DataPayloadFormatIndicator string                // 85
	ApplicationTemplates       []ApplicationTemplate // 61
	CommonDataTemplates        []CommonDataTemplate  // 62
}

// ApplicationTemplate ...
type ApplicationTemplate struct {
	BERTLV
	ApplicationSpecificTransparentTemplates []ApplicationSpecificTransparentTemplate // 63
}

// CommonDataTemplate ...
type CommonDataTemplate struct {
	BERTLV
	CommonDataTransparentTemplates []CommonDataTransparentTemplate // 64
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
	DataIssuerApplicationData         string // "9F10"
	DataTokenRequestorID              string // "9F19"
	DataPaymentAccountReference       string // "9F24"
	DataLast4DigitsOfPAN              string // "9F25"
	DataApplicationCryptogram         string // "9F26"
	DataApplicationTransactionCounter string // "9F36"
	DataUnpredictableNumber           string // "9F37"
}

// GeneratePayload ...
func (c *EMVQR) GeneratePayload() (string, error) {
	s := ""
	if c.DataPayloadFormatIndicator != "" {
		s += format(IDPayloadFormatIndicator, toHex(c.DataPayloadFormatIndicator))
	} else {
		return "", fmt.Errorf("DataPayloadFormatIndicator is mandatory")
	}
	if len(c.ApplicationTemplates) > 0 {
		for _, t := range c.ApplicationTemplates {
			template := formattingTemplate((t.BERTLV))
			if len(t.ApplicationSpecificTransparentTemplates) > 0 {
				for _, tt := range t.ApplicationSpecificTransparentTemplates {
					ttemplate := formattingTemplate((tt.BERTLV))
					template += format(IDApplicationSpecificTransparentTemplate, ttemplate)
				}
			}
			s += format(IDApplicationTemplate, template)
		}
	}
	if len(c.CommonDataTemplates) > 0 {
		for _, t := range c.CommonDataTemplates {
			template := formattingTemplate(t.BERTLV)
			if len(t.CommonDataTransparentTemplates) > 0 {
				for _, tt := range t.CommonDataTransparentTemplates {
					ttemplate := formattingTemplate(tt.BERTLV)
					template += format(IDCommonDataTransparentTemplate, ttemplate)
				}
			}
			s += format(IDCommonDataTemplate, template)
		}
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
	lengthStr = "00" + fmt.Sprintf("%X", length)
	return id + lengthStr[len(lengthStr)-2:] + value
}

func toHex(s string) string {
	src := []byte(s)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return string(dst)
}

func formattingTemplate(t BERTLV) string {
	template := ""
	if t.DataApplicationDefinitionFileName != "" {
		template += format(TagApplicationDefinitionFileName, t.DataApplicationDefinitionFileName)
	}
	if t.DataApplicationLabel != "" {
		template += format(TagApplicationLabel, toHex(t.DataApplicationLabel))
	}
	if t.DataTrack2EquivalentData != "" {
		template += format(TagTrack2EquivalentData, t.DataTrack2EquivalentData)
	}
	if t.DataApplicationPAN != "" {
		template += format(TagApplicationPAN, t.DataApplicationPAN)
	}
	if t.DataCardholderName != "" {
		template += format(TagCardholderName, toHex(t.DataCardholderName))
	}
	if t.DataLanguagePreference != "" {
		template += format(TagLanguagePreference, toHex(t.DataLanguagePreference))
	}
	if t.DataIssuerURL != "" {
		template += format(TagIssuerURL, toHex(t.DataIssuerURL))
	}
	if t.DataApplicationVersionNumber != "" {
		template += format(TagApplicationVersionNumber, t.DataApplicationVersionNumber)
	}
	if t.DataIssuerApplicationData != "" {
		template += format(TagIssuerApplicationData, t.DataIssuerApplicationData)
	}
	if t.DataTokenRequestorID != "" {
		template += format(TagTokenRequestorID, t.DataTokenRequestorID)
	}
	if t.DataPaymentAccountReference != "" {
		template += format(TagPaymentAccountReference, t.DataPaymentAccountReference)
	}
	if t.DataLast4DigitsOfPAN != "" {
		template += format(TagLast4DigitsOfPAN, t.DataLast4DigitsOfPAN)
	}
	if t.DataApplicationCryptogram != "" {
		template += format(TagApplicationCryptogram, t.DataApplicationCryptogram)
	}
	if t.DataApplicationTransactionCounter != "" {
		template += format(TagApplicationTransactionCounter, t.DataApplicationTransactionCounter)
	}
	if t.DataUnpredictableNumber != "" {
		template += format(TagUnpredictableNumber, t.DataUnpredictableNumber)
	}
	return template
}
