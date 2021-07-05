package views

import (
	"fmt"

	"github.com/ewenquim/horkruxes/model"
	"github.com/ewenquim/horkruxes/service"
)

type PageData struct {
	Server   ServerData
	Messages []model.Message
	PageInfo PageInfo
}

type ServerData struct {
	Name string
	IP   string
	Info string
}

type PageInfo struct {
	Title           string
	SubTitle        string
	PostToMessageID string
}

// Get Local and online messages, checks validity and return view
func GetMessagesAndMainPageInfo(s service.Service) PageData {

	// Get local messages
	messages := model.GetMessagesFromDB(s)

	// Get other pods messages
	remoteMessages := getMessagesFrom("/api/message", s.ServerConfig.PublicPods...)
	messages = append(messages, remoteMessages...)

	messages = SortByDate(messages)

	// Inject view
	return PageData{
		Messages: CleanMessagesClientSide(messages),
		Server:   ServerData{Name: s.ServerConfig.Name, IP: s.ServerConfig.URL, Info: s.ServerConfig.Info},
		PageInfo: PageInfo{Title: "All Messages"},
	}
}

// Get Local and online messages, checks validity and return view
func GetAuthorMessagesAndMainPageInfo(s service.Service, pubKey string) PageData {

	// Get local messages
	messages := model.GetMessagesFromAuthor(s, pubKey)

	// Get other pods messages
	remoteMessages := getMessagesFrom("/api/message/user/"+pubKey, s.ServerConfig.PublicPods...)
	messages = append(messages, remoteMessages...)

	messages = SortByDate(messages)

	// Inject view
	return PageData{
		Messages: CleanMessagesClientSide(messages),
		Server:   ServerData{Name: s.ServerConfig.Name, IP: s.ServerConfig.URL, Info: s.ServerConfig.Info},
		PageInfo: PageInfo{Title: "Author", SubTitle: pubKey},
	}
}

// Get Local and online messages, checks validity and return view
func GetCommentsAndMainPageInfo(s service.Service, messageID string) PageData {

	messages := []model.Message{}

	// Get local messages
	messages = append(messages, model.GetCommentsTo(s, messageID)...)

	// Get other pods messages
	remoteMessages := getMessagesFrom("/api/message/comments/"+messageID, s.ServerConfig.PublicPods...)
	messages = append(messages, remoteMessages...)

	messages = SortByDate(messages)

	op, err := model.GetMessageFromDB(s, messageID)
	if err != nil {
		fmt.Println("err:", err)
	}
	messages = append(messages, CleanSingleMessageClientSide(op))
	messages = CleanMessagesClientSide(messages)

	// Inject view
	return PageData{
		Messages: messages,
		Server:   ServerData{Name: s.ServerConfig.Name, IP: s.ServerConfig.URL, Info: s.ServerConfig.Info},
		PageInfo: PageInfo{Title: "Comments", SubTitle: messageID, PostToMessageID: messageID},
	}
}