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

func TestParseAdditionalDataFieldTemplate(t *testing.T) {
	type args struct {
		data map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    AdditionalDataFieldTemplate
		wantErr bool
	}{
		{
			name: "bill number",
			args: args{
				data: map[string]string{
					"id":    "62",
					"value": "01041234",
				},
			},
			want: AdditionalDataFieldTemplate{
				BillNumber: "1234",
			},
			wantErr: false,
		},
		{
			name: "mobile number",
			args: args{
				data: map[string]string{
					"id":    "62",
					"value": "021109012345678",
				},
			},
			want: AdditionalDataFieldTemplate{
				MobileNumber: "09012345678",
			},
			wantErr: false,
		},
		{
			name: "store label",
			args: args{
				data: map[string]string{
					"id":    "62",
					"value": "03041234",
				},
			},
			want: AdditionalDataFieldTemplate{
				StoreLabel: "1234",
			},
			wantErr: false,
		},
		{
			name: "loyalty number",
			args: args{
				data: map[string]string{
					"id":    "62",
					"value": "0403***",
				},
			},
			want: AdditionalDataFieldTemplate{
				LoyaltyNumber: "***",
			},
			wantErr: false,
		},
		{
			name: "reference label",
			args: args{
				data: map[string]string{
					"id":    "62",
					"value": "0503***",
				},
			},
			want: AdditionalDataFieldTemplate{
				ReferenceLabel: "***",
			},
			wantErr: false,
		},
		{
			name: "customer label",
			args: args{
				data: map[string]string{
					"id":    "62",
					"value": "0603***",
				},
			},
			want: AdditionalDataFieldTemplate{
				CustomerLabel: "***",
			},
			wantErr: false,
		},
		{
			name: "terminal label",
			args: args{
				data: map[string]string{
					"id":    "62",
					"value": "0708A6008667",
				},
			},
			want: AdditionalDataFieldTemplate{
				TerminalLabel: "A6008667",
			},
			wantErr: false,
		},
		{
			name: "purpose label",
			args: args{
				data: map[string]string{
					"id":    "62",
					"value": "0803***",
				},
			},
			want: AdditionalDataFieldTemplate{
				PurposeTransaction: "***",
			},
			wantErr: false,
		},
		{
			name: "additional consumer data request",
			args: args{
				data: map[string]string{
					"id":    "62",
					"value": "0902ME",
				},
			},
			want: AdditionalDataFieldTemplate{
				AdditionalConsumerDataRequest: "ME",
			},
			wantErr: false,
		},
		{
			name: "failed readNext",
			args: args{
				data: map[string]string{
					"id":    "62",
					"value": "00aa",
				},
			},
			want:    AdditionalDataFieldTemplate{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAdditionalDataFieldTemplate(tt.args.data)
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

func Test_readNext(t *testing.T) {
	type args struct {
		inputText string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		want1   string
		wantErr bool
	}{
		{
			name: "not remain inputText",
			args: args{
				inputText: "0002ab",
			},
			want: map[string]string{
				"id":    "00",
				"value": "ab",
			},
			want1:   "",
			wantErr: false,
		},
		{
			name: "remain inputText",
			args: args{
				inputText: "0002ab0102cd",
			},
			want: map[string]string{
				"id":    "00",
				"value": "ab",
			},
			want1:   "0102cd",
			wantErr: false,
		},
		{
			name: "length is not number",
			args: args{
				inputText: "00aa",
			},
			want:    nil,
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := readNext(tt.args.inputText)
			if (err != nil) != tt.wantErr {
				t.Errorf("readNext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readNext() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readNext() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_substring(t *testing.T) {
	type args struct {
		str    string
		start  int
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				str:    "0102ab",
				start:  0,
				length: 2,
			},
			want: "01",
		},
		{
			name: "start is a negative",
			args: args{
				str:    "abc",
				start:  -1,
				length: 0,
			},
			want: "abc",
		},
		{
			name: "length is zero",
			args: args{
				str:    "abc",
				start:  0,
				length: 0,
			},
			want: "abc",
		},
		{
			name: "over length of str",
			args: args{
				str:    "abc",
				start:  0,
				length: 5,
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := substring(tt.args.str, tt.args.start, tt.args.length); got != tt.want {
				t.Errorf("substring() = %v, want %v", got, tt.want)
			}
		})
	}
}
