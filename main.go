package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Message struct {
	AuthorPubKey  []byte
	AuthorRep     string
	Name          string
	Content       string
	Signature     []byte
	SignatureRep  string
	SignedCorrect bool
	Color         string
}

type SignMessage struct {
	AuthorPubKey string // base64-encoded
	AuthorSecKey string // base64-encoded
	Content      string
	Signature    string
}

type Data struct {
	PageTitle string
	Messages  []Message
	KeyGen    KeyGen
}

func main() {
	// One time
	// Generates public and secret key
	pub, sec, _ := ed25519.GenerateKey(nil)
	pub2, sec2, _ := ed25519.GenerateKey(nil)
	// Sign every message TODO: possibility to do it in the site + advices
	// Needs the message, the pub and secret key. Outputs the signature
	// When the pod receives the message (content, public and signature) possibility to verify
	signature := ed25519.Sign(sec, []byte("hey guys, hello world"+string(pub)))
	signature2 := ed25519.Sign(sec, []byte("my first secure tweet"+string(pub)))
	signature3 := ed25519.Sign(sec2, []byte("lorem <strong>ipsum</strong>i skip\n lines"+string(pub2)))
	println("--------")
	data := Data{
		PageTitle: "Horkruxes",
		Messages: []Message{
			{Name: "ewen", AuthorPubKey: pub, Content: "hey guys, hello world", Signature: signature},
			{Name: "chloe.sa", AuthorPubKey: pub, Content: "my first secure tweet", Signature: signature2},
			{Name: "seraph", AuthorPubKey: []byte("2eb1ek2ed9g"), Content: `lorem https://ewen.quimerch.com/ <strong>ipsum</strong>i skip\n lines`},
			{Name: "marius", AuthorPubKey: pub2, Content: "lorem <strong>ipsum</strong>i skip\n lines", Signature: signature3},
		},
		KeyGen: KeyGen{},
	}
	// fmt.Printf("%+v", data.Messages)

	// serv
	mainTmpl := template.Must(template.ParseFiles("templates/root.html", "templates/main/_base.html", "templates/main/main.html", "templates/main/pods.html", "templates/main/keys.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		for i, message := range data.Messages {
			data.Messages[i].SignedCorrect = message.verifyOwnerShip()
			fmt.Println(data.Messages[i].SignedCorrect)
			data.Messages[i].AuthorRep = base64.StdEncoding.EncodeToString(message.AuthorPubKey)
			data.Messages[i].Color = colorFromString(string(message.AuthorPubKey))
			data.Messages[i].SignatureRep = base64.StdEncoding.EncodeToString(message.Signature)
		}
		mainTmpl.Execute(w, data)
	})

	// keygen
	keyGenTmpl := template.Must(template.ParseFiles("templates/root.html", "templates/keys/keys.html"))
	http.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("----------")
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		outputData := genKeys()
		
		// POST Form
		form := SignMessage{
			AuthorSecKey: r.FormValue("secret-key"),
			AuthorPubKey: r.FormValue("public-key"),
			Content:      r.FormValue("message"),
			Signature:    r.FormValue("signature"),
		}

		outputData.Content = form.Content
		fmt.Printf("FORM %+v\n", form)

		// GET Page
		if r.Method != http.MethodPost {
			keyGenTmpl.Execute(w, outputData)
			return
		}

		// Form "generate"
		if form.Signature == "" {
			outputData.Sig = signMessage(form.AuthorSecKey, form.AuthorPubKey, form.Content)
			form.AuthorSecKey = ""
			keyGenTmpl.Execute(w, outputData)
		} else {
			// Form "verify"
			outputData.Valid = verifyFromString(form.AuthorPubKey, form.Signature, form.Content)
			fmt.Println("sig valid", outputData.Valid)
			form.Signature = ""
			keyGenTmpl.Execute(w, outputData)
		}
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)
}

func colorFromString(s string) string {
	// java String#hashCode
	str := strings.ToLower(s)
	hash := 0
	for i := 0; i < len(str); i++ {
		num, _ := strconv.Atoi(string(str[i]))
		hash = num + ((hash << 5) - hash)
	}
	base := int(math.Abs(float64(hash)))

	colors := []string{
		"pink",
		"#9b88ee",
		"GainsBoRo",
		"yellowGreen",
		"skyBlue",
		"salmon",
		"sandyBrown",
		"paleGreen",
		"paleTurquoise",
		"red",
	}

	return "fill:" + strings.ToLower(colors[trueMod(base, len(colors))])
}

func trueMod(n int, m int) int {
	return ((n % m) + m) % m
}
