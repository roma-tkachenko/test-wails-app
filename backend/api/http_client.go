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
			TLSClientConfig: &tls.Config{},
			QUICConfig:      &quic.Config{},
		}

		client = &http.Client{
			Transport: tr,
			Jar:       jar,
		}
	})
}

// SendGETRequest виконує GET-запит
func SendGETRequest(requestURL string) (string, error) {
	initHttpClient() // Ініціалізація клієнта (якщо ще не ініціалізований)

	resp, err := client.Get(requestURL)
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

	resp, err := client.PostForm(requestURL, data)
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

func extractDleLoginHash(body string) {
	hashRegex := regexp.MustCompile(`var dle_login_hash\s*=\s*['"]([a-zA-Z0-9]+)['"]`)
	matches := hashRegex.FindSubmatch([]byte(body))
	if len(matches) >= 2 {
		configs.UserHash = string(matches[1])
		log.Printf("Extracted dle_login_hash: %s\n", configs.UserHash)
	}
}
