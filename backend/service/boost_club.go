package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"test-app/backend/api"
	"test-app/backend/configs"
	"time"
)

const (
	ActionRefresh  = "refresh"
	ActionSubmit   = "submit"
	ActionRetry    = "retry"
	ActionNoAction = "none"
)

type JsonResponse struct {
	BoostNo          string `json:"boost_no"`
	BoostHTML        string `json:"boost_html"`
	BoostHTMLChanged string `json:"boost_html_changed"`
	Error            string `json:"error"`
}

var (
	cardID string
	action string
)

var appContext context.Context

func SetAppContext(ctx context.Context) {
	appContext = ctx
}

func StartProcessing(ctx context.Context) {
	SetAppContext(ctx)
	var wg sync.WaitGroup

	const interval = configs.RefreshCardInterval * time.Millisecond

	// первинний запит на сторінку внесення карт
	html, err := api.SendGETRequest(configs.BoostClubURL)
	if err != nil {
		log.Printf("Помилка первинного запиту: %v\n", err)
		return
	}

	// отримання ID карти на сторінці та події яку варто виконати
	cardID, action, err := getCardActionFromHTML(html)
	if err != nil {
		log.Printf("Помилка парсингу HTML: %v\n", err)
		return
	}

	// channel for incoming messages
	resetTimerSignal := make(chan bool, 1)
	defer close(resetTimerSignal)

	// канал для сигналу зупинки циклів
	stopSignal := make(chan bool, 1)

	// Відкриття каналу для сигналів внесення карти
	submitSignal := make(chan bool, 1)
	defer close(submitSignal)

	// Відправлення сигналу на внесення карти
	if action == ActionSubmit {
		submitSignal <- true
	}

	// Цикл для запитів оновлення карти
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Refresh loop stopped")
				return
			case <-stopSignal: // Отримання сигналу для зупинки
				log.Println("Received stop signal for refresh loop")
				return
			case <-time.After(interval):
				// Робимо запит оновлення карти
				cardID, action, err = makeServerRequest(cardID, ActionRefresh)
				if err != nil {
					log.Printf("Refresh error: %v", err)
					continue
				}

				log.Printf("Refresh response action: %s", action)
				if action == ActionSubmit {
					// Відправлення сигналу на внесення карти
					submitSignal <- true
				}
			case <-resetTimerSignal:
				log.Println("Handle reset timer income message and move to the next iteration")
			}
		}
	}()

	// Цикл для внесення карти
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Card submission loop stopped")
				return
			case <-stopSignal: // Отримання сигналу для зупинки
				log.Println("Received stop signal for refresh loop")
				return
			case <-submitSignal:
				log.Println("Submit loop:")
				for {
					// Робимо запит внесення карти
					cardID, action, err = makeServerRequest(cardID, ActionSubmit)
					if err != nil {
						// Вихід з циклу при виникненні помилки
						log.Printf("Submit error: %v", err)
						break
					}

					// Після успішного виконання запиту на внесення карти скидаємо таймер оновлення карти
					resetTimerSignal <- true

					// Повторюємо запит на внесення карти (умови для повторного внесення карти залежать від відповіді сервера)
					log.Printf("Submit response action: %s", action)
					if action == ActionRetry {
						action = ActionSubmit
						// Додаємо очікування перед повторним запитом
						time.Sleep(configs.RetryBoostCardInterval * time.Millisecond)
						log.Println("Retrying card submission by retry action...")
						continue
					} else if action == ActionSubmit {
						// Додаємо очікування перед новим запитом
						time.Sleep(configs.RetryBoostCardInterval * time.Millisecond)
						log.Println("Submission triggered another submission")
						continue
					}
					break
				}
			}
		}
	}()

	// Очікуємо завершення всіх горутин
	wg.Wait()
}

func getCardActionFromHTML(body string) (cardID string, action string, err error) {
	// Регулярний вираз для кнопки з класом "club__boost-btn"
	boostBtnRegex := regexp.MustCompile(`<button[^>]*class="[^"]*club__boost-btn[^"]*"[^>]*data-card-id="([^"]*)"`)
	boostMatches := boostBtnRegex.FindStringSubmatch(body)

	// Регулярний вираз для кнопки з класом "club__boost__refresh-btn"
	refreshBtnRegex := regexp.MustCompile(`<button[^>]*class="[^"]*club__boost__refresh-btn[^"]*"[^>]*data-card-id="([^"]*)"`)
	refreshMatches := refreshBtnRegex.FindStringSubmatch(body)

	// Регулярний вираз для пошуку тега <img> всередині <div class="club-boost__image">
	imageRegex := regexp.MustCompile(`<div[^>]*class="[^"]*club-boost__image[^"]*"[^>]*>\s*<img[^>]*src="([^"]*)"`)
	imageMatches := imageRegex.FindAllStringSubmatch(body, -1)

	// Якщо знайдено результати
	if len(imageMatches) > 0 {
		// Перевірка та додавання домену до Image
		if !isAbsoluteURL(imageMatches[0][1]) {
			// Видаляємо початковий слеш, якщо він є
			cleanedImagePath := strings.TrimPrefix(imageMatches[0][1], "/")
			imgUrl := configs.BaseURL + cleanedImagePath

			runtime.EventsEmit(appContext, "boostImage", imgUrl)
		}
	} else {
		fmt.Println("Картинки не знайдено.")
	}

	if len(boostMatches) > 1 {
		// Якщо знайдено кнопку "club__boost-btn"
		cardID = boostMatches[1]
		action = "submit"
		return cardID, action, nil
	} else if len(refreshMatches) > 1 {
		// Якщо знайдено кнопку "club__boost__refresh-btn"
		cardID = refreshMatches[1]
		action = "refresh"
		return cardID, action, nil
	}

	// Якщо кнопки не знайдено
	return "", "", errors.New("не знайдено відповідних кнопок")
}

func makeServerRequest(cardID string, action string) (string, string, error) {
	var jsonResponseData JsonResponse
	var err error
	var jsonStr string

	// Виконуємо відповідний запит залежно від дії
	switch action {
	case ActionRefresh:
		jsonStr, err = performRefreshRequest(cardID)
		jsonResponseData, err = parseJsonResponse(jsonStr)
	case ActionSubmit:
		jsonStr, err = performBoostRequest(cardID)
		jsonResponseData, err = parseJsonResponse(jsonStr)
	default:
		log.Printf("Unknown action: %s", action)
		return "", "", errors.New("unknown action")
	}

	// Перевіряємо помилку запиту
	if err != nil {
		log.Printf("Server request failed: %v", err)
		return "", "", err
	}

	// Обробляємо json відповідь після виконання запитів з картами
	cardID, action, err = processJsonResponseData(jsonResponseData, cardID, action)
	if err != nil {
		log.Printf("Failed to process JSON response: %v", err)
		return "", "", err
	}

	return cardID, action, nil
}

func processJsonResponseData(jsonResponseData JsonResponse, cardID string, action string) (string, string, error) {
	// Обробляємо поле `error`, якщо воно не порожнє
	if jsonResponseData.Error != "" {
		log.Printf("Action: %s -> Error from server: %s \n", action, jsonResponseData.Error)

		// Використання регулярного виразу для перевірки тексту
		reRetry := regexp.MustCompile(`^Следующую карту можно сдать клубу через -?\d+ секунд$`)

		switch {
		case reRetry.MatchString(jsonResponseData.Error):
			return cardID, ActionRetry, nil
		case jsonResponseData.Error == "Ваша карта заблокирована, для пожертвования клубу разблокируйте её":
			return cardID, ActionNoAction, nil
		case jsonResponseData.Error == "Достигнут дневной лимит пожертвований в клуб, подождите до завтра":
			return cardID, ActionNoAction, nil
		default:
			log.Printf("Unknown Error: %s\n", jsonResponseData.Error)
			return cardID, ActionRefresh, nil
		}
	}

	// Обробляємо поле `boost_no`, якщо воно не порожнє
	if jsonResponseData.BoostNo != "" {
		log.Printf("Action: %s -> BoostNo from server: %s", action, jsonResponseData.BoostNo)
		switch jsonResponseData.BoostNo {
		case "Нужная клубу карта не менялась":
			return cardID, ActionRefresh, nil
		default:
			log.Printf("Unknown BoostNo: %s\n", jsonResponseData.BoostNo)
			return cardID, ActionRefresh, nil
		}
	}

	// Вибір HTML-вмісту для парсингу
	var html string
	if jsonResponseData.BoostHTMLChanged != "" {
		html = jsonResponseData.BoostHTMLChanged
		log.Println("Using `boost_html_changed` from response.")
	} else if jsonResponseData.BoostHTML != "" {
		html = jsonResponseData.BoostHTML
		log.Println("Using `boost_html` from response.")
	} else {
		log.Println("No valid HTML content in the response.")
		return cardID, ActionRefresh, errors.New("no valid HTML content")
	}

	// Парсинг HTML для отримання ID карти та екшену
	parsedCardID, parsedAction, err := getCardActionFromHTML(html)
	if err != nil {
		log.Printf("Помилка парсингу HTML: %v\n", err)
		return cardID, ActionRefresh, err
	}

	log.Printf("Відповідь парсингу HTML: %s -> %s\n", parsedAction, parsedCardID)
	return parsedCardID, parsedAction, nil
}

func performBoostRequest(cardID string) (string, error) {
	// Параметри POST-запиту
	postParams := url.Values{
		"action":    {"boost"},
		"card_id":   {cardID},
		"user_hash": {configs.UserHash},
	}

	// Виконуємо HTTP/3 POST-запит
	jsonStr, err := api.SendPOSTRequest(configs.BoostClubURL, postParams)
	if err != nil {
		log.Printf("HTTP/3 POST request failed: %v", err)
		return "", err
	}

	return jsonStr, nil
}

func performRefreshRequest(cardID string) (string, error) {
	urlParams := url.Values{
		"mod": {"clubs_ajax"},
	}

	finalURL := configs.СontrollerURL + "?" + urlParams.Encode()

	// Параметри POST-запиту
	postParams := url.Values{
		"action":    {"boost_refresh"},
		"card_id":   {cardID},
		"user_hash": {configs.UserHash},
	}

	// Perform an HTTP/3 GET request
	jsonStr, err := api.SendPOSTRequest(finalURL, postParams)
	if err != nil {
		log.Printf("HTTP/3 POST request failed: %v", err)
		return "", err
	}

	return jsonStr, nil
}

func parseJsonResponse(jsonStr string) (JsonResponse, error) {
	// Перевіряємо, чи наданий параметр не порожній
	if jsonStr == "" {
		log.Println("Empty response from server")
		return JsonResponse{}, errors.New("empty response")
	}

	// Парсинг JSON відповіді
	var jsonResponseData JsonResponse
	err := json.Unmarshal([]byte(jsonStr), &jsonResponseData)
	if err != nil {
		log.Printf("Error parsing JSON response: %v", err)
		return JsonResponse{}, err
	}

	// Повертаємо заповнений JsonResponse і nil як помилку
	return jsonResponseData, nil
}
