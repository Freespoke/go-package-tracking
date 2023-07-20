package parcel_test

import (
	"fmt"
	"testing"

	parcel "dev.freespoke.com/go-package-tracking"
	"dev.freespoke.com/go-package-tracking/internal"
)

// Testing for all positive and negative test cases embedded in the shipping
// service json files.
func TestAllServices(t *testing.T) {
	if len(internal.Services) == 0 {
		t.Fatal("TestAllServices: No services loaded.")
	}
	for _, service := range internal.Services {
		for _, test := range service.TestNumbers.Valid {
			name := fmt.Sprintf("valid %s:%s", service.ID, test)
			t.Run(name, func(t *testing.T) {
				got, err := parcel.Track(test)
				if err != nil {
					t.Errorf("parcel.Track() error %v", err)
					return
				}
				if len(got) == 0 {
					t.Error("parcel.Track() expected a result. Got none.")
					return
				}
				s := make([]string, 0)
				for _, v := range got {
					if v.Service == service.Name {
						return
					}
					s = append(s, v.Service)
				}
				t.Errorf("parcel.Track() unexpected service. Wanted %s, got [%v]", service.Name, s)
			})
		}

		for _, test := range service.TestNumbers.Invalid {
			name := fmt.Sprintf("invalid %s:%s", service.ID, test)
			t.Run(name, func(t *testing.T) {
				got, err := parcel.Track(test)
				if err != nil {
					t.Errorf("parcel.Track() error %v", err)
					return
				}
				for _, v := range got {
					if v.Service == service.Name {
						t.Errorf("parcelTrack() should not match service %s", service.Name)
						return
					}
				}
			})
		}
	}
}
