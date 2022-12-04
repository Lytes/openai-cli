package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// Define a struct to represent the configuration data.
type Config struct {
	SessionToken  string `json:"sessionToken"`
	MixpanelToken string `json:"mixpanelToken"`
}

func get_req_info(prompt string) (map[string]string, map[string]interface{}) {
	var header = make(map[string]string)
	header["Host"] = "chat.openai.com"
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:107.0) Gecko/20100101 Firefox/107.0"
	header["Accept"] = "*/*"
	header["Accept-Language"] = "en-US,en;q=0.5"
	header["Accept-Encoding"] = "gzip, deflate"
	header["Referer"] = "https://chat.openai.com/chat"
	header["Sec-Fetch-Dest"] = "empty"
	header["Sec-Fetch-Mode"] = "cors"
	header["Sec-Fetch-Site"] = "same-origin"
	header["Te"] = "trailers"

	// generate new uuid
	uuid := uuid.NewV4()
	uuidStr := uuid.String()

	// Create a JSON object
	jsonData := map[string]interface{}{
		"action":            "next",
		"parent_message_id": "ac93cda0-5d2d-489c-9ce3-4262e83a1481",
		"model":             "text-davinci-002-render",
	}

	// Create a JSON array of messages
	messages := []interface{}{
		map[string]interface{}{
			"id":   uuidStr,
			"role": "user",
			"content": map[string]interface{}{
				"content_type": "text",
				"parts": []interface{}{
					prompt,
				},
			},
		},
	}
	jsonData["messages"] = messages
	return header, jsonData
}

func get_config() string {
	// Read the configuration file into memory.
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	// Parse the JSON data into the Config struct.
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	cookie := fmt.Sprintf("__Secure-next-auth.session-token=%s;mp_d7d7628de9d5e6160010b84db960a7ee_mixpanel=%s", config.SessionToken, config.MixpanelToken)
	// Return the values of the config variables.
	return cookie
}

// Send a request to retreive access tokens
func get_cookie(cookie string, headers map[string]string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://chat.openai.com/api/auth/session", nil)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var respo map[string]interface{}
	err = json.Unmarshal(bodyText, &respo)
	if err != nil {
		panic(err)
	}
	authToken, ok := respo["accessToken"].(string)
	if !ok {
		panic(ok)
	}
	authToken = fmt.Sprintf("Bearer %s", authToken)
	return authToken
}

func send_prompt(cookie string, authToken string, headers map[string]string, jsonData map[string]interface{}) string {
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", "https://chat.openai.com/backend-api/conversation", bytes.NewBuffer(jsonString))
	if err != nil {
		panic(err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Authorization", authToken)
	client := &http.Client{}
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		// Handle the error
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var retdata string
	for scanner.Scan() {
		data := strings.TrimPrefix(scanner.Text(), "data:")
		if len(data) < 8 {
			continue
		}
		retdata = data
	}
	if len(retdata) < 1 {
		panic(err)
	}
	return retdata
}

func main() {
	prompt := flag.String("prompt", "", "Prompt String")
	flag.Parse()
	promp := *prompt
	cookie := get_config()
	headers, jsonData := get_req_info(promp)
	authToken := get_cookie(cookie, headers)
	body := send_prompt(cookie, authToken, headers, jsonData)
	fmt.Println(body)
}
