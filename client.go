package sfc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// Client is an SFC client.
type Client struct {
	C       *http.Client
	BaseURL *url.URL
}

type TrackResponse struct {
	Count     int       `json:"count"`
	OrderCode string    `json:"order_code"`
	OrderInfo OrderInfo `json:"orderInfo"`
	Status    int       `json:"track_status"`
	Len       int       `json:"tracking_len"`
	Events    []Event   `json:"trackingList"`

	// raw body bytes
	body []byte
}

func (res TrackResponse) Body() []byte {
	b := make([]byte, len(res.body))
	copy(b, res.body)
	return b
}

type OrderInfo struct {
	ShippingCode       string  `json:"ship_type_code"`
	TrackingNumber     string  `json:"tracking_number"`
	TrackingNumberUSPS string  `json:"tracking_number_usps"`
	Quantity           int     `json:"numbers"`
	Weight             string  `json:"weight"`
	OrderID            string  `json:"order_id"`
	OrderCode          string  `json:"order_code"`
	CustomerOrderCode  string  `json:"customer_order_code"`
	SenderCountry      Country `json:"sender_country"`
	Country            Country `json:"country"`
}

type Country string

var _ json.Unmarshaler = func(c Country) *Country { return &c }("")

func (c *Country) UnmarshalJSON(b []byte) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("unmarshal: %w", err)
		}
	}()
	var out struct {
		Name string `json:"cn_name"`
	}
	if err := json.Unmarshal(b, &out); err == nil {
		*c = Country(out.Name)
	}
	return err
}

type Event struct {
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
	Status   string    `json:"statu"` // this is not a typo
}

var _ json.Unmarshaler = &Event{}

func (e *Event) UnmarshalJSON(b []byte) error {
	type NOPEvent Event
	var out struct {
		NOPEvent
		Date string `json:"date"`
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	t, err := time.Parse("01/02/2006 15:04:05", out.Date)
	if err != nil {
		return fmt.Errorf("parse time: %w", err)
	}

	*e = Event(out.NOPEvent)
	e.Date = t

	return nil
}

func (c *Client) Track(ctx context.Context, trackingNumber string) (*TrackResponse, error) {
	u := *c.BaseURL
	u.Path = path.Join(u.Path, "/track/track/get-track-for-web")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(url.Values{
		"tracknumber": {trackingNumber},
	}.Encode()))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := c.C.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	var (
		buf bytes.Buffer
		out TrackResponse
	)
	if err := json.NewDecoder(io.TeeReader(res.Body, &buf)).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	out.body = buf.Bytes()

	return &out, nil
}
