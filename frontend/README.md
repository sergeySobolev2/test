# Лаба 7: Redux Toolkit + Axios + Swagger

Краткое руководство для демонстрации и сдачи.

## Порядок показа
- Авторизация: страницы Вход/Регистрация. После входа логин в меню.
- Добавление симптомов: на карточке «Добавить» — формируется черновик.
- Страница черновика: редактирование интенсивности, основной, комментарий; «Сформировать» переводит в статус.
- Список заявок: таблица заявок пользователя с статусами и счётчиком комментариев.
- Показ Application → Local Storage/cookies: `app:token`, `session_id`.
- Использование токена/cookie в Insomnia/Postman: просмотреть `/api/requests`.

## Архитектура состояния
- Redux Toolkit: slices `auth`, `filters`, `requests`.
  - `authSlice`: thunks `loginUser`, `registerUser`, `logoutUser`, `fetchProfile`, `updateProfile`, храним `token` и `user` в localStorage.
  - `filtersSlice`: хранит на фронте фильтры симптомов; сохраняем в localStorage.
  - `requestsSlice`: `fetchCart`, `addSymptomToDraft`, `removeSymptomFromDraft`, `listRequests`, `formRequest`.
- Store: инициализация из `localStorage`; при выходе — сброс фильтров и конструктора.

## Axios и кодогенерация
- Axios-клиент: [src/modules/services/http.ts](src/modules/services/http.ts), `withCredentials`, `Authorization: Bearer` из `localStorage`.
- Мини-клиент по swagger: [src/api/generated.ts](src/api/generated.ts).
  - Методы: `RequestsService.cart/list/get/update/form`, `RequestSymptomsService.add/remove/update`.

## Контрольные вопросы
- Redux Toolkit схема: reducer, store, middleware (thunk встроен), actions созданы слайсами.
- useContext: используется редко, в данном проекте управление состоянием — через Redux.
- axios: http-клиент с интерсептором токена, методы в `generated.ts`.
- Local Storage: хранение `app:token`, `app:user`, фильтров `app:filters`.

## Быстрый запуск
```bash
cd frontend
npm install
npm run dev
```

Бэкенд:
```powershell
task shell: Go mod tidy and build
task shell: Run migrate main.go
task shell: Run seed main.go
task shell: Run Go server (dev)
```
# Frontend (React + Vite)

## Quick start
- Install deps: `npm install`
- Dev server: `npm run dev` (http://localhost:5176)
- Backend proxy: requests to `/api` go to `http://localhost:8082`

## Env toggles
- `VITE_USE_MOCK`: `true` to use local mock data (no network)
- `VITE_API_BASE`: absolute API base for production (e.g. `https://example.com`); leave empty in dev to rely on proxy

Examples:
```
# .env.development
VITE_API_BASE=
VITE_USE_MOCK=false

# .env.production - GitHub Pages or remote API
VITE_API_BASE=https://your-host-or-lan-ip
VITE_USE_MOCK=false
```

## PWA
- Build: `npm run build`; Preview: `npm run preview`
- Place icons at `public/icons/icon-192.png` and `public/icons/icon-512.png`
- Service Worker auto-updates (via vite-plugin-pwa)

## Notes
- In mock mode: filters work client-side, images fallback to `public/placeholder.svg`
- In backend mode: ensure API is reachable by the browser (CORS/proxy)
