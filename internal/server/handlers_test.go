package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
	"melvadTestTask/internal/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRedisIncrement(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{{
		name: "OK",
		body: gin.H{
			"key":   "age",
			"value": 10,
		},
		checkResponse: func(recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			var response map[string]int
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("Failed to unmarshal response JSON: %v", err)
			}
			expectedValue := 10
			require.Equal(t, expectedValue, response["value"])
		},
	},
		{
			name: "InvalidRequest",
			body: gin.H{
				"key":   "",
				"value": 1,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		redisMock, mr := testRedis()
		defer mr.Close()
		router.POST("/redis/incr", handleRedisIncrement(redisMock))
		rr := httptest.NewRecorder()
		data, err := json.Marshal(tc.body)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/redis/incr", bytes.NewReader(data))
		require.NoError(t, err)

		router.ServeHTTP(rr, req)
		tc.checkResponse(rr)
	}

}

func TestHandleHMACSHA512(t *testing.T) {

	testCases := []struct {
		name          string
		body          gin.H
		checkResponce func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"key":  "mysecretkey",
				"text": "sometext",
			},
			checkResponce: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "EmptyKey",
			body: gin.H{
				"key":  "",
				"text": "sometext",
			},
			checkResponce: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.POST("/sign/hmacsha512", handleHMACSHA512())
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, _ := http.NewRequest("POST", "/sign/hmacsha512", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			tc.checkResponce(recorder)
		})
	}
}

func TestHandlePostgresUser(t *testing.T) {
	mockDB := mocks.MockPostgres{
		ExecFunc: func(query string, args ...interface{}) (sql.Result, error) {
			// Возвращаем успешный результат, чтобы симулировать успешное выполнение запроса
			return nil, nil
		},
		GetFunc: func(dest interface{}, query string, args ...interface{}) error {
			return nil
		},
	}
	testCases := []struct {
		name          string
		body          gin.H
		mockDB        mocks.MockPostgres
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name": "Alex",
				"age":  10,
			},
			mockDB: mockDB,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response map[string]int
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response JSON: %v", err)
				}
				expectedValue := 0
				require.Equal(t, expectedValue, response["id"])
			},
		},
		{
			name: "InvalidRequest",
			body: gin.H{
				"name": "",
				"age":  1,
			},
			mockDB: mockDB,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError1",
			body: gin.H{
				"name": "age",
				"age":  1,
			},
			mockDB: mocks.MockPostgres{
				ExecFunc: func(query string, args ...interface{}) (sql.Result, error) {
					// Возвращаем ошибку, чтобы симулировать внутреннюю ошибку
					return nil, errors.New("Internal error")
				},
				GetFunc: func(dest interface{}, query string, args ...interface{}) error {
					return nil
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalError2",
			body: gin.H{
				"name": "age",
				"age":  1,
			},
			mockDB: mocks.MockPostgres{
				ExecFunc: func(query string, args ...interface{}) (sql.Result, error) {
					return nil, nil
				},
				GetFunc: func(dest interface{}, query string, args ...interface{}) error {
					return errors.New("Internal error")
				},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		router.POST("/postgres/users", handlePostgresUser(tc.mockDB))
		rr := httptest.NewRecorder()
		data, err := json.Marshal(tc.body)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/postgres/users", bytes.NewReader(data))
		require.NoError(t, err)

		router.ServeHTTP(rr, req)
		tc.checkResponse(rr)
	}
}

func testRedis() (*redismock.ClientMock, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return redismock.NewNiceMock(client), mr
}
