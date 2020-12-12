package database

import (
    "context"
)

type ServerGRPC struct {
	DB   Database
}

func (s ServerGRPC) GetAccount(c context.Context, req *GetAccountRequest) (*GetAccountReply, error) {
	a, err := s.DB.AccountByID(req.AccountId)
	if err != nil {
		return &GetAccountReply{}, nil
	}

	acct := GetAccountReply{}
	acct.AccountId = a.AccountID.String()
	acct.RootApiKey = a.RootAPIKey.String()
	acct.AlertType = a.AlertType
	acct.AdminEmail = a.AdminEmail
	acct.AlertConfig, err = a.AlertConfig.JSON()
	return &acct, err
}

func (s ServerGRPC) mustEmbedUnimplementedApiServer() {
	return
}
