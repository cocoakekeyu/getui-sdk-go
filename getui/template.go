package getui

type Message struct {
	AppKey            string `json:"appkey"`
	IsOffline         bool   `json:"is_offline"`
	OfflineExpireTime uint32 `json:"offline_expire_time"`
	MsgType           string `json:"msgtype"`
}

type Style struct {
	Type        int    `json:"type"`
	Text        string `json:"text"`
	Title       string `json:"title"`
	Logo        string `json:"logo"`
	LogoUrl     string `json:"logourl,omitempty"`
	IsRing      string `json:"is_ring,omitempty"`
	IsVibrate   string `json:"is_virate,omitempty"`
	IsClearable string `json:"is_clearable,omitempty"`
}

type Notification struct {
	Style               Style  `json:"style"`
	TransmissionType    bool   `json:"transmission_type"`
	TransmissionContent string `json:"transmission_content"`
}

type Transmission struct {
	TransmissionType    bool   `json:"transmission_type"`
	TransmissionContent string `json:"transmission_content"`
}

type Link struct {
	Style Style  `json:"style"`
	Url   string `json:"url"`
}

type PushInfo struct {
	Aps struct {
		Alert struct {
			Title string `json:"title,omitempty"`
			Body  string `json:"body,omitempty"`
		} `json:"alert"`
		AutoBadge        string `json:"autoBadge,omitempty"`
		ContentAvailable int    `json:"content-available,omitempty"`
	} `json:"aps"`
	Payload string `json:"payload"`
}

type Condition struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	OptValue string `json:"opt_type"`
}

type TemplateInterface interface {
	TemplateMap() map[string]interface{}
	EnsureTemplateValue(string)
}

type Template struct {
	Message Message `json:"message"`
}

func (t *Template) EnsureTemplateValue(AppKey string) {
	if len(t.Message.AppKey) == 0 {
		t.Message.AppKey = AppKey
	}

	if t.Message.OfflineExpireTime == 0 {
		t.Message.OfflineExpireTime = 2000000
	}
}

type NotificationTemplate struct {
	Template
	Notification Notification `json:"notification"`
}

func (t *NotificationTemplate) EnsureTemplateValue(AppKey string) {
	if len(t.Message.MsgType) == 0 {
		t.Message.MsgType = "notification"
	}
	t.Template.EnsureTemplateValue(AppKey)
}

func (t *NotificationTemplate) TemplateMap() (result map[string]interface{}) {
	result = make(map[string]interface{})
	result["message"] = t.Message
	result["notification"] = t.Notification
	return
}

type TransmissionTemplate struct {
	Template
	Transmission Transmission `json:"transmission"`
	PushInfo     PushInfo     `json:"push_info"`
}

func (t *TransmissionTemplate) TemplateMap() (result map[string]interface{}) {
	result = make(map[string]interface{})
	result["message"] = t.Message
	result["transmission"] = t.Transmission
	result["push_info"] = t.PushInfo
	return
}

func (t *TransmissionTemplate) EnsureTemplateValue(AppKey string) {
	if len(t.Message.MsgType) == 0 {
		t.Message.MsgType = "transmission"
	}
	t.Template.EnsureTemplateValue(AppKey)
}

type LinkTemplate struct {
	Template
	Link Link `json:"link"`
}

func (t *LinkTemplate) TemplateMap() (result map[string]interface{}) {
	result = make(map[string]interface{})
	result["message"] = t.Message
	result["link"] = t.Link
	return
}

func (t *LinkTemplate) EnsureTemplateValue(AppKey string) {
	if len(t.Message.MsgType) == 0 {
		t.Message.MsgType = "link"
	}
	t.Template.EnsureTemplateValue(AppKey)
}

type BatchMessageTemplate struct {
	Template TemplateInterface
	CID      string
}

func (t *BatchMessageTemplate) TemplateMap() (result map[string]interface{}) {
	result = t.Template.TemplateMap()
	result["cid"] = t.CID
	return result
}

func (t *BatchMessageTemplate) EnsureTemplateValue(AppKey string) {
	t.Template.EnsureTemplateValue(AppKey)
}

func NewNotificationTemplate(AppKey string) *NotificationTemplate {
	t := new(NotificationTemplate)
	t.EnsureTemplateValue(AppKey)
	return t
}

func NewTransmissionTemplate(AppKey string) *TransmissionTemplate {
	t := new(TransmissionTemplate)
	t.EnsureTemplateValue(AppKey)
	return t
}

func NewLinkTemplate(AppKey string) *LinkTemplate {
	t := new(LinkTemplate)
	t.EnsureTemplateValue(AppKey)
	return t
}

func NewBatchMessageTemplate(AppKey string, Template TemplateInterface, CID string) *BatchMessageTemplate {
	t := new(BatchMessageTemplate)
	t.Template = Template
	t.CID = CID
	t.EnsureTemplateValue(AppKey)
	return t
}
