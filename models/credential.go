package models

import "gorm.io/gorm"

type CredentialOperations interface {
	GetTokenKey() (string, error)
	Insert(*Credential) error
}
type Credential struct {
	gorm.Model
	Credential string `json:"credential"`
	Type       string `json:"type" gorm:"unique"`
}

type CredentialOperationsImpl struct {
	DB *gorm.DB
}

func (c *CredentialOperationsImpl) GetTokenKey() (string, error) {
	credential := &Credential{}
	err := c.DB.Where("type = ?", "jwt-token-key").First(credential).Error
	if err != nil {
		return "", err
	}
	return credential.Credential, nil
}
