package views

import (
	"fmt"

	"github.com/ewenquim/horkruxes/model"
	"github.com/ewenquim/horkruxes/service"
)

type PageData struct {
	Server     ServerData
	TopMessage model.Message
	Messages   []model.Message
	PageInfo   PageInfo
}

type ServerData struct {
	Name string
	IP   string
	Info string
}

type PageInfo struct {
	MainPage        bool
	Title           string
	SubTitle        string
	PostToMessageID string
}

// Get Local and online messages, checks validity and return view
func GetMessagesAndMainPageInfo(s service.Service) PageData {

	// Get local messages
	messages := model.GetMessagesFromDB(s)

	// Get other pods messages
	remoteMessages := getMessagesFrom(s, "/api/message")
	messages = append(messages, remoteMessages...)

	messages = model.SortByDate(messages)

	// Inject view
	return PageData{
		Messages: model.CleanMessagesClientSide(messages),
		Server:   ServerData{Name: s.ServerConfig.Name, IP: s.ServerConfig.URL, Info: s.ServerConfig.Info},
		PageInfo: PageInfo{MainPage: true, Title: "All Messages"},
	}
}

// Get Local and online messages, checks validity and return view
// Pubkey is in base64 form
func GetAuthorMessagesAndMainPageInfo(s service.Service, pubKey string) PageData {

	// Get local messages
	messages := model.GetMessagesFromAuthor(s, pubKey)

	// Get other pods messages
	remoteMessages := getMessagesFrom(s, "/api/user/"+service.Base64ToSafeURL(pubKey))
	messages = append(messages, remoteMessages...)

	messages = model.SortByDate(messages)

	// Inject view
	return PageData{
		Messages: model.CleanMessagesClientSide(messages),
		Server:   ServerData{Name: s.ServerConfig.Name, IP: s.ServerConfig.URL, Info: s.ServerConfig.Info},
		PageInfo: PageInfo{Title: "Author", SubTitle: pubKey},
	}
}

// Get Local and online messages, checks validity and return view
func GetCommentsAndMainPageInfo(s service.Service, messageID string) PageData {

	messages := []model.Message{}

	// Get local comments
	messages = append(messages, model.GetCommentsTo(s, messageID)...)

	// Get other pods comments
	remoteMessages := getMessagesFrom(s, "/api/comments/"+messageID)
	messages = append(messages, remoteMessages...)

	messages = model.SortByDate(messages)
	messages = model.CleanMessagesClientSide(messages)

	// Try to get local OP
	op, err := model.GetMessageFromDB(s, messageID)
	if err != nil {
		fmt.Println("err:", err)
		// Asks other pods to get comment
		remoteOPs := getSingleMessageFromEachPod(s, "/api/message/"+messageID)
		fmt.Println(remoteOPs)
		if len(remoteOPs) > 0 {
			op = remoteOPs[0]
		}
	}
	op = model.CleanSingleMessageClientSide(op)

	// Inject view
	return PageData{
		TopMessage: op,
		Messages:   messages,
		Server:     ServerData{Name: s.ServerConfig.Name, IP: s.ServerConfig.URL, Info: s.ServerConfig.Info},
		PageInfo:   PageInfo{Title: "Comments", SubTitle: messageID, PostToMessageID: messageID},
	}
}