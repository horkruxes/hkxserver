package views

import (
	"fmt"

	"github.com/horkruxes/hkxserver/client"
	"github.com/horkruxes/hkxserver/model"
	"github.com/horkruxes/hkxserver/service"
)

type ClientData struct {
	PublicPods  bool
	PrivatePods bool
	PodsString  string
}

// PageData represent information sent to the page
type PageData struct {
	Server     service.GeneralConfig
	Client     ClientData // Filters to feed back loop
	TopMessage model.Message
	Messages   []model.Message
	PageInfo   PageInfo
}

type PageInfo struct {
	MainPage        bool
	CommentsPage    bool
	Title           string
	SubTitle        string
	PostToMessageID string
}

// Get Local and online messages, checks validity and return view
func GetMessagesAndMainPageInfo(s service.Service) PageData {

	// Get local messages
	messages := model.GetMessagesFromDB(s)

	// Get other pods messages
	if s.ClientConfig.PublicPods {
		remoteMessages := client.GetMessagesFrom(s.GeneralConfig.TrustedPods, "/api/message")
		messages = append(messages, remoteMessages...)
	}

	messages = model.SortByDate(messages)

	messages = client.CleanMessagesClientSide(messages)
	for i, msg := range messages {
		if len(msg.Content) > 250 {
			messages[i].Content = msg.Content[:250] + "..."
		}
	}

	// Inject view
	return PageData{
		Messages: messages,
		Server:   s.GeneralConfig,
		// Client:   ClientData{PublicPods: s.ClientConfig.PublicPods, PrivatePods: s.ClientConfig.SpecificPods, PodsString: s.ClientConfig.SpecificPodsListString},
		PageInfo: PageInfo{MainPage: true, Title: "All Messages"},
	}
}

// Get Local and online messages, checks validity and return view
// Pubkey is in base64 form
func GetAuthorMessagesAndMainPageInfo(s service.Service, pubKey string) PageData {

	// Get local messages
	messages := model.GetMessagesFromAuthor(s, pubKey)

	// Get other pods messages
	remoteMessages := client.GetMessagesFrom(s.GeneralConfig.TrustedPods, "/api/user/"+pubKey)
	messages = append(messages, remoteMessages...)

	messages = model.SortByDate(messages)

	// Inject view
	return PageData{
		Messages: client.CleanMessagesClientSide(messages),
		Server:   s.GeneralConfig,
		PageInfo: PageInfo{Title: "Author", SubTitle: pubKey},
	}
}

// Get Local and online messages, checks validity and return view
func GetCommentsAndMainPageInfo(s service.Service, messageID string) PageData {

	messages := []model.Message{}

	// Get local comments
	messages = append(messages, model.GetCommentsTo(s, messageID)...)

	// Get other pods comments
	remoteMessages := client.GetMessagesFrom(s.GeneralConfig.TrustedPods, "/api/comments/"+messageID)
	messages = append(messages, remoteMessages...)

	messages = model.SortByDate(messages)
	messages = client.CleanMessagesClientSide(messages)

	// Try to get local OP
	op, err := model.GetMessageFromDB(s, messageID)
	if err != nil {
		fmt.Println("err:", err)
		// Asks other pods to get comment
		remoteOPs := client.GetSingleMessageFromEachPod(s.GeneralConfig.TrustedPods, "/api/message/"+messageID)
		fmt.Println(remoteOPs)
		if len(remoteOPs) > 0 {
			op = remoteOPs[0]
		}
	}
	op = client.CleanSingleMessageClientSide(op)

	// Inject view
	return PageData{
		TopMessage: op,
		Messages:   messages,
		Server:     s.GeneralConfig,
		PageInfo:   PageInfo{Title: "Comments", SubTitle: messageID, PostToMessageID: messageID, CommentsPage: true},
	}
}
