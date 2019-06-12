package mpm

import (
	"reflect"
	"testing"
)

func TestParsePayload(t *testing.T) {
	type args struct {
		emvString string
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
				emvString: "",
			},
			want:    &EMVQR{},
			wantErr: false,
		},
		{
			name: "payload format indicator",
			args: args{
				emvString: "000201",
			},
			want: &EMVQR{
				PayloadFormatIndicator: "01",
			},
			wantErr: false,
		},
		{
			name: "point of initiation method",
			args: args{
				emvString: "010211",
			},
			want: &EMVQR{
				PointOfInitiationMethod: "11",
			},
			wantErr: false,
		},
		{
			name: "merchant account information",
			args: args{
				emvString: "02081234abcd",
			},
			want: &EMVQR{
				MerchantAccountInformation: []MerchantAccountInformation{
					MerchantAccountInformation{
						Tag:   "02",
						Value: "1234abcd",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple merchant account information",
			args: args{
				emvString: "02081234abcd26085678efgh",
			},
			want: &EMVQR{
				MerchantAccountInformation: []MerchantAccountInformation{
					MerchantAccountInformation{
						Tag:   "02",
						Value: "1234abcd",
					},
					MerchantAccountInformation{
						Tag:   "26",
						Value: "5678efgh",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "merchant category code",
			args: args{
				emvString: "52044111",
			},
			want: &EMVQR{
				MerchantCategoryCode: "4111",
			},
			wantErr: false,
		},
		{
			name: "transaction currency",
			args: args{
				emvString: "5303156",
			},
			want: &EMVQR{
				TransactionCurrency: "156",
			},
			wantErr: false,
		},
		{
			name: "transaction amount",
			args: args{
				emvString: "540523.72",
			},
			want: &EMVQR{
				TransactionAmount: 23.72,
			},
			wantErr: false,
		},
		{
			name: "tip or convenience indicator",
			args: args{
				emvString: "550201",
			},
			want: &EMVQR{
				TipOrConvenienceIndicator: "01",
			},
			wantErr: false,
		},
		{
			name: "value of convenience fee fixed",
			args: args{
				emvString: "5603500",
			},
			want: &EMVQR{
				ValueOfConvenienceFeeFixed: "500",
			},
			wantErr: false,
		},
		{
			name: "value of convenience fee percentage",
			args: args{
				emvString: "57015",
			},
			want: &EMVQR{
				ValueOfConvenienceFeePercentage: "5",
			},
			wantErr: false,
		},
		{
			name: "country code",
			args: args{
				emvString: "5802CN",
			},
			want: &EMVQR{
				CountryCode: "CN",
			},
			wantErr: false,
		},
		{
			name: "merchant name",
			args: args{
				emvString: "5914BEST TRANSPORT",
			},
			want: &EMVQR{
				MerchantName: "BEST TRANSPORT",
			},
			wantErr: false,
		},
		{
			name: "merchant city",
			args: args{
				emvString: "6007BEIJING",
			},
			want: &EMVQR{
				MerchantCity: "BEIJING",
			},
			wantErr: false,
		},
		{
			name: "postal code",
			args: args{
				emvString: "61071234567",
			},
			want: &EMVQR{
				PostalCode: "1234567",
			},
			wantErr: false,
		},
		{
			name: "additional data field template",
			args: args{
				emvString: "6233030412340603***0708A60086670902ME",
			},
			want: &EMVQR{
				AdditionalDataFieldTemplate: AdditionalDataFieldTemplate{
					StoreLabel:                    "1234",
					CustomerLabel:                 "***",
					TerminalLabel:                 "A6008667",
					AdditionalConsumerDataRequest: "ME",
				},
			},
			wantErr: false,
		},
		{
			name: "merchant information language template",
			args: args{
				emvString: "64200002ZH0104最佳运输0202北京",
			},
			want: &EMVQR{
				MerchantInformationLanguageTemplate: MerchantInformationLanguageTemplate{
					LanguagePreference: "ZH",
					MerchantName:       "最佳运输",
					MerchantCity:       "北京",
				},
			},
			wantErr: false,
		},
		{
			name: "failed readNext",
			args: args{
				emvString: "00aa00",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed parse id",
			args: args{
				emvString: "bb0200",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed parse additional data field template",
			args: args{
				emvString: "620401cc",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed parse merchant information language template",
			args: args{
				emvString: "64100002JA01dd",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed parse transaction amount",
			args: args{
				emvString: "5405abcde",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePayload(tt.args.emvString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMerchantAccountInformation(t *testing.T) {
	type args struct {
		data map[string]string
	}
	tests := []struct {
		name string
		args args
		want MerchantAccountInformation
	}{
		{
			name: "ok",
			args: args{
				data: map[string]string{
					"id":    "26",
					"value": "1234",
				},
			},
			want: MerchantAccountInformation{
				Tag:   "26",
				Value: "1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseMerchantAccountInformation(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMerchantAccountInformation() = %v, want %v", got, tt.want)
			}
		})
	}
}
