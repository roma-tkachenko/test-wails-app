package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"test-app/backend/api"
	"test-app/backend/configs"
	"time"
)

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

	var iterations int

	rewardLimit, err := fetchRewardCards(ctx)
	if err != nil {
		log.Printf("Помилка при отриманні карт винагород: %v", err)
		runtime.EventsEmit(ctx, "rewardError", err.Error())
		StopClaimRewards()
		return
	}

	iterations = rewardLimit

	for i := 0; i < iterations || iterations == 0; i++ { // Якщо iterations == 0, то потрібно отримати перше значення
		select {
		case <-stopMonitorCh:
			log.Println("Зупинено фоновий процес для отримання карт винагород. 2")
			StopClaimRewards()
			return
		case <-ticker.C:
			// Виконуємо запит на сервер
			rewardLimit, err := fetchRewardCards(ctx)
			if err != nil {
				log.Printf("Помилка при отриманні карт винагород: %v", err)
				runtime.EventsEmit(ctx, "rewardError", err.Error())
				StopClaimRewards() // Зупиняємо процес
				return
			}

			if iterations == 0 {
				iterations = rewardLimit
			}
		}
	}

	log.Println("Фоновий процес перевірки винагород завершено.")
	claimActive = false
}

func fetchRewardCards(ctx context.Context) (int, error) {
	urlParams := url.Values{
		"mod":       {"reward_card"},
		"action":    {"check_reward"},
		"user_hash": {configs.UserHash},
	}

	response, err := api.SendGETRequest(configs.СontrollerURL + "?" + urlParams.Encode())
	log.Printf("Response CHECK reward: %s\n", response)
	if err != nil {
		return 0, err
	}

	var data RewardCards
	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		return 0, err
	}

	if data.StopReward == "yes" {
		log.Printf("StopReward is value: %d", data.StopReward)
		StopClaimRewards()
		return 0, nil
	}

	// Перевірка вмісту `cards`
	if string(data.Cards) == `""` {
		log.Println("Cards is an empty string")
		StopClaimRewards()
		return 0, nil
	}

	var card Card
	err = json.Unmarshal(data.Cards, &card)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal cards: %w", err)
	}
	log.Printf("Parsed card: %+v", card)

	if !isAbsoluteURL(card.Image) {
		cleanedImagePath := strings.TrimPrefix(card.Image, "/")
		card.Image = configs.BaseURL + cleanedImagePath
	}

	rewardCards = append(rewardCards, card)

	if card.OwnerID != 0 {
		urlParams = url.Values{
			"mod": {"cards_ajax"},
		}

		postParams := url.Values{
			"action":   {"take_card"},
			"owner_id": {strconv.Itoa(card.OwnerID)},
		}

		log.Printf("OwnerID is not empty: %d", card.OwnerID)
		response, err = api.SendPOSTRequest(configs.СontrollerURL+"?"+urlParams.Encode(), postParams)
		log.Printf("Response GET reward: %s", response)
		if err != nil {
			return 0, err
		}
	}

	runtime.EventsEmit(ctx, "rewardCards", rewardCards)

	//fmt.Printf("reward limit: %v", data.RewardLimit)
	return data.RewardLimit, nil
}

func StopClaimRewards() {
	claimMutex.Lock()
	defer claimMutex.Unlock()

	if claimActive {
		claimActive = false
		log.Println("Фоновий процес отримання карт винагород зупинено. 1")
		close(stopMonitorCh)
		stopMonitorCh = make(chan struct{})
	}
}

func isAbsoluteURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return parsedURL.IsAbs()
}
