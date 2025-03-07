package voicebuilder

import (
	"RestChatBot/src/config"
	"RestChatBot/src/model"
	"RestChatBot/src/mp3player"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/wujunwei928/edge-tts-go/edge_tts"
)

func ConvertToJsonByte(obj model.VoicePayload) []byte {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil
	}

	return jsonData
}

func ConvertToMp3(Message, VoiceName, FileName string) error {
	config := config.GetConfig()
	if config.EnableVoice {
		if config.UseEdgeTTS {
			return ConvertToMp3_EdgeTTS(Message, VoiceName, FileName)
		}
		return ConvertToMp3_LLM(Message, VoiceName, FileName)
	}
	return nil
}

func ConvertToMp3_EdgeTTS(Message, VoiceName, FileName string) error {
	//kk, err := edge_tts.ListVoices("")
	//for _, voice := range kk {
	//	if voice.Locale == "zh-CN" || voice.Locale == "zh-TW" {
	//		fmt.Println(voice.ShortName)
	//	}
	//}

	//Voice name list
	//zh-CN-XiaoxiaoNeural
	//zh-CN-XiaoyiNeural
	//zh-CN-YunjianNeural
	//zh-CN-YunxiNeural
	//zh-CN-YunxiaNeural
	//zh-CN-YunyangNeural
	//zh-TW-HsiaoChenNeural
	//zh-TW-YunJheNeural
	//zh-TW-HsiaoYuNeural

	connOptions := []edge_tts.CommunicateOption{
		edge_tts.SetVoice(VoiceName),
		edge_tts.SetRate("+20%"),
		edge_tts.SetVolume("+0%"),
		edge_tts.SetPitch("+0Hz"),
		edge_tts.SetReceiveTimeout(20), //20秒 timeout
	}

	conn, err := edge_tts.NewCommunicate(
		Message,
		connOptions...,
	)

	if err != nil {
		return err
	}

	audioData, err := conn.Stream()
	if err != nil {
		return err
	}

	err = os.WriteFile(FileName, audioData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ConvertToMp3_LLM(Message, VoiceName, FileName string) error {
	config := config.GetConfig()
	apiKey := config.OPENAI_API_KEY
	baseURL := config.OPENAI_BASE_URL

	voicepayload := model.VoicePayload{
		Model:           "tts-1-hd",
		Voice:           VoiceName,
		Response_format: "mp3",
		Speed:           1.0,
		Input:           Message,
	}

	jsonPayload := ConvertToJsonByte(voicepayload)

	url := baseURL + "audio/speech"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBuffer(jsonPayload))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		audioData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response: %v", err)
		}

		err = os.WriteFile(FileName, audioData, os.ModePerm)
		if err != nil {
			log.Fatalf("Error writing audio file: %v", err)
		}

		fmt.Println("Audio saved as", FileName)
		return nil
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
			return err
		}
		msg := fmt.Sprintf("Failed to convert text to speech: %v - %v", resp.StatusCode, string(body))
		return errors.New(msg)
	}
}

var chanMp3 = make(chan bool)
var isPlaying = false

func XXPlayMp3File(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening audio file: %v", err)
	}
	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		log.Fatalf("Error creating MP3 decoder: %v", err)
	}

	context, err := oto.NewContext(decoder.SampleRate(), 2, 2, 4096)
	if err != nil {
		log.Fatalf("Error creating audio context: %v", err)
	}
	defer context.Close()

	player := context.NewPlayer()
	defer player.Close()

	fmt.Println("Playing audio:", filename)
	buffer := make([]byte, 1)
	for {
		_, err := decoder.Read(buffer)
		if err != nil {
			break
		}
		player.Write(buffer)
	}

	time.Sleep(1 * time.Second) // 等待一秒確保播放完成
	fmt.Println("Audio playback finished:", filename)
}

func XXPlayMp3FileInline(filename string) {
	if filename == "" {
		time.Sleep(time.Second) //無條件等一秒，在同時排隊，要最後才對
		<-chanMp3
		return
	}
	needKeepinLine := true
	if !isPlaying {
		isPlaying = true
		needKeepinLine = false
	}

	go func(isNeedWait bool) {
		if isNeedWait {
			//排隊等訊號
			fmt.Println("排隊：", filename)
			<-chanMp3
			XXPlayMp3File(filename)
		} else {
			XXPlayMp3File(filename)
		}
		chanMp3 <- true
	}(needKeepinLine)
}

func PlayVoiceSample() {
	player := mp3player.NewMP3Player()
	defer player.Close()
	voices := []string{"nova", "shimmer", "echo", "fable", "alloy", "onyx"}
	for _, v := range voices {
		filename := fmt.Sprintf("%s_Sample.mp3", v)
		player.Add(filename, nil, nil)
	}
}

func BuildSampleVoice() {
	player := mp3player.NewMP3Player()
	defer player.Close()

	voices := []string{"nova", "shimmer", "echo", "fable", "alloy", "onyx"}
	message := "This is a test voice, glad to serve you. 這是測試語音，很高興能為您服務，我是。"
	for _, v := range voices {
		filename := fmt.Sprintf("%s_Sample.mp3", v)
		newMsg := fmt.Sprint(message, v)
		err := ConvertToMp3(newMsg, v, filename)
		if err != nil {
			return
		}
		player.Add(filename,
			func() { log.Println("開始播放...", filename) },
			func() { log.Println("結束播放...", filename) },
		)
	}
}
