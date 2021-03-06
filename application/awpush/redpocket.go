package awpush

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"subcenter/infra"
	"subcenter/infra/conf"
	"subcenter/infra/dto"
	"subcenter/infra/log"
	"time"
)

func enterLiveRoom(roomId, cookie, csrf string) error {
	rawUrl := "https://api.live.bilibili.com/xlive/web-room/v1/index/roomEntryAction"
	data := url.Values{
		"room_id":    []string{roomId},
		"platform":   []string{"pc"},
		"csrf":       []string{csrf},
		"csrf_token": []string{csrf},
		"visit_id":   []string{""},
	}
	body, err := infra.PostFormWithCookie(rawUrl, cookie, data)
	if err != nil {
		log.Error("PostFormWithCookie error: %v, raw data: %v", err, data)
		return err
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
		return err
	}
	if resp.Code != 0 {
		err = errors.New("response error")
		log.Error("EnterLiveRoom error: %v, resp: %v", err, resp)
		return err
	}
	return nil
}

// joinRedPocket refers to bilibili live lottery
func joinRedPocket(client *AWPushClient, redPocket dto.RedPocketMsg) {
	rawUrl := "https://api.live.bilibili.com/xlive/lottery-interface/v1/popularityRedPocket/RedPocketDraw"
	var roomId string
	switch val := redPocket.Data.RoomID.(type) {
	case string:
		roomId = val
	default:
		data := val.(float64)
		roomId = fmt.Sprintf("%d", int64(data))
	}
	data := url.Values{
		"ruid":       []string{fmt.Sprint(redPocket.Data.UID)},
		"room_id":    []string{roomId},
		"lot_id":     []string{fmt.Sprint(redPocket.Data.LotteryID)},
		"spm_id":     []string{"444.8.red_envelope.extract"},
		"session_id": []string{""},
		"jump_from":  []string{""},
	}
	for _, user := range conf.BiliConf.Users {
		if err := enterLiveRoom(roomId, user.Cookie, user.Csrf); err != nil {
			log.Info("User %d enter live room %s error", user.Uid, roomId)
			continue
		}
		timer := time.NewTimer(time.Second)
		<-timer.C
		body, err := infra.PostFormWithCookie(rawUrl, user.Cookie, data)
		if err != nil {
			log.Error("PostFormWithCookie error: %v, raw data: %v", err, data)
			continue
		}
		var resp dto.BiliBaseResp
		if err = json.Unmarshal(body, &resp); err != nil {
			log.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
		}
		if resp.Code == 0 {
			log.Info("User %d join redpocket %d success",
				user.Uid, redPocket.Data.LotteryID)
		} else {
			log.Info("User %d join redpocket %d failed because %s",
				user.Uid, redPocket.Data.LotteryID, resp.Message)
		}
	}
	go listenRoom(roomId)
}

// HandleRedPocket deal with red pocket message
func HandleRedPocket(client *AWPushClient, msg []byte) error {
	var redPocket dto.RedPocketMsg
	if err := json.Unmarshal(msg, &redPocket); err != nil {
		log.Error("Unmarshal RedPocketMsg error: %v, raw data: %s", err, string(msg))
		client.sleep.Reset(time.Microsecond)
		return err
	}
	client.sleep.Reset(time.Microsecond)
	// go joinRedPocket(client, redPocket)
	return nil
}
