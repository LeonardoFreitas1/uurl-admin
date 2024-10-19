package models

import "time"

type LanguageTag struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	IsoCode1 string     `json:"isoCode1"`
	IsoCode2 string     `json:"isoCode2"`
	Variants []Variants `json:"variants"`
}

type Variants struct {
	ID                      int       `json:"id"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`
	VariantTag              string    `json:"variantTag"`
	Description             string    `json:"description"`
	IsIANALanguageSubTag    bool      `json:"isIANALanguageSubTag"`
	InstancesOnDomainsCount int       `json:"instancesOnDomainsCount"`
	LanguageTagID           int       `json:"languageTagId"`
}
