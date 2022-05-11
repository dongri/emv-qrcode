package mpm

import (
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		emvqr *EMVQR
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty EMVQR",
			args: args{
				emvqr: &EMVQR{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				emvqr: &EMVQR{
					PayloadFormatIndicator: TLV{
						Tag:    "00",
						Length: "02",
						Value:  "01",
					},
					PointOfInitiationMethod: TLV{
						Tag:    IDPointOfInitiationMethod,
						Length: "02",
						Value:  PointOfInitiationMethodDynamic,
					},
					MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
						ID("29"): {
							Tag:    "29",
							Length: "30",
							Value: &MerchantAccountInformation{
								GloballyUniqueIdentifier: TLV{
									Tag:    "00",
									Length: "12",
									Value:  "D15600000000",
								},
								PaymentNetworkSpecific: []TLV{
									{
										Tag:    "05",
										Length: "10",
										Value:  "A93FO3230Q",
									},
								},
							},
						},
					},
					MerchantCategoryCode: TLV{
						Tag:    IDMerchantCategoryCode,
						Length: "04",
						Value:  "4111",
					},
					TransactionCurrency: TLV{
						Tag:    IDTransactionCurrency,
						Length: "03",
						Value:  "156",
					},
					TransactionAmount: TLV{
						Tag:    IDTransactionAmount,
						Length: "05",
						Value:  "23.72",
					},
					CountryCode: TLV{
						Tag:    IDCountryCode,
						Length: "02",
						Value:  "CN",
					},
					MerchantName: TLV{
						Tag:    IDMerchantName,
						Length: "14",
						Value:  "BEST TRANSPORT",
					},
					MerchantCity: TLV{
						Tag:    IDMerchantCity,
						Length: "07",
						Value:  "BEIJING",
					},
					MerchantInformationLanguageTemplate: &MerchantInformationLanguageTemplate{
						LanguagePreference: TLV{
							Tag:    MerchantInformationIDLanguagePreference,
							Length: "02",
							Value:  "ZH",
						},
						MerchantName: TLV{
							Tag:    MerchantInformationIDMerchantName,
							Length: "04",
							Value:  "最佳运输",
						},
						MerchantCity: TLV{
							Tag:    MerchantInformationIDMerchantCity,
							Length: "02",
							Value:  "北京",
						},
					},
					TipOrConvenienceIndicator: TLV{
						Tag:    IDTipOrConvenienceIndicator,
						Length: "02",
						Value:  "01",
					},
					AdditionalDataFieldTemplate: &AdditionalDataFieldTemplate{
						StoreLabel: TLV{
							Tag:    AdditionalIDStoreLabel,
							Length: "04",
							Value:  "1234",
						},
						CustomerLabel: TLV{
							Tag:    AdditionalIDCustomerLabel,
							Length: "03",
							Value:  "***",
						},
						TerminalLabel: TLV{
							Tag:    AdditionalIDTerminalLabel,
							Length: "08",
							Value:  "A6008667",
						},
						AdditionalConsumerDataRequest: TLV{
							Tag:    "09",
							Length: "02",
							Value:  "ME",
						},
					},
					UnreservedTemplates: map[ID]UnreservedTemplateTLV{
						ID("91"): {
							Tag:    "91",
							Length: "32",
							Value: &UnreservedTemplate{
								GloballyUniqueIdentifier: TLV{
									Tag:    "00",
									Length: "16",
									Value:  "A011223344998877",
								},
								ContextSpecificData: []TLV{
									{
										Tag:    "07",
										Length: "08",
										Value:  "12345678",
									},
								},
							},
						},
					},
				},
			},
			want:    "00020101021229300012D156000000000510A93FO3230Q520441115303156540523.725502015802CN5914BEST TRANSPORT6007BEIJING6233030412340603***0708A60086670902ME64200002ZH0104最佳运输0202北京91320016A011223344998877070812345678" + formatCrc("00020101021229300012D156000000000510A93FO3230Q520441115303156540523.725502015802CN5914BEST TRANSPORT6007BEIJING6233030412340603***0708A60086670902ME64200002ZH0104最佳运输0202北京91320016A011223344998877070812345678"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.emvqr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
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
			want:    &EMVQR{},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				payload: "00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING64200002ZH0104最佳运输0202北京540523.7253031565502016233030412340603***0708A60086670902ME91320016A0112233449988770708123456786304A13A",
			},
			want: &EMVQR{
				PayloadFormatIndicator: TLV{
					Tag:    IDPayloadFormatIndicator,
					Length: "02",
					Value:  "01",
				},
				PointOfInitiationMethod: TLV{
					Tag:    IDPointOfInitiationMethod,
					Length: "02",
					Value:  PointOfInitiationMethodDynamic,
				},
				MerchantAccountInformation: map[ID]MerchantAccountInformationTLV{
					ID("29"): {
						Tag:    "29",
						Length: "30",
						Value: &MerchantAccountInformation{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "12",
								Value:  "D15600000000",
							},
							PaymentNetworkSpecific: []TLV{
								{
									Tag:    "05",
									Length: "10",
									Value:  "A93FO3230Q",
								},
							},
						},
					},
					ID("31"): {
						Tag:    "31",
						Length: "28",
						Value: &MerchantAccountInformation{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "12",
								Value:  "D15600000001",
							},
							PaymentNetworkSpecific: []TLV{
								{
									Tag:    "03",
									Length: "08",
									Value:  "12345678",
								},
							},
						},
					},
				},
				MerchantCategoryCode: TLV{
					Tag:    IDMerchantCategoryCode,
					Length: "04",
					Value:  "4111",
				},
				TransactionCurrency: TLV{
					Tag:    IDTransactionCurrency,
					Length: "03",
					Value:  "156",
				},
				TransactionAmount: TLV{
					Tag:    IDTransactionAmount,
					Length: "05",
					Value:  "23.72",
				},
				CountryCode: TLV{
					Tag:    IDCountryCode,
					Length: "02",
					Value:  "CN",
				},
				MerchantName: TLV{
					Tag:    IDMerchantName,
					Length: "14",
					Value:  "BEST TRANSPORT",
				},
				MerchantCity: TLV{
					Tag:    IDMerchantCity,
					Length: "07",
					Value:  "BEIJING",
				},
				MerchantInformationLanguageTemplate: &MerchantInformationLanguageTemplate{
					LanguagePreference: TLV{
						Tag:    MerchantInformationIDLanguagePreference,
						Length: "02",
						Value:  "ZH",
					},
					MerchantName: TLV{
						Tag:    MerchantInformationIDMerchantName,
						Length: "04",
						Value:  "最佳运输",
					},
					MerchantCity: TLV{
						Tag:    MerchantInformationIDMerchantCity,
						Length: "02",
						Value:  "北京",
					},
				},
				TipOrConvenienceIndicator: TLV{
					Tag:    IDTipOrConvenienceIndicator,
					Length: "02",
					Value:  "01",
				},
				AdditionalDataFieldTemplate: &AdditionalDataFieldTemplate{
					StoreLabel: TLV{
						Tag:    AdditionalIDStoreLabel,
						Length: "04",
						Value:  "1234",
					},
					CustomerLabel: TLV{
						Tag:    AdditionalIDCustomerLabel,
						Length: "03",
						Value:  "***",
					},
					TerminalLabel: TLV{
						Tag:    AdditionalIDTerminalLabel,
						Length: "08",
						Value:  "A6008667",
					},
					AdditionalConsumerDataRequest: TLV{
						Tag:    "09",
						Length: "02",
						Value:  "ME",
					},
				},
				UnreservedTemplates: map[ID]UnreservedTemplateTLV{
					ID("91"): {
						Tag:    "91",
						Length: "32",
						Value: &UnreservedTemplate{
							GloballyUniqueIdentifier: TLV{
								Tag:    "00",
								Length: "16",
								Value:  "A011223344998877",
							},
							ContextSpecificData: []TLV{
								{
									Tag:    "07",
									Length: "08",
									Value:  "12345678",
								},
							},
						},
					},
				},
				CRC: TLV{
					Tag:    IDCRC,
					Length: "04",
					Value:  "A13A",
				},
			},
			wantErr: false,
		},
		{
			name: "failed parse",
			args: args{
				payload: "00020", // not enough length
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
