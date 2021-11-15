package vwap_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/blewater/zh/log"
	"github.com/blewater/zh/types"
	"github.com/blewater/zh/vwap"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type Ticker struct {
	Price     *big.Float
	Volume    *big.Float
	ProductID string
}

// VWAPTestSuite is a test suite that runs a series of tests comparing the
// computed VWAP float64 and big.Float results.
type VWAPTestSuite struct {
	logger *zap.Logger
	ctx    context.Context
	suite.Suite
}

func (suite *VWAPTestSuite) SetupTest() {
	suite.logger = log.New(true, "Test Suite")
	suite.ctx = log.ContextWithLogger(context.Background(), suite.logger)
}

func (suite *VWAPTestSuite) TestVMAPResults() {
	type args struct {
		quotes       []Ticker
		productsVWAP *vwap.ProductsVwap
	}

	testCases := []struct {
		name        string
		args        args
		results     []types.VWAPResult
		expectPass  bool
		expectedErr string
	}{
		{
			name: "Single Product Default Zero",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(0),
						Volume:    big.NewFloat(0),
						ProductID: "0",
					},
				},
				productsVWAP: vwap.New([]string{"0"}, 1),
			},
			results: []types.VWAPResult{
				{
					ProductID: "0",
					Vwap:      big.NewFloat(0),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Single Product Default One",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(1),
						ProductID: "1",
					},
				},
				productsVWAP: vwap.New([]string{"1"}, 1),
			},
			results: []types.VWAPResult{
				{
					ProductID: "1",
					Vwap:      big.NewFloat(1),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Single Product Single Price Two",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(2),
						Volume:    big.NewFloat(1),
						ProductID: "Prod",
					},
				},
				productsVWAP: vwap.New([]string{"Prod"}, 1),
			},
			results: []types.VWAPResult{
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(2),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Single Product Half Window Zero",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(0),
						Volume:    big.NewFloat(0),
						ProductID: "0",
					},
					{
						Price:     big.NewFloat(0),
						Volume:    big.NewFloat(0),
						ProductID: "0",
					},
				},
				productsVWAP: vwap.New([]string{"0"}, 1),
			},
			results: []types.VWAPResult{
				{
					ProductID: "0",
					Vwap:      big.NewFloat(0),
				},
				{
					ProductID: "0",
					Vwap:      big.NewFloat(0),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Single Product 2 Items Full Window",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(2.5),
						Volume:    big.NewFloat(1),
						ProductID: "1",
					},
					{
						Price:     big.NewFloat(4.5),
						Volume:    big.NewFloat(3),
						ProductID: "1",
					},
				},
				productsVWAP: vwap.New([]string{"1"}, 2),
			},
			results: []types.VWAPResult{
				{
					ProductID: "1",
					Vwap:      big.NewFloat(2.5),
				},
				{
					ProductID: "1",
					Vwap:      big.NewFloat(4),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Single Product 2 Items Half Window One",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(2.5),
						Volume:    big.NewFloat(1),
						ProductID: "1",
					},
					{
						Price:     big.NewFloat(4.5),
						Volume:    big.NewFloat(3),
						ProductID: "1",
					},
				},
				productsVWAP: vwap.New([]string{"1"}, 1),
			},
			results: []types.VWAPResult{
				{
					ProductID: "1",
					Vwap:      big.NewFloat(2.5),
				},
				{
					ProductID: "1",
					Vwap:      big.NewFloat(4.5),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Single Product 3 Items Full Window",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(2.5),
						Volume:    big.NewFloat(1),
						ProductID: "1",
					},
					{
						Price:     big.NewFloat(4.5),
						Volume:    big.NewFloat(3),
						ProductID: "1",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(1),
						ProductID: "1",
					},
				},
				productsVWAP: vwap.New([]string{"1"}, 3),
			},
			results: []types.VWAPResult{
				{
					ProductID: "1",
					Vwap:      big.NewFloat(2.5),
				},
				{
					ProductID: "1",
					Vwap:      big.NewFloat(4),
				},
				{
					ProductID: "1",
					Vwap:      big.NewFloat(3.4),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Single Product 3 Items Window of One",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(2.5),
						Volume:    big.NewFloat(1),
						ProductID: "1",
					},
					{
						Price:     big.NewFloat(4.5),
						Volume:    big.NewFloat(3),
						ProductID: "1",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(1),
						ProductID: "1",
					},
				},
				productsVWAP: vwap.New([]string{"1"}, 1),
			},
			results: []types.VWAPResult{
				{
					ProductID: "1",
					Vwap:      big.NewFloat(2.5),
				},
				{
					ProductID: "1",
					Vwap:      big.NewFloat(4.5),
				},
				{
					ProductID: "1",
					Vwap:      big.NewFloat(1),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Single Product 3 Items Window of Two",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(2.5),
						Volume:    big.NewFloat(1),
						ProductID: "1",
					},
					{
						Price:     big.NewFloat(4.5),
						Volume:    big.NewFloat(3),
						ProductID: "1",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(1),
						ProductID: "1",
					},
				},
				productsVWAP: vwap.New([]string{"1"}, 2),
			},
			results: []types.VWAPResult{
				{
					ProductID: "1",
					Vwap:      big.NewFloat(2.5),
				},
				{
					ProductID: "1",
					Vwap:      big.NewFloat(4),
				},
				{
					ProductID: "1",
					Vwap:      big.NewFloat(3.625),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Multi Product Default Zero",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(0),
						Volume:    big.NewFloat(0),
						ProductID: "0",
					},
					{
						Price:     big.NewFloat(0),
						Volume:    big.NewFloat(0),
						ProductID: "1",
					},
					{
						Price:     big.NewFloat(0),
						Volume:    big.NewFloat(0),
						ProductID: "2",
					},
				},
				productsVWAP: vwap.New([]string{"0", "1", "2"}, 1),
			},
			results: []types.VWAPResult{
				{
					ProductID: "0",
					Vwap:      big.NewFloat(0),
				},
				{
					ProductID: "1",
					Vwap:      big.NewFloat(0),
				},
				{
					ProductID: "2",
					Vwap:      big.NewFloat(0),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Single Product 10 Tickers With Full sized Window",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(0),
						Volume:    big.NewFloat(0),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(0),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(1),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(1),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(2),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(2),
						Volume:    big.NewFloat(1),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(3),
						Volume:    big.NewFloat(1),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(5),
						Volume:    big.NewFloat(1),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(2),
						Volume:    big.NewFloat(2),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(20),
						Volume:    big.NewFloat(18),
						ProductID: "Prod",
					},
				},
				productsVWAP: vwap.New([]string{"Prod"}, 10),
			},
			results: []types.VWAPResult{
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(0),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(0),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(1),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(1),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(1),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(1.2),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(1.5),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(2),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(2),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(14),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
		{
			name: "Single Product 10 Tickers with Window Size of Two",
			args: args{
				quotes: []Ticker{
					{
						Price:     big.NewFloat(0),
						Volume:    big.NewFloat(0),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(0),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(1),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(1),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(1),
						Volume:    big.NewFloat(2),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(2),
						Volume:    big.NewFloat(2),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(3),
						Volume:    big.NewFloat(3),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(5),
						Volume:    big.NewFloat(1),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(2),
						Volume:    big.NewFloat(2),
						ProductID: "Prod",
					},
					{
						Price:     big.NewFloat(20),
						Volume:    big.NewFloat(18),
						ProductID: "Prod",
					},
				},
				productsVWAP: vwap.New([]string{"Prod"}, 2),
			},
			results: []types.VWAPResult{
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(0),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(0),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(1),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(1),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(1),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(1.5),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(2.6),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(3.5),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(3),
				},
				{
					ProductID: "Prod",
					Vwap:      big.NewFloat(18.2),
				},
			},
			expectPass:  true,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		suite.Run(
			tc.name, func() {
				suite.logger.Info("Running Test", zap.String("Name", tc.name))
				for i, quote := range tc.args.quotes {
					err := tc.args.productsVWAP.ProduceVwap(
						suite.ctx, quote.ProductID, quote.Price, quote.Volume,
					)
					if !tc.expectPass {
						suite.Require().Error(err, tc.expectedErr)
						continue
					}

					prodRes := tc.results[i].ProductID
					res := <-tc.args.productsVWAP.GetResultsQ()
					suite.Require().Equal(
						prodRes, res.ProductID,
						"ProductID not matching",
					)
					suite.Require().Equal(
						tc.results[i].Vwap.String(),
						res.Vwap.String(), "VWAP Result not matching",
					)
				}
			},
		)
	}
}

func TestVWAPTestSuite(t *testing.T) {
	suite.Run(t, new(VWAPTestSuite))
}
