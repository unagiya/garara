package garara

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	*http.Client
	v1User    string
	v1Pass    string
	SiteID    string
	ServiceID string
}

func NewDefaultClient() *Client {
	return &Client{Client: http.DefaultClient, v1User: "", v1Pass: "", SiteID: "", ServiceID: ""}
}
func (c *Client) SetV1User(user string) {
	c.v1User = user
}

func (c *Client) SetV1Pass(password string) {
	c.v1Pass = password
}

func (c *Client) V1SendQueueMode(r V1MailRequest) (string, error) {
	rb, err := xml.Marshal(r)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://bvam001.am.arara.com/tm/lpmail_qmode.php", strings.NewReader(string(rb)))
	req.Header.Set("Content-Type", "application/xml")
	req.Header["X-AutomailUser"] = []string{c.v1User}
	req.Header["X-AutomailPassword"] = []string{c.v1Pass}
	req.Header["X-AutomailUseSite"] = []string{c.SiteID}
	req.Header["X-AutomailUseService"] = []string{c.ServiceID}

	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil

}
