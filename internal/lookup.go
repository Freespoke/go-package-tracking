package internal

type Lookup struct {
	Matches          string      `json:"matches"`
	MatchesRegex     RegexParser `json:"matches_regex"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	CountryCode      string      `json:"country_code"`
	CountryShortCode string      `json:"country_short_code"`
	Country          string      `json:"country"`
	Courier          string      `json:"courier"`
	CourierURL       string      `json:"courier_url"`
	UPUReferenceURL  string      `json:"upu_reference_url"`
}
