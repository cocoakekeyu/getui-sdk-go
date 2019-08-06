package getui

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/cocoakekeyu/getui-sdk-go/utils"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	AUTH_TOKEN_URL_TPL        = "https://restapi.getui.com/v1/%s/auth_sign"
	CLOSE_AUTH_URL_TPL        = "https://restapi.getui.com/v1/%s/auth_close"
	PUSH_SINGLE_URL_TPL       = "https://restapi.getui.com/v1/%s/push_single"
	PUSH_APP_URL_TPL          = "https://restapi.getui.com/v1/%s/push_app"
	USER_STATUS_URL_TPL       = "https://restapi.getui.com/v1/%s/user_status/%s"
	STOP_TASK_URL_TPL         = "https://restapi.getui.com/v1/%s/stop_task/%s"
	SVAE_LIST_BODY_URL_TPL    = "https://restapi.getui.com/v1/%s/save_list_body"
	PUSH_LIST_URL_TPL         = "https://restapi.getui.com/v1/%s/push_list"
	PUSH_SINGEL_BATCH_URL_TPL = "https://restapi.getui.com/v1/%s/push_single_batch"
)

type Client struct {
	AppID                string
	AppKey               string
	MasterSecret         string
	authToken            string
	lastRefreshTokenTime time.Time
}

func NewGeTuiClient(AppID, AppKey, MasterSecret string) (*Client, error) {
	var c = Client{
		AppID:        AppID,
		AppKey:       AppKey,
		MasterSecret: MasterSecret,
	}
	err := c.RefreshAuthToken()

	if err == nil {
		// Ticker to refresh authtoken
		go func() {
			timer := time.NewTicker(20 * time.Hour)
			for range timer.C {
				c.RefreshAuthToken()
			}
		}()
	}

	return &c, err
}

func (c *Client) String() string {
	return fmt.Sprintf("GeTuiClient(AppId: %s, authToken: %s)", c.AppID, c.authToken)
}

func (c *Client) RefreshAuthToken() error {
	var ts = fmt.Sprintf("%d", time.Now().UnixNano()/int64(1000000))
	var sign = fmt.Sprintf("%x", sha256.Sum256([]byte(c.AppKey+ts+c.MasterSecret)))

	data := map[string]string{
		"appkey":    c.AppKey,
		"timestamp": ts,
		"sign":      sign,
	}

	payload, _ := json.Marshal(data)

	url := fmt.Sprintf(AUTH_TOKEN_URL_TPL, c.AppID)

	req, err := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewReader(payload)))
	if err != nil {
		return fmt.Errorf("[RefreshAuthToken] create http request failed, err: %s", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("[RefreshAuthToken] request auth token failed, err: %s", err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("[RefreshAuthToken] read response body failed, err: %s", err)
	}

	var result map[string]string

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return fmt.Errorf("[RefreshAuthToken] parse response body failed, err: %s", err)
	}

	var ok bool
	c.authToken, ok = result["auth_token"]
	if !ok {
		return fmt.Errorf("[RefreshAuthToken] get auth token failed, err: %s", respBody)
	}
	c.lastRefreshTokenTime = time.Now()

	return nil
}

func (c *Client) CloseAuth() (map[string]string, error) {
	url := fmt.Sprintf(CLOSE_AUTH_URL_TPL, c.AppID)
	return c.httpRequestPost(url, nil)
}

func (c *Client) PushToSingle(t TemplateInterface, CID string) (map[string]string, error) {
	t.EnsureTemplateValue(c.AppKey)

	body := t.TemplateMap()
	body["cid"] = CID
	body["requestid"] = utils.GenerateRequestID()

	payload, _ := json.Marshal(body)

	url := fmt.Sprintf(PUSH_SINGLE_URL_TPL, c.AppID)

	return c.httpRequestPost(url, ioutil.NopCloser(bytes.NewReader(payload)))
}

func (c *Client) PushToApp(t TemplateInterface, condition []Condition) (map[string]string, error) {
	t.EnsureTemplateValue(c.AppKey)

	body := t.TemplateMap()
	body["requestid"] = utils.GenerateRequestID()

	if len(condition) > 0 {
		body["condition"] = condition
	}

	payload, _ := json.Marshal(body)

	url := fmt.Sprintf(PUSH_APP_URL_TPL, c.AppID)

	return c.httpRequestPost(url, ioutil.NopCloser(bytes.NewReader(payload)))
}

func (c *Client) SaveListBody(t TemplateInterface) (map[string]string, error) {
	body := t.TemplateMap()
	payload, _ := json.Marshal(body)
	url := fmt.Sprintf(SVAE_LIST_BODY_URL_TPL, c.AppID)
	return c.httpRequestPost(url, ioutil.NopCloser(bytes.NewReader(payload)))
}

func (c *Client) PushToList(CID []string, TaskID string, NeedDetail bool) (map[string]string, error) {
	body := make(map[string]interface{})
	body["cid"] = CID
	body["taskid"] = TaskID
	body["need_detail"] = NeedDetail

	payload, _ := json.Marshal(body)
	url := fmt.Sprintf(PUSH_LIST_URL_TPL, c.AppID)
	return c.httpRequestPost(url, ioutil.NopCloser(bytes.NewReader(payload)))
}

func (c *Client) PushSingleBatch(templates []TemplateInterface, NeedDetail bool) (map[string]string, error) {
	body := make(map[string]interface{})
	MsgList := make([]interface{}, 0)
	for _, t := range templates {
		t.EnsureTemplateValue(c.AppKey)
		m := t.TemplateMap()
		m["requestid"] = utils.GenerateRequestID()
		MsgList = append(MsgList, m)
	}
	body["msg_list"] = MsgList
	body["need_detail"] = NeedDetail
	payload, _ := json.Marshal(body)
	fmt.Println(string(payload))
	url := fmt.Sprintf(PUSH_SINGEL_BATCH_URL_TPL, c.AppID)
	return c.httpRequestPost(url, ioutil.NopCloser(bytes.NewReader(payload)))
}

func (c *Client) UserStatus(CID string) (map[string]string, error) {
	url := fmt.Sprintf(USER_STATUS_URL_TPL, c.AppID, CID)
	return c.httpRequestGet(url)
}

func (c *Client) StopTask(TaskID string) (map[string]string, error) {
	url := fmt.Sprintf(STOP_TASK_URL_TPL, c.AppID, TaskID)
	return c.httpRequestDelete(url)
}

func (c *Client) httpRequest(method string, url string, body io.Reader) (map[string]string, error) {
	var req *http.Request
	var err error
	switch method {
	case "GET":
		req, err = http.NewRequest("GET", url, nil)
	case "POST":
		req, err = http.NewRequest("POST", url, body)
	case "DELETE":
		req, err = http.NewRequest("DELETE", url, nil)
	}

	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["authtoken"] = []string{c.authToken}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[HttpRequest] request url(%s) failed, err: %s", url, err)
	}

	defer resp.Body.Close()

	var unmarshal map[string]interface{}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[HttpRequest] read response body failed, err: %s", err)
	}

	err = json.Unmarshal(data, &unmarshal)
	if err != nil {
		return nil, fmt.Errorf("[HttpRequest] json umarshal response body failed, err: %s", err)
	}

	var result = make(map[string]string)
	for key, value := range unmarshal {
		result[key] = fmt.Sprintf("%v", value)
	}

	return result, nil
}

func (c *Client) httpRequestGet(url string) (map[string]string, error) {
	return c.httpRequest("GET", url, nil)
}

func (c *Client) httpRequestDelete(url string) (map[string]string, error) {
	return c.httpRequest("DELETE", url, nil)
}

func (c *Client) httpRequestPost(url string, body io.Reader) (map[string]string, error) {
	return c.httpRequest("POST", url, body)
}
