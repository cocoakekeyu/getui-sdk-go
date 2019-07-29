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
	client, err := getui.NewGeTuiClient(AppID, AppKey, MasterSecret)
	if err != nil {
		fmt.Printf("[Main] init getui client failed, err: %s\n", err)
	}

	fmt.Printf("init a getui client: %v\n", client)

	fmt.Println("transmission message")
	t := getui.NewTransmissionTemplate(client.AppKey)
	t.Transmission.TransmissionContent = ""
	t.PushInfo.Aps.Alert.Title = "test"
	t.PushInfo.Aps.Alert.Body = "test"
	t.PushInfo.Aps.AutoBadge = "+1"
	t.PushInfo.Aps.ContentAvailable = 1

	CID := "cid"
	result, err := client.PushToSingle(t, CID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("push to single result: %v\n", result)

	result, err = client.PushToApp(t, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("push to app result: %v\n", result)

	fmt.Println("notification message")

	n := getui.NewNotificationTemplate(client.AppKey)
	n.Notification.Style.Text = "test"
	n.Notification.Style.Title = "test"

	result, err = client.PushToSingle(n, CID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("push to single result: %v\n", result)
}
