package vwap

import (
	"math/big"
	"reflect"
	"sync"
	"testing"
)

func TestNewWindowQueue(t *testing.T) {
	type args struct {
		size uint16
	}
	tests := []struct {
		name string
		args args
		want *WindowQueue
	}{
		{
			name: "Create",
			args: args{
				size: 10,
			},
			want: &WindowQueue{
				Mutex:   sync.Mutex{},
				content: make([]*vwapCache, 10),
				size:    10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewWindowQueue(tt.args.size); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewWindowQueue() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestWindowQueue_Peek(t *testing.T) {
	type fields struct {
		content   []*vwapCache
		readHead  uint16
		writeHead uint16
		len       uint16
		size      uint16
	}
	type args struct {
		pos uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *vwapCache
		ok     bool
	}{
		{
			name: "Peek",
			fields: fields{
				content: []*vwapCache{
					{
						TPV:  big.NewFloat(0),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{
						TPV:  big.NewFloat(1),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{
						TPV:  big.NewFloat(2),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
				},
				readHead:  0,
				writeHead: 0,
				len:       3,
				size:      3,
			},
			args: args{
				pos: 2,
			},
			want: &vwapCache{
				TPV:  big.NewFloat(2),
				TVol: big.NewFloat(0),
				PV:   big.NewFloat(0),
				Vol:  big.NewFloat(0),
			},
			ok: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				q := &WindowQueue{
					content:   tt.fields.content,
					readHead:  tt.fields.readHead,
					writeHead: tt.fields.writeHead,
					len:       tt.fields.len,
					size:      tt.fields.size,
				}
				got, got1 := q.Peek(tt.args.pos)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Peek() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.ok {
					t.Errorf("Peek() got1 = %v, want %v", got1, tt.ok)
				}
			},
		)
	}
}

func TestWindowQueue_Pop(t *testing.T) {
	type fields struct {
		content   []*vwapCache
		readHead  uint16
		writeHead uint16
		len       uint16
		size      uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   *vwapCache
		ok     bool
	}{
		{
			name: "Pop1",
			fields: fields{
				content: []*vwapCache{
					{
						TPV:  big.NewFloat(0),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{
						TPV:  big.NewFloat(1),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{
						TPV:  big.NewFloat(2),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
				},
				readHead:  0,
				writeHead: 0,
				len:       3,
				size:      3,
			},
			want: &vwapCache{
				TPV:  big.NewFloat(0),
				TVol: big.NewFloat(0),
				PV:   big.NewFloat(0),
				Vol:  big.NewFloat(0),
			},
			ok: true,
		},
		{
			name: "Pop2",
			fields: fields{
				content: []*vwapCache{
					{
						TPV:  big.NewFloat(1),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{
						TPV:  big.NewFloat(2),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
				},
				readHead:  0,
				writeHead: 0,
				len:       2,
				size:      10,
			},
			want: &vwapCache{
				TPV:  big.NewFloat(1),
				TVol: big.NewFloat(0),
				PV:   big.NewFloat(0),
				Vol:  big.NewFloat(0),
			},
			ok: true,
		},
		{
			name: "Pop3",
			fields: fields{
				content: []*vwapCache{
					{
						TPV:  big.NewFloat(2),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
				},
				readHead:  0,
				writeHead: 0,
				len:       1,
				size:      10,
			},
			want: &vwapCache{
				TPV:  big.NewFloat(2),
				TVol: big.NewFloat(0),
				PV:   big.NewFloat(0),
				Vol:  big.NewFloat(0),
			},
			ok: true,
		},
		{
			name: "Pop4",
			fields: fields{
				content:   []*vwapCache{},
				readHead:  0,
				writeHead: 0,
				len:       0,
				size:      3,
			},
			want: nil,
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				q := &WindowQueue{
					content:   tt.fields.content,
					readHead:  tt.fields.readHead,
					writeHead: tt.fields.writeHead,
					len:       tt.fields.len,
					size:      tt.fields.size,
				}
				got, got1 := q.Pop()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Pop() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.ok {
					t.Errorf("Pop() got1 = %v, want %v", got1, tt.ok)
				}
			},
		)
	}
}

func TestWindowQueue_Push(t *testing.T) {
	type fields struct {
		content   []*vwapCache
		readHead  uint16
		writeHead uint16
		len       uint16
		size      uint16
	}
	type args struct {
		e *vwapCache
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Push when empty",
			fields: fields{
				content:   make([]*vwapCache, 3),
				readHead:  0,
				writeHead: 0,
				len:       0,
				size:      3,
			},
			args: args{
				e: &vwapCache{
					TPV:  big.NewFloat(0),
					TVol: big.NewFloat(0),
					PV:   big.NewFloat(0),
					Vol:  big.NewFloat(0),
				},
			},
			want: true,
		},
		{
			name: "Push when one is empty",
			fields: fields{
				content: []*vwapCache{
					{
						TPV:  big.NewFloat(0),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{},
				},
				readHead:  0,
				writeHead: 1,
				len:       1,
				size:      2,
			},
			args: args{
				e: &vwapCache{
					TPV:  big.NewFloat(0),
					TVol: big.NewFloat(0),
					PV:   big.NewFloat(0),
					Vol:  big.NewFloat(0),
				},
			},
			want: true,
		},
		{
			name: "Push when full",
			fields: fields{
				content: []*vwapCache{
					{
						TPV:  big.NewFloat(0),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{
						TPV:  big.NewFloat(1),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{
						TPV:  big.NewFloat(2),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
				},
				readHead:  0,
				writeHead: 0,
				len:       3,
				size:      3,
			},
			args: args{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				q := &WindowQueue{
					content:   tt.fields.content,
					readHead:  tt.fields.readHead,
					writeHead: tt.fields.writeHead,
					len:       tt.fields.len,
					size:      tt.fields.size,
				}
				if got := q.Push(tt.args.e); got != tt.want {
					t.Errorf("Push() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestWindowQueue_PeekLast(t *testing.T) {
	type fields struct {
		content   []*vwapCache
		readHead  uint16
		writeHead uint16
		len       uint16
		size      uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   *vwapCache
		want1  bool
	}{
		{
			name: "PeekLast when full of two to get first item",
			fields: fields{
				content: []*vwapCache{
					{
						TPV:  big.NewFloat(0),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{
						TPV:  big.NewFloat(1),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
				},
				readHead:  0,
				writeHead: 0,
				len:       2,
				size:      2,
			},
			want: &vwapCache{
				TPV:  big.NewFloat(1),
				TVol: big.NewFloat(0),
				PV:   big.NewFloat(0),
				Vol:  big.NewFloat(0),
			},
			want1: true,
		},
		{
			name: "PeekLast when just before full to get mid item",
			fields: fields{
				content: []*vwapCache{
					{
						TPV:  big.NewFloat(0),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{
						TPV:  big.NewFloat(1),
						TVol: big.NewFloat(0),
						PV:   big.NewFloat(0),
						Vol:  big.NewFloat(0),
					},
					{},
				},
				readHead:  0,
				writeHead: 2,
				len:       3,
				size:      3,
			},
			want: &vwapCache{
				TPV:  big.NewFloat(1),
				TVol: big.NewFloat(0),
				PV:   big.NewFloat(0),
				Vol:  big.NewFloat(0),
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				q := &WindowQueue{
					content:   tt.fields.content,
					readHead:  tt.fields.readHead,
					writeHead: tt.fields.writeHead,
					len:       tt.fields.len,
					size:      tt.fields.size,
				}
				got, got1 := q.PeekLast()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("PeekLast() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1 {
					t.Errorf("PeekLast() got1 = %v, want %v", got1, tt.want1)
				}
			},
		)
	}
}
