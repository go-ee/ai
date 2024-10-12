package api

import "context"

type Workflow struct {
	Chain []Plugin
}

func (o *Workflow) Execute() (err error) {
	session := &Session{}
	for _, plugin := range o.Chain {
		switch plugin.GetType() {
		case PluginTypeInput:
			input := plugin.(Input)
			var messages []*Message
			if messages, err = input.GetMessages(); err != nil {
				return
			}
			session.Append(messages...)
		case PluginTypeChatter:
			chatter := plugin.(Chatter)
			var messages []*Message
			if messages, err = chatter.Chat(context.Background(), session.GetChatMessages(), nil); err != nil {
				return
			}
			session.Append(messages...)
		case PluginTypeTransformer:
			messages := session.GetLastMessages()
			transformer := plugin.(Transformer)
			if messages, err = transformer.Transform(messages); err != nil {
				return
			}
			session.ReplaceLastMessages(messages)
		case PluginTypeOutput:
			output := plugin.(Output)
			if err = output.Output(session.Messages); err != nil {
				return
			}
		}
	}
	return
}
