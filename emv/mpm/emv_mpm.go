package mpm

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
	if err := emvqr.Validate(); err != nil {
		return nil, err
	}
	return emvqr, nil
}
