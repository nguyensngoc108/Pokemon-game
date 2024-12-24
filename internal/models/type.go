package models

type When struct {
	/*
		"multiplier":0.5,"name":"fire"
	*/
	Multiplier interface{} `json:"multiplier"`
	Name       string      `json:"name"`
}

type Types struct {
	WhenDefending []When `json:"whenDefending"`
	WhenAttacking []When `json:"whenAttacking"`
	ID            string `json:"_id"`
	Rev           string `json:"_rev"`
}
