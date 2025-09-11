package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// 用 sync.Map 替换 Redis，用于事件去重和过期

// 带过期自动清理（简单实现：存储时间戳，定期清理）
type eventEntry struct {
	CreatedAt time.Time
}

// 启动后台清理已过期事件
func startCacheCleanup(expire time.Duration) {
	go func() {
		for {
			time.Sleep(time.Hour)
			now := time.Now()
			eventCache.Range(func(key, value interface{}) bool {
				entry := value.(eventEntry)
				if now.Sub(entry.CreatedAt) > expire {
					eventCache.Delete(key)
				}
				return true
			})
		}
	}()
}

func CheckEventID(eventID string) (bool, error) {
	key := fmt.Sprintf("event:%s", eventID)
	_, exists := eventCache.Load(key)
	if exists {
		return false, nil
	}
	eventCache.Store(key, eventEntry{CreatedAt: time.Now()})
	return true, nil
}

func GenerateEventID() string {
	return uuid.New().String()
}

func AffilikaPostback(action, trackingID string) error {
	var url string
	switch action {
	case "reg":
		url = fmt.Sprintf("https://pwa.market/pwa-pb/?action=reg&tracking_id=%s", trackingID)
	case "dep":
		url = fmt.Sprintf("https://pwa.market/pwa-pb/?action=dep&tracking_id=%s", trackingID)
	default:
		return nil
	}
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func FbConversionAPI(action, eventID, trackingID, ip string) error {
	accessToken := os.Getenv("FB_TOKEN")
	pixelID := os.Getenv("FB_PIXEL_ID")
	eventName := ""
	switch action {
	case "reg":
		eventName = "CompleteRegistration"
	case "dep":
		eventName = "Purchase"
	default:
		return nil
	}
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/events?access_token=%s", pixelID, accessToken)
	userData := map[string]interface{}{
		"external_id":       Sha256(trackingID),
		"client_ip_address": ip,
	}
	payload := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"event_name":    eventName,
				"event_time":    time.Now().Unix(),
				"event_id":      eventID,
				"action_source": "website",
				"user_data":     userData,
			},
		},
	}
	body, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func Sha256(str string) string {
	sum := sha256.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}

// 在 main() 入口处调用：startCacheCleanup(24*time.Hour)
