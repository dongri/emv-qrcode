package mpm

import (
	"reflect"
	"testing"
)

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
			want:    nil,
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				payload: "00020101021229300012D156000000000510A93FO3230Q31280012D15600000001030812345678520441115802CN5914BEST TRANSPORT6007BEIJING64200002ZH0104最佳运输0202北京540523.7253031565502016233030412340603***0708A60086670902ME91320016A0112233449988770708123456786304A13A",
			},
			want: &EMVQR{
				PayloadFormatIndicator:  "01",
				PointOfInitiationMethod: PointOfInitiationMethodDynamic,
				MerchantAccountInformation: map[ID]*MerchantAccountInformation{
					ID("29"): &MerchantAccountInformation{
						Value: "0012D156000000000510A93FO3230Q",
					},
					ID("31"): &MerchantAccountInformation{
						Value: "0012D15600000001030812345678",
					},
				},
				MerchantCategoryCode: "4111",
				TransactionCurrency:  "156",
				TransactionAmount:    "23.72",
				CountryCode:          "CN",
				MerchantName:         "BEST TRANSPORT",
				MerchantCity:         "BEIJING",
				MerchantInformationLanguageTemplate: &MerchantInformationLanguageTemplate{
					LanguagePreference: "ZH",
					MerchantName:       "最佳运输",
					MerchantCity:       "北京",
				},
				TipOrConvenienceIndicator: "01",
				AdditionalDataFieldTemplate: &AdditionalDataFieldTemplate{
					StoreLabel:                    "1234",
					CustomerLabel:                 "***",
					TerminalLabel:                 "A6008667",
					AdditionalConsumerDataRequest: "ME",
				},
				UnreservedTemplates: map[ID]*UnreservedTemplate{
					ID("91"): &UnreservedTemplate{
						Value: "0016A011223344998877070812345678",
					},
				},
				CRC: "A13A",
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
					PayloadFormatIndicator:  "01",
					PointOfInitiationMethod: PointOfInitiationMethodDynamic,
					MerchantAccountInformation: map[ID]*MerchantAccountInformation{
						ID("29"): &MerchantAccountInformation{
							Value: "0012D156000000000510A93FO3230Q",
						},
					},
					MerchantCategoryCode: "4111",
					TransactionCurrency:  "156",
					TransactionAmount:    "23.72",
					CountryCode:          "CN",
					MerchantName:         "BEST TRANSPORT",
					MerchantCity:         "BEIJING",
					MerchantInformationLanguageTemplate: &MerchantInformationLanguageTemplate{
						LanguagePreference: "ZH",
						MerchantName:       "最佳运输",
						MerchantCity:       "北京",
					},
					TipOrConvenienceIndicator: "01",
					AdditionalDataFieldTemplate: &AdditionalDataFieldTemplate{
						StoreLabel:                    "1234",
						CustomerLabel:                 "***",
						TerminalLabel:                 "A6008667",
						AdditionalConsumerDataRequest: "ME",
					},
					UnreservedTemplates: map[ID]*UnreservedTemplate{
						ID("91"): &UnreservedTemplate{
							Value: "0016A011223344998877070812345678",
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
