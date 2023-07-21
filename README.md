# go-package-tracking

[![Go Reference](https://pkg.go.dev/badge/dev.freespoke.com/go-package-tracking.svg)](https://pkg.go.dev/dev.freespoke.com/go-package-tracking)

A go package for matching tracking numbers to couriers, or for discovering tracking numbers in text.

## Install

```sh
$ go get dev.freespoke.com/go-package-tracking@latest
```

## Usage

```go
// Track a number.
tracking, err := parcel.Track("986578788855")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Service: %s\nURL: %s\n\n", tracking.Service, tracking.TrackingURL)

// Find tracking numbers in a string.
// Tracking numbers are grouped by service name.
tracks, err := parcel.Find("track 986578788855 and JVGL0999999990")
if err != nil {
    log.Fatal(err)
}

for k, v := range tracks {
    for _, track := range v {
        fmt.Printf("Service: %s\nURL: %s\n", k, track.TrackingURL)
    }
}
```

## Resources

* [tracking number data](https://github.com/jkeen/tracking_number_data)

## Notes

The original package specifies its regex as PCRE which is not fully supported by Go. To prevent breaking changes from impacting this package, the dependency to jkeen/tracking_number_data is currently tied to the specific commit which includes go.mod. Once a compatible release is tagged, the dependency can be
updated.

A small helper function "fixes" the regex provided.

Tests are generated to run against the test cases embedded in the courier json file. Separate tests also
validate the check digit functions.

The signatures for the two exposed functions are

`func Track(string) (parcel.Tracking, error)`

`func Find(string) (map[string][]parcelTracking, error)`

## License

MIT (See LICENSE.md).
