## Телеграм бот

### Описание
Получает курсы валют по API, сохраняет их в базе данных и отдаёт пользователям посредством REST API или Telegram бота
### Стек
- Docker
- PostgreSQL
### Алгоритм работы
- При запуске сервиса и затем каждые 5 минут происходит обновление курса валют для основных монет (BTC, ETH) в фоне
- Клиент отправляет команды на получение курса валют или отправляет запрос на сервис
- Клиент может включить / выключить автоматическую отправку данных по курсам через Telegram бот

##### Команды REST API
- GET /rates
- GET /rates/{cryptocurrency}
##### Команды Telegram бота
- /start
- /rates
- /rates {cryptocurrency}
- /start-auto {minutes_count} (пример /start-auto 10, что значит отправлять каждые 10 минут) 
- /stop-auto

### Настройка 
1. Обновите файл .env вставив в строку TOKEN=[YOUR TELEGRAM TOKEN] Ваш телеграм токен
2. Соберите контейнер Телеграм бота через Docker командой _build-image_ из файла Makefile
3. Запустите контейнер PostgreSQL через Docker командой _build-postgres_ из файла Makefile
4. Создайте базу данных _telegram_bot_
5. Установите файлы миграции на базу данных командой _migrations_up_ из файла Makefile
6. Запустите контейнер Телеграм бота командой _migrations_up_ из файла Makefile
7. Вы прекрасны :)
