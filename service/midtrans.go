package service

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"os"
	"payment-service/model"
)

type MidtransService interface {
	GenerateSnapURL(payment model.PaymentRequest) (string, error)
	VerifyPayment(orderID string) (string, error)
}

type MidtransServiceImpl struct {
	snapClient snap.Client
	coreClient coreapi.Client
}

func NewMidtransService() *MidtransServiceImpl {
	var snapClient snap.Client
	var coreClient coreapi.Client
	snapClient.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)
	coreClient.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)
	return &MidtransServiceImpl{snapClient: snapClient, coreClient: coreClient}
}

func (service *MidtransServiceImpl) GenerateSnapURL(payment model.PaymentRequest) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  payment.OrderID,
			GrossAmt: payment.GrossAmt,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items:           &payment.ItemDetails,
	}
	response, _ := service.snapClient.CreateTransaction(req)
	return response.RedirectURL, nil
}

func (service *MidtransServiceImpl) VerifyPayment(orderID string) (string, error) {
	transactionStatusResp, err := service.coreClient.CheckTransaction(orderID)
	if err != nil {
		return "", err
	} else {
		if transactionStatusResp != nil {
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
				} else if transactionStatusResp.FraudStatus == "accept" {
					return transactionStatusResp.PaymentType, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				return transactionStatusResp.PaymentType, nil
			} else if transactionStatusResp.TransactionStatus == "deny" {
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
			} else if transactionStatusResp.TransactionStatus == "pending" {
			}
		}
	}
	return "", nil
}
