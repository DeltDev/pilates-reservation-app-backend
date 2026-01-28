package services

import (
	"fmt"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService struct {
	snapClient snap.Client
}

func NewPaymentService(serverKey, environment string) *PaymentService {
	var env midtrans.EnvironmentType
	if environment == "production" {
		env = midtrans.Production
	} else {
		env = midtrans.Sandbox
	}

	snapClient := snap.Client{}
	snapClient.New(serverKey, env)

	return &PaymentService{
		snapClient: snapClient,
	}
}

type CreateTransactionRequest struct {
	OrderID       string
	GrossAmount   int64
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	ItemName      string
	ItemPrice     int64
	ItemQuantity  int32
}

type CreateTransactionResponse struct {
	Token       string
	RedirectURL string
}

func (s *PaymentService) CreateTransaction(req CreateTransactionRequest) (*CreateTransactionResponse, error) {

	if req.OrderID == "" {
		req.OrderID = fmt.Sprintf("ORDER-%d", time.Now().Unix())
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.GrossAmount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.CustomerName,
			Email: req.CustomerEmail,
			Phone: req.CustomerPhone,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "PILATES-SESSION",
				Name:  req.ItemName,
				Price: req.ItemPrice,
				Qty:   req.ItemQuantity,
			},
		},
		EnabledPayments: snap.AllSnapPaymentType,
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	snapResp, err := s.snapClient.CreateTransaction(snapReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create midtrans transaction: %w", err)
	}

	return &CreateTransactionResponse{
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}
