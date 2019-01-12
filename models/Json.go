package models

type Json struct {
	Name   string
	Schema []Schema
}

type Schema struct {
	Field string
	Type  string
}
