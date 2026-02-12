package eswagger

import (
	_ "embed"
)

var (
	//go:embed doc/data.bin
	databin []byte

	CryptoKey = []byte("x76cgqt36i9c863bzmotuf8626dxiwu0")
)

const apijson = "/api.json"

type SwaggerInfo struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Version        string `json:"version"`
	TermsOfService string `json:"termsOfService"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	Email          string `json:"email"`
	Auths          []Auth `json:"auths"`
}

type Auth struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}
