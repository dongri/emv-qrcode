package mpm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dongri/emv-qrcode/crc16"
)

// Encode ...
func Encode(emvqr *EMVQR) (string, error) {
	if err := emvqr.Validate(); err != nil {
		return "", err
	}
	return emvqr.GeneratePayload(), nil
}

// Decode ...
func Decode(payload string) (*EMVQR, error) {
	emvqr, err := ParseEMVQR(payload)
	if err != nil {
		return nil, err
	}
	if err := checkCRC(payload, emvqr.CRC.Value); err != nil {
		return nil, err
	}
	if err := emvqr.Validate(); err != nil {
		return emvqr, err
	}
	return emvqr, nil
}

func checkCRC(payload string, crc string) error {
	if crc == "" {
		return errors.New("CRC is mandatory")
	}
	if len(payload) < 8 {
		return errors.New("payload too short to contain CRC field")
	}

	crc = strings.ToUpper(crc)

	// Spec (EMV QRCODE 4.1): CRC shall be the last data object.
	// checkCRC relies on this — it hashes payload[:len-4] assuming the
	// last 4 chars are always the CRC hex digits.
	if strings.ToUpper(payload[len(payload)-8:len(payload)-4]) != "6304" {
		return errors.New("CRC field (6304) must be the last field in the payload")
	}
	if strings.ToUpper(payload[len(payload)-4:]) != crc {
		return errors.New("CRC field (6304) must be the last field in the payload")
	}

	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	computed := crc16.Checksum([]byte(payload[:len(payload)-4]), table)
	computedS := fmt.Sprintf("%04X", computed)

	if computedS != crc {
		return fmt.Errorf("invalid CRC: computed %s, got %s", computedS, crc)
	}

	return nil
}
