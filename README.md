   # DAO Vote

   `dao_vote` - это система голосования для децентрализованных автономных организаций (DAO), реализованная на языке Go. Она предоставляет конечные точки для обработки голосов пользователей, расчета силы голоса и хранения голосов.

   ## Структура проекта

   ```plaintext
   dao_vote/
   │
   ├── bin/
   ├── cmd/
   │   └── main.go
   ├── dsc-example/
   ├── frontend/
   ├── internal/
   │   ├── handlers/
   │   │   ├── auth_handler.go
   │   │   ├── dao_team_vote_handler.go
   │   │   ├── vote_handler.go
   │   │   └── withdraw_handler.go
   │   ├── models/
   │   │   └── vote.go
   │   ├── repository/
   │   │   ├── user_repository.go
   │   │   └── vote_repository.go
   │   ├── services/
   │   │   ├── dao_team_vote_service.go
   │   │   └── vote_service.go
   │   └── utils/
   ├── swagger/
   │   ├── favicon-16x16.png
   │   ├── favicon-32x32.png
   │   ├── index.css
   │   ├── index.html
   │   ├── oauth2-redirect.html
   │   ├── swagger.yaml
   │   ├── swagger-initializer.js
   │   ├── swagger-ui.css
   │   └── swagger-ui.js
   ```

   ## Начало работы

   ### Предварительные требования

   Убедитесь, что у вас установлен Go. Скачать и установить его можно с [golang.org](https://golang.org/dl/).

   ### Установка

   Клонируйте репозиторий:

   ```bash
   git clone https://github.com/GODAOSOFTWARE/GODAO.git
   cd GODAO
   ```

   ### Запуск приложения

   Для сборки и запуска приложения:

   ```bash
   go build ./cmd
   ./cmd
   ```

   ## Эндпоинты

   ## Обзор кода

   ### Обработчики (Handlers)

   - **vote_handler.go**: Содержит HTTP обработчики для действий, связанных с голосованием.
   - **auth_handler.go**: Управляет аутентификацией пользователей.
   - **dao_team_vote_handler.go**: Обрабатывает операции голосования команды DAO.
   - **withdraw_handler.go**: Управляет выводом голосов.

   ### Модели (Models)

   - **vote.go**: Определяет структуру `UserVote` и связанные модели данных.

   ### Репозиторий (Repository)

   - **user_repository.go**: Управляет взаимодействием с данными пользователей.
   - **vote_repository.go**: Управляет взаимодействием с данными голосов.

   ### Сервисы (Services)

   - **dao_team_vote_service.go**: Содержит логику для голосования команды DAO.
   - **vote_service.go**: Содержит логику для общих операций голосования.

   ### Утилиты (Utils)

   Вспомогательные функции, используемые в приложении.

   ## Вклад

   Вклады приветствуются! Пожалуйста, создайте issue или отправьте pull request.

   ## Лицензия

   Этот проект лицензирован по лицензии MIT - см. файл [LICENSE](LICENSE) для подробностей.
   ```
