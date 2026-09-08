package _go

import (
	"testing"
)

func TestConnectRequest(t *testing.T) {
	req := &ConnectRequest{
		User:     12345678,
		Password: "test_password",
		Host:     "127.0.0.1",
		Port:     443,
	}

	if req.GetUser() != 12345678 {
		t.Errorf("Expected User=12345678, got %d", req.GetUser())
	}
	if req.GetPassword() != "test_password" {
		t.Errorf("Expected Password='test_password', got %s", req.GetPassword())
	}
	if req.GetHost() != "127.0.0.1" {
		t.Errorf("Expected Host='127.0.0.1', got %s", req.GetHost())
	}
	if req.GetPort() != 443 {
		t.Errorf("Expected Port=443, got %d", req.GetPort())
	}
}

func TestAccountSummaryData(t *testing.T) {
	data := &AccountSummaryData{
		AccountLogin:    12345678,
		AccountBalance:  10000.50,
		AccountEquity:   10250.75,
		AccountCurrency: "USD",
		AccountLeverage: 100,
	}

	if data.GetAccountLogin() != 12345678 {
		t.Errorf("Expected Login=12345678, got %d", data.GetAccountLogin())
	}
	if data.GetAccountBalance() != 10000.50 {
		t.Errorf("Expected Balance=10000.50, got %f", data.GetAccountBalance())
	}
	if data.GetAccountEquity() != 10250.75 {
		t.Errorf("Expected Equity=10250.75, got %f", data.GetAccountEquity())
	}
	if data.GetAccountCurrency() != "USD" {
		t.Errorf("Expected Currency='USD', got %s", data.GetAccountCurrency())
	}
	if data.GetAccountLeverage() != 100 {
		t.Errorf("Expected Leverage=100, got %d", data.GetAccountLeverage())
	}
}

func TestGetIdRequest(t *testing.T) {
	req := &GetIdRequest{
		User:     "12345678",
		Password: "test_password",
	}
	if req.GetUser() != "12345678" {
		t.Errorf("Expected User='12345678', got %s", req.GetUser())
	}
	if req.GetPassword() != "test_password" {
		t.Errorf("Expected Password='test_password', got %s", req.GetPassword())
	}

	reply := &GetIdReply{
		Response: &GetIdReply_Data{
			Data: &GetIdData{
				Id: "68c935ee-a2b1-4f3e-bb36-3982845cfa85",
			},
		},
	}
	if reply.GetData().GetId() != "68c935ee-a2b1-4f3e-bb36-3982845cfa85" {
		t.Errorf("Expected Id='68c935ee-a2b1-4f3e-bb36-3982845cfa85', got %s", reply.GetData().GetId())
	}
}
