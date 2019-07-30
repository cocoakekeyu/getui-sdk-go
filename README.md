## 个推 Go SDK

个推简单的服务端 Go 语言实现的 SDK, 已可以使用单推、推送 APP 等功能。


## 用法

获取个推应用的 `appid`,  `appkey`, `mastersecret`

#### 创建个推 Client

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

#### 选择推送模版并填充数据

```go
//透传模版消息

t := getui.NewTransmissionTemplate(client.AppKey)
t.Message.IsOffline = true
t.Message.OfflineExpireTime = 2000000
t.Transmission.TransmissionContent = "test"
t.PushInfo.Aps.Alert.Title = "test"
t.PushInfo.Aps.Alert.Body = "test"
t.PushInfo.Aps.AutoBadge = "+1"
t.PushInfo.Aps.ContentAvailable = 1

// 通知模版消息
n := getui.NewNotificationTemplate(client.AppKey)
n.Message.IsOffline = false
n.Notification.Style.Text = "test"
n.Notification.Style.Title = "test"

// 链接模版消息
l := getui.NewLinkTemplate(client.AppKey)
l.Link.Style.Text = "test"
l.Link.Style.Title = "test"
l.Link.Url = "http://www.baidu.com"



```

#### 推送方式

```go

var CID = "cid"
var result map[string]string

// 单推 cid
result, err = client.PushToSingle(t, CID)
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

// 推送 CID 列表
result, err = client.SaveListBody(t)
if err != nil {
 	fmt.Println(err)
}
fmt.Printf("save list result: %v\n", result)
taskid := result["taskid"]

var cids = []string{CID}
result, err = client.PushToList(cids, taskid, false)
if err != nil {
 	fmt.Println(err)
}
fmt.Printf("push to list result: %v\n", result)

```

#### 其他接口

```go
// 查询用户状态
result, err = client.UserStatus(CID)
if err != nil {
 	fmt.Println(err)
}
fmt.Printf("user status result: %v\n", result)

// 停止群推任务
result, err = client.StopTask(CID)
if err != nil {
	fmt.Println(err)
}
fmt.Printf("stop task result: %v\n", result)

```

## 示例

见 example
