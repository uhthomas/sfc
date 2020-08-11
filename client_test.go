package sfc_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/uhthomas/sfc"
)

func TestClient_Track(t *testing.T) {
	t.Run("should track", func(t *testing.T) {
		const someTrackingNumber = "some tracking number"

		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/track/track/get-track-for-web"; got != want {
				t.Fatalf("unexpected url path: got %s, want %s", got, want)
			}
			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if got, want := r.FormValue("tracknumber"), someTrackingNumber; got != want {
				t.Fatalf("unexpected tracking number: got %s, want %s", got, want)
			}

			w.Write([]byte(`{
				"count": 1,
				"order_code": "some order code",
				"orderInfo": {
					"ship_type_code": "some ship type code",
					"tracking_number": "some tracking number",
					"tracking_number_usps": "some usps tracking number",
					"numbers": 1,
					"weight": "",
					"order_id": "",
					"order_code": "WW6404007290132",
					"customer_order_code": "#29516",
					"sender_country": {
						"cn_name": "\u4e2d\u56fd"
					},
					"country": {
						"cn_name": "\u82f1\u56fd"
					}
				},
				"track_status": 1,
				"tracking_len": 8,
				"trackingList": [
					{
						"date": "08\/06\/2020 14:45:00",
						"location": " ",
						"statu": "Flight has arrived"
					},
					{
						"date": "08\/06\/2020 11:22:02",
						"location": "gb",
						"statu": "pre-advice"
					},
					{
						"date": "08\/03\/2020 13:12:00",
						"location": " ",
						"statu": "Flight has taken off"
					},
					{
						"date": "07\/31\/2020 11:04:43",
						"location": "Shenzhen",
						"statu": "Departured from SFC warehouse"
					},
					{
						"date": "07\/30\/2020 18:20:09",
						"location": "shenzhen",
						"statu": "Arrive SFC warehouse in processing"
					},
					{
						"date": "07\/30\/2020 09:12:33",
						"location": "Shenzhen",
						"statu": "SFC driver pick-up"
					},
					{
						"date": "07\/29\/2020 20:06:56",
						"location": " ",
						"statu": "Shipment information sent to SFC"
					},
					{
						"date": "07\/29\/2020 20:06:55",
						"location": " ",
						"statu": "SHIPMENT INFORMATION SUBMITTED"
					}
				]
			}`))
		}))
		defer s.Close()

		u, err := url.Parse(s.URL)
		if err != nil {
			t.Fatal(err)
		}

		res, err := (&sfc.Client{
			C:       s.Client(),
			BaseURL: u,
		}).Track(context.Background(), someTrackingNumber)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := *res, (sfc.TrackResponse{
			Count:     1,
			OrderCode: "some order code",
			OrderInfo: sfc.OrderInfo{
				ShippingCode:       "some ship type code",
				TrackingNumber:     "some tracking number",
				TrackingNumberUSPS: "some usps tracking number",
				Quantity:           1,
				Weight:             "",
				OrderID:            "",
				OrderCode:          "WW6404007290132",
				CustomerOrderCode:  "#29516",
				SenderCountry:      "\u4e2d\u56fd",
				Country:            "\u82f1\u56fd",
			},
			Status: 1,
			Len:    8,
			Events: []sfc.Event{
				{
					Date:     time.Date(2020, 8, 6, 14, 45, 0, 0, time.UTC),
					Location: " ",
					Status:   "Flight has arrived",
				},
				{
					Date:     time.Date(2020, 8, 6, 11, 22, 2, 0, time.UTC),
					Location: "gb",
					Status:   "pre-advice",
				},
				{
					Date:     time.Date(2020, 8, 3, 13, 12, 0, 0, time.UTC),
					Location: " ",
					Status:   "Flight has taken off",
				},
				{
					Date:     time.Date(2020, 7, 31, 11, 4, 43, 0, time.UTC),
					Location: "Shenzhen",
					Status:   "Departured from SFC warehouse",
				},
				{
					Date:     time.Date(2020, 7, 30, 18, 20, 9, 0, time.UTC),
					Location: "shenzhen",
					Status:   "Arrive SFC warehouse in processing",
				},
				{
					Date:     time.Date(2020, 7, 30, 9, 12, 33, 0, time.UTC),
					Location: "Shenzhen",
					Status:   "SFC driver pick-up",
				},
				{
					Date:     time.Date(2020, 7, 29, 20, 6, 56, 0, time.UTC),
					Location: " ",
					Status:   "Shipment information sent to SFC",
				},
				{
					Date:     time.Date(2020, 7, 29, 20, 6, 55, 0, time.UTC),
					Location: " ",
					Status:   "SHIPMENT INFORMATION SUBMITTED",
				},
			},
		}); !cmp.Equal(got, want, cmpopts.IgnoreUnexported(sfc.TrackResponse{})) {
			t.Fatalf("unexpected response: %s", cmp.Diff(got, want, cmpopts.IgnoreUnexported(sfc.TrackResponse{})))
		}
	})
}
