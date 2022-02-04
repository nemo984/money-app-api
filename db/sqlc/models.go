// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"fmt"
)

type DateFrequency string

const (
	DateFrequencyDay   DateFrequency = "day"
	DateFrequencyWeek  DateFrequency = "week"
	DateFrequencyMonth DateFrequency = "month"
	DateFrequencyYear  DateFrequency = "year"
)

func (e *DateFrequency) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = DateFrequency(s)
	case string:
		*e = DateFrequency(s)
	default:
		return fmt.Errorf("unsupported scan type for DateFrequency: %T", src)
	}
	return nil
}

type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityMedium NotificationPriority = "medium"
	NotificationPriorityHigh   NotificationPriority = "high"
)

func (e *NotificationPriority) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = NotificationPriority(s)
	case string:
		*e = NotificationPriority(s)
	default:
		return fmt.Errorf("unsupported scan type for NotificationPriority: %T", src)
	}
	return nil
}

type Budget struct {
	BudgetID   int32         `json:"budget_id"`
	CategoryID sql.NullInt32 `json:"category_id"`
	Percentage int32         `json:"percentage"`
	StartDate  sql.NullTime  `json:"start_date"`
	EndDate    sql.NullTime  `json:"end_date"`
	UserID     int32         `json:"user_id"`
}

type Category struct {
	CategoryID int32  `json:"category_id"`
	Name       string `json:"name"`
}

type Expense struct {
	ExpenseID  int32          `json:"expense_id"`
	CategoryID int32          `json:"category_id"`
	Amount     string         `json:"amount"`
	CreatedAt  sql.NullTime   `json:"created_at"`
	Frequency  DateFrequency  `json:"frequency"`
	Note       sql.NullString `json:"note"`
	UserID     int32          `json:"user_id"`
}

type Income struct {
	IncomeID     int32          `json:"income_id"`
	IncomeTypeID string         `json:"income_type_id"`
	Description  sql.NullString `json:"description"`
	Amount       string         `json:"amount"`
	CreatedAt    sql.NullTime   `json:"created_at"`
	Frequency    DateFrequency  `json:"frequency"`
	UserID       int32          `json:"user_id"`
}

type IncomesType struct {
	IncomeTypeID int32  `json:"income_type_id"`
	Type         string `json:"type"`
}

type Notification struct {
	NotificationID int32                `json:"notification_id"`
	UserID         int32                `json:"user_id"`
	Description    sql.NullString       `json:"description"`
	Type           string               `json:"type"`
	Priority       NotificationPriority `json:"priority"`
	Read           sql.NullBool         `json:"read"`
	CreatedAt      sql.NullTime         `json:"created_at"`
}

type User struct {
	UserID     int32          `json:"user_id"`
	Username   string         `json:"username"`
	Name       sql.NullString `json:"name"`
	Password   string         `json:"password"`
	ProfileUrl sql.NullString `json:"profile_url"`
	CreatedAt  sql.NullTime   `json:"created_at"`
}
