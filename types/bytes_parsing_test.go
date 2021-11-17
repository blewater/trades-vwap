package types

import "testing"

func TestParseString(t *testing.T) {
	type ParserFunc func([]byte)(string, int)
	type args struct {
		tokenFunc     ParserFunc
		msg      []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Parse Type",
			args: args{
				tokenFunc: ParseType,
				msg:      []byte(`"type":"match"`),
			},
			want: "match",
		},
		{
			name: "Parse product id",
			args: args{
				tokenFunc: ParseProductID,
				msg:       []byte(`"type":"match","trade_id":178622422,"maker_order_id":"253c56b0-f115-4364-9e06-65ffd2412f3b","taker_order_id":"928f8eb1-b6b4-4735-b12a-a512a0da684f","side":"sell","size":"0.00269988","price":"4606.8","product_id":"ETH-USD","sequence":22394045199,"time":"2021-11-10T21:37:07.988255Z"`),
			},
			want: "ETH-USD",
		},
		{
			name: "Parse time",
			args: args{
				tokenFunc: func(msg []byte) (string, int) {
					return ParseString(tokenSep, 35, msg)
				},
				msg:      []byte(`"type":"match","trade_id":178622422,"maker_order_id":"253c56b0-f115-4364-9e06-65ffd2412f3b","taker_order_id":"928f8eb1-b6b4-4735-b12a-a512a0da684f","side":"sell","size":"0.00269988","price":"4606.8","product_id":"ETH-USD","sequence":22394045199,"time":"2021-11-10T21:37:07.988255Z"`),
			},
			want: "2021-11-10T21:37:07.988255Z",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if gotVal, _ := tt.args.tokenFunc(tt.args.msg); gotVal != tt.want {
					t.Errorf("ParseString() = %v, want %v", gotVal, tt.want)
				}
			},
		)
	}
}

func TestParseF64(t *testing.T) {
	type ParserFunc func([]byte)(float64, int)
	type args struct {
		tokenFunc     ParserFunc
		msg      []byte
	}
	tests := []struct {
		name  string
		args  args
		want  float64
	}{
		{
			name: "Parse volume",
			args: args{
				tokenFunc: ParseVolume,
				msg:      []byte(`"type":"match","trade_id":178622422,"maker_order_id":"253c56b0-f115-4364-9e06-65ffd2412f3b","taker_order_id":"928f8eb1-b6b4-4735-b12a-a512a0da684f","side":"sell","size":"0.00269988","price":"4606.8","product_id":"ETH-USD","sequence":22394045199,"time":"2021-11-10T21:37:07.988255Z"`),
			},
			want: 0.00269988,
		},
		{
			name: "Parse price",
			args: args{
				tokenFunc: ParsePrice,
				msg:      []byte(`"type":"match","trade_id":178622422,"maker_order_id":"253c56b0-f115-4364-9e06-65ffd2412f3b","taker_order_id":"928f8eb1-b6b4-4735-b12a-a512a0da684f","side":"sell","size":"0.00269988","price":"4606.8","product_id":"ETH-USD","sequence":22394045199,"time":"2021-11-10T21:37:07.988255Z"`),
			},
			want: 4606.8,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, _ := tt.args.tokenFunc(tt.args.msg)
				if got != tt.want {
					t.Errorf("ParseF64() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}