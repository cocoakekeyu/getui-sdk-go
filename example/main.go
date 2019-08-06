package main

import (
	"fmt"
	"github.com/cocoakekeyu/getui-sdk-go/getui"
)

func main() {
	// example
	var (
		AppID        = "getui appid"
		AppKey       = "getui appkey"
		MasterSecret = "getui mastersecret"
	)
	var err error
	var result map[string]string
	var CID = "8e42279b3dfdbd335079b18f10b7ce9c"

	client, err := getui.NewGeTuiClient(AppID, AppKey, MasterSecret)
	if err != nil {
		fmt.Printf("[Main] init getui client failed, err: %s\n", err)
	}

	fmt.Printf("init a getui client: %v\n", client)

	// create transmission message template
	t := getui.NewTransmissionTemplate(client.AppKey)
	t.Message.IsOffline = true
	t.Transmission.TransmissionContent = "test"
	t.PushInfo.Aps.Alert.Title = "test"
	t.PushInfo.Aps.Alert.Body = "test"
	t.PushInfo.Aps.AutoBadge = "+1"
	t.PushInfo.Aps.ContentAvailable = 1

	// push to single
	result, err = client.PushToSingle(t, CID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("push to single result: %v\n", result)

	// push to app
	result, err = client.PushToApp(t, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("push to app result: %v\n", result)

	// create notification message template
	n := getui.NewNotificationTemplate(client.AppKey)
	n.Notification.Style.Text = "test"
	n.Notification.Style.Title = "test"

	// link message template
	// l := getui.NewLinkTemplate(client.AppKey)
	// l.Link.Style.Text = "test"
	// l.Link.Style.Title = "test"
	// l.Link.Url = "http://www.baidu.com"

	// push to single
	result, err = client.PushToSingle(n, CID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("push to single result: %v\n", result)

	// user status
	result, err = client.UserStatus(CID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("user status result: %v\n", result)

	// save list body
	result, err = client.SaveListBody(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("save list result: %v\n", result)
	taskid := result["taskid"]

	// push to list
	result, err = client.PushToList([]string{CID}, taskid, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("push to list result: %v\n", result)

	// push single batch
	b1 := getui.NewBatchMessageTemplate(client.AppKey, t, CID)
	b2 := getui.NewBatchMessageTemplate(client.AppKey, n, CID)

	templates := []getui.TemplateInterface{b1, b2}

	result, err = client.PushSingleBatch(templates, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("push single list result: %v\n", result)

}
