package api

import "fmt"

const ChatMessageRoleMeta = "meta"

type Session struct {
	Name     string
	Messages []*Message

	chatMessages []*Message

	lastMessages []*Message
}

func (o *Session) IsEmpty() bool {
	return len(o.Messages) == 0
}

func (o *Session) Append(messages ...*Message) {
	o.lastMessages = messages
	if o.chatMessages != nil {
		for _, message := range messages {
			o.Messages = append(o.Messages, message)
			o.appendChatMessage(message)
		}
	} else {
		o.Messages = append(o.Messages, messages...)
	}
}

func (o *Session) GetChatMessages() (ret []*Message) {
	if o.chatMessages == nil {
		o.chatMessages = []*Message{}
		for _, message := range o.Messages {
			o.appendChatMessage(message)
		}
	}
	ret = o.chatMessages
	return
}

func (o *Session) appendChatMessage(message *Message) {
	if message.Role != ChatMessageRoleMeta {
		o.chatMessages = append(o.chatMessages, message)
	}
}

func (o *Session) GetLastMessages() (ret []*Message) {
	ret = o.lastMessages
	return
}

func (o *Session) ReplaceLastMessages(newMessages []*Message) {
	o.lastMessages = newMessages
	lastMessagesLen := len(o.lastMessages)

	if lastMessagesLen > 0 {
		if len(o.Messages) >= lastMessagesLen {
			o.Messages = o.Messages[:len(o.Messages)-lastMessagesLen]
		} else {
			o.Messages = []*Message{}
		}

		if len(o.chatMessages) >= lastMessagesLen {
			o.chatMessages = o.chatMessages[:len(o.chatMessages)-lastMessagesLen]
		} else {
			o.chatMessages = []*Message{}
		}
	}

	o.Append(newMessages...)
}

func (o *Session) String() (ret string) {
	for _, message := range o.Messages {
		ret += fmt.Sprintf("\n--- \n[%v]\n\n%v", message.Role, message.Content)
	}
	return
}
