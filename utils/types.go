package utils

type BwCounter struct {
	Limit float64 `json:"monthly_bw_limit_b, omitempty"`
	Used  float64 `json:"bw_counter_b, omitempty"`
	Reset int32   `json:"bw_reset_day_of_month, omitempty"`
}
