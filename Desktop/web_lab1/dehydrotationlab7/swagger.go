package docs

// @title Dehydration Lab API
// @version 4.0
// @description API для управления заявками на расчет обезвоживания пациентов
// @description
// @description Без авторизации доступны:
// @description - Регистрация и аутентификация
// @description - Просмотр симптомов (GET)
// @description
// @description С авторизацией (пользователь):
// @description - Управление своими заявками
// @description - Добавление симптомов в заявки
// @description
// @description С ролью модератора:
// @description - Все методы пользователя
// @description - CRUD операции с симптомами
// @description - Завершение/отклонение заявок

// @contact.name API Support
// @contact.email support@dehydrationlab.com

// @host localhost:8082
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите JWT токен в формате: Bearer {token}

// @securityDefinitions.apikey CookieAuth
// @in header
// @name Cookie
// @description Сессия через cookie (session_id)

// @tag.name auth
// @tag.description Аутентификация и авторизация

// @tag.name symptoms
// @tag.description Управление симптомами (чтение доступно всем, изменение - только модераторам)

// @tag.name requests
// @tag.description Управление заявками (требуется авторизация)

// @tag.name request-symptoms
// @tag.description Управление симптомами в заявках (требуется авторизация)

var _ = ""
