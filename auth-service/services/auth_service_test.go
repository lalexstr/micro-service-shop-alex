package services

import (
	"os"
	"testing"

	"auth-service/config"
	"auth-service/db"
	"auth-service/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	// Устанавливаем переменные окружения перед инициализацией config
	os.Setenv("SQLITE_PATH", ":memory:")
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-only")
	os.Setenv("JWT_TTL_MINUTES", "60")

	// Переустанавливаем конфигурацию для тестов
	config.SQLitePath = ":memory:"
	config.JWTSecret = []byte("test-secret-key-for-testing-only")
	config.JWTTTLMin = 60

	// Инициализируем тестовую БД напрямую
	var err error
	db.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}

	// Миграция схемы
	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate test database")
	}
}

func TestAuthService_Register(t *testing.T) {
	setupTestDB()

	svc := NewAuthService()

	tests := []struct {
		name      string
		email     string
		password  string
		fullName  string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "Успешная регистрация нового пользователя",
			email:     "test@example.com",
			password:  "password123",
			fullName:  "Test User",
			wantError: false,
		},
		{
			name:      "Ошибка при регистрации с существующим email",
			email:     "test@example.com",
			password:  "password123",
			fullName:  "Another User",
			wantError: true,
			errorMsg:  "email already used",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.Register(tt.email, tt.password, tt.fullName)

			if tt.wantError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но её не было")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Ожидалась ошибка '%s', получена '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
					return
				}
				if user == nil {
					t.Errorf("Пользователь не был создан")
					return
				}
				if user.Email != tt.email {
					t.Errorf("Email не совпадает. Ожидалось: %s, получено: %s", tt.email, user.Email)
				}
				if user.FullName != tt.fullName {
					t.Errorf("FullName не совпадает. Ожидалось: %s, получено: %s", tt.fullName, user.FullName)
				}
				if user.Role != "user" {
					t.Errorf("Роль должна быть 'user', получено: %s", user.Role)
				}
				if user.Password == "" {
					t.Errorf("Пароль должен быть установлен")
				}
				// Проверяем, что пароль захеширован
				if user.Password == tt.password {
					t.Errorf("Пароль не должен храниться в открытом виде")
				}
			}
		})
	}
}

func TestAuthService_Authenticate(t *testing.T) {
	setupTestDB()

	svc := NewAuthService()

	// Сначала регистрируем пользователя
	email := "auth@example.com"
	password := "testpass123"
	fullName := "Auth User"

	_, err := svc.Register(email, password, fullName)
	if err != nil {
		t.Fatalf("Не удалось зарегистрировать пользователя для теста: %v", err)
	}

	tests := []struct {
		name      string
		email     string
		password  string
		wantError bool
	}{
		{
			name:      "Успешная аутентификация с правильными данными",
			email:     email,
			password:  password,
			wantError: false,
		},
		{
			name:      "Ошибка при неправильном пароле",
			email:     email,
			password:  "wrongpassword",
			wantError: true,
		},
		{
			name:      "Ошибка при несуществующем email",
			email:     "nonexistent@example.com",
			password:  password,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.Authenticate(tt.email, tt.password)

			if tt.wantError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но её не было")
				}
				if user != nil {
					t.Errorf("Пользователь не должен быть возвращен при ошибке")
				}
			} else {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
					return
				}
				if user == nil {
					t.Errorf("Пользователь не был найден")
					return
				}
				if user.Email != tt.email {
					t.Errorf("Email не совпадает. Ожидалось: %s, получено: %s", tt.email, user.Email)
				}
			}
		})
	}
}
