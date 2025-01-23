package service

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"net/url"
	"regexp"
	"test-app/backend/api"
	"test-app/backend/configs"
	"time"
)

var (
	authActive = false
)

// Authorize виконує авторизацію користувача
func Authorize(ctx context.Context, username string, password string) (bool, error) {
	data := url.Values{
		"login_name":     {username},
		"login_password": {password},
		"login":          {"submit"},
	}

	// Виконання POST-запиту для авторизації
	response, err := api.SendPOSTRequest(configs.BaseURL, data)
	if err != nil {
		log.Printf("Авторизація не вдалася: %v", err)
		return false, err
	}

	_ = checkAuthStatus(ctx, response)

	if authActive {
		configs.UserName = string(username)
		// Запускаємо фоновий процес для перевірки статусу авторизації
		go monitorAuthStatus(ctx)
	}

	return authActive, nil
}

// monitorAuthStatus перевіряє статус авторизації з певним інтервалом
//func monitorAuthStatus(ctx context.Context) {
//	log.Println("Запущено фоновий процес для перевірки авторизації...")
//	ticker := time.NewTicker(30 * time.Second) // Інтервал перевірки (30 секунд)
//	defer ticker.Stop()
//
//	for {
//		select {
//		case <-ticker.C:
//			// Перевіряємо авторизацію, якщо вона активна
//			if authActive {
//				err := checkAuthStatus(ctx, "")
//				if err != nil {
//					log.Println("Авторизація завершена або помилка перевірки:", err)
//					authActive = false
//		            StopClaimRewards()
//					runtime.EventsEmit(ctx, "authStatus", authActive)
//					return // Завершуємо фоновий процес
//				}
//			} else {
//				log.Println("Авторизація більше не активна. Завершуємо перевірку.")
//				return
//			}
//		}
//	}
//}

func monitorAuthStatus(ctx context.Context) {
	log.Println("Запущено фоновий процес для перевірки авторизації...")
	ticker := time.NewTicker(60 * time.Second) // Інтервал перевірки (30 секунд)
	defer ticker.Stop()

	maxRetries := 2 // Максимальна кількість повторних спроб

	for {
		select {
		case <-ticker.C:
			// Перевіряємо авторизацію, якщо вона активна
			if authActive {
				var err error
				retries := 0

				for retries <= maxRetries {
					err = checkAuthStatus(ctx, "")
					if err != nil {
						retries++
						log.Printf("Помилка перевірки авторизації. Спроба #%d: %v", retries, err)
						// Якщо це була остання спроба, виходимо із циклу
						if retries > maxRetries {
							break
						}
					} else {
						// Якщо запит успішний, виходимо з циклу повторів
						break
					}
				}

				// Якщо після повторів авторизацію не вдалось підтвердити
				if err != nil {
					log.Println("Авторизація завершена або помилка перевірки:", err)
					authActive = false
					StopClaimRewards()
					runtime.EventsEmit(ctx, "authStatus", authActive)
					return // Завершуємо фоновий процес
				}
			} else {
				log.Println("Авторизація більше не активна. Завершуємо перевірку.")
				return
			}
		}
	}
}

// checkAuthStatus виконує запит для перевірки активності авторизації
func checkAuthStatus(ctx context.Context, body string) error {
	// Якщо body — це порожній рядок, виконуємо запит до сервера
	if body == "" {
		response, err := api.SendGETRequest(configs.BaseURL)
		if err != nil {
			return err
		}
		body = response // Призначаємо отриманий результат змінній body
	}

	//log.Printf("Responce: %s", body)

	// Регулярний вираз для пошуку рядка OneSignal.sendTag("userId", "...");
	hashRegex := regexp.MustCompile(`OneSignal\.sendTag\("userId",\s*"(\d+)"\);`)
	matches := hashRegex.FindSubmatch([]byte(body))

	if len(matches) > 1 {
		newUserId := string(matches[1])
		if newUserId != configs.UserId {
			log.Printf("UserID changed: %s -> %s\n", configs.UserId, newUserId)
			configs.UserId = newUserId // Оновлюємо значення
		}

		authActive = true
	} else {
		authActive = false
		StopClaimRewards()
		log.Println("userId not found")
	}

	log.Printf("Check authorization status: authActive = %v", authActive)

	runtime.EventsEmit(ctx, "authStatus", authActive)

	body = ""
	return nil
}
