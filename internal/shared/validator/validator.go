package validator

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// Глобальный экземпляр валидатора
var validate *validator.Validate

// Регулярное выражение для username: только латинские буквы (lowercase), цифры, дефис, подчеркивание
var usernameRegex = regexp.MustCompile(`^[a-z0-9_-]+$`)

func init() {
	validate = validator.New()

	// Регистрируем кастомный валидатор для username
	validate.RegisterValidation("username", validateUsername)

	// Регистрируем кастомный валидатор для сложного пароля
	validate.RegisterValidation("strong_password", validateStrongPassword)
}

// validateUsername проверяет формат username согласно FR-003:
// - 3-50 символов (проверяется через тег min=3,max=50)
// - только латинские буквы (lowercase), цифры, дефис, подчеркивание
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// Проверка длины
	if len(username) < 3 || len(username) > 50 {
		return false
	}

	// Проверка формата через регулярное выражение
	return usernameRegex.MatchString(username)
}

// validateStrongPassword проверяет сложность пароля согласно FR-005:
// - минимум 12 символов (проверяется через тег min=12)
// - минимум 1 заглавная буква
// - минимум 1 строчная буква
// - минимум 1 цифра
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Проверка минимальной длины
	if len(password) < 12 {
		return false
	}

	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)

	// Проверяем наличие заглавной буквы, строчной буквы и цифры
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}

		// Если все условия выполнены, можем завершить проверку раньше
		if hasUpper && hasLower && hasDigit {
			return true
		}
	}

	// Возвращаем true только если все три условия выполнены
	return hasUpper && hasLower && hasDigit
}

// Validate выполняет валидацию структуры
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// ValidationError представляет ошибку валидации конкретного поля
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// GetValidationErrors преобразует ошибки валидатора в удобочитаемый формат
func GetValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	// Проверяем, является ли ошибка ошибкой валидации
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			ve := ValidationError{
				Field:   fieldErr.Field(),
				Message: getErrorMessage(fieldErr),
			}
			validationErrors = append(validationErrors, ve)
		}
	}

	return validationErrors
}

// getErrorMessage возвращает понятное сообщение об ошибке на основе типа ошибки валидации
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Это поле обязательно для заполнения"
	case "username":
		return "Логин должен содержать от 3 до 50 символов и состоять только из латинских букв, цифр, дефиса и подчеркивания"
	case "strong_password":
		return "Пароль должен содержать минимум 12 символов, включая заглавные и строчные буквы, а также цифры"
	case "eqfield":
		return "Значение должно совпадать с полем " + fe.Param()
	case "min":
		return "Минимальная длина: " + fe.Param() + " символов"
	case "max":
		return "Максимальная длина: " + fe.Param() + " символов"
	case "email":
		return "Некорректный формат email"
	default:
		return "Некорректное значение поля"
	}
}
