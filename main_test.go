package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

// тест на количества записей
func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		countOk int
	}{
		// варианты проверки
		{"/cafe?city=moscow&count=0", 0},
		{"/cafe?city=tula&count=0", 0},
		{"/cafe?city=moscow&count=1", 1},
		{"/cafe?city=tula&count=1", 1},
		{"/cafe?city=moscow&count=2", 2},
		{"/cafe?city=tula&count=2", 2},
		{"/cafe?city=moscow&count=100", len(cafeList["moscow"])},
		{"/cafe?city=tula&count=100", len(cafeList["tula"])},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		// конвертируем ответ
		countResp := strings.Split(response.Body.String(), ",")

		// для случая с 0
		if (len(countResp) == 1) && (countResp[0] == "") {

			assert.Equal(t, v.countOk, 0)
			return
		}

		// для всех вариантов которые > 0
		assert.Equal(t, v.countOk, len(countResp))

	}
}

// тест на параметр поиска
func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		// варианты проверки
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
		{"", len(cafeList["moscow"])},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+v.search, nil)
		handler.ServeHTTP(response, req)

		// конвертируем ответ
		cafeSearch := strings.Split(response.Body.String(), ",")

		found := 0

		// проверяем ответ на количество совпадений
		for _, cafe := range cafeSearch {
			if strings.Contains(strings.ToLower(cafe), strings.ToLower(v.search)) {
				found++
			}
		}

		assert.Equal(t, v.wantCount, found)

	}
}
