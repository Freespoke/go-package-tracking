package internal

import (
	"embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	tracking "github.com/jkeen/tracking_number_data"
)

const jsonPath = "couriers"

var (
	Services []Service
)

func init() {
	Services = loadServices(tracking.Couriers)
}

// Courier is the data from a single json file
type Courier struct {
	Name        string    `json:"name"`
	CourierCode string    `json:"courier_code"`
	Services    []Service `json:"tracking_numbers"`
}

// Service is a single service from a courier json file. Couriers may provide
// more than one service with uniquely formatted tracking numbers.
type Service struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CourierName string `json:"-"` // from flattening courier json file
	CourierCode string `json:"-"` // from flattening courier json file
	Description string `json:"description"`
	TrackingURL string `json:"tracking_url"`

	Regex RegexParser `json:"regex"`

	Validation `json:"validation,omitempty"`

	// Test numbers drive the package level tests.
	TestNumbers struct {
		Valid   []string `json:"valid"`
		Invalid []string `json:"invalid"`
	} `json:"test_numbers"`

	Additional []Additional `json:"additional,omitempty"`

	// TODO
	Partners []struct {
		Description string `json:"description"`
		PartnerType string `json:"partner_type"`
		PartnerID   string `json:"partner_id"`
		Validation  struct {
			MatchesAll []struct {
				RegexGroupName string `json:"regex_group_name"`
				Matches        string `json:"matches"`
			} `json:"matches_all"`
		} `json:"validation"`
	} `json:"partners,omitempty"`
}

// ValidateAdditionalExists validates lookup values from the additional data.
// The values must be defined in the service to pass this validation.
func (s Service) ValidateAdditionalExists(matchKey string, keys map[string]string) bool {
	for _, sa := range s.Additional {
		if sa.Name != matchKey {
			continue
		}
		if val, ok := keys[sa.RegexGroupName]; ok {
			for _, lookup := range sa.Lookups {
				if lookup.Matches == val {
					return true
				}
			}
		}
	}

	return false
}

// Additional describes lookup values that may be encoded into the tracking number.
type Additional struct {
	Name           string   `json:"name"`
	RegexGroupName string   `json:"regex_group_name"`
	Lookups        []Lookup `json:"lookup"`
}

// Validation defines the validations to be applied to a tracking number.
type Validation struct {
	CheckDigitOpts `json:"checksum"`
	Validator      CheckDigit `json:"-"`
	Additional     struct {
		// Exists verifies a portions of the additional data extracted matches
		// an entry in the corresponding list of lookups defined with the
		// service.
		Exists []string `json:"exists"`
	} `json:"additional"`

	// SerialNumberFormat conditionally modifies the serial number prior
	// to validating.
	SerialNumberFormat struct {
		PrependIf PrependIf `json:"prepend_if"`
	} `json:"serial_number_format"`
}

// SetValidator applies the appropriate validation function for check digits.
func (val *Validation) SetValidator() {
	c := val.CheckDigitOpts

	switch c.Name {
	case "mod7":
		val.Validator = NewMod7()
	case "mod10":
		val.Validator = NewMod10(c.EvensMultiplier, c.OddsMultiplier)
	case "s10":
		val.Validator = NewS10(c.Weightings, val.Additional.Exists)
	case "sum_product_with_weightings_and_modulo":
		val.Validator = NewSumProductWithWeightingsAndModulo(c.Weightings, c.Modulo1, c.Modulo2)
	case "mod_37_36":
		val.Validator = NewMod3736()
	default:
		val.Validator = NewNoop()
	}

}

// CheckDigitOpts contains the configuration values for check digit validations
type CheckDigitOpts struct {
	Name            string `json:"name"`
	EvensMultiplier int    `json:"evens_multiplier"`
	OddsMultiplier  int    `json:"odds_multiplier"`
	Weightings      []int  `json:"weightings"`
	Modulo1         int    `json:"modulo1"`
	Modulo2         int    `json:"modulo2"`
}

// PrependIf contains the definition for tracking number prefixes that may be
// required to validate properly.
type PrependIf struct {
	Regex   RegexParser `json:"matches_regex"`
	Content string      `json:"content"`
}

// RegexParser is a helper type to convert the PCRE regex to compatible Regex.
type RegexParser struct {
	Regex *regexp.Regexp
}

func (r *RegexParser) UnmarshalJSON(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}

	var out []string
	if err := json.Unmarshal(buf, &out); err != nil {
		var s string
		if err := json.Unmarshal(buf, &s); err != nil {
			return err
		}
		out = []string{s}
	}

	str := strings.Join(out, "")

	// Fix PCRE Named Groups for RE2
	str = strings.ReplaceAll(str, "(?<", "(?P<")

	// Fix PrependIf Regex since negated look aheads aren't supported
	str = strings.ReplaceAll(str, "^(?!", "^(")

	var err error
	r.Regex, err = regexp.Compile(str)

	return err
}

// loadServices loads the json definitions for all courier services when this
// package is initialized.
// It makes a best effort to load the tracking service json files and skips any
// that can't be loaded.
func loadServices(fs embed.FS) []Service {
	services := make([]Service, 0)
	files, err := fs.ReadDir(jsonPath)
	if err != nil {
		return services
	}

	for _, f := range files {
		file, err := fs.Open(fmt.Sprintf("%s/%s", jsonPath, f.Name()))
		if err != nil {
			continue
		}
		defer func() { file.Close() }()
		var courier Courier
		if err := json.NewDecoder(file).Decode(&courier); err != nil {
			continue
		}
		for _, service := range courier.Services {
			service.CourierCode = courier.CourierCode
			service.CourierName = courier.Name
			service.Validation.SetValidator()
			services = append(services, service)
		}
	}

	return services
}
