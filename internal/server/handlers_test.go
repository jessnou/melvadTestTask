package server

import (
	"bytes"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
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
		{
			name: "InternalError",
			body: gin.H{
				"key":   "age",
				"value": 10,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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
