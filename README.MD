## 个推 Go SDK

个推简单的服务端 Go 语言实现的 SDK, 已可以使用单推、推送 APP 等功能。


## 用法

获取个推应用的 `appid`,  `appkey`, `mastersecret`

#### 创建一个个推 Client

```go
var (
  AppID        = "getui appid"
  AppKey       = "getui appkey"
  MasterSecret = "getui mastersecret"
)
client, err := getui.NewGeTuiClient(AppID, AppKey, MasterSecret)
if err != nil {
  fmt.Printf("[Main] init getui client failed, err: %s\n", err)
}

```

#### 选择模版并填充数据

```
//透传模版消息

t := getui.NewTransmissionTemplate(client.AppKey)
t.Transmission.TransmissionContent = "test"
t.PushInfo.Aps.Alert.Title = "test"
t.PushInfo.Aps.Alert.Body = "test"
t.PushInfo.Aps.AutoBadge = "+1"
t.PushInfo.Aps.ContentAvailable = 1

// 通知模版消息
n := getui.NewNotificationTemplate(client.AppKey)
n.Notification.Style.Text = "test"
n.Notification.Style.Title = "test"


```

### 发送单推或者推APP

```

// 单推 cid
CID := "cid"
result, err := client.PushToSingle(t, CID)
if err != nil {
  fmt.Println(err)
}
fmt.Printf("push to single result: %v\n", result)

// 推 APP
result, err = client.PushToApp(t, nil)
if err != nil {
  fmt.Println(err)
}
fmt.Printf("push to app result: %v\n", result)


```

## 示例

见 example
