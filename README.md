# openai-cli
A tool written in golang that allows access to chat.openai.com from your command line interface. 

## Installation
Compile from source
```
git clone https://github.com/Lytes/openai-cli
cd openai-cli
go build -o openai-cli main.go
```
OR use the openai-cli executable (ELF file for Linux Users)

### Setup
A `config.json` file has to be placed in the same directory as the executable. The file contains both `__Secure-next-auth.session-token` (Session Token) and `mixpaneltoken`, both are cookies that should be copied from user browser after logging in. 

Here's an example of the json file
```
{"sessiontoken":"eyJhbGciOiJkaXIi...",
"mixpaneltoken":"%7B%22distinct_id%22%3A%20%22user-..."
}
```
If unfamiliar with getting cookies from your browser, I suggest asking the AI on chat.openai.com

### Usage
Example
```
./openai-cli -prompt "Is it pronounced tomato or tomato"
```
Output
```
{"message": {"id": "f42f9493-9f6f-4650-ae2b-066e3c60202f", "role": "assistant", "user": null, "create_time": null, "update_time": null, "content": {"content_type": "text", "parts": ["Tomato is the correct spelling for the fruit (yes, tomato is a fruit, not a vegetable) that is commonly eaten raw or used as an ingredient in a variety of dishes. The word comes from the Nahuatl word \"tomatl,\" which means \"the swelling fruit.\" The proper way to write it is with a lowercase \"t\" at the beginning, followed by an uppercase \"T\" in the middle and lowercase letters for the rest of the word.\n\n"]}, "end_turn": null, "weight": 1.0, "metadata": {}, "recipient": "all"}, "conversation_id": "d340cd50-4a9d-4915-ab90-2f38d38b73ed", "error": null}
```
