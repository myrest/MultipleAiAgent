# Langchain 使用 Golang 語言

## 系統設定檔 `systemconfig.json`

以下是系統設定檔中的各項配置說明：

- **OPENAI_API_KEY**: 
  - 當你使用 OpenAI 的服務時，需要提供 API key。詳情請參考 [OpenAI 官網](https://openai.com)。
  
- **OPENAI_BASE_URL**: 
  - 根據你所使用的 AI 來源，資料會傳送到此位置。

- **DefaultMaxResponseLength**: 
  - 預設限制 AI 回答的內容長度。

- **BotSettingFileName**: 
  - AI 機器人設定檔。預設包含兩個範例檔案：
    - `botsetting.json`：辯論比賽設定
    - `storybotsetting.json`：說故事設定

- **AiSource**: 
  - 表示你的 AI 來源，可能的值為：
    - `OpenAI`
    - `Ollama`
    - `OpenAIStory`
    - `OllamaStory`
  - 如果使用 OpenAI，則需要填寫 `OPENAI_API_KEY`；如果使用 Ollama，則不需要提供。

- **ModelName**: 
  - 預計使用的模型名稱，例如：`llama3.1`, `llama3-8b`, `gpt-4o-mini`。

- **Sequential**: 
  - 表示 AI 參與的情況，可以是：
    - `Random`（隨機）
    - `Sequential`（依序）

- **EnableVoice**: 
  - 設定是否啟用產生 MP3 檔案：
    - `true`：啟用產生 MP3 檔
    - `false`：不啟用（使用 Ollama 時需啟用UseEdgeTTS）

- **UseEdgeTTS**: 
  - 設定是否使用Microsoft Server Speech Text to Speech Voice產生 MP3 ：
    - `true`：使用Edge TTS
    - `false`：使用LLM baseURL + "audio/speech"的方式產生

## AI 設定檔

以下是 AI 設定檔中的各項配置說明：

- **Name**: 
  - AIBot 的名字。

- **SystemPrompt**: 
  - 用於 System prompt 的提示語。

- **BotMaxResponseLength**: 
  - 限制 AI 回答的內容長度。
  - 若未設定或其實為0，將參考系統設定中的DefaultMaxResponseLength值。

- **Voice**: 
  - Speech 聲音名稱：
    - `nova`
    - `shimmer`
    - `echo`
    - `fable`
    - `alloy`
    - `onyx`
  - Edge TTS 支援中文的聲音名稱：
    - `zh-CN-XiaoxiaoNeural`
    - `zh-CN-XiaoyiNeural`
    - `zh-CN-YunjianNeural`
    - `zh-CN-YunxiNeural`
    - `zh-CN-YunxiaNeural`
    - `zh-CN-YunyangNeural`
    - `zh-TW-HsiaoChenNeural`
    - `zh-TW-YunJheNeural`
    - `zh-TW-HsiaoYuNeural `

## 如何使用
1. **初始化系統**
   在終端中執行以下命令以初始化系統，`RestChatBot` 是本專案名稱：
   ```bash
   go mod init RestChatBot
2. **下載必要套件**
   執行以下命令來下載專案所需的所有依賴套件：
   ```bash
   go mod tidy
3. **建置及編譯執行檔**
   使用以下命令來編譯並建置執行檔：
   ```bash
   go build ./src/main.go
4. **設定執行期環境變數**
   使用以下命令來編譯並建置執行檔：
   ```bash
   export CONFIGFOLDER=/你的/設定檔/路徑
4. **執行**
   使用以下命令來啟動執行檔：
   ```bash
   ./main
