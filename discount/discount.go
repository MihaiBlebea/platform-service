package discount

import "time"

// Discount model
type Discount struct {
	ID         int
	ProductID  int
	Code       string
	Percentage float64
	Active     bool
	Expires    time.Time
	Created    time.Time
	Updated    time.Time
}

// IsValid returns true if the discount code is valid
func (d *Discount) IsValid() bool {
	// Check active
	if d.Active == false {
		return false
	}

	// Check data
	if d.Expires.Before(time.Now()) {
		return false
	}

	return true
}
