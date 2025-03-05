package mp3player

import (
	"container/list"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

type MP3Player struct {
	queue   *list.List
	mutex   sync.Mutex
	playing bool
	wg      sync.WaitGroup // 用於等待播放完成
}

// 建構式
func NewMP3Player() *MP3Player {
	return &MP3Player{
		queue:   list.New(),
		playing: false,
	}
}

func (p *MP3Player) Add(fileName string, funcBefore, funcAfter func()) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.queue.PushBack(struct {
		fileName   string
		funcBefore func()
		funcAfter  func()
	}{fileName, funcBefore, funcAfter})

	// 如果播放器未在播放，則啟動播放
	if !p.playing {
		p.playing = true
		p.wg.Add(1) // 增加 WaitGroup 計數
		go p.Play() // 開始播放
	}
}

func (p *MP3Player) Play() {
	for {
		p.mutex.Lock()
		if p.queue.Len() == 0 {
			p.playing = false
			p.mutex.Unlock()
			break
		}

		item := p.queue.Remove(p.queue.Front()).(struct {
			fileName   string
			funcBefore func()
			funcAfter  func()
		})
		p.mutex.Unlock()

		// 執行播放前的函式
		if item.funcBefore != nil {
			item.funcBefore()
		}

		// 播放 MP3 檔案
		p.playFile(item.fileName)

		// 執行播放後的函式
		if item.funcAfter != nil {
			item.funcAfter()
		}
	}
	p.wg.Done() // 減少 WaitGroup 計數
}

func (p *MP3Player) playFile(fileName string) {
	//讀檔
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close() // 確保檔案關閉

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		log.Println("Error creating decoder:", err)
		return
	}

	context, err := oto.NewContext(decoder.SampleRate(), 2, 2, 2048)
	if err != nil {
		log.Println("Error creating NewContext:", err)
		return
	}
	defer context.Close()

	player := context.NewPlayer()
	defer player.Close() // 確保播放器關閉

	buffer := make([]byte, 1)
	for {
		n, err := decoder.Read(buffer)
		if err != nil {
			break
		}
		if _, err := player.Write(buffer[:n]); err != nil {
			log.Println("Error writing to player:", err)
			break
		}
	}

	// 確保音頻播放完成
	time.Sleep(200 * time.Millisecond)
}

func (p *MP3Player) Close() {
	p.wg.Wait() // 等待所有播放完成
}
