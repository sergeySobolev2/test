# Контрольные ответы (шпаргалка)

- Flux: однонаправленный поток данных — действие → редьюсер → новое состояние → UI.
- Redux Toolkit: упрощает редьюсеры/стор, иммутабельность через Immer, слайсы и DevTools.
- Persist: сохранение части стора (у нас — фильтры) в `localStorage`.
- PWA vs Tauri: PWA — браузер/Pages, офлайн-режим, без нативного доступа; Tauri — нативная оболочка, системные возможности, офлайн и общий код фронта.
- GitHub Pages: статический хостинг, деплой через Actions, `base=/medical_dehydration` в Vite.
- Прокси Vite: dev-запросы `/api` направляются на `http://localhost:8082`.
- Env: `VITE_API_BASE` для прод/Pages/Tauri; `VITE_USE_MOCK` — принудительный mock.
- Offline: Workbox `NetworkFirst` для `/api/symptoms`; статика кешируется на билд-шаге.
- Адаптивность: Bootstrap гриды `xs, sm, md, lg, xl` для карточек симптомов.
