# Shipping and Tracking Numbers for Testing

From *Tracking Web Service Developer Guide*, page 7

The latest version can be obtained from [Tracking.zip](https://www.ups.com/gec/techdocs/pdf/Tracking.zip).

---

The following tables contain shipping and tracking numbers for use in the testing environment and describes the results that developers should expect when attempting a TrackRequest.

When using tracking numbers other than those designated for testing in CIE, the numbers will be evaluated as production tracking numbers.

NOTE: For Web Service only. Signature image details in CIE are returned for only the tracking numbers below.

| Type | Value | Option | Service | Response |
| --- | --- | --- | --- | --- |
| Tracking  | 1Z12345E0205271688 (Signature Availability) | Activity (All) | 2nd Day Air | Delivered |
| Tracking  | 1Z12345E6605272234 | None (Last) | WorldWide Express | Delivered |
| Tracking  | 1Z12345E1305277940 | None (Last) | Next Day Air Saver |  ORIGIN SCAN |
| Tracking  | 1Z12345E6205277936 | Activity (All) | Next Day Air Saver |  2nd Delivery Attempt |
| Tracking  | 1Z12345E020527079 | None (Last) |  | Invalid Tracking Number |
| Tracking  | 1Z12345E1505270452 | None (Last) |  | No Tracking Information Available |
| Tracking  | 990728071 | Activity (All) | UPS Freight LTL | In Transit |
| Tracking  | 3251026119 | Activity (All) |  | Delivered Origin CFS |
MI Tracking Number | 9102084383041101186729 | None (Last) |  |  |
MI Reference Number | cgish000116630 | None (Last) |  |  |
| Tracking  | 1Z648616E192760718 | Activity | UPS Worldwide Express Freight | Order Process by UPS |
| Tracking  | 5548789114 | Activity | UPS Express Freight | Response for UPS Air Freight |
| Tracking  | ER751105042015062 | Activity | UPS Ocean | Response for UPS Ocean Freight |
| Tracking  | 1ZWX0692YP40636269 | Activity | UPS SUREPOST | Respnose for UPS SUREPOST |
