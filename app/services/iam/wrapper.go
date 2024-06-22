package iam

import (
	"net/http"
	"saas-billing/config"
)

type IamService struct {
	httpClient *http.Client
	iamHost    string
}

func NewIamService() *IamService {
	httpClient := &http.Client{}
	return &IamService{
		httpClient: httpClient,
		iamHost:    config.IAM_HOST,
	}
}

func (s *IamService) GetUserOrgById(userId string) error {
	url := s.iamHost + "/users/" + userId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+config.JWT_SECRET)
	s.httpClient.Do(req)

	return nil
}
