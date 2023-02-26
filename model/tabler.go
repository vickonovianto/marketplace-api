package model

// for overriding gorm table name
type Tabler interface {
	TableName() string
}
