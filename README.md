# DAO Vote

`dao_vote` - это система голосования для децентрализованных автономных организаций (DAO), реализованная на языке Go. Она предоставляет конечные точки для обработки голосов пользователей, расчета силы голоса и хранения голосов.

## Эндпоинты

## Обзор кода

### Обработчики (Handlers)

- `auth_handler.go`: Управляет аутентификацией пользователей.
- `dao_team_vote_handler.go`: Обрабатывает операции голосования команды DAO.
- `vote_handler.go`: Содержит HTTP обработчики для действий, связанных с голосованием.
- `withdraw_handler.go`: Управляет выводом голосов.

### Модели (Models)

- `common.go`: Содержит общие функции и переменные, используемые в различных частях приложения.
- `vote.go`: Определяет структуру UserVote и связанные модели данных.

### Репозиторий (Repository)

- `user_repository.go`: Управляет взаимодействием с данными пользователей.
- `vote_repository.go`: Управляет взаимодействием с данными голосов.

### Сервисы (Services)

- `dao_team_vote_service.go`: Содержит логику для голосования команды DAO.
- `vote_service.go`: Содержит логику для общих операций голосования.

### Утилиты (Utils)

- `response.go`: Вспомогательные функции для работы с HTTP ответами.

## Условия предоставления ПО

- **Базовая лицензия**: 10 USDT в месяц
- **Профессиональная лицензия**: 25 USDT в месяц
- **Корпоративная лицензия**: 2000 USDT единоразово
- **Корпоративная поддержка**: 1500$ ежемесячно

[Подробнее](#)

## Лицензия

Этот проект лицензирован по лицензии MIT. Для легального использования программного обеспечения GO DAO необходимо иметь лицензию, дающую права на платформе. Подробности см. в файле LICENSE.

## Инструкция

### Как все происходит

1. **Покупка лицензии**: Лицензия покупается в магазине.
2. **Получение ключа и установочного файла**: 
    - После оплаты вы получите лицензионный ключ и ссылку на установочный файл приложения.
    - Установочный файл необходимо скачать после того, как система зачислит оплату.
    - Ключ сохраняйте, чтобы можно было получить доступ с любого устройства.
3. **Установка программы**: 
    - Запустите скачанный установочный файл и следуйте инструкциям по установке.
4. **Регистрация учетной записи**:
    - После установки программы необходимо зарегистрировать учетную запись для работы с программой.
    - Переходя по ссылке с единоразовой активацией, регистрация произойдет автоматически под личным ключом из магазина.

## Запуск приложения

Для сборки и запуска приложения:

```sh
go build -o dao_vote main.go
./dao_vote
```

Программное обеспечение распространяется как компьютерная программа, задействующая как локальные, так и удаленные серверы, для открытия клиент-серверного взаимодействия.

## OPEN API

После команды `go build` в main.go на локальном порту развернется Swagger. Документация доступна по адресу: [http://localhost:8080/swagger/](http://localhost:8080/swagger/).
