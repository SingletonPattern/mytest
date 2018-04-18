/*
 *    Copyright 2016-2018 Li ZongZe
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package wechat

import "github.com/labstack/gommon/log"

// EventType 事件类型
type EventType string

// RecvEvent 事件消息
type RecvEvent interface {
	RecvMsg
}

// EventBase 事件基础类
type EventBase struct {
	RecvMsg
	ToUserName   string    // 开发者微信号
	FromUserName string    // 发送方帐号（一个OpenID）
	CreateTime   string    // 消息创建时间（整型）
	MsgType      MsgType   // 消息类型，event
	Event        EventType // 事件类型
}

// NewRecvEvent 把通用 struct 转化成相应类型的 struct
func NewRecvEvent(msg *Message) RecvEvent {
	switch msg.Event {
	case EventTypeSubscribe:
		return NewEventSubscribe(msg)
	case EventTypeUnsubscribe:
		return NewEventSubscribe(msg)
	case EventTypeLocation:
		return NewEventLocation(msg)
	case EventTypeClick:
		return NewEventClick(msg)
	case EventTypeView:
		return NewEventView(msg)
	case EventTypeScancodePush:
		return NewEventScancodePush(msg)
	case EventTypeScancodeWaitmsg:
		return NewEventScancodeWaitmsg(msg)
	case EventTypePicSysphoto:
		return NewEventPicSysphoto(msg)
	case EventTypePicPhotoOrAlbum:
		return NewEventPicPhotoOrAlbum(msg)
	case EventTypePicWeixin:
		return NewEventPicWeixin(msg)
	case EventTypeLocationSelect:
		return NewEventLocationSelect(msg)
	case EventTypeQualificationVerifySuccess:
		return NewEventQualificationVerifySuccess(msg)
	case EventTypeQualificationVerifyFail:
		return NewEventQualificationVerifyFail(msg)
	case EventTypeNamingVerifySuccess:
		return NewEventNamingVerifySuccess(msg)
	case EventTypeNamingVerifyFail:
		return NewEventNamingVerifyFail(msg)
	case EventTypeAnnualRenew:
		return NewEventAnnualRenew(msg)
	case EventTypeVerifyExpired:
		return NewEventVerifyExpired(msg)
	default:
		log.Errorf("unexpected receive EventType: %s", msg.Event)
		return nil
	}
}
