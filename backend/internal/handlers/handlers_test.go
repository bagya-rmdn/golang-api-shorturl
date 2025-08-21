package handlers
import (
"encoding/json"
"net/http"
"net/http/httptest"
"strings"
"testing"
"github.com/gofiber/fiber/v2"
"gorm.io/driver/sqlite"
"gorm.io/gorm"
"urlshortener/internal/models"
"urlshortener/internal/routes"
)
type cfg struct{ Port, DatabaseURL, AppBaseURL string }
func setupTestApp(t *testing.T) (*fiber.App,
*gorm.DB) {
db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"),
&gorm.Config{})
if err != nil { t.Fatal(err) }
if err := db.AutoMigrate(&models.URLMapping{}); err != nil {
t.Fatal(err) }
app := fiber.New()
routes.Register(app, db, &cfg{AppBaseURL: "http://localhost:8080"})
return app, db
}
func TestShorten_Positive(t *testing.T) {
app, _
:= setupTestApp(t)
req := httptest.NewRequest(http.MethodPost, "/shorten",
strings.NewReader(`{"url":"https://example.com/page"}`))
req.Header.Set("Content-Type", "application/json")
resp, _
:= app.Test(req)
if resp.StatusCode != http.StatusCreated { t.Fatalf("want 201, got %d",
resp.StatusCode) }
var body map[string]string
json.NewDecoder(resp.Body).Decode(&body)
if body["token"] == "" { t.Fatal("token empty") }
}
func TestShorten_MissingURL(t *testing.T) {
app, _
:= setupTestApp(t)
req := httptest.NewRequest(http.MethodPost, "/shorten",
strings.NewReader(`{}`))
req.Header.Set("Content-Type", "application/json")
resp, _
:= app.Test(req)
if resp.StatusCode != http.StatusBadRequest { t.Fatalf("want 400, got
%d", resp.StatusCode) }
}
func TestShorten_InvalidURL(t *testing.T) {
app, _
:= setupTestApp(t)
req := httptest.NewRequest(http.MethodPost, "/shorten",
strings.NewReader(`{"url":"::not-a-url"}`))
req.Header.Set("Content-Type", "application/json")
resp, _
:= app.Test(req)
if resp.StatusCode != http.StatusUnprocessableEntity { t.Fatalf("want
422, got %d", resp.StatusCode) }
}
func TestIdempotent_SameToken(t *testing.T) {
app, _
:= setupTestApp(t)
b := strings.NewReader(`{"url":"https://EXAMPLE.com/Page"}`)
req1 := httptest.NewRequest(http.MethodPost, "/shorten", b)
req1.Header.Set("Content-Type", "application/json")
resp1, _
:= app.Test(req1)
var a map[string]string
json.NewDecoder(resp1.Body).Decode(&a)
req2 := httptest.NewRequest(http.MethodPost, "/shorten",
strings.NewReader(`{"url":"https://example.com/Page"}`))
req2.Header.Set("Content-Type", "application/json")
resp2, _
:= app.Test(req2)
var bdy map[string]string
json.NewDecoder(resp2.Body).Decode(&bdy)
if a["token"] != bdy["token"] { t.Fatal("tokens not equal for same
URL") }
}
func TestRedirect_And_Stats(t *testing.T) {
app, _
:= setupTestApp(t)
// create
req := httptest.NewRequest(http.MethodPost, "/shorten",
strings.NewReader(`{"url":"https://example.com"}`))
req.Header.Set("Content-Type", "application/json")
resp, _
:= app.Test(req)
var body map[string]string
json.NewDecoder(resp.Body).Decode(&body)
token := body["token"]
// redirect 3 times
for i := 0; i < 3; i++ {
redir := httptest.NewRequest(http.MethodGet, "/"+token, nil)
resp, _
:= app.Test(redir)
if resp.StatusCode != http.StatusFound { t.Fatalf("want 302, got
%d", resp.StatusCode) }
}
// stats
statsReq := httptest.NewRequest(http.MethodGet, "/stats/"+token, nil)
statsResp, _
:= app.Test(statsReq)
if statsResp.StatusCode != http.StatusOK { t.Fatalf("want 200, got %d",
statsResp.StatusCode) }
var stats models.URLMapping
json.NewDecoder(statsResp.Body).Decode(&stats)
if stats.Clicks != 3 { t.Fatalf("want 3 clicks, got %d", stats.Clicks) }
}
func TestRedirect_UnknownToken(t *testing.T) {
app, _
:= setupTestApp(t)
req := httptest.NewRequest(http.MethodGet, "/doesnotexist", nil)
resp, _
:= app.Test(req)
if resp.StatusCode != http.StatusNotFound { t.Fatalf("want 404, got %d",
resp.StatusCode) }
}
func TestStats_UnknownToken(t *testing.T) {
app, _
:= setupTestApp(t)
req := httptest.NewRequest(http.MethodGet, "/stats/unknown", nil)
resp, _
:= app.Test(req)
if resp.StatusCode != http.StatusNotFound { t.Fatalf("want 404, got %d",
resp.StatusCode) }
}