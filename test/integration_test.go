package test

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"spectre-go/controllers"
	"spectre-go/repo"
	"spectre-go/utility/rest"
	"strings"
	"testing"
)

func runTestServer() *httptest.Server {
	//host := os.Getenv("REDIS_HOST")
	//passwd := os.Getenv("REDIS_PASSWD")

	host := "localhost:6379"
	passwd := "pass"
	log.Printf("Host %s , Passwd %s", host, passwd)
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: passwd,
		DB:       0,
	})

	//if err := client.Ping(context.Background()); err != nil {
	//	log.Fatal(err)
	//}

	log.Printf("Redis client initialized")

	siteResultRepo := repo.NewSiteResultRepo(client)

	h := controllers.NewBaseHandler(siteResultRepo)

	return httptest.NewServer(controllers.HandleRequests(h))
}

func Test_genPassword(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	t.Run("It should return 200 when health ok", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/health", ts.URL))

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("It should return correct password result", func(t *testing.T) {
		url := fmt.Sprintf("%s/api/getPassword", ts.URL)
		contentType := "application/json"
		body := generateBody()
		resp, err := http.Post(url, contentType, strings.NewReader(body))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		respBody := string(bodyBytes)
		respBody = strings.TrimSuffix(fmt.Sprintf("%s", respBody), "\n")
		log.Printf("Isi body " + respBody)
		assert.Equal(t, respBody, `{"result":"Mat4;Noq"}`)
	})
}

func Benchmark_intTest_genPassword_withCache(b *testing.B) {
	ts := runTestServer()
	defer ts.Close()
	url := fmt.Sprintf("%s/api/getPassword", ts.URL)
	contentType := "application/json"

	headers := map[string]interface{}{
		"Content-Type":  contentType,
		"Utilize-Cache": "true",
	}

	//log.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		body := generateFakeBody()
		http.Post(url, contentType, strings.NewReader(body))
		(rest.RestClient{}).Post(
			rest.RestRequest{
				Path:    url,
				Headers: headers,
				Body:    []byte(body),
				Method:  "POST",
			},
		)
	}
}

func Benchmark_intTest_genPassword_withOutCache(b *testing.B) {
	ts := runTestServer()
	defer ts.Close()
	url := fmt.Sprintf("%s/api/getPassword", ts.URL)
	contentType := "application/json"

	headers := map[string]interface{}{
		"Content-Type":  contentType,
		"Utilize-Cache": "false",
	}

	//log.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		body := generateFakeBody()
		http.Post(url, contentType, strings.NewReader(body))
		(rest.RestClient{}).Post(
			rest.RestRequest{
				Path:    url,
				Headers: headers,
				Body:    []byte(body),
				Method:  "POST",
			},
		)
	}
}

func generateBody() string {
	return fmt.Sprintf(`{
		"username":"a",
		"password":"acde",
		"site":"twitter.com",
		"keyCounter":1,
		"keyPurpose":"com.lyndir.masterpassword",
		"keyType":"med" }`)
}

func generateFakeBody() string {
	id := uuid.New()
	idstr := id.String()
	return fmt.Sprintf(`{
		"username":"%s",
		"password":"%s",
		"site":"%s",
		"keyCounter":1,
		"keyPurpose":"com.lyndir.masterpassword",
		"keyType":"med" }`, idstr, idstr, idstr)
}

//compare performance pake uuid aja.
