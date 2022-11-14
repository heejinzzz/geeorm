package model

// Model is a model of a table
type Model interface {
	TableName() string
}
