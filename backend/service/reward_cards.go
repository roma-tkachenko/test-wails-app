package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"net/url"
	"strconv"
	"sync"
	"test-app/backend/api"
	"test-app/backend/configs"
	"time"
)

type RewardCards struct {
	Cards       Card   `json:"cards"`
	IfReward    string `json:"if_reward"`
	RewardLimit int    `json:"reward_limit"`
	StopReward  string `json:"stop_reward"`
}

type Card struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	NewsID    string `json:"news_id"`
	Rank      string `json:"rank"`
	Image     string `json:"image"`
	Source    string `json:"source"`
	Approve   string `json:"approve"`
	Reward    string `json:"reward"`
	Comment   string `json:"comment"`
	NewAuthor string `json:"new_author"`
	OwnerID   int    `json:"owner_id"`
}

var (
	rewardCards   []Card // Масив для збереження отриманих карт
	claimActive   = false
	claimMutex    sync.Mutex
	stopMonitorCh = make(chan struct{})
)

// ClaimRewards запускає процес отримання карт винагород
func ClaimRewards(ctx context.Context) (bool, error) {
	claimMutex.Lock()
	defer claimMutex.Unlock()

	// Якщо процес вже запущений, повертаємо успішний статус
	if claimActive {
		log.Println("Фоновий процес для отримання винагород вже запущено.")
		return false, nil
	}

	claimActive = true
	go monitorClaimRewards(ctx)
	return true, nil
}

// monitorClaimRewards перевіряє винагороди з певним інтервалом
func monitorClaimRewards(ctx context.Context) {
	log.Println("Запущено фоновий процес для отримання карт винагород...")
	ticker := time.NewTicker(configs.ClaimRewardCardInterval * time.Minute) // Інтервал перевірки (30 секунд)
	defer ticker.Stop()

	var iterations int // Кількість ітерацій, визначена сервером

	// Перший запит одразу після запуску
	rewardLimit, err := fetchRewardCards(ctx)
	if err != nil {
		log.Printf("Помилка при отриманні карт винагород: %v", err)
		claimActive = false
		runtime.EventsEmit(ctx, "rewardError", err.Error())
		close(stopMonitorCh) // Зупиняємо процес
		return
	}

	// Встановлюємо кількість ітерацій після першого запиту
	iterations = rewardLimit

	for i := 0; i < iterations || iterations == 0; i++ { // Якщо iterations == 0, то потрібно отримати перше значення
		select {
		case <-stopMonitorCh:
			log.Println("Зупинено фоновий процес для отримання карт винагород.")
			claimActive = false
			return
		case <-ticker.C:
			// Виконуємо запит на сервер
			rewardLimit, err := fetchRewardCards(ctx)
			if err != nil {
				log.Printf("Помилка при отриманні карт винагород: %v", err)
				claimActive = false
				runtime.EventsEmit(ctx, "rewardError", err.Error())
				close(stopMonitorCh) // Зупиняємо процес
				return
			}

			// Встановлюємо кількість ітерацій після першого запиту
			if iterations == 0 {
				iterations = rewardLimit
			}
		}
	}

	log.Println("Фоновий процес перевірки винагород завершено.")
	claimActive = false
}

// fetchRewardCards виконує запит до сервера та повертає reward_limit
func fetchRewardCards(ctx context.Context) (int, error) {
	// Формування URL із параметрами
	urlParams := url.Values{
		"mod":       {"reward_card"},
		"action":    {"check_reward"},
		"user_hash": {configs.UserHash},
	}

	log.Printf("get cart url : %s\n", configs.СontrollerURL+"?"+urlParams.Encode())

	response, err := api.SendGETRequest(configs.СontrollerURL + "?" + urlParams.Encode())
	log.Printf("Response CHECK reward: %s\n", response)
	if err != nil {
		return 0, err
	}

	// Парсинг JSON відповіді
	var data RewardCards
	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		return 0, err
	}

	if data.StopReward == "yes" {
		log.Printf("StopReward is value: %d", data.StopReward)
		close(stopMonitorCh) // Зупиняємо процес
	}

	// Зберігаємо отримані карти в масив
	rewardCards = append(rewardCards, data.Cards)

	if data.Cards.OwnerID != 0 {
		urlParams = url.Values{
			"mod": {"cards_ajax"},
		}

		postParams := url.Values{
			"action":   {"take_card"},
			"owner_id": {strconv.Itoa(data.Cards.OwnerID)},
		}

		log.Printf("OwnerID is not empty: %d", data.Cards.OwnerID)
		// Виконання HTTP/3 POST-запиту
		response, err = api.SendPOSTRequest(configs.СontrollerURL+"?"+urlParams.Encode(), postParams)
		log.Printf("Response GET reward: %s", response)
		if err != nil {
			return 0, err
		}
	}

	// Тригеримо івент для передачі отриманої карти на фронтенд
	runtime.EventsEmit(ctx, "rewardCards", data.Cards)

	fmt.Printf("reward limit: %v", data.RewardLimit)
	return data.RewardLimit, nil
}

// StopClaimRewards зупиняє фоновий процес отримання винагород
func StopClaimRewards() {
	claimMutex.Lock()
	defer claimMutex.Unlock()

	if claimActive {
		close(stopMonitorCh)
		stopMonitorCh = make(chan struct{}) // Оновлюємо канал для наступного запуску
		claimActive = false
		log.Println("Фоновий процес отримання карт винагород зупинено.")
	}
}
