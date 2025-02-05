package domain

type Port struct {
	ID          *string   `json:"id,omitempty" db:"id"`
	Name        string    `json:"name" db:"name"`
	City        string    `json:"city" db:"city"`
	Country     string    `json:"country" db:"country"`
	Alias       []string  `json:"alias" db:"alias"`
	Regions     []string  `json:"regions" db:"regions"`
	Coordinates []float64 `json:"coordinates" db:"coordinates"`
	Province    string    `json:"province" db:"province"`
	Timezone    string    `json:"timezone" db:"timezone"`
	Unlocs      []string  `json:"unlocs" db:"unlocs"`
	Code        string    `json:"code" db:"code"`
}
