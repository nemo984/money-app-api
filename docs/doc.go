package docs

// date frequency
// swagger:enum DateFrequency
type DateFrequency string

const (
	DateFrequencyDay   DateFrequency = "day"
	DateFrequencyWeek  DateFrequency = "week"
	DateFrequencyMonth DateFrequency = "month"
	DateFrequencyYear  DateFrequency = "year"
)
