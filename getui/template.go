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

type Template struct {
	Message Message `json:"message"`
}

type NotificationTemplate struct {
	Template
	Notification Notification `json:"notification"`
}

type TransmissionTemplate struct {
	Template
	Transmission Transmission `json:"transmission"`
	PushInfo     PushInfo     `json:"push_info"`
}

func (t *Template) EnsureMessageAppKey(AppKey string) {
	if len(t.Message.AppKey) == 0 {
		t.Message.AppKey = AppKey
	}
}

func (t *Template) EnsureMsgtype(msgtype string) {
	if len(t.Message.MsgType) == 0 {
		t.Message.MsgType = msgtype
	}
}

func (t *Template) EnsureOfflineExpireTime() {
	if t.Message.OfflineExpireTime == 0 {
		t.Message.OfflineExpireTime = 2000000
	}
}
