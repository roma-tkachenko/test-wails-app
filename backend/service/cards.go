package service

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"strconv"
	"strings"
	"test-app/backend/api"
	"test-app/backend/configs"
)

var cards []Card
var stopSync chan bool

func SyncCards(ctx context.Context) (bool, error) {

	// page := 1

	// url := fmt.Sprintf(configs.BaseCardsURL, configs.UserName, page)

	// log.Printf("Sync url: %v", url)

	// response, err := api.SendGETRequest(url)
	// if err != nil {
	// 	log.Printf("Sync error: %v", err)
	// 	return false, err
	// }

	// cards, err := ParseCards(response)

	// runtime.EventsEmit(ctx, "allCards", cards)

	// канал для сигналу зупинки циклів
	stopSync := make(chan bool, 1)

	page := 1 // Початкова сторінка
	for {
		// Формуємо URL
		url := fmt.Sprintf(configs.BaseCardsURL, configs.UserName, page)
		log.Printf("Sync url: %v", url)

		// Відправляємо GET-запит
		response, err := api.SendGETRequest(url)
		if err != nil {
			log.Printf("Sync error: %v", err)
			return false, err
		}

		// Парсимо карти
		cards, err := ParseCards(response)
		if err != nil {
			log.Printf("Parsing error: %v", err)
			return false, err
		}

		// Якщо картки порожні — виходимо з циклу
		if len(cards) == 0 {
			log.Println("No more cards to sync. Exiting...")
			break
		}

		// Відправляємо подію з картками
		runtime.EventsEmit(ctx, "allCards", cards)

		// Умова для виходу з циклу (наприклад, перевірка контексту)
		select {
		case <-stopSync:
			log.Println("Sync interrupted by context. Exiting...")
			return false, err
		default:
		}

		if page == 52 {
			stopSync <- true
		}

		// Збільшуємо номер сторінки
		page++
	}

	return true, nil
}

func ParseCards(html string) ([]Card, error) {
	// Парсимо HTML-документ
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// Знаходимо всі елементи з класом "anime-cards__item-wrapper"
	doc.Find(".anime-cards__item-wrapper").Each(func(i int, wrapper *goquery.Selection) {
		// Знаходимо елемент з класом "anime-cards__item" всередині wrapper
		item := wrapper.Find(".anime-cards__item")
		if item.Length() > 0 {
			// Створюємо карту зі значень data-атрибутів
			card := Card{
				Name:       item.AttrOr("data-name", ""),
				ID:         item.AttrOr("data-id", ""),
				Rank:       item.AttrOr("data-rank", ""),
				AnimeName:  item.AttrOr("data-anime-name", ""),
				AnimeLink:  item.AttrOr("data-anime-link", ""),
				Author:     item.AttrOr("data-author", ""),
				Image:      item.AttrOr("data-image", ""),
				CanTrade:   item.AttrOr("data-can-trade", ""),
				IsFavorite: item.AttrOr("data-favourite", ""),
			}

			// Конвертація OwnerID з рядка в int
			ownerIDStr := item.AttrOr("data-owner-id", "")
			ownerID, err := strconv.Atoi(ownerIDStr)
			if err != nil {
				ownerID = 0 // Якщо виникла помилка, використовуємо значення за замовчуванням
			}

			card.OwnerID = ownerID

			// Перевірка та додавання домену до Image
			if !isAbsoluteURL(card.Image) {
				// Видаляємо початковий слеш, якщо він є
				cleanedImagePath := strings.TrimPrefix(card.Image, "/")
				card.Image = configs.BaseURL + cleanedImagePath
			}

			// Додаємо карту в список
			cards = append(cards, card)
		} else {
			stopSync <- true
		}
	})

	return cards, nil
}

func GetAllcards(ctx context.Context) (bool, error) {

	return true, nil
}
