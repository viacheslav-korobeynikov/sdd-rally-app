package middleware

import "net/http"

// Обработчик CORS
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем Origin из запроса
		origin := r.Header.Get("Origin")
		// Если origin не передан - работаем как обычно
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}
		// Если в origin передат какой-то домен
		header := w.Header()
		// Добавляем в хедер Access-Control-Allow-Origin значение из Origin
		header.Set("Access-Control-Allow-Origin", origin)
		header.Set("Access-Control-Allow-Credentials", "true")
		// Если был вызван метод OPTIONS
		if r.Method == http.MethodOptions {
			// Описываем разрешённые HTTP-методы
			header.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD,PATCH")
			// Описываем разрешённые хедеры для отправки
			header.Set("Access-Control-Allow-Headers", "authorization,content-type,content-length")
			// Описываем кэширование предварительного запроса на N секунд
			header.Set("Access-Control-Max-Age", "86400")
		}
		next.ServeHTTP(w, r)
	})
}
