package tests

import (
	"bytes"
	"dating-apps/app"
	"dating-apps/app/api/initialization"
	"dating-apps/app/model/entity"
	"dating-apps/app/model/response"
	"dating-apps/helper/config"
	"dating-apps/helper/database"
	"dating-apps/helper/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"
)

var (
	db  database.Database
	cfg *config.Config
	gdb *gorm.DB
	mux *http.ServeMux
)

// responseHttp is the expected structure of our JSON response.
type responseHttp struct {
	Meta   meta        `json:"meta"`
	Data   data        `json:"data"`
	Errors interface{} `json:"errors,omitempty"`
}

type meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type data struct {
	Record response.LoginResponse `json:"record,omitempty"`
}

// Sample token (adjust if needed)
var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6Im1hdWx2aUBleGFtcGxlLmNvbSIsImV4cCI6MTczOTU0MDU5MiwiaWF0IjoxNzM5NDU0MTkyfQ.TbvrUUmSsx-im2_QvX8aaEYVLgMHPkqapY3lMcGKdoo"

// TestMain initializes configuration, database, logger and routes.
func TestMain(m *testing.M) {
	cfg = config.Init(config.WithConfigPath("../"))
	log := logger.NewLogger(&cfg.ServerConfig.LogConfig)

	// Initialize DB (using your helper / initialization code)
	db = initDatabase(cfg)
	gdb, _ = db.Client().(*gorm.DB)

	// Pass infra to your application.
	infra := &app.Infra{Db: &db, Log: log, Config: cfg}

	mux = initialization.InitRouting(infra)
	code := m.Run()
	os.Exit(code)
}

// initDatabase wraps your initialization function.
func initDatabase(cfg *config.Config) database.Database {
	if db != nil {
		return db
	}
	db, err := initialization.InitDatabase(&cfg.DBConfig)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// executeRequest sends a request through the mux and returns the recorder.
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

//------------------------------
// Integration tests for HTTP endpoints using table-driven tests
//------------------------------

func TestSignUp(t *testing.T) {
	tests := []struct {
		name           string
		payload        []byte
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "Valid SignUp",
			payload:        []byte(`{"email":"maulvi@example.com","password":"secret","name":"New User","gender":"M"}`),
			expectedStatus: http.StatusOK,
			expectedCode:   200,
			expectedMsg:    "Success",
		},
		{
			name:           "Duplicate Email",
			payload:        []byte(`{"email":"maulvi@example.com","password":"secret","name":"Maulvi","gender":"M"}`),
			expectedStatus: http.StatusBadRequest,
			// Adjust expected error code/message as defined in your message package.
			expectedCode: 400,
			expectedMsg:  "Data already exists",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, cfg.UrlWithPrefix("auth/signup"), bytes.NewBuffer(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			response := executeRequest(req)

			if response.Code != tc.expectedStatus {
				t.Errorf("Expected response code %d. Got %d", tc.expectedStatus, response.Code)
			}

			var resp responseHttp
			if err := json.Unmarshal(response.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if resp.Meta.Code != tc.expectedCode {
				t.Errorf("Expected Meta Code %d. Got %d", tc.expectedCode, resp.Meta.Code)
			}

			if resp.Meta.Message != tc.expectedMsg {
				t.Errorf("Expected Meta Message '%s'. Got '%s'", tc.expectedMsg, resp.Meta.Message)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		payload        []byte
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "Valid Login",
			payload:        []byte(`{"email":"maulvi@example.com","password":"secret"}`),
			expectedStatus: http.StatusOK,
			expectedCode:   200,
			expectedMsg:    "Success",
		},
		{
			name:           "Invalid Password",
			payload:        []byte(`{"email":"maulvi@example.com","password":"wrong"}`),
			expectedStatus: http.StatusUnauthorized,
			// Adjust expected error code/message as defined in your message package.
			expectedCode: 401,
			expectedMsg:  "Invalid credentials",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, cfg.UrlWithPrefix("auth/login"), bytes.NewBuffer(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			response := executeRequest(req)

			if response.Code != tc.expectedStatus {
				t.Errorf("Expected response code %d. Got %d", tc.expectedStatus, response.Code)
			}

			var resp responseHttp
			if err := json.Unmarshal(response.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if resp.Meta.Code != tc.expectedCode {
				t.Errorf("Expected Meta Code %d. Got %d", tc.expectedCode, resp.Meta.Code)
			}

			if resp.Meta.Message != tc.expectedMsg {
				t.Errorf("Expected Meta Message '%s'. Got '%s'", tc.expectedMsg, resp.Meta.Message)
			}

			if resp.Meta.Code == 200 && resp.Meta.Message == "Success" {
				token = resp.Data.Record.Token
			}
		})
	}
}

func TestSwipe(t *testing.T) {
	// For Swipe tests, we include the Authorization header.
	tests := []struct {
		name           string
		payload        []byte
		authToken      string
		setupFunc      func()
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "Valid Swipe",
			payload:        []byte(`{"targetUserID": 2, "action": "like"}`),
			authToken:      token,
			expectedStatus: http.StatusOK,
			expectedCode:   200,
			expectedMsg:    "Success",
		},
		{
			name:           "Invalid Action",
			payload:        []byte(`{"targetUserID": 2, "action": "invalid"}`),
			authToken:      token,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
			expectedMsg:    "action must be 'like' or 'pass'",
		},
		{
			name:           "Missing Token",
			payload:        []byte(`{"targetUserID": 2, "action": "like"}`),
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
			// You can define expected error meta code/message for missing token if your error encoder returns one.
			expectedCode: 401,
			expectedMsg:  "No authorization token was found",
		},
		{
			name:      "Already Swiped Today",
			payload:   []byte(`{"targetUserID": 2, "action": "like"}`),
			authToken: token,
			setupFunc: func() {
				// Assume the user ID decoded from the token is 1.
				// Insert a swipe record for user 1 swiping on target 2 today.
				today := time.Now().Truncate(24 * time.Hour)
				// Clean any previous record first.
				gdb.Where("user_id = ? AND target_user_id = ? AND swipe_date = ?", 1, 2, today).Delete(&entity.Swipe{})
				gdb.Create(&entity.Swipe{
					UserID:       1,
					TargetUserID: 2,
					Action:       "like",
					SwipeDate:    today,
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
			expectedMsg:    "already swiped on this profile today",
		},
		{
			name:      "Daily Swipe Limit Reached",
			payload:   []byte(`{"targetUserID": 3, "action": "like"}`),
			authToken: token,
			setupFunc: func() {
				today := time.Now().Truncate(24 * time.Hour)
				gdb.Where("user_id = ? AND swipe_date = ?", 1, today).Delete(&entity.Swipe{})
				limit := int(cfg.AppConfig.SwipeLimit)
				for i := 0; i < limit; i++ {
					if err := gdb.Create(&entity.Swipe{
						UserID:       1,
						TargetUserID: uint(i + 10),
						Action:       "like",
						SwipeDate:    today,
					}).Error; err != nil {
						t.Fatalf("Failed to create swipe record: %v", err)
					}
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   400,
			expectedMsg:    "daily swipe limit reached",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFunc != nil {
				tc.setupFunc()
			}
			req, _ := http.NewRequest(http.MethodPost, cfg.UrlWithPrefix("user/swipe"), bytes.NewBuffer(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			if tc.authToken != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.authToken))
			}
			response := executeRequest(req)

			if response.Code != tc.expectedStatus {
				t.Errorf("Expected response code %d. Got %d", tc.expectedStatus, response.Code)
			}

			var resp responseHttp
			if err := json.Unmarshal(response.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if resp.Meta.Code != tc.expectedCode {
				t.Errorf("Expected Meta Code %d. Got %d", tc.expectedCode, resp.Meta.Code)
			}

			if resp.Meta.Message != tc.expectedMsg {
				t.Errorf("Expected Meta Message '%s'. Got '%s'", tc.expectedMsg, resp.Meta.Message)
			}
		})
	}
}
