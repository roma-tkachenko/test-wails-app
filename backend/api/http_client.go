package api

import (
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"test-app/backend/configs"
)

// Глобальна змінна для клієнта
var (
	client     *http.Client
	clientOnce sync.Once // Використовується для ініціалізації один раз
)

// initHttpClient ініціалізує HTTP/3 клієнт
func initHttpClient() {
	clientOnce.Do(func() {
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalf("Failed to create cookie jar: %v", err)
		}

		tr := &http3.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			QUICConfig: &quic.Config{
				//HandshakeIdleTimeout: 10 * time.Second, // Тайм-аут під час рукопотискання
			},
		}

		client = &http.Client{
			Transport: tr,
			Jar:       jar,
			//Timeout:   30 * time.Second, // Загальний тайм-аут для всіх запитів
		}
	})
}

// SendGETRequest виконує GET-запит
func SendGETRequest(requestURL string) (string, error) {
	initHttpClient() // Ініціалізація клієнта (якщо ще не ініціалізований)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return "", err
	}

	// Використання функції для встановлення заголовків
	setRequestHeaders(req)

	resp, err := client.Do(req)

	//resp, err := client.Get(requestURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Перевірка статусу відповіді
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		log.Println(err)
		return "", err
	}

	// Читання тіла відповіді
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//Extract dle_login_hash
	extractDleLoginHash(string(body))

	return string(body), nil
}

// SendPOSTRequest виконує POST-запит
func SendPOSTRequest(requestURL string, data url.Values) (string, error) {
	initHttpClient() // Ініціалізація клієнта (якщо ще не ініціалізований)

	req, err := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Використання функції для встановлення заголовків
	setRequestHeaders(req)

	resp, err := client.Do(req)

	//resp, err := client.PostForm(requestURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Перевірка статусу відповіді
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return "", nil
	}

	// Читання тіла відповіді
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//Extract dle_login_hash
	extractDleLoginHash(string(body))

	return string(body), nil
}

// setRequestHeaders встановлює стандартні заголовки для запитів
func setRequestHeaders(req *http.Request) {
	req.Header.Set("Alt-Svc", `h3=":443"; ma=86400`)
	req.Header.Set("documentLifecycle", "active")
	req.Header.Set("frameType", "outermost_frame")
	req.Header.Set("initiator", "https://animestars.org")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://animestars.org/")
	req.Header.Set("sec-ch-ua", `"Not A(Brand";v="8", "Chromium";v="132", "Brave";v="132"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
}

func extractDleLoginHash(body string) {
	hashRegex := regexp.MustCompile(`var dle_login_hash\s*=\s*['"]([a-zA-Z0-9]+)['"]`)
	matches := hashRegex.FindSubmatch([]byte(body))
	if len(matches) >= 2 {
		newHash := string(matches[1]) // Нове значення dle_login_hash
		if newHash != configs.UserHash {
			log.Printf("dle_login_hash changed: %s -> %s\n", configs.UserHash, newHash)
			configs.UserHash = newHash // Оновлюємо значення
		}
	}
}
