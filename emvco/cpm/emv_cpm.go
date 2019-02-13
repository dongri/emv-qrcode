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
			s += format(IDApplicationTemplate, template)
		}
	}

	if len(c.CommonDataTemplates) > 0 {
		for _, t := range c.CommonDataTemplates {
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

			if len(t.CommonDataTransparentTemplates) > 0 {
				for _, tt := range t.CommonDataTransparentTemplates {
					ttemplate := ""
					if tt.DataApplicationDefinitionFileName != "" {
						ttemplate += format(TagApplicationDefinitionFileName, tt.DataApplicationDefinitionFileName)
					}
					if tt.DataApplicationLabel != "" {
						ttemplate += format(TagApplicationLabel, toHex(tt.DataApplicationLabel))
					}
					if tt.DataTrack2EquivalentData != "" {
						ttemplate += format(TagTrack2EquivalentData, tt.DataTrack2EquivalentData)
					}
					if tt.DataApplicationPAN != "" {
						ttemplate += format(TagApplicationPAN, tt.DataApplicationPAN)
					}
					if tt.DataCardholderName != "" {
						ttemplate += format(TagCardholderName, toHex(tt.DataCardholderName))
					}
					if tt.DataLanguagePreference != "" {
						ttemplate += format(TagLanguagePreference, toHex(tt.DataLanguagePreference))
					}
					if tt.DataIssuerURL != "" {
						ttemplate += format(TagIssuerURL, toHex(tt.DataIssuerURL))
					}
					if tt.DataApplicationVersionNumber != "" {
						ttemplate += format(TagApplicationVersionNumber, tt.DataApplicationVersionNumber)
					}
					if tt.DataIssuerApplicationData != "" {
						ttemplate += format(TagIssuerApplicationData, tt.DataIssuerApplicationData)
					}
					if tt.DataTokenRequestorID != "" {
						ttemplate += format(TagTokenRequestorID, tt.DataTokenRequestorID)
					}
					if tt.DataPaymentAccountReference != "" {
						ttemplate += format(TagPaymentAccountReference, tt.DataPaymentAccountReference)
					}
					if tt.DataLast4DigitsOfPAN != "" {
						ttemplate += format(TagLast4DigitsOfPAN, tt.DataLast4DigitsOfPAN)
					}
					if tt.DataApplicationCryptogram != "" {
						ttemplate += format(TagApplicationCryptogram, tt.DataApplicationCryptogram)
					}
					if tt.DataApplicationTransactionCounter != "" {
						ttemplate += format(TagApplicationTransactionCounter, tt.DataApplicationTransactionCounter)
					}
					if tt.DataUnpredictableNumber != "" {
						ttemplate += format(TagUnpredictableNumber, tt.DataUnpredictableNumber)
					}
					template += format(IDCommonDataTransparentTemplate, ttemplate)
				}
				s += format(IDCommonDataTemplate, template)
			}
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

//fmt.Println(template)
//ts := format(IDApplicationSpecificTransparentTemplate, template)
//fmt.Println(ts)
