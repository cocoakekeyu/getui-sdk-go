package getui

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/cocoakekeyu/getui-sdk-go/getui/utils"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	AUTH_TOKEN_URL_TPL  = "https://restapi.getui.com/v1/%s/auth_sign"
	CLOSE_AUTH_URL_TPL  = "https://restapi.getui.com/v1/%s/auth_close"
	PUSH_SINGLE_URL_TPL = "https://restapi.getui.com/v1/%s/push_single"
	PUSH_APP_URL_TPL    = "https://restapi.getui.com/v1/%s/push_app"
	USER_STATUS_URL_TPL = "https://restapi.getui.com/v1/%s/user_status/%s"
	STOP_TASK_URL_TPL   = "https://restapi.getui.com/v1/%s/stop_task/%s"
)

func EnsureTemplateValue(t interface{}, AppKey string) error {

	vv := reflect.ValueOf(t)
	if vv.Kind() != reflect.Ptr {
		return fmt.Errorf("t is not a pointer value")
	}
	value := vv.Elem()
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("the value that t point to is not a struct value")
	}
	type_ := value.Type()

	// ensure Message AppKey
	vv.MethodByName("EnsureMessageAppKey").Call([]reflect.Value{reflect.ValueOf(AppKey)})

	msgtype := strings.ToLower(strings.TrimRight(type_.Name(), "Template"))
	vv.MethodByName("EnsureMsgtype").Call([]reflect.Value{reflect.ValueOf(msgtype)})
	vv.MethodByName("EnsureOfflineExpireTime").Call([]reflect.Value{})

	return nil

}

func BuildTemplateMap(t interface{}) map[string]interface{} {
	vv := reflect.ValueOf(t)
	value := vv.Elem()
	type_ := value.Type()

	var result = map[string]interface{}{}

	switch type_.Name() {
	case "NotificationTemplate":
		{
			tem := t.(*NotificationTemplate)
			result["message"] = tem.Message
			result["notification"] = tem.Notification
		}
	case "TransmissionTemplate":
		{
			tem := t.(*TransmissionTemplate)
			result["message"] = tem.Message
			result["transmission"] = tem.Transmission
			result["push_info"] = tem.PushInfo
		}
	}
	return result
}

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

func (c *Client) PushToSingle(t interface{}, CID string) (map[string]string, error) {
	err := EnsureTemplateValue(t, c.AppKey)
	if err != nil {
		return nil, err
	}

	body := BuildTemplateMap(t)
	body["cid"] = CID
	body["requestid"] = utils.GenerateRequestID()

	payload, _ := json.Marshal(body)

	url := fmt.Sprintf(PUSH_SINGLE_URL_TPL, c.AppID)

	return c.httpRequestPost(url, ioutil.NopCloser(bytes.NewReader(payload)))
}

func (c *Client) PushToApp(t interface{}, condition []Condition) (map[string]string, error) {
	err := EnsureTemplateValue(t, c.AppKey)
	if err != nil {
		return nil, err
	}
	body := BuildTemplateMap(t)
	body["requestid"] = utils.GenerateRequestID()

	if len(condition) > 0 {
		body["condition"] = condition
	}

	payload, _ := json.Marshal(body)

	url := fmt.Sprintf(PUSH_APP_URL_TPL, c.AppID)

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
		{
			req, err = http.NewRequest("GET", url, nil)
		}
	case "POST":
		{
			req, err = http.NewRequest("POST", url, body)
		}
	case "DELETE":
		{
			req, err = http.NewRequest("DELETE", url, nil)
		}
	}

	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["authtoken"] = []string{c.authToken}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[HttpRequest] request url(%s) failed, err: %s", url, err)
	}

	defer resp.Body.Close()

	var result map[string]string

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[HttpRequest] read response body failed, err: %s", err)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("[HttpRequest] json umarshal response body failed, err: %s", err)
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
