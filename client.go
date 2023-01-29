package garara

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const attachMaxFileSize = 3

type V1Client struct {
	*http.Client
	apiUser     string
	apiPassword string
	mgrUser     string
	mgrPassword string
	SiteID      int
	ServiceID   int
}

func NewDefaultV1Client(apiUser, apiPassword, mgrUser, mgrPassword string, siteID, serviceID int) *V1Client {
	return &V1Client{
		Client:      http.DefaultClient,
		apiUser:     apiUser,
		apiPassword: apiPassword,
		mgrUser:     mgrUser,
		mgrPassword: mgrPassword,
		SiteID:      siteID,
		ServiceID:   serviceID,
	}
}

func NewV1Client(client *http.Client, apiUser, apiPassword, mgrUser, mgrPassword string, siteID, serviceID int) *V1Client {
	return &V1Client{
		Client:      client,
		apiUser:     apiUser,
		apiPassword: apiPassword,
		mgrUser:     mgrUser,
		mgrPassword: mgrPassword,
		SiteID:      siteID,
		ServiceID:   serviceID,
	}
}

func (c *V1Client) SetApiUser(user string) {
	c.apiUser = user
}

func (c *V1Client) SetApiPassword(password string) {
	c.apiPassword = password
}

func (c *V1Client) SetMgrUser(user string) {
	c.mgrUser = user
}

func (c *V1Client) SetMgrPassword(password string) {
	c.mgrPassword = password
}

func (c *V1Client) getAuthMail() AuthMail {
	return AuthMail{
		Auth: Auth{
			Site:    AttrID{c.SiteID},
			Service: AttrID{c.ServiceID},
			Name:    CDATAString{c.apiUser},
			Pass:    CDATAString{c.apiPassword},
		},
	}
}

func (c *V1Client) getStatus(ctx context.Context, r V1AuthRequest, host string) ([]ResDelivery, error) {
	x, err := xml.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(infoEndpoint, host), strings.NewReader(string(x)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	return c.getResDeliveries(req)
}

func (c *V1Client) getResDeliveries(r *http.Request) ([]ResDelivery, error) {
	res, err := c.Do(r)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Status: %d, http error: %s", res.StatusCode, string(b)))
	}

	resp := &V1MailResult{}
	err = xml.Unmarshal(b, resp)
	if err != nil {
		return nil, err
	}

	if resp.Result.Code == "-1" {
		messages := make([]string, 0, len(resp.Errors.Error))
		for _, e := range resp.Errors.Error {
			messages = append(messages, fmt.Sprintf("{\"code\":\"%s\", \"kind\":\"%s\", \"msg\":\"%s\"}", e.Code, e.Kind, e.Value))
		}

		return nil, errors.New(fmt.Sprintf("[%s]", strings.Join(messages, ",")))
	}

	size := len(resp.Delivery)
	for i := 0; i < size; i++ {
		if resp.Delivery[i].Result.Code != "-1" {
			continue
		}

		messages := make([]string, 0, len(resp.Delivery[i].Errors.Error))
		for _, e := range resp.Delivery[i].Errors.Error {
			messages = append(messages, fmt.Sprintf("{\"code\":\"%s\", \"kind\":\"%s\", \"msg\":\"%s\"}", e.Code, e.Kind, e.Value))
		}

		resp.Delivery[i].Errors.Err = errors.New(
			fmt.Sprintf(
				"{\"deliver_id\":\"%s\",\"request_id\":\"%s\",\"from\":\"%s\",\"to\":\"%s\",\"errors\":[%s]}",
				resp.Delivery[i].DeliverID,
				resp.Delivery[i].RequestID,
				resp.Delivery[i].From,
				resp.Delivery[i].To,
				strings.Join(messages, ","),
			),
		)
	}

	return resp.Delivery, nil
}

func (c *V1Client) ResDeliverySplitErrors(r []ResDelivery) (res []ResDelivery, errs []error) {
	res = make([]ResDelivery, 0, len(r))
	errs = make([]error, 0, len(r))
	for _, v := range r {
		if v.Errors.Err != nil {
			errs = append(errs, v.Errors.Err)
			continue
		}
		res = append(res, v)
	}

	if len(errs) == 0 {
		return r, nil
	}
	return res, errs
}

func (c *V1Client) ResponseCSVToClickLog(body []byte) ([]ClickLog, error) {
	rb := strings.NewReader(string(body))
	r := csv.NewReader(rb)

	var list []ClickLog
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}

		if len(row) != 10 {
			return nil, errors.New("this csv does not meet the required requirements")
		}

		s := ClickLog{
			ClickDate:   row[0],
			URL:         row[1],
			UrlID:       row[2],
			CarrierCode: row[3],
			DeviceCode:  row[4],
			MailAddress: row[5],
			UserAgent:   row[6],
			AddInfo:     row[7],
			DeliverID:   row[8],
			OpenedFlag:  row[9],
		}
		list = append(list, s)
	}
	return list, nil
}

func (c *V1Client) ErrorDeliverySplitErrors(e []ErrorDelivery) (res []ErrorDelivery, errs []error) {
	res = make([]ErrorDelivery, 0, len(e))
	errs = make([]error, 0, len(e))

	for _, v := range e {
		if v.Errors.Err != nil {
			errs = append(errs, v.Errors.Err)
			continue
		}
		res = append(res, v)
	}

	if len(errs) == 0 {
		return e, nil
	}
	return res, errs

}

func (c *V1Client) getErrorDeliveries(ctx context.Context, r V1AuthRequest, host string) ([]ErrorDelivery, error) {

	x, err := xml.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(infoEndpoint, host), bytes.NewBuffer(x))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("status: %d, response: %s", res.StatusCode, string(b)))
	}

	resp := ErrorList{}
	err = xml.Unmarshal(b, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Result.Code == "-1" {
		return nil, errors.New(fmt.Sprintf("Error code: %s, message: %s", resp.Result.Code, resp.Result.Value))
	}

	size := len(resp.Delivery)
	for i := 0; i < size; i++ {
		if resp.Delivery[i].Result.Code != "-1" {
			continue
		}

		messages := make([]string, 0, len(resp.Delivery[i].Errors.Error))
		for _, e := range resp.Delivery[i].Errors.Error {
			messages = append(messages, fmt.Sprintf("{\"code\":\"%s\", \"kind\":\"%s\", \"msg\":\"%s\"}", e.Code, e.Kind, e.Value))
		}

		resp.Delivery[i].Errors.Err = errors.New(
			fmt.Sprintf(
				"{\"delivery_id\":\"%s\",\"from\":\"%s\",\"to\":\"%s\",\"errors\":[%s]}",
				resp.Delivery[i].DeliveryID,
				resp.Delivery[i].From,
				resp.Delivery[i].To,
				strings.Join(messages, ","),
			),
		)
	}

	return resp.Delivery, nil

}

func (c *V1Client) deleteReservation(ctx context.Context, r V1AuthRequest, host string) (resp V1MailResult, err error) {
	x, err := xml.Marshal(r)
	if err != nil {
		return resp, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(infoEndpoint, host), bytes.NewBuffer(x))
	if err != nil {
		return resp, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	res, err := c.Do(req)
	if err != nil {
		return resp, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}

	if res.StatusCode != http.StatusOK {
		return resp, errors.New(fmt.Sprintf("status: %d, response: %s", res.StatusCode, string(b)))
	}

	err = xml.Unmarshal(b, &resp)
	if err != nil {
		return resp, err
	}

	if resp.Result.Code == "-1" {
		return resp, errors.New(fmt.Sprintf("error code: %s, message: %s", resp.Result.Code, resp.Result.Value))
	}

	return resp, err

}
