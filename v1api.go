package garara

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	mailEndpoint        = "https://%s/tm/lpmail_qmode.php"
	infoEndpoint        = "https://%s/tm/lpmail.php"
	clickCountEndpoint  = "https://%s/click/get_click_log.php"
	errorFilterEndpoint = "https://%s/errflt/getErrorFilterList.php"
)

// SendQueueMode is 配信予約APIの実行
func (c *V1Client) SendQueueMode(ctx context.Context, r V1MailRequest, host string) ([]ResDelivery, error) {
	rb, err := xml.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(mailEndpoint, host), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header["X-AutomailUser"] = []string{c.apiUser}
	req.Header["X-AutomailPassword"] = []string{c.apiPassword}
	req.Header["X-AutomailUseSite"] = []string{strconv.Itoa(c.SiteID)}
	req.Header["X-AutomailUseService"] = []string{strconv.Itoa(c.ServiceID)}

	return c.getResDeliveries(req)
}

// GetStatusByDeliverIDs is DeliveryID指定の配信ステータス取得APIの実行
func (c *V1Client) GetStatusByDeliverIDs(ctx context.Context, deliverIDs []string, host string) ([]ResDelivery, error) {
	deliveries := make([]Delivery, 0, len(deliverIDs))
	for _, v := range deliverIDs {
		d := Delivery{
			Action:    GET_STATUS,
			DeliverID: v,
		}
		deliveries = append(deliveries, d)
	}

	r := V1AuthRequest{
		AuthMail: c.getAuthMail(),
		Delivery: deliveries,
	}

	x, err := xml.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(infoEndpoint, host), strings.NewReader(string(x)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	return c.getStatus(ctx, r, host)
}

// GetStatusByRequestIDs is RequestID指定の配信ステータス取得APIの実行
func (c *V1Client) GetStatusByRequestIDs(ctx context.Context, requestIDs []string, host string) ([]ResDelivery, error) {
	deliveries := make([]Delivery, 0, len(requestIDs))
	for _, v := range requestIDs {
		d := Delivery{
			Action:    GET_STATUS,
			RequestID: v,
		}
		deliveries = append(deliveries, d)

	}
	r := V1AuthRequest{
		AuthMail: c.getAuthMail(),
		Delivery: deliveries,
	}

	return c.getStatus(ctx, r, host)
}

// GetStatusByTerms is 期間指定の配信ステータス取得APIの実行
func (c *V1Client) GetStatusByTerms(ctx context.Context, terms []Term, host string) ([]ResDelivery, error) {
	deliveries := make([]Delivery, 0, len(terms))
	for _, t := range terms {
		d := Delivery{
			Action: GET_STATUS,
			Term:   t,
		}
		deliveries = append(deliveries, d)
	}

	r := V1AuthRequest{
		AuthMail: c.getAuthMail(),
		Delivery: deliveries,
	}

	return c.getStatus(ctx, r, host)
}

// GetClickLog is クリックカウントログ取得APIの実行
func (c *V1Client) GetClickLog(ctx context.Context, deliverID, termFrom, termTo, mode, host string) ([]byte, error) {
	form := url.Values{}
	form.Add("site_id", strconv.Itoa(c.SiteID))
	form.Add("user_id", c.mgrUser)
	form.Add("password", c.mgrPassword)
	form.Add("mode", mode)
	form.Add("deliver_id", deliverID)
	form.Add("term_from", termFrom)
	form.Add("term_to", termTo)

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(clickCountEndpoint, host), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		return nil, errors.New(fmt.Sprintf("Status: %d, http error response: %s", res.StatusCode, string(b)))
	}

	return b, nil
}

// GetErrorFilter is エラーフィルターリスト取得APIの実行
func (c *V1Client) GetErrorFilter(ctx context.Context, host string) (*ErrorFilter, error) {
	form := url.Values{}
	form.Add("site_id", strconv.Itoa(c.SiteID))
	form.Add("user_name", c.mgrUser)
	form.Add("user_passwd", c.mgrPassword)

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(errorFilterEndpoint, host), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	resp := &ErrorFilter{}
	err = xml.Unmarshal(b, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetErrorListByDeliverIDs is デリバリーID指定のエラーリスト取得API実行
func (c *V1Client) GetErrorListByDeliverIDs(ctx context.Context, deliverIDs []string, host string, isAllError bool) ([]ErrorDelivery, error) {
	deliveries := make([]Delivery, 0, len(deliverIDs))
	for k, id := range deliverIDs {
		var d Delivery
		d.ID = k
		d.Action = GET_ERROR
		d.DeliverID = id
		if isAllError {
			d.Mode = "all"
		}
		deliveries = append(deliveries, d)
	}

	r := V1AuthRequest{
		AuthMail: c.getAuthMail(),
		Delivery: deliveries,
	}

	return c.getErrorDeliveries(ctx, r, host)
}

// GetErrorListByTerms is 期間指定のエラーリスト取得API実行
func (c *V1Client) GetErrorListByTerms(ctx context.Context, terms []Term, host string, isAllError bool) ([]ErrorDelivery, error) {
	deliveries := make([]Delivery, 0, len(terms))
	for k, term := range terms {
		var d Delivery
		d.ID = k
		d.Action = GET_ERROR
		d.Term = term
		if isAllError {
			d.Mode = "all"
		}
		deliveries = append(deliveries, d)
	}

	r := V1AuthRequest{
		AuthMail: c.getAuthMail(),
		Delivery: deliveries,
	}

	return c.getErrorDeliveries(ctx, r, host)
}

// GetResultListByDeliverIDs is 配信結果リスト取得APIを複数DeliverIDにて取得する
func (c *V1Client) GetResultListByDeliverIDs(ctx context.Context, deliverIDs []string, host string) ([]ResultDelivery, error) {
	deliveries := make([]Delivery, 0, len(deliverIDs))
	for k, id := range deliverIDs {
		var d Delivery
		d.ID = k
		d.Action = GET_RESULT
		d.DeliverID = id
		deliveries = append(deliveries, d)
	}
	r := V1AuthRequest{
		AuthMail: c.getAuthMail(),
		Delivery: deliveries,
	}

	x, err := xml.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(infoEndpoint, host), strings.NewReader(string(x)))
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

	resp := ResultList{}
	err = xml.Unmarshal(b, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Delivery, nil
}

func (c *V1Client) DeleteReservationByDeliverIDs(ctx context.Context, deliverIDs []string, host string) (resp V1MailResult, err error) {
	deliveries := make([]Delivery, 0, len(deliverIDs))
	for k, v := range deliverIDs {
		var d Delivery
		d.ID = k
		d.DeliverID = v
		deliveries = append(deliveries, d)
	}

	r := V1AuthRequest{
		AuthMail: c.getAuthMail(),
		Delivery: deliveries,
	}

	return c.deleteReservation(ctx, r, host)
}

func (c *V1Client) DeleteReservationByRequestIDs(ctx context.Context, requestIDs []string, host string) (resp V1MailResult, err error) {
	deliveries := make([]Delivery, 0, len(requestIDs))
	for k, v := range requestIDs {
		var d Delivery
		d.ID = k
		d.RequestID = v
		deliveries = append(deliveries, d)
	}

	r := V1AuthRequest{
		AuthMail: c.getAuthMail(),
		Delivery: deliveries,
	}

	return c.deleteReservation(ctx, r, host)
}
