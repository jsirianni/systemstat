package database

import (
    "context"
)

func (s Server) GetAccount(c context.Context, req *GetAccountRequest) (*GetAccountReply, error) {
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

func (s Server) mustEmbedUnimplementedApiServer() {
	return
}
