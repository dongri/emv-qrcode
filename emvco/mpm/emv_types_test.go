package mpm

import (
	"reflect"
	"testing"
)

func TestID_String(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want string
	}{
		{
			id:   ID("01"),
			want: "01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.String(); got != tt.want {
				t.Errorf("ID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestID_ParseInt(t *testing.T) {
	tests := []struct {
		name    string
		id      ID
		want    int64
		wantErr bool
	}{
		{
			name:    "ok",
			id:      ID("00"),
			want:    0,
			wantErr: false,
		},
		{
			name:    "not number",
			id:      ID("ab"),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.id.ParseInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("ID.ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ID.ParseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestID_Equal(t *testing.T) {
	type args struct {
		val ID
	}
	tests := []struct {
		name string
		id   ID
		args args
		want bool
	}{
		{
			name: "equal",
			id:   ID("01"),
			args: args{
				val: ID("01"),
			},
			want: true,
		},
		{
			name: "not equal",
			id:   ID("01"),
			args: args{
				val: ID("02"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.Equal(tt.args.val); got != tt.want {
				t.Errorf("ID.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestID_Between(t *testing.T) {
	type args struct {
		start ID
		end   ID
	}
	tests := []struct {
		name    string
		id      ID
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "within",
			id:   ID("01"),
			args: args{
				start: ID("00"),
				end:   ID("02"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "within equal start",
			id:   ID("01"),
			args: args{
				start: ID("01"),
				end:   ID("02"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "within equal end",
			id:   ID("01"),
			args: args{
				start: ID("00"),
				end:   ID("01"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "within equal start and end",
			id:   ID("01"),
			args: args{
				start: ID("01"),
				end:   ID("01"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "without smaller than start",
			id:   ID("00"),
			args: args{
				start: ID("01"),
				end:   ID("02"),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "without larger than end",
			id:   ID("03"),
			args: args{
				start: ID("01"),
				end:   ID("02"),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "not number id",
			id:   ID("ab"),
			args: args{
				start: ID("01"),
				end:   ID("02"),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "not number start",
			id:   ID("00"),
			args: args{
				start: ID("ab"),
				end:   ID("02"),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "not number end",
			id:   ID("00"),
			args: args{
				start: ID("01"),
				end:   ID("ab"),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.id.Between(tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("ID.Between() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ID.Between() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseEMVQR(t *testing.T) {
	type args struct {
		payload string
	}
	tests := []struct {
		name    string
		args    args
		want    *EMVQR
		wantErr bool
	}{
		{
			name: "empty payload",
			args: args{
				payload: "",
			},
			want: &EMVQR{},
		},
		{
			name: "id parse error",
			args: args{
				payload: "ab",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "value parse error",
			args: args{
				payload: "00020",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parse payload format indicator",
			args: args{
				payload: "000201",
			},
			want: &EMVQR{
				PayloadFormatIndicator: TLV{"00", "02", "01"},
			},
			wantErr: false,
		},
		{
			name: "parse point of initiation method",
			args: args{
				payload: "010211",
			},
			want: &EMVQR{
				PointOfInitiationMethod: TLV{"01", "02", PointOfInitiationMethodStatic},
			},
			wantErr: false,
		},
		{
			name: "parse merchant category code",
			args: args{
				payload: "52044111",
			},
			want: &EMVQR{
				MerchantCategoryCode: TLV{IDMerchantCategoryCode, "04", "4111"},
			},
			wantErr: false,
		},
		{
			name: "parse transaction currency",
			args: args{
				payload: "5303156",
			},
			want: &EMVQR{
				TransactionCurrency: TLV{IDTransactionCurrency, "03", "156"},
			},
			wantErr: false,
		},
		{
			name: "parse transaction amount",
			args: args{
				payload: "540523.72",
			},
			want: &EMVQR{
				TransactionAmount: TLV{IDTransactionAmount, "05", "23.72"},
			},
			wantErr: false,
		},
		{
			name: "parse tip or convenience indicator",
			args: args{
				payload: "550201",
			},
			want: &EMVQR{
				TipOrConvenienceIndicator: TLV{IDTipOrConvenienceIndicator, "02", "01"},
			},
			wantErr: false,
		},
		{
			name: "parse value of convenience fee fixed",
			args: args{
				payload: "5603500",
			},
			want: &EMVQR{
				ValueOfConvenienceFeeFixed: TLV{IDValueOfConvenienceFeeFixed, "03", "500"},
			},
			wantErr: false,
		},
		{
			name: "parse value of convenience fee percentage",
			args: args{
				payload: "57015",
			},
			want: &EMVQR{
				ValueOfConvenienceFeePercentage: TLV{IDValueOfConvenienceFeePercentage, "01", "5"},
			},
			wantErr: false,
		},
		{
			name: "parse country code",
			args: args{
				payload: "5802CN",
			},
			want: &EMVQR{
				CountryCode: TLV{IDCountryCode, "02", "CN"},
			},
			wantErr: false,
		},
		{
			name: "parse merchant name",
			args: args{
				payload: "5914BEST TRANSPORT",
			},
			want: &EMVQR{
				MerchantName: TLV{IDMerchantName, "14", "BEST TRANSPORT"},
			},
			wantErr: false,
		},
		{
			name: "parse merchant city",
			args: args{
				payload: "6007BEIJING",
			},
			want: &EMVQR{
				MerchantCity: TLV{IDMerchantCity, "07", "BEIJING"},
			},
			wantErr: false,
		},
		{
			name: "parse postal code",
			args: args{
				payload: "61071234567",
			},
			want: &EMVQR{
				PostalCode: TLV{IDPostalCode, "07", "1234567"},
			},
			wantErr: false,
		},
		{
			name: "parse crc",
			args: args{
				payload: "6304A13A",
			},
			want: &EMVQR{
				CRC: TLV{IDCRC, "04", "A13A"},
			},
			wantErr: false,
		},
		{
			name: "parse additional data field template",
			args: args{
				payload: "6233030412340603***0708A60086670902ME",
			},
			want: &EMVQR{
				AdditionalDataFieldTemplate: &AdditionalDataFieldTemplate{
					StoreLabel:                    TLV{AdditionalIDStoreLabel, "04", "1234"},
					CustomerLabel:                 TLV{AdditionalIDCustomerLabel, "03", "***"},
					TerminalLabel:                 TLV{AdditionalIDTerminalLabel, "08", "A6008667"},
					AdditionalConsumerDataRequest: TLV{AdditionalIDAdditionalConsumerDataRequest, "02", "ME"},
				},
			},
			wantErr: false,
		},
		{
			name: "parse failed additional data field template",
			args: args{
				payload: "6231030412340603***0708A60086670902", // not enough length
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parse merchant information language template",
			args: args{
				payload: "64200002ZH0104最佳运输0202北京",
			},
			want: &EMVQR{
				MerchantInformationLanguageTemplate: &MerchantInformationLanguageTemplate{
					LanguagePreference: TLV{MerchantInformationIDLanguagePreference, "02", "ZH"},
					MerchantName:       TLV{MerchantInformationIDMerchantName, "04", "最佳运输"},
					MerchantCity:       TLV{MerchantInformationIDMerchantCity, "02", "北京"},
				},
			},
			wantErr: false,
		},
		{
			name: "parse failed merchant information language template",
			args: args{
				payload: "64180002ZH0104最佳运输0202", // not enough length
			},
			want:    nil,
			wantErr: true,
		},
		// {
		// 	name: "parse merchant account information",
		// 	args: args{
		// 		payload: "02160004hoge0104abcd",
		// 	},
		// 	want: &EMVQR{
		// 		MerchantAccountInformation: map[ID]*MerchantAccountInformation{
		// 			ID("02"): &MerchantAccountInformation{
		// 				Value: "0004hoge0104abcd",
		// 			},
		// 		},
		// 	},
		// 	wantErr: false,
		// },
		{
			name: "parse failed merchant account information",
			args: args{
				payload: "02140004hoge0104", // not enough length
			},
			want:    nil,
			wantErr: true,
		},
		// {
		// 	name: "parse multiple merchant account information",
		// 	args: args{
		// 		payload: "02160004hoge0104abcd26160004fuga0204efgh",
		// 	},
		// 	want: &EMVQR{
		// 		MerchantAccountInformation: map[ID]*MerchantAccountInformation{
		// 			ID("02"): &MerchantAccountInformation{
		// 				Value: "0004hoge0104abcd",
		// 			},
		// 			ID("26"): &MerchantAccountInformation{
		// 				Value: "0004fuga0204efgh",
		// 			},
		// 		},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "parse RFU for EMVCo",
		// 	args: args{
		// 		payload: "6504abcd",
		// 	},
		// 	want: &EMVQR{
		// 		RFUforEMVCo: map[ID]*RFUforEMVCo{
		// 			ID("65"): &RFUforEMVCo{
		// 				Value: "abcd",
		// 			},
		// 		},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "parse multiple RFU for EMVCo",
		// 	args: args{
		// 		payload: "6504abcd7904efgh",
		// 	},
		// 	want: &EMVQR{
		// 		RFUforEMVCo: map[ID]*RFUforEMVCo{
		// 			ID("65"): &RFUforEMVCo{
		// 				Value: "abcd",
		// 			},
		// 			ID("79"): &RFUforEMVCo{
		// 				Value: "efgh",
		// 			},
		// 		},
		// 	},
		// 	wantErr: false,
		// },
		// 	{
		// 		name: "parse unreserved templates",
		// 		args: args{
		// 			payload: "8004abcd",
		// 		},
		// 		want: &EMVQR{
		// 			UnreservedTemplates: map[ID]*UnreservedTemplate{
		// 				ID("80"): &UnreservedTemplate{
		// 					Value: "abcd",
		// 				},
		// 			},
		// 		},
		// 		wantErr: false,
		// 	},
		// 	{
		// 		name: "parse multiple unreserved templates",
		// 		args: args{
		// 			payload: "8004abcd9904efgh",
		// 		},
		// 		want: &EMVQR{
		// 			UnreservedTemplates: map[ID]*UnreservedTemplate{
		// 				ID("80"): &UnreservedTemplate{
		// 					Value: "abcd",
		// 				},
		// 				ID("99"): &UnreservedTemplate{
		// 					Value: "efgh",
		// 				},
		// 			},
		// 		},
		// 		wantErr: false,
		// 	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEMVQR(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEMVQR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseEMVQR() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestEMVQR_GeneratePayload(t *testing.T) {
// 	type fields struct {
// 		PayloadFormatIndicator              string
// 		PointOfInitiationMethod             PointOfInitiationMethod
// 		MerchantAccountInformation          map[ID]*MerchantAccountInformation
// 		MerchantCategoryCode                string
// 		TransactionCurrency                 string
// 		TransactionAmount                   string
// 		TipOrConvenienceIndicator           string
// 		ValueOfConvenienceFeeFixed          string
// 		ValueOfConvenienceFeePercentage     string
// 		CountryCode                         string
// 		MerchantName                        string
// 		MerchantCity                        string
// 		PostalCode                          string
// 		AdditionalDataFieldTemplate         *AdditionalDataFieldTemplate
// 		CRC                                 string
// 		MerchantInformationLanguageTemplate *MerchantInformationLanguageTemplate
// 		RFUforEMVCo                         map[ID]*RFUforEMVCo
// 		UnreservedTemplates                 map[ID]*UnreservedTemplate
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			name:    "empty EMVQR",
// 			fields:  fields{},
// 			want:    "0000" + formatCrc("0000"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify payload format indicator",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 			},
// 			want:    "000201" + formatCrc("000201"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify point of initiation method",
// 			fields: fields{
// 				PointOfInitiationMethod: PointOfInitiationMethodStatic,
// 			},
// 			want:    "0000" + "010211" + formatCrc("0000"+"010211"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify merchant category code",
// 			fields: fields{
// 				MerchantCategoryCode: "4111",
// 			},
// 			want:    "0000" + "52044111" + formatCrc("0000"+"52044111"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify transaction currency",
// 			fields: fields{
// 				TransactionCurrency: "156",
// 			},
// 			want:    "0000" + "5303156" + formatCrc("0000"+"5303156"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify transaction amount",
// 			fields: fields{
// 				TransactionAmount: "23.72",
// 			},
// 			want:    "0000" + "540523.72" + formatCrc("0000"+"540523.72"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify tip or convenience indicator",
// 			fields: fields{
// 				TipOrConvenienceIndicator: "01",
// 			},
// 			want:    "0000" + "550201" + formatCrc("0000"+"550201"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify value of convenience fee fixed",
// 			fields: fields{
// 				ValueOfConvenienceFeeFixed: "500",
// 			},
// 			want:    "0000" + "5603500" + formatCrc("0000"+"5603500"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify value of convenience fee percentage",
// 			fields: fields{
// 				ValueOfConvenienceFeePercentage: "5",
// 			},
// 			want:    "0000" + "57015" + formatCrc("0000"+"57015"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify country code",
// 			fields: fields{
// 				CountryCode: "CN",
// 			},
// 			want:    "0000" + "5802CN" + formatCrc("0000"+"5802CN"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify merchant name",
// 			fields: fields{
// 				MerchantName: "BEST TRANSPORT",
// 			},
// 			want:    "0000" + "5914BEST TRANSPORT" + formatCrc("0000"+"5914BEST TRANSPORT"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify merchant city",
// 			fields: fields{
// 				MerchantCity: "BEIJING",
// 			},
// 			want:    "0000" + "6007BEIJING" + formatCrc("0000"+"6007BEIJING"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify postal code",
// 			fields: fields{
// 				PostalCode: "1234567",
// 			},
// 			want:    "0000" + "61071234567" + formatCrc("0000"+"61071234567"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify additional data field template",
// 			fields: fields{
// 				AdditionalDataFieldTemplate: &AdditionalDataFieldTemplate{
// 					StoreLabel:                    "1234",
// 					CustomerLabel:                 "***",
// 					TerminalLabel:                 "A6008667",
// 					AdditionalConsumerDataRequest: "ME",
// 				},
// 			},
// 			want:    "0000" + "6233030412340603***0708A60086670902ME" + formatCrc("0000"+"6233030412340603***0708A60086670902ME"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify merchant information language template",
// 			fields: fields{
// 				MerchantInformationLanguageTemplate: &MerchantInformationLanguageTemplate{
// 					LanguagePreference: "ZH",
// 					MerchantName:       "最佳运输",
// 					MerchantCity:       "北京",
// 				},
// 			},
// 			want:    "0000" + "64200002ZH0104最佳运输0202北京" + formatCrc("0000"+"64200002ZH0104最佳运输0202北京"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify merchant account information",
// 			fields: fields{
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{
// 						Value: "0004hoge0104abcd",
// 					},
// 				},
// 			},
// 			want:    "0000" + "02160004hoge0104abcd" + formatCrc("0000"+"02160004hoge0104abcd"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify RFU for EMVCo",
// 			fields: fields{
// 				RFUforEMVCo: map[ID]*RFUforEMVCo{
// 					ID("65"): &RFUforEMVCo{
// 						Value: "abcd",
// 					},
// 				},
// 			},
// 			want:    "0000" + "6504abcd" + formatCrc("0000"+"6504abcd"),
// 			wantErr: false,
// 		},
// 		{
// 			name: "stringify unreserved templates",
// 			fields: fields{
// 				UnreservedTemplates: map[ID]*UnreservedTemplate{
// 					ID("80"): &UnreservedTemplate{
// 						Value: "abcd",
// 					},
// 				},
// 			},
// 			want:    "0000" + "8004abcd" + formatCrc("0000"+"8004abcd"),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &EMVQR{
// 				PayloadFormatIndicator:          tt.fields.PayloadFormatIndicator,
// 				PointOfInitiationMethod:         tt.fields.PointOfInitiationMethod,
// 				MerchantAccountInformation:      tt.fields.MerchantAccountInformation,
// 				MerchantCategoryCode:            tt.fields.MerchantCategoryCode,
// 				TransactionCurrency:             tt.fields.TransactionCurrency,
// 				TransactionAmount:               tt.fields.TransactionAmount,
// 				TipOrConvenienceIndicator:       tt.fields.TipOrConvenienceIndicator,
// 				ValueOfConvenienceFeeFixed:      tt.fields.ValueOfConvenienceFeeFixed,
// 				ValueOfConvenienceFeePercentage: tt.fields.ValueOfConvenienceFeePercentage,
// 				CountryCode:                     tt.fields.CountryCode,
// 				MerchantName:                    tt.fields.MerchantName,
// 				MerchantCity:                    tt.fields.MerchantCity,
// 				PostalCode:                      tt.fields.PostalCode,
// 				AdditionalDataFieldTemplate:     tt.fields.AdditionalDataFieldTemplate,
// 				CRC: tt.fields.CRC,
// 				MerchantInformationLanguageTemplate: tt.fields.MerchantInformationLanguageTemplate,
// 				RFUforEMVCo:                         tt.fields.RFUforEMVCo,
// 				UnreservedTemplates:                 tt.fields.UnreservedTemplates,
// 			}
// 			got, err := c.GeneratePayload()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("EMVQR.GeneratePayload() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("EMVQR.GeneratePayload() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestEMVQR_Validate(t *testing.T) {
// 	type fields struct {
// 		PayloadFormatIndicator              string
// 		PointOfInitiationMethod             PointOfInitiationMethod
// 		MerchantAccountInformation          map[ID]*MerchantAccountInformation
// 		MerchantCategoryCode                string
// 		TransactionCurrency                 string
// 		TransactionAmount                   string
// 		TipOrConvenienceIndicator           string
// 		ValueOfConvenienceFeeFixed          string
// 		ValueOfConvenienceFeePercentage     string
// 		CountryCode                         string
// 		MerchantName                        string
// 		MerchantCity                        string
// 		PostalCode                          string
// 		AdditionalDataFieldTemplate         *AdditionalDataFieldTemplate
// 		CRC                                 string
// 		MerchantInformationLanguageTemplate *MerchantInformationLanguageTemplate
// 		RFUforEMVCo                         map[ID]*RFUforEMVCo
// 		UnreservedTemplates                 map[ID]*UnreservedTemplate
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{
// 			name: "minimum ok",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode: "1443",
// 				TransactionCurrency:  "354",
// 				CountryCode:          "JP",
// 				MerchantName:         "Sample",
// 				MerchantCity:         "tokyo",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:    "lack of PayloadFormatIndicator",
// 			fields:  fields{},
// 			wantErr: true,
// 		},
// 		{
// 			name: "lack of MerchantAccountInformation",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "lack of MerchantCategoryCode",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "lack of TransactionCurrency",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode: "1443",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "lack of CountryCode",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode: "1443",
// 				TransactionCurrency:  "354",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "lack of MerchantName",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode: "1443",
// 				TransactionCurrency:  "354",
// 				CountryCode:          "JP",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "lack of MerchantCity",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode: "1443",
// 				TransactionCurrency:  "354",
// 				CountryCode:          "JP",
// 				MerchantName:         "Sample",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "PointOfInitiationMethod is unknown",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode:    "1443",
// 				TransactionCurrency:     "354",
// 				CountryCode:             "JP",
// 				MerchantName:            "Sample",
// 				MerchantCity:            "tokyo",
// 				PointOfInitiationMethod: "00", // should be 11 or 12
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "PointOfInitiationMethod is unknown",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode:    "1443",
// 				TransactionCurrency:     "354",
// 				CountryCode:             "JP",
// 				MerchantName:            "Sample",
// 				MerchantCity:            "tokyo",
// 				PointOfInitiationMethod: "00", // should be 11 or 12
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "failed validate merchant information language template",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode:                "1443",
// 				TransactionCurrency:                 "354",
// 				CountryCode:                         "JP",
// 				MerchantName:                        "Sample",
// 				MerchantCity:                        "tokyo",
// 				MerchantInformationLanguageTemplate: &MerchantInformationLanguageTemplate{},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "exist merchant account information",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode: "1443",
// 				TransactionCurrency:  "354",
// 				CountryCode:          "JP",
// 				MerchantName:         "Sample",
// 				MerchantCity:         "tokyo",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "exist additional data field template",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode:        "1443",
// 				TransactionCurrency:         "354",
// 				CountryCode:                 "JP",
// 				MerchantName:                "Sample",
// 				MerchantCity:                "tokyo",
// 				AdditionalDataFieldTemplate: &AdditionalDataFieldTemplate{},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "exist RFU for EMVCo",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode: "1443",
// 				TransactionCurrency:  "354",
// 				CountryCode:          "JP",
// 				MerchantName:         "Sample",
// 				MerchantCity:         "tokyo",
// 				RFUforEMVCo: map[ID]*RFUforEMVCo{
// 					ID("79"): &RFUforEMVCo{},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "exist unreserved template",
// 			fields: fields{
// 				PayloadFormatIndicator: "01",
// 				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
// 					ID("02"): &MerchantAccountInformation{},
// 				},
// 				MerchantCategoryCode: "1443",
// 				TransactionCurrency:  "354",
// 				CountryCode:          "JP",
// 				MerchantName:         "Sample",
// 				MerchantCity:         "tokyo",
// 				UnreservedTemplates: map[ID]*UnreservedTemplate{
// 					ID("99"): &UnreservedTemplate{},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &EMVQR{
// 				PayloadFormatIndicator:          tt.fields.PayloadFormatIndicator,
// 				PointOfInitiationMethod:         tt.fields.PointOfInitiationMethod,
// 				MerchantAccountInformation:      tt.fields.MerchantAccountInformation,
// 				MerchantCategoryCode:            tt.fields.MerchantCategoryCode,
// 				TransactionCurrency:             tt.fields.TransactionCurrency,
// 				TransactionAmount:               tt.fields.TransactionAmount,
// 				TipOrConvenienceIndicator:       tt.fields.TipOrConvenienceIndicator,
// 				ValueOfConvenienceFeeFixed:      tt.fields.ValueOfConvenienceFeeFixed,
// 				ValueOfConvenienceFeePercentage: tt.fields.ValueOfConvenienceFeePercentage,
// 				CountryCode:                     tt.fields.CountryCode,
// 				MerchantName:                    tt.fields.MerchantName,
// 				MerchantCity:                    tt.fields.MerchantCity,
// 				PostalCode:                      tt.fields.PostalCode,
// 				AdditionalDataFieldTemplate:     tt.fields.AdditionalDataFieldTemplate,
// 				CRC: tt.fields.CRC,
// 				MerchantInformationLanguageTemplate: tt.fields.MerchantInformationLanguageTemplate,
// 				RFUforEMVCo:                         tt.fields.RFUforEMVCo,
// 				UnreservedTemplates:                 tt.fields.UnreservedTemplates,
// 			}
// 			if err := c.Validate(); (err != nil) != tt.wantErr {
// 				t.Errorf("EMVQR.Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestPointOfInitiationMethod_IsStaticMethod(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		m    PointOfInitiationMethod
// 		want bool
// 	}{
// 		{
// 			name: "is static",
// 			m:    PointOfInitiationMethod("11"),
// 			want: true,
// 		},
// 		{
// 			name: "is not static",
// 			m:    PointOfInitiationMethod("12"),
// 			want: false,
// 		},
// 		{
// 			name: "undefined",
// 			m:    PointOfInitiationMethod(""),
// 			want: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.m.IsStaticMethod(); got != tt.want {
// 				t.Errorf("PointOfInitiationMethod.IsStaticMethod() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPointOfInitiationMethod_IsDynamicMethod(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		m    PointOfInitiationMethod
// 		want bool
// 	}{
// 		{
// 			name: "is dynamic",
// 			m:    PointOfInitiationMethod("12"),
// 			want: true,
// 		},
// 		{
// 			name: "is not dynamic",
// 			m:    PointOfInitiationMethod("11"),
// 			want: false,
// 		},
// 		{
// 			name: "undefined",
// 			m:    PointOfInitiationMethod(""),
// 			want: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.m.IsDynamicMethod(); got != tt.want {
// 				t.Errorf("PointOfInitiationMethod.IsDynamicMethod() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPointOfInitiationMethod_GeneratePayload(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		m    PointOfInitiationMethod
// 		want string
// 	}{
// 		{
// 			name: "ok",
// 			m:    PointOfInitiationMethod("11"),
// 			want: "010211",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.m.GeneratePayload(); got != tt.want {
// 				t.Errorf("PointOfInitiationMethod.GeneratePayload() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPointOfInitiationMethod_Validate(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		c       PointOfInitiationMethod
// 		wantErr bool
// 	}{
// 		{
// 			name:    "unknown method",
// 			c:       PointOfInitiationMethod("00"),
// 			wantErr: true,
// 		},
// 		{
// 			name:    "static method",
// 			c:       PointOfInitiationMethod("11"),
// 			wantErr: false,
// 		},
// 		{
// 			name:    "dynamic method",
// 			c:       PointOfInitiationMethod("12"),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := tt.c.Validate(); (err != nil) != tt.wantErr {
// 				t.Errorf("PointOfInitiationMethod.Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func Test_ParseRFUforEMVCo(t *testing.T) {
// 	type args struct {
// 		value string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *RFUforEMVCo
// 	}{
// 		{
// 			args: args{
// 				value: "abcd",
// 			},
// 			want: &RFUforEMVCo{
// 				Value: "abcd",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := ParseRFUforEMVCo(tt.args.value); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ParseRFUforEMVCo() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestRFUforEMVCo_Validate(t *testing.T) {
// 	type fields struct {
// 		Value string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{
// 			fields: fields{
// 				Value: "abcd",
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &RFUforEMVCo{
// 				Value: tt.fields.Value,
// 			}
// 			if err := c.Validate(); (err != nil) != tt.wantErr {
// 				t.Errorf("RFUforEMVCo.Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func Test_ParseUnreservedTemplate(t *testing.T) {
// 	type args struct {
// 		id    ID
// 		value string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *UnreservedTemplate
// 	}{
// 		{
// 			args: args{
// 				value: "abcd",
// 			},
// 			want: &UnreservedTemplate{
// 				Value: "abcd",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := ParseUnreservedTemplate(tt.args.value); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ParseUnreservedTemplate() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestUnreservedTemplate_Validate(t *testing.T) {
// 	type fields struct {
// 		Value string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{
// 			fields: fields{
// 				Value: "abcd",
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &UnreservedTemplate{
// 				Value: tt.fields.Value,
// 			}
// 			if err := c.Validate(); (err != nil) != tt.wantErr {
// 				t.Errorf("UnreservedTemplate.Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func Test_ParseMerchantAccountInformationTemplate(t *testing.T) {
// 	type args struct {
// 		payload string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *MerchantAccountInformation
// 	}{
// 		{
// 			name: "empty payload",
// 			args: args{
// 				payload: "",
// 			},
// 			want: &MerchantAccountInformation{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := ParseMerchantAccountInformation(tt.args.payload)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ParseMerchantAccountInformation() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMerchantAccountInformation_Validate(t *testing.T) {
// 	type fields struct {
// 		Value string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{
// 			fields: fields{
// 				Value: "abcd",
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &MerchantAccountInformation{
// 				Value: tt.fields.Value,
// 			}
// 			if err := c.Validate(); (err != nil) != tt.wantErr {
// 				t.Errorf("MerchantAccountInformation.Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func Test_ParseAdditionalDataFieldTemplate(t *testing.T) {
// 	type args struct {
// 		payload string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *AdditionalDataFieldTemplate
// 		wantErr bool
// 	}{
// 		{
// 			name: "empty payload",
// 			args: args{
// 				payload: "",
// 			},
// 			want:    &AdditionalDataFieldTemplate{},
// 			wantErr: false,
// 		},
// 		{
// 			name: "id parse error",
// 			args: args{
// 				payload: "ab",
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "value parse error",
// 			args: args{
// 				payload: "00020",
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "parse bill number",
// 			args: args{
// 				payload: "010200",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				BillNumber: "00",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse mobile number",
// 			args: args{
// 				payload: "020200",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				MobileNumber: "00",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse store label",
// 			args: args{
// 				payload: "030200",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				StoreLabel: "00",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse loyalty nubmer",
// 			args: args{
// 				payload: "040200",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				LoyaltyNumber: "00",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse reference label",
// 			args: args{
// 				payload: "050200",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				ReferenceLabel: "00",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse customer label",
// 			args: args{
// 				payload: "060200",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				CustomerLabel: "00",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse terminal label",
// 			args: args{
// 				payload: "070200",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				TerminalLabel: "00",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse purpose transaction",
// 			args: args{
// 				payload: "080200",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				PurposeTransaction: "00",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse additional consumer data request",
// 			args: args{
// 				payload: "0902ME",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				AdditionalConsumerDataRequest: "ME",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse payment system specific",
// 			args: args{
// 				payload: "50160004hoge0104abcd",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				PaymentSystemSpecific: map[ID]*PaymentSystemSpecific{
// 					ID("50"): &PaymentSystemSpecific{
// 						Value: "0004hoge0104abcd",
// 					},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse failed payment system specific",
// 			args: args{
// 				payload: "50140004hoge0104", // not enough length
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "parse multiple payment system specific",
// 			args: args{
// 				payload: "50160004hoge0104abcd99160004fuga0304efgh",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				PaymentSystemSpecific: map[ID]*PaymentSystemSpecific{
// 					ID("50"): &PaymentSystemSpecific{
// 						Value: "0004hoge0104abcd",
// 					},
// 					ID("99"): &PaymentSystemSpecific{
// 						Value: "0004fuga0304efgh",
// 					},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse RFU for EMVCo",
// 			args: args{
// 				payload: "1004abcd",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				RFUforEMVCo: map[ID]*RFUforEMVCo{
// 					ID("10"): &RFUforEMVCo{Value: "abcd"},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse multiple RFU for EMVCo",
// 			args: args{
// 				payload: "1004abcd4904efgh",
// 			},
// 			want: &AdditionalDataFieldTemplate{
// 				RFUforEMVCo: map[ID]*RFUforEMVCo{
// 					ID("10"): &RFUforEMVCo{Value: "abcd"},
// 					ID("49"): &RFUforEMVCo{Value: "efgh"},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := ParseAdditionalDataFieldTemplate(tt.args.payload)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ParseAdditionalDataFieldTemplate() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ParseAdditionalDataFieldTemplate() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestAdditionalDataFieldTemplate_GeneratePayload(t *testing.T) {
// 	type fields struct {
// 		BillNumber                    string
// 		MobileNumber                  string
// 		StoreLabel                    string
// 		LoyaltyNumber                 string
// 		ReferenceLabel                string
// 		CustomerLabel                 string
// 		TerminalLabel                 string
// 		PurposeTransaction            string
// 		AdditionalConsumerDataRequest string
// 		RFUforEMVCo                   map[ID]*RFUforEMVCo
// 		PaymentSystemSpecific         map[ID]*PaymentSystemSpecific
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   string
// 	}{
// 		{
// 			name: "stringify bill number",
// 			fields: fields{
// 				BillNumber: "00",
// 			},
// 			want: "010200",
// 		},
// 		{
// 			name: "stringify mobile number",
// 			fields: fields{
// 				MobileNumber: "00",
// 			},
// 			want: "020200",
// 		},
// 		{
// 			name: "stringify store label",
// 			fields: fields{
// 				StoreLabel: "00",
// 			},
// 			want: "030200",
// 		},
// 		{
// 			name: "stringify loyalty nubmer",
// 			fields: fields{
// 				LoyaltyNumber: "00",
// 			},
// 			want: "040200",
// 		},
// 		{
// 			name: "stringify reference label",
// 			fields: fields{
// 				ReferenceLabel: "00",
// 			},
// 			want: "050200",
// 		},
// 		{
// 			name: "stringify customer label",
// 			fields: fields{
// 				CustomerLabel: "00",
// 			},
// 			want: "060200",
// 		},
// 		{
// 			name: "stringify terminal label",
// 			fields: fields{
// 				TerminalLabel: "00",
// 			},
// 			want: "070200",
// 		},
// 		{
// 			name: "stringify purpose transaction",
// 			fields: fields{
// 				PurposeTransaction: "00",
// 			},
// 			want: "080200",
// 		},
// 		{
// 			name: "stringify additional consumer data request",
// 			fields: fields{
// 				AdditionalConsumerDataRequest: "ME",
// 			},
// 			want: "0902ME",
// 		},
// 		{
// 			name: "stringify payment system specific",
// 			fields: fields{
// 				PaymentSystemSpecific: map[ID]*PaymentSystemSpecific{
// 					ID("50"): &PaymentSystemSpecific{
// 						Value: "0004hoge0104abcd",
// 					},
// 				},
// 			},
// 			want: "50160004hoge0104abcd",
// 		},
// 		{
// 			name: "stringify RFU for EMVCo",
// 			fields: fields{
// 				RFUforEMVCo: map[ID]*RFUforEMVCo{
// 					ID("10"): &RFUforEMVCo{Value: "abcd"},
// 				},
// 			},
// 			want: "1004abcd",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &AdditionalDataFieldTemplate{
// 				BillNumber:                    tt.fields.BillNumber,
// 				MobileNumber:                  tt.fields.MobileNumber,
// 				StoreLabel:                    tt.fields.StoreLabel,
// 				LoyaltyNumber:                 tt.fields.LoyaltyNumber,
// 				ReferenceLabel:                tt.fields.ReferenceLabel,
// 				CustomerLabel:                 tt.fields.CustomerLabel,
// 				TerminalLabel:                 tt.fields.TerminalLabel,
// 				PurposeTransaction:            tt.fields.PurposeTransaction,
// 				AdditionalConsumerDataRequest: tt.fields.AdditionalConsumerDataRequest,
// 				RFUforEMVCo:                   tt.fields.RFUforEMVCo,
// 				PaymentSystemSpecific:         tt.fields.PaymentSystemSpecific,
// 			}
// 			if got := c.GeneratePayload(); got != tt.want {
// 				t.Errorf("AdditionalDataFieldTemplate.GeneratePayload() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestAdditionalDataFieldTemplate_Validate(t *testing.T) {
// 	type fields struct {
// 		BillNumber                    string
// 		MobileNumber                  string
// 		StoreLabel                    string
// 		LoyaltyNumber                 string
// 		ReferenceLabel                string
// 		CustomerLabel                 string
// 		TerminalLabel                 string
// 		PurposeTransaction            string
// 		AdditionalConsumerDataRequest string
// 		RFUforEMVCo                   map[ID]*RFUforEMVCo
// 		PaymentSystemSpecific         map[ID]*PaymentSystemSpecific
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{
// 			fields:  fields{},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &AdditionalDataFieldTemplate{
// 				BillNumber:                    tt.fields.BillNumber,
// 				MobileNumber:                  tt.fields.MobileNumber,
// 				StoreLabel:                    tt.fields.StoreLabel,
// 				LoyaltyNumber:                 tt.fields.LoyaltyNumber,
// 				ReferenceLabel:                tt.fields.ReferenceLabel,
// 				CustomerLabel:                 tt.fields.CustomerLabel,
// 				TerminalLabel:                 tt.fields.TerminalLabel,
// 				PurposeTransaction:            tt.fields.PurposeTransaction,
// 				AdditionalConsumerDataRequest: tt.fields.AdditionalConsumerDataRequest,
// 				RFUforEMVCo:                   tt.fields.RFUforEMVCo,
// 				PaymentSystemSpecific:         tt.fields.PaymentSystemSpecific,
// 			}
// 			if err := c.Validate(); (err != nil) != tt.wantErr {
// 				t.Errorf("AdditionalDataFieldTemplate.Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func Test_ParsePaymentSystemSpecific(t *testing.T) {
// 	type args struct {
// 		value string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *PaymentSystemSpecific
// 	}{
// 		{
// 			args: args{
// 				value: "abcd",
// 			},
// 			want: &PaymentSystemSpecific{
// 				Value: "abcd",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := ParsePaymentSystemSpecific(tt.args.value); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ParsePaymentSystemSpecific() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPaymentSystemSpecific_GeneratePayload(t *testing.T) {
// 	type fields struct {
// 		Value string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   string
// 	}{
// 		{
// 			fields: fields{
// 				Value: "abcd",
// 			},
// 			want: "abcd",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			v := &PaymentSystemSpecific{
// 				Value: tt.fields.Value,
// 			}
// 			if got := v.GeneratePayload(); got != tt.want {
// 				t.Errorf("PaymentSystemSpecific.GeneratePayload() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_ParseMerchantInformationLanguageTemplate(t *testing.T) {
// 	type args struct {
// 		payload string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *MerchantInformationLanguageTemplate
// 		wantErr bool
// 	}{
// 		{
// 			name: "empty payload",
// 			args: args{
// 				payload: "",
// 			},
// 			want:    &MerchantInformationLanguageTemplate{},
// 			wantErr: false,
// 		},
// 		{
// 			name: "id parse error",
// 			args: args{
// 				payload: "ab",
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "value parse error",
// 			args: args{
// 				payload: "00020",
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "parse language preference",
// 			args: args{
// 				payload: "0002JP",
// 			},
// 			want: &MerchantInformationLanguageTemplate{
// 				LanguagePreference: "JP",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse merchant name",
// 			args: args{
// 				payload: "0106sample",
// 			},
// 			want: &MerchantInformationLanguageTemplate{
// 				MerchantName: "sample",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse merchant name",
// 			args: args{
// 				payload: "0205TOKYO",
// 			},
// 			want: &MerchantInformationLanguageTemplate{
// 				MerchantCity: "TOKYO",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "parse RFU for EMVCo",
// 			args: args{
// 				payload: "1004abcd",
// 			},
// 			want: &MerchantInformationLanguageTemplate{
// 				RFUForEMVCo: map[ID]*RFUforEMVCo{
// 					ID("10"): &RFUforEMVCo{
// 						Value: "abcd",
// 					},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := ParseMerchantInformationLanguageTemplate(tt.args.payload)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ParseMerchantInformationLanguageTemplate() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ParseMerchantInformationLanguageTemplate() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMerchantInformationLanguageTemplate_GeneratePayload(t *testing.T) {
// 	type fields struct {
// 		LanguagePreference string
// 		MerchantName       string
// 		MerchantCity       string
// 		RFUForEMVCo        map[ID]*RFUforEMVCo
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   string
// 	}{
// 		{
// 			name: "stringify language preference",
// 			fields: fields{
// 				LanguagePreference: "JP",
// 			},
// 			want: "0002JP",
// 		},
// 		{
// 			name: "stringify merchant name",
// 			fields: fields{
// 				MerchantName: "sample",
// 			},
// 			want: "0106sample",
// 		},
// 		{
// 			name: "stringify merchant city",
// 			fields: fields{
// 				MerchantCity: "TOKYO",
// 			},
// 			want: "0205TOKYO",
// 		},
// 		{
// 			name: "stringify RFU for EMVCo",
// 			fields: fields{
// 				RFUForEMVCo: map[ID]*RFUforEMVCo{
// 					ID("10"): &RFUforEMVCo{
// 						Value: "abcd",
// 					},
// 				},
// 			},
// 			want: "1004abcd",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &MerchantInformationLanguageTemplate{
// 				LanguagePreference: tt.fields.LanguagePreference,
// 				MerchantName:       tt.fields.MerchantName,
// 				MerchantCity:       tt.fields.MerchantCity,
// 				RFUForEMVCo:        tt.fields.RFUForEMVCo,
// 			}
// 			if got := c.GeneratePayload(); got != tt.want {
// 				t.Errorf("MerchantInformationLanguageTemplate.GeneratePayload() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMerchantInformationLanguageTemplate_Validate(t *testing.T) {
// 	type fields struct {
// 		LanguagePreference string
// 		MerchantName       string
// 		MerchantCity       string
// 		RFUForEMVCo        map[ID]*RFUforEMVCo
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{
// 			name: "minumum ok",
// 			fields: fields{
// 				LanguagePreference: "JA",
// 				MerchantName:       "sample",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:    "lack of LanguagePreference",
// 			fields:  fields{},
// 			wantErr: true,
// 		},
// 		{
// 			name: "lack of MechantName",
// 			fields: fields{
// 				LanguagePreference: "JA",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "exist RFU for EMVCo",
// 			fields: fields{
// 				LanguagePreference: "JA",
// 				MerchantName:       "sample",
// 				RFUForEMVCo: map[ID]*RFUforEMVCo{
// 					ID("99"): &RFUforEMVCo{},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &MerchantInformationLanguageTemplate{
// 				LanguagePreference: tt.fields.LanguagePreference,
// 				MerchantName:       tt.fields.MerchantName,
// 				MerchantCity:       tt.fields.MerchantCity,
// 				RFUForEMVCo:        tt.fields.RFUForEMVCo,
// 			}
// 			if err := c.Validate(); (err != nil) != tt.wantErr {
// 				t.Errorf("MerchantInformationLanguageTemplate.Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
