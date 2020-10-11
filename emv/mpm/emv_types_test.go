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

func TestEMVQR_GeneratePayload(t *testing.T) {
	type fields struct {
		PayloadFormatIndicator              TLV
		PointOfInitiationMethod             TLV
		MerchantAccountInformation          map[ID]MerchantAccountInformationTLV
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
		MerchantInformationLanguageTemplate *MerchantInformationLanguageTemplate
		RFUforEMVCo                         []TLV
		UnreservedTemplates                 map[ID]UnreservedTemplateTLV
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "stringify payload format indicator",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    "00",
					Length: "02",
					Value:  "01",
				},
			},
			want:    "000201" + formatCrc("000201"),
			wantErr: false,
		},
		{
			name: "stringify point of initiation method",
			fields: fields{
				PointOfInitiationMethod: TLV{
					Tag:    "01",
					Length: "02",
					Value:  "11",
				},
			},
			want:    "010211" + formatCrc("010211"),
			wantErr: false,
		},
		{
			name: "stringify merchant category code",
			fields: fields{
				MerchantCategoryCode: TLV{
					Tag:    "52",
					Length: "04",
					Value:  "4111",
				},
			},
			want:    "52044111" + formatCrc("52044111"),
			wantErr: false,
		},
		{
			name: "stringify transaction currency",
			fields: fields{
				TransactionCurrency: TLV{
					Tag:    "53",
					Length: "03",
					Value:  "156",
				},
			},
			want:    "5303156" + formatCrc("5303156"),
			wantErr: false,
		},
		{
			name: "stringify transaction amount",
			fields: fields{
				TransactionAmount: TLV{
					Tag:    "54",
					Length: "05",
					Value:  "23.72",
				},
			},
			want:    "540523.72" + formatCrc("540523.72"),
			wantErr: false,
		},
		{
			name: "stringify tip or convenience indicator",
			fields: fields{
				TipOrConvenienceIndicator: TLV{
					Tag:    "55",
					Length: "02",
					Value:  "01",
				},
			},
			want:    "550201" + formatCrc("550201"),
			wantErr: false,
		},
		{
			name: "stringify value of convenience fee fixed",
			fields: fields{
				ValueOfConvenienceFeeFixed: TLV{
					Tag:    "56",
					Length: "03",
					Value:  "500",
				},
			},
			want:    "5603500" + formatCrc("5603500"),
			wantErr: false,
		},
		{
			name: "stringify value of convenience fee percentage",
			fields: fields{
				ValueOfConvenienceFeePercentage: TLV{
					Tag:    "57",
					Length: "01",
					Value:  "5",
				},
			},
			want:    "57015" + formatCrc("57015"),
			wantErr: false,
		},
		{
			name: "stringify country code",
			fields: fields{
				CountryCode: TLV{
					Tag:    "58",
					Length: "02",
					Value:  "CN",
				},
			},
			want:    "5802CN" + formatCrc("5802CN"),
			wantErr: false,
		},
		{
			name: "stringify merchant name",
			fields: fields{
				MerchantName: TLV{
					Tag:    "59",
					Length: "14",
					Value:  "BEST TRANSPORT",
				},
			},
			want:    "5914BEST TRANSPORT" + formatCrc("5914BEST TRANSPORT"),
			wantErr: false,
		},
		{
			name: "stringify merchant city",
			fields: fields{
				MerchantCity: TLV{
					Tag:    "60",
					Length: "07",
					Value:  "BEIJING",
				},
			},
			want:    "6007BEIJING" + formatCrc("6007BEIJING"),
			wantErr: false,
		},
		{
			name: "stringify postal code",
			fields: fields{
				PostalCode: TLV{
					Tag:    "61",
					Length: "07",
					Value:  "1234567",
				},
			},
			want:    "61071234567" + formatCrc("61071234567"),
			wantErr: false,
		},
		{
			name: "stringify additional data field template",
			fields: fields{
				AdditionalDataFieldTemplate: &AdditionalDataFieldTemplate{
					StoreLabel: TLV{
						Tag:    "03",
						Length: "04",
						Value:  "1234",
					},
					CustomerLabel: TLV{
						Tag:    "06",
						Length: "03",
						Value:  "***",
					},
					TerminalLabel: TLV{
						Tag:    "07",
						Length: "08",
						Value:  "A6008667",
					},
					AdditionalConsumerDataRequest: TLV{
						Tag:    "09",
						Length: "02",
						Value:  "ME",
					},
				},
			},
			want:    "6233030412340603***0708A60086670902ME" + formatCrc("6233030412340603***0708A60086670902ME"),
			wantErr: false,
		},
		{
			name: "stringify merchant information language template",
			fields: fields{
				MerchantInformationLanguageTemplate: &MerchantInformationLanguageTemplate{
					LanguagePreference: TLV{
						Tag:    "00",
						Length: "02",
						Value:  "ZH",
					},
					MerchantName: TLV{
						Tag:    "01",
						Length: "04",
						Value:  "最佳运输",
					},
					MerchantCity: TLV{
						Tag:    "02",
						Length: "02",
						Value:  "北京",
					},
				},
			},
			want:    "64200002ZH0104最佳运输0202北京" + formatCrc("64200002ZH0104最佳运输0202北京"),
			wantErr: false,
		},
		{
			name: "stringify merchant account information",
			fields: fields{
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{
						Tag:    "02",
						Length: "32",
						Value: &MerchantAccountInformation{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "04",
								Value:  "hoge",
							},
							PaymentNetworkSpecific: []TLV{
								TLV{
									Tag:    "01",
									Length: "04",
									Value:  "abcd",
								},
								TLV{
									Tag:    "15",
									Length: "04",
									Value:  "efgh",
								},
								TLV{
									Tag:    "13",
									Length: "04",
									Value:  "ijkl",
								},
							},
						},
					},
				},
			},
			want:    "02320004hoge0104abcd1304ijkl1504efgh" + formatCrc("02320004hoge0104abcd1304ijkl1504efgh"),
			wantErr: false,
		},
		{
			name: "stringify RFU for EMVCo",
			fields: fields{
				RFUforEMVCo: []TLV{
					TLV{
						Tag:    "65",
						Length: "04",
						Value:  "abcd",
					},
				},
			},
			want:    "6504abcd" + formatCrc("6504abcd"),
			wantErr: false,
		},
		{
			name: "stringify unreserved templates",
			fields: fields{
				UnreservedTemplates: map[ID]UnreservedTemplateTLV{
					ID("80"): UnreservedTemplateTLV{
						Tag:    "80",
						Length: "16",
						Value: &UnreservedTemplate{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "04",
								Value:  "abcd",
							},
							ContextSpecificData: []TLV{
								TLV{
									Tag:    "01",
									Length: "04",
									Value:  "efgh",
								},
							},
						},
					},
				},
			},
			want:    "80160004abcd0104efgh" + formatCrc("80160004abcd0104efgh"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EMVQR{
				PayloadFormatIndicator:          tt.fields.PayloadFormatIndicator,
				PointOfInitiationMethod:         tt.fields.PointOfInitiationMethod,
				MerchantAccountInformation:      tt.fields.MerchantAccountInformation,
				MerchantCategoryCode:            tt.fields.MerchantCategoryCode,
				TransactionCurrency:             tt.fields.TransactionCurrency,
				TransactionAmount:               tt.fields.TransactionAmount,
				TipOrConvenienceIndicator:       tt.fields.TipOrConvenienceIndicator,
				ValueOfConvenienceFeeFixed:      tt.fields.ValueOfConvenienceFeeFixed,
				ValueOfConvenienceFeePercentage: tt.fields.ValueOfConvenienceFeePercentage,
				CountryCode:                     tt.fields.CountryCode,
				MerchantName:                    tt.fields.MerchantName,
				MerchantCity:                    tt.fields.MerchantCity,
				PostalCode:                      tt.fields.PostalCode,
				AdditionalDataFieldTemplate:     tt.fields.AdditionalDataFieldTemplate,
				CRC: tt.fields.CRC,
				MerchantInformationLanguageTemplate: tt.fields.MerchantInformationLanguageTemplate,
				RFUforEMVCo:                         tt.fields.RFUforEMVCo,
				UnreservedTemplates:                 tt.fields.UnreservedTemplates,
			}
			got := c.GeneratePayload()
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("EMVQR.GeneratePayload() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			if got != tt.want {
				t.Errorf("EMVQR.GeneratePayload() = %v, want %v", got, tt.want)
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
		{
			name: "parse merchant account information",
			args: args{
				payload: "02160004hoge0104abcd",
			},
			want: &EMVQR{
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{
						Tag:    "02",
						Length: "16",
						Value: &MerchantAccountInformation{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "04",
								Value:  "hoge",
							},
							PaymentNetworkSpecific: []TLV{
								TLV{
									Tag:    "01",
									Length: "04",
									Value:  "abcd",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "parse failed merchant account information",
			args: args{
				payload: "02140004hoge0104", // not enough length
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parse multiple merchant account information",
			args: args{
				payload: "02160004hoge0104abcd26160004fuga0204efgh",
			},
			want: &EMVQR{
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{
						Tag:    "02",
						Length: "16",
						Value: &MerchantAccountInformation{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "04",
								Value:  "hoge",
							},
							PaymentNetworkSpecific: []TLV{
								TLV{
									Tag:    "01",
									Length: "04",
									Value:  "abcd",
								},
							},
						},
					},
					ID("26"): MerchantAccountInformationTLV{
						Tag:    "26",
						Length: "16",
						Value: &MerchantAccountInformation{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "04",
								Value:  "fuga",
							},
							PaymentNetworkSpecific: []TLV{
								TLV{
									Tag:    "02",
									Length: "04",
									Value:  "efgh",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "parse RFU for EMVCo",
			args: args{
				payload: "6504abcd",
			},
			want: &EMVQR{
				RFUforEMVCo: []TLV{
					TLV{
						Tag:    "65",
						Length: "04",
						Value:  "abcd",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "parse multiple RFU for EMVCo",
			args: args{
				payload: "6504abcd7904efgh",
			},
			want: &EMVQR{
				RFUforEMVCo: []TLV{
					TLV{
						Tag:    "65",
						Length: "04",
						Value:  "abcd",
					},
					TLV{
						Tag:    "79",
						Length: "04",
						Value:  "efgh",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "parse unreserved templates",
			args: args{
				payload: "80160004hoge0104abcd",
			},
			want: &EMVQR{
				UnreservedTemplates: map[ID]UnreservedTemplateTLV{
					ID("80"): UnreservedTemplateTLV{
						Tag:    "80",
						Length: "16",
						Value: &UnreservedTemplate{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "04",
								Value:  "hoge",
							},
							ContextSpecificData: []TLV{
								TLV{
									Tag:    "01",
									Length: "04",
									Value:  "abcd",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "parse multiple unreserved templates",
			args: args{
				payload: "80240004hoge0104abcd0204efgh",
			},
			want: &EMVQR{
				UnreservedTemplates: map[ID]UnreservedTemplateTLV{
					ID("80"): UnreservedTemplateTLV{
						Tag:    "80",
						Length: "24",
						Value: &UnreservedTemplate{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "04",
								Value:  "hoge",
							},
							ContextSpecificData: []TLV{
								TLV{
									Tag:    "01",
									Length: "04",
									Value:  "abcd",
								},
								TLV{
									Tag:    "02",
									Length: "04",
									Value:  "efgh",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
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

func TestEMVQR_Validate(t *testing.T) {
	type fields struct {
		PayloadFormatIndicator              TLV
		PointOfInitiationMethod             TLV
		MerchantAccountInformation          map[ID]MerchantAccountInformationTLV
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
		MerchantInformationLanguageTemplate *MerchantInformationLanguageTemplate
		RFUforEMVCo                         []TLV
		UnreservedTemplates                 map[ID]UnreservedTemplateTLV
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "minimum ok",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "354",
				},
				CountryCode: TLV{
					Tag:    IDCountryCode,
					Length: "02",
					Value:  "JP",
				},
				MerchantName: TLV{
					Tag:    IDMerchantName,
					Length: "06",
					Value:  "Sample",
				},
				MerchantCity: TLV{
					Tag:    IDMerchantCity,
					Length: "05",
					Value:  "tokyo",
				},
			},
			wantErr: false,
		},
		{
			name:    "lack of PayloadFormatIndicator",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "lack of MerchantAccountInformation",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
			},
			wantErr: true,
		},
		{
			name: "lack of MerchantCategoryCode",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
			},
			wantErr: true,
		},
		{
			name: "lack of TransactionCurrency",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
			},
			wantErr: true,
		},
		{
			name: "lack of CountryCode",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "354",
				},
			},
			wantErr: true,
		},
		{
			name: "lack of MerchantName",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "354",
				},
				CountryCode: TLV{
					Tag:    IDCountryCode,
					Length: "02",
					Value:  "JP",
				},
			},
			wantErr: true,
		},
		{
			name: "lack of MerchantCity",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "354",
				},
				CountryCode: TLV{
					Tag:    IDCountryCode,
					Length: "02",
					Value:  "JP",
				},
				MerchantName: TLV{
					Tag:    IDMerchantName,
					Length: "06",
					Value:  "Sample",
				},
			},
			wantErr: true,
		},
		{
			name: "PointOfInitiationMethod is unknown",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "354",
				},
				CountryCode: TLV{
					Tag:    IDCountryCode,
					Length: "02",
					Value:  "JP",
				},
				MerchantName: TLV{
					Tag:    IDMerchantName,
					Length: "06",
					Value:  "Sample",
				},
				PointOfInitiationMethod: TLV{
					Tag:    IDPointOfInitiationMethod,
					Length: "02",
					Value:  "00", // should be 11 or 12
				},
			},
			wantErr: true,
		},
		{
			name: "failed validate merchant information language template",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "354",
				},
				CountryCode: TLV{
					Tag:    IDCountryCode,
					Length: "02",
					Value:  "JP",
				},
				MerchantName: TLV{
					Tag:    IDMerchantName,
					Length: "06",
					Value:  "Sample",
				},
				MerchantCity: TLV{
					Tag:    IDMerchantCity,
					Length: "05",
					Value:  "tokyo",
				},
				MerchantInformationLanguageTemplate: &MerchantInformationLanguageTemplate{},
			},
			wantErr: true,
		},
		{
			name: "exist merchant account information",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "354",
				},
				CountryCode: TLV{
					Tag:    IDCountryCode,
					Length: "02",
					Value:  "JP",
				},
				MerchantName: TLV{
					Tag:    IDMerchantName,
					Length: "06",
					Value:  "Sample",
				},
				MerchantCity: TLV{
					Tag:    IDMerchantCity,
					Length: "05",
					Value:  "tokyo",
				},
			},
			wantErr: false,
		},
		{
			name: "exist additional data field template",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "354",
				},
				CountryCode: TLV{
					Tag:    IDCountryCode,
					Length: "02",
					Value:  "JP",
				},
				MerchantName: TLV{
					Tag:    IDMerchantName,
					Length: "06",
					Value:  "Sample",
				},
				MerchantCity: TLV{
					Tag:    IDMerchantCity,
					Length: "05",
					Value:  "tokyo",
				},
				AdditionalDataFieldTemplate: &AdditionalDataFieldTemplate{},
			},
			wantErr: false,
		},
		{
			name: "exist RFU for EMVCo",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "354",
				},
				CountryCode: TLV{
					Tag:    IDCountryCode,
					Length: "02",
					Value:  "JP",
				},
				MerchantName: TLV{
					Tag:    IDMerchantName,
					Length: "06",
					Value:  "Sample",
				},
				MerchantCity: TLV{
					Tag:    IDMerchantCity,
					Length: "05",
					Value:  "tokyo",
				},
				RFUforEMVCo: []TLV{},
			},
			wantErr: false,
		},
		{
			name: "exist unreserved template",
			fields: fields{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("02"): MerchantAccountInformationTLV{},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "1443",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "354",
				},
				CountryCode: TLV{
					Tag:    IDCountryCode,
					Length: "02",
					Value:  "JP",
				},
				MerchantName: TLV{
					Tag:    IDMerchantName,
					Length: "06",
					Value:  "Sample",
				},
				MerchantCity: TLV{
					Tag:    IDMerchantCity,
					Length: "05",
					Value:  "tokyo",
				},
				UnreservedTemplates: map[ID]UnreservedTemplateTLV{
					ID("80"): UnreservedTemplateTLV{
						Tag:    "80",
						Length: "16",
						Value: &UnreservedTemplate{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "04",
								Value:  "abcd",
							},
							ContextSpecificData: []TLV{
								TLV{
									Tag:    "01",
									Length: "04",
									Value:  "efgh",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EMVQR{
				PayloadFormatIndicator:          tt.fields.PayloadFormatIndicator,
				PointOfInitiationMethod:         tt.fields.PointOfInitiationMethod,
				MerchantAccountInformation:      tt.fields.MerchantAccountInformation,
				MerchantCategoryCode:            tt.fields.MerchantCategoryCode,
				TransactionCurrency:             tt.fields.TransactionCurrency,
				TransactionAmount:               tt.fields.TransactionAmount,
				TipOrConvenienceIndicator:       tt.fields.TipOrConvenienceIndicator,
				ValueOfConvenienceFeeFixed:      tt.fields.ValueOfConvenienceFeeFixed,
				ValueOfConvenienceFeePercentage: tt.fields.ValueOfConvenienceFeePercentage,
				CountryCode:                     tt.fields.CountryCode,
				MerchantName:                    tt.fields.MerchantName,
				MerchantCity:                    tt.fields.MerchantCity,
				PostalCode:                      tt.fields.PostalCode,
				AdditionalDataFieldTemplate:     tt.fields.AdditionalDataFieldTemplate,
				CRC: tt.fields.CRC,
				MerchantInformationLanguageTemplate: tt.fields.MerchantInformationLanguageTemplate,
				RFUforEMVCo:                         tt.fields.RFUforEMVCo,
				UnreservedTemplates:                 tt.fields.UnreservedTemplates,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("EMVQR.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ParseMerchantAccountInformationTemplate(t *testing.T) {
	type args struct {
		payload string
	}
	tests := []struct {
		name string
		args args
		want *MerchantAccountInformation
	}{
		{
			name: "empty payload",
			args: args{
				payload: "",
			},
			want: &MerchantAccountInformation{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ParseMerchantAccountInformation(tt.args.payload)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMerchantAccountInformation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseAdditionalDataFieldTemplate(t *testing.T) {
	type args struct {
		payload string
	}
	tests := []struct {
		name    string
		args    args
		want    *AdditionalDataFieldTemplate
		wantErr bool
	}{
		{
			name: "empty payload",
			args: args{
				payload: "",
			},
			want:    &AdditionalDataFieldTemplate{},
			wantErr: false,
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
			name: "parse bill number",
			args: args{
				payload: "010200",
			},
			want: &AdditionalDataFieldTemplate{
				BillNumber: TLV{
					Tag:    AdditionalIDBillNumber,
					Length: "02",
					Value:  "00",
				},
			},
			wantErr: false,
		},
		{
			name: "parse mobile number",
			args: args{
				payload: "020200",
			},
			want: &AdditionalDataFieldTemplate{
				MobileNumber: TLV{
					Tag:    AdditionalIDMobileNumber,
					Length: "02",
					Value:  "00",
				},
			},
			wantErr: false,
		},
		{
			name: "parse store label",
			args: args{
				payload: "030200",
			},
			want: &AdditionalDataFieldTemplate{
				StoreLabel: TLV{
					Tag:    AdditionalIDStoreLabel,
					Length: "02",
					Value:  "00",
				},
			},
			wantErr: false,
		},
		{
			name: "parse loyalty nubmer",
			args: args{
				payload: "040200",
			},
			want: &AdditionalDataFieldTemplate{
				LoyaltyNumber: TLV{
					Tag:    AdditionalIDLoyaltyNumber,
					Length: "02",
					Value:  "00",
				},
			},
			wantErr: false,
		},
		{
			name: "parse reference label",
			args: args{
				payload: "050200",
			},
			want: &AdditionalDataFieldTemplate{
				ReferenceLabel: TLV{
					Tag:    AdditionalIDReferenceLabel,
					Length: "02",
					Value:  "00",
				},
			},
			wantErr: false,
		},
		{
			name: "parse customer label",
			args: args{
				payload: "060200",
			},
			want: &AdditionalDataFieldTemplate{
				CustomerLabel: TLV{
					Tag:    AdditionalIDCustomerLabel,
					Length: "02",
					Value:  "00",
				},
			},
			wantErr: false,
		},
		{
			name: "parse terminal label",
			args: args{
				payload: "070200",
			},
			want: &AdditionalDataFieldTemplate{
				TerminalLabel: TLV{
					Tag:    AdditionalIDTerminalLabel,
					Length: "02",
					Value:  "00",
				},
			},
			wantErr: false,
		},
		{
			name: "parse purpose transaction",
			args: args{
				payload: "080200",
			},
			want: &AdditionalDataFieldTemplate{
				PurposeTransaction: TLV{
					Tag:    AdditionalIDPurposeTransaction,
					Length: "02",
					Value:  "00",
				},
			},
			wantErr: false,
		},
		{
			name: "parse additional consumer data request",
			args: args{
				payload: "0902ME",
			},
			want: &AdditionalDataFieldTemplate{
				AdditionalConsumerDataRequest: TLV{
					Tag:    "09",
					Length: "02",
					Value:  "ME",
				},
			},
			wantErr: false,
		},
		{
			name: "parse payment system specific",
			args: args{
				payload: "50160004hoge0104abcd",
			},
			want: &AdditionalDataFieldTemplate{
				PaymentSystemSpecific: []TLV{
					TLV{
						Tag:    "50",
						Length: "16",
						Value:  "0004hoge0104abcd",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "parse RFU for EMVCo",
			args: args{
				payload: "1004abcd",
			},
			want: &AdditionalDataFieldTemplate{
				RFUforEMVCo: []TLV{
					TLV{
						Tag:    "10",
						Length: "04",
						Value:  "abcd",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "parse multiple RFU for EMVCo",
			args: args{
				payload: "1004abcd4904efgh",
			},
			want: &AdditionalDataFieldTemplate{
				RFUforEMVCo: []TLV{
					TLV{
						Tag:    "10",
						Length: "04",
						Value:  "abcd",
					},
					TLV{
						Tag:    "49",
						Length: "04",
						Value:  "efgh",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAdditionalDataFieldTemplate(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAdditionalDataFieldTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAdditionalDataFieldTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseMerchantInformationLanguageTemplate(t *testing.T) {
	type args struct {
		payload string
	}
	tests := []struct {
		name    string
		args    args
		want    *MerchantInformationLanguageTemplate
		wantErr bool
	}{
		{
			name: "empty payload",
			args: args{
				payload: "",
			},
			want:    &MerchantInformationLanguageTemplate{},
			wantErr: false,
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
			name: "parse language preference",
			args: args{
				payload: "0002JP",
			},
			want: &MerchantInformationLanguageTemplate{
				LanguagePreference: TLV{
					Tag:    MerchantInformationIDLanguagePreference,
					Length: "02",
					Value:  "JP",
				},
			},
			wantErr: false,
		},
		{
			name: "parse merchant name",
			args: args{
				payload: "0106sample",
			},
			want: &MerchantInformationLanguageTemplate{
				MerchantName: TLV{
					Tag:    MerchantInformationIDMerchantName,
					Length: "06",
					Value:  "sample",
				},
			},
			wantErr: false,
		},
		{
			name: "parse merchant name",
			args: args{
				payload: "0205TOKYO",
			},
			want: &MerchantInformationLanguageTemplate{
				MerchantCity: TLV{
					Tag:    MerchantInformationIDMerchantCity,
					Length: "05",
					Value:  "TOKYO",
				},
			},
			wantErr: false,
		},
		{
			name: "parse RFU for EMVCo",
			args: args{
				payload: "1004abcd",
			},
			want: &MerchantInformationLanguageTemplate{
				RFUforEMVCo: []TLV{
					TLV{
						Tag:    "10",
						Length: "04",
						Value:  "abcd",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMerchantInformationLanguageTemplate(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMerchantInformationLanguageTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMerchantInformationLanguageTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
