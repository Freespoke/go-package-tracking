package parcel

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	"dev.freespoke.com/go-package-tracking/internal"
)

var (
	ErrBadString  = errors.New("extended characters present")
	ErrNoServices = errors.New("no tracking services loaded")
)

// Tracking contains results extracted from a valid tracking number.
type Tracking struct {
	// Always returned
	Courier        string
	Service        string
	TrackingNumber string
	SerialNumber   string

	// Always populated if used
	CheckDigit  string
	TrackingURL string

	// Extra details that may be encoded into the tracking number.
	Details map[string]string
}

// Track identifies valid package tracking codes.
// If a valid code is identified, it returns encoded tracking information.
// There may be multiple matches.
func Track(in string) ([]Tracking, error) {
	// Exit early if not a simple string
	if len(in) != utf8.RuneCountInString(in) {
		return nil, ErrBadString
	}

	// Exit early if no tracking services available to check.
	if len(internal.Services) == 0 {
		return nil, ErrNoServices
	}

	in = strings.ReplaceAll(strings.ToUpper(in), " ", "")
	res := make([]Tracking, 0)

	for _, service := range internal.Services {
		// initialize a single, empty result
		tracker := Tracking{
			Details: map[string]string{},
		}

		// Identify potential matches
		if service.Regex.Regex.MatchString(in) {
			matches := service.Regex.Regex.FindStringSubmatch(in)
			for i, val := range matches {
				if key := service.Regex.Regex.SubexpNames()[i]; key != "" {
					val = strings.TrimSpace(val)
					if val != "" {
						tracker.Details[key] = val
					}
				}
			}

			if v, ok := tracker.Details["SerialNumber"]; ok {
				prepend := service.Validation.SerialNumberFormat.PrependIf
				if prepend.Regex.Regex != nil && !prepend.Regex.Regex.MatchString(v) {
					v = prepend.Content + v
				}
				tracker.SerialNumber = v
				delete(tracker.Details, "SerialNumber")
			}

			if v, ok := tracker.Details["CheckDigit"]; ok {
				tracker.CheckDigit = v
				delete(tracker.Details, "CheckDigit")
			}

			// Confirm match
			if service.Validation.Validator.Validate(tracker.SerialNumber, tracker.CheckDigit) {
				match := true
				// If additional validations exist, check them
				for _, matchKey := range service.Validation.Additional.Exists {
					if ok := service.ValidateAdditionalExists(matchKey, tracker.Details); !ok {
						match = false
						break
					}
				}
				if !match {
					continue
				}

				tracker.TrackingNumber = in
				tracker.Courier = service.CourierCode
				tracker.Service = service.Name
				if service.TrackingURL != "" {
					tracker.TrackingURL = fmt.Sprintf(service.TrackingURL, in)
				}

				// Populate additional details
				tracker.populate(service.Additional)

				res = append(res, tracker)
			}
		}
	}

	return res, nil
}

// Find extracts detected tracking numbers from a string based on word boundaries.
// It returns a list of tracking results with the corresponding tracking number.
func Find(in string) (map[string][]Tracking, error) {
	// Exit early if no tracking services available to check.
	if len(internal.Services) == 0 {
		return nil, ErrNoServices
	}

	out := make(map[string][]Tracking)
	terms := strings.Fields(in)

	wg := new(sync.WaitGroup)
	wg.Add(len(terms))
	mutex := new(sync.Mutex)

	for _, v := range terms {
		go func(v string, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			track, err := Track(v)
			if err == nil && len(track) != 0 {
				mu.Lock()
				defer mu.Unlock()
				out[v] = track
			}
		}(v, wg, mutex)
	}

	wg.Wait()

	// If no results yet, try handling the case with a single tracking number
	// containing white space.
	if len(out) == 0 {
		in = regexp.MustCompile(`\s`).ReplaceAllString(in, "")
		track, err := Track(in)
		if err == nil && len(track) != 0 {
			out[in] = track
		}
	}

	return out, nil
}

// populate extracts lookup values to enrich the response.
func (t *Tracking) populate(details []internal.Additional) {
	deets := make(map[string]string)
	for _, detail := range details {
		var found bool
		for k, v := range t.Details {
			if k == detail.RegexGroupName {
				for _, lookup := range detail.Lookups {
					if v == lookup.Matches {
						found = true
					}
					if r := lookup.MatchesRegex.Regex; r != nil {
						if r.MatchString(v) {
							found = true
						}
					}
					if found {
						switch k {
						case "ServiceType":
							if lookup.Name != "" {
								deets["ServiceName"] = lookup.Name
							}
							if lookup.Description != "" {
								deets["ServiceDesc"] = lookup.Description
							}
						case "CountryCode":
							if lookup.Country != "" {
								deets["Country"] = lookup.Country
							}
							if lookup.Courier != "" {
								deets["Courier"] = lookup.Courier
							}
							if lookup.CourierURL != "" {
								deets["CourierURL"] = lookup.CourierURL
							}
							if lookup.UPUReferenceURL != "" {
								deets["UPURefURL"] = lookup.UPUReferenceURL
							}
						case "ShippingContainerType":
							if lookup.Name != "" {
								deets[k] = fmt.Sprintf("%s %s", k, lookup.Name)
							}
						default:
							deets[k+":other"] = fmt.Sprintf("%s", lookup)
						}
						break
					}
				}
				break
			}
		}
	}

	for k, v := range deets {
		t.Details[k] = v
	}
}
