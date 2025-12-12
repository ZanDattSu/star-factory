![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/ZanDattSu/f617c6f48434b663c898480190e13075/raw/coverage.json)

Стек: Golang, PosgreSQL, Kafka, Docker, Redis, MongoDB, Taskfile, REST API, gRPC, Git, Linux.

---
## Требования

- **Go** ≥ `1.24`
- **protoc** ≥ `3.21`
- **Buf CLI** (`task install-buf`) - для генерации protobuf
- **Node.js + npm** ≥ `18` - нужны только для сборки OpenAPI через **Redocly**
- Docker
- **Taskfile CLI** → [инструкция по установке](https://taskfile.dev/#/installation)

---

## Быстрый старт
```
# Клонируем репозиторий
git clone https://github.com/{yourusername}/star-factory.git
cd star-factory

# Установка зависимостей
go work sync

# Поднять окружение в docker-compose и запустить сервисы
task run

# Запускаем тесты с подсчетом покрытия
task test-coverage
```

---

## Архитектура проекта
![architecture.png](readme/architecture.png)

Полную версию можно посмотреть по ссылке: https://excalidraw.com/#json=ZFRc_9KZJyyPMIi4lOCbb,1ucK2qdt1BvLS9mC1LiK-g

В проекте используется монорепозиторий с Go Workspaces. Каждый сервис — в своём модуле с `go.mod`, подключённом в `go.work`.

--- 

## Общая структура монорепозитория

```
.
├── assembly/ — асинхронная сборка заказа через Kafka
├── auth/ — авторизация и хранения пользователей
├── notification/ — отправка уведомлений о собранных заказах
├── inventory/  — каталог деталей (gRPC + MongoDB)
├── order/ — управление заказами (OpenAPI HTTP + gRPC + PostgreSQL + Kafka)
├── payment/ — симуляция платёжного шлюза (gRPC)
├── platform/ — платформенная библиотека
├── shared/ — контракты и сгенерированный код
├── deploy/ — инфраструктура и развертывание 
├── .github
├── ├── workflows - CI
├── Taskfile.yml
└── go.work
```

### shared

- OpenAPI контракты для OrderService

- Protobuf контракты для Auth/Inventory/Payment/Events

- Автоcгенерированный код клиентов и серверов


### platform

Общая библиотека для всех сервисов:

- Собственная обёртка над `uber/zap`, с пробросом traceIDKey и userIDKey для возможности добавления трейсинга в будущем

- Обертка над `IBM/sarama` для kafka producer/consumer (авто-ребалансировка, graceful shutdown)

- Обертка над `gomodule/redigo` (контроль пула соединений)

- Компоненты для работы с `Goose` миграциями

- Компоненты для `health check` проверки состояния сервиса 

- Компонент `closer`, отвечающий за корректное закрытие ресурсов в порядке `LIFO`

- Универсальная структура gRPC сервера

- gRPC интерцепторы для логирования, валидации, аутентификации

- Http мидлвар для аутентификации между сервисами

## Ключевые особенности

---
### Чистая архитектура

Проект построен в соответствии с принципами чистой архитектуры и разделен на следующие слои:

#### 1. App
- Располагается в `internal/app`
- Точка входа в приложение
- Проброс конфигов
- Автоматически инициализация зависимостей через Dependency Injection контейнер

#### 2. API слой (Адаптеры)
- Располагается в `internal/api`
- Отвечает за обработку внешних запросов (gRPC, HTTP)
- Преобразует данные из внешнего формата во внутренние модели
- Делегирует выполнение бизнес-логики в сервисный слой

#### 3. Сервисный слой (Use Cases)
- Располагается в `internal/service`
- Содержит бизнес-логику приложения
- Не зависит от внешних деталей (БД, протоколы и т.д.)
- Оперирует только сервисными моделями

#### 3. Репозиторный слой (Адаптеры)
- Располагается в `internal/repository`
- Oтвечает за доступ к данным (SQL, NoSQL, InMemory)
- Скрывает детали хранения данных от остальных слоев 
- Имеет собственные модели данных (`internal/repository/model`), отличные от доменных моделей
- Использует конвертеры для преобразования между моделями репозитория и доменными моделями

#### 4. Модели (Entities)
- Располагаются в `internal/model`
- Представляют основные бизнес-сущности
- Не зависят от других слоев

#### 5. Конвертеры
- Располагаются в `internal/converter` и `internal/repository/converter`
- Отвечают за преобразование данных между различными форматами
- Обеспечивают изоляцию между слоями
- Включают:
    - Конвертеры между Proto/Ogen и доменными моделями (в `internal/converter`)
    - Конвертеры между доменными моделями и моделями репозитория (в `internal/repository/converter`)

### Преимущества данной архитектуры

1. **Разделение ответственностей** - каждый слой имеет четко определенную ответственность
2. **Изоляция зависимостей** - зависимости направлены внутрь (к ядру приложения)
3. **Тестируемость** - слои можно тестировать изолированно
4. **Гибкость** - можно легко заменить конкретные реализации (например, базу данных)
5. **Устойчивость к изменениям** - изменения в одном слое минимально влияют на другие

Все сервисы подчиняются единой структуре:

```
order 
│   ├── cmd
│   ├── internal
│   │   ├── api - хендлеры (gRPC, HTTP)
│   │   ├── app - точка входа в приложение и dependency injection
│   │   ├── client (опционально) - внешние сервисы (например, Inventory)
│   │   ├── config - конфиги
│   │   │   ├── env
│   │   ├── model - сущности сервисного слоя
│   │   ├── repository - хранилища
│   │   │   ├── model - сущности и конверторы repo слоя
│   │   │   ├── order
│   │   │   │   └── postgresql
│   │   ├── server - обертка на http сервером
│   │   └── service - бизнес-логика (интерфейсы и реализации)
│   │       ├── consumer
│   │       ├── order
│   │       ├── produser
│   └── migrations (опционально)
```

---

### Тесты

Unit-тесты покрывают бизнес-логику **без внешних зависимостей**.

Все внешние вызовы — **через интерфейсы и моки** (автоматическая генерация моков с помощью Mockery).

Тестовые данные подготавливаются через библиотеку gofakeit

**Моки и тесты** для всех слоёв размещаются рядом с реализациями.

По всех unit тестах используется паттерн Arrange–Act–Assert (AAA) через **Test Suite** из `testify/suite`:

```go
type ServiceSuite struct {
    suite.Suite  
	  
	ctx context.Context //nolint:containedctx  
	  
	orderRepository      *mocks.OrderRepository  
	paymentClient        *clientMocks.PaymentClient  
	inventoryClient      *clientMocks.InventoryClient  
	orderProducerService *serviceMocks.OrderProducerService  
	  
	service *service
}

func (s *SuiteService) TestPayOrderSuccess() { 
	// Arrange 
	
    order := &model.Order{}  
    paymentMethod := RandomPaymentMethod()  
    expectedTransactionUUID := gofakeit.UUID()  
	  
    s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).  
       Return(order, nil).Once()  
	  
    s.orderRepository.On("UpdateOrder", ...})).Return(nil).Once()  
	  
    s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID, order.UserUUID, paymentMethod).  
       Return(expectedTransactionUUID, nil).Once()  
	  
    s.orderProducerService.On("ProduceOrderPaid", s.ctx, mock.Anything).Return(nil)  
	  
	// Act
	
    transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod, order.OrderUUID)  
	  
	// Assert
	
    s.Require().NoError(err)  
    s.Require().Equal(expectedTransactionUUID, transactionUUID)  
}
```

---

### CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Выполняется автоматическое извлечение версий из Taskfile.yml
  - Запуск линтера golangci-lint
  - Тестирование и подсчёт процента тестового покрытия
  - Процент покрытия обновляется в gist

---

### Единая конфигурация через переменные окружения

Описывается в `deploy/env/.env`, генерируется в `deploy/compose/{service_name}.env`:

Переменные окружения автоматически загружаются в сервисы в `{service_name}/internal/config`: реализовано через интерфейсы, можно подменить реализацию конфигов, например на yaml или json

---

### Docker Compose

Полностью автоматизирует поднятие инфраструктуры
##### TL;DR
**Поднять всю инфраструктуру**: `task up-all`

Используются отдельные compose файлы для описания зависимостей каждого сервиса и core для описание зависимостей всего приложения:

Docker-network для связи compose файлов

- Использует плейсхолдеры для подстановки переменных окружения из .env
- Создает volumes для серсисов которые должны хранить состояние.
- Задаёт переменные окружения которые вычитываются из .env файлов
- Делает healthcheck для каждой зависимости, с политикой ретраев и автоматически перезапускает контейнер при сбоях

---

### Kafka инфраструктура

Развернута в KRaft-режиме (без Zookeeper).

Конфигурация в `deploy/compose/core/docker-compose.yml`:

- 1 брокер

- Один узел выполняет роли:
  - **controller** (управляет метаданными)

  - **broker** (принимает и отдает сообщения)

- Собственный Volume `kafka_data` (сообщения в топиках 7 дней)

- Авто создание топиков

- Гарантия доставки: At-least-once для всех сервисов

##### В docker-compose поднят также Kafka UI:

- Образ: `provectuslabs/kafka-ui`

- Доступен по адресу: `http://localhost:8090`


Там можно смотреть:

- топики

- ключи/значения сообщений

- consumer groups

- offset'ы


#### Взаимодействие сервисов с kafka

##### OrderService

- **Producer →** `order.paid`

- **Consumer ←** `ship.assembled`

- Consumer group: `order-group-order-assembled`


##### AssemblyService

- **Consumer ←** `order.paid`

- Consumer group: `assembly-group-order-paid`

- **Producer →** `ship.assembled`

##### NotificationService

- **Consumer ←** `order.paid`

- Consumer group: `notification-group-order-paid`

- **Consumer ←** `ship.assembled`

- Consumer group: `notification-group-ship-assembled`

Сервисы подключаются через единый KafkaConfig из ENV.

---

### Taskfile

Файл с готовыми командами для генерации кода, моков, форматирования, линтинга и поднятия окружения, прогона тестов и других задач.
Полный список команд см. [Taskfile](readme/taskfile.md)

---
### Так же для Order, Inventory, Payment сервисов имеется OpenApi Swagger документация.

При старте сервисов по адресу (можно изменить в конфиге)

Inventory: localhost:8081

Payment: : localhost:8082

Order: если у вас IDE Goland => ./shared/api/order/v1/order.openapi.yaml можно открыть через неё

Если нет => ./shared/api/bundles/order.openapi.v1.bundle.yaml копируем содержимое, заходим на сайт https://editor.swagger.io/ и вставляем его туда

---

## OrderService
**Центральный сервис для оформления заказов**

Общается с Inventory и Payment по gRPC, пишет данные в PostgreSQL, публикует события в Kafka. Доступен по HTTP.
#### Архитектурные особенности:

- HTTP API(chi роутер) строго по OpenAPI контракту (Ogen)

- Защита от Slowloris атак, через readHeaderTimeout = 5s

- middlewares: 
  - Logger 
  - Recoverer
  - Response timeout = 10s

- gRPC клиенты InventoryService и PaymentService

- PostgreSQL с миграциями Goose

- Kafka producer/consumer

- DI-контейнер

- чистая архитектура: api → service → repository → postgres
#### Основные эндпоинты:

1. `POST /api/v1/orders` — создание заказа

   Создаёт новый заказ на основе выбранных пользователем деталей.

   **Поведение:**
   - Получает детали через `InventoryService.ListParts`.
   - Проверяет, что все детали существуют. Если хотя бы одной нет — возвращает ошибку.
   - Считает `total_price`.
   - Генерирует `order_uuid`.
   - Сохраняет заказ со статусом `PENDING_PAYMENT`.

2. `POST /api/v1/orders/{order_uuid}/pay` — оплата заказа

   Проводит оплату ранее созданного заказа.

   **Поведение:**
   - Находит заказ по `order_uuid`. Если не существует — возвращает 404 Not Found.
   - Вызывает `PaymentService.PayOrder`, передаёт `user_uuid`, `order_uuid` и `payment_method`. Получает`transaction_uuid`.
   - Обновляет заказ: статус → `PAID`, сохраняет `transaction_uuid`, `payment_method`.
   - Публикует события в топик `order.paid` в Kafka.
   - По ходу работу сервиса асинхронно слушает топик `ship.Assembled`, вычитывает событие ShipAssembled и обновляет статус в БД

3. `GET /api/v1/orders/{order_uuid}`  —  получить заказ по UUID

   Возвращает информацию о заказе.

   **Поведение:**
   - Ищет заказ по UUID.
   - Если найден — возвращает.
   - Если не найден — 404 Not Found.

4. `POST /api/v1/orders/{order_uuid}/cancel` — отменить заказ

   Отменяет заказ.

   **Ответы:**
   - `204 No Content` — заказ успешно отменён
   - `404 Not Found` — заказ не найден
   - `409 Conflict` — заказ уже оплачен и не может быть отменён

   **Поведение:**
   - Проверяет статус заказа.
   - Если `PENDING_PAYMENT` — меняет статус на `CANCELLED`.
   - Если `PAID` — возвращает ошибку 409.

## InventoryService
**Сервис хранения и поиска деталей**

gRPC-сервис с MongoDB, предоставляет OrderService информацию о деталях при оформлении заказов.

#### Архитектурные особенности:
- gRPC методы сгенерированные через protobuf по proto-контракту
- HTTP Gateway с Swagger UI (OpenAPI, сгенерировано из proto)
- MongoDB с коллекцией `parts`, индекс по `uuid`
- Кастомные интерцепторы для логирования и валидации запросов

#### Основные ручки:

1. `GetPart(uuid string) Part` — возврат информации о детали по её UUID

   **Поведение:**
    - Валидирует формат uuid.
    - Запрашивает документ в MongoDB по индексу uuid.
    - Если не найден — возвращает ошибку NotFound.

2. `ListParts(filter PartsFilter) []Part` — получение списка деталей с фильтрацией

    **Поведение:**
    - Если все поля фильтра пусты — возвращаются все детали.
    - Фильтрация происходит по принципу:
        - *логическое ИЛИ внутри одного поля фильтра* (например, имя `"main"` **или** `"main booster"`)
        - *логическое И между различными полями* (например, категория = `ENGINE` **и** страна = `"Germany"`)
    - Фильтрация выполняется за 1 проход.
    - Возвращает массив найденных деталей.

---

## PaymentService
**Сервис обработки оплаты заказов**

Принимает gRPC-запросы от Order, не имеет своей базы

#### Архитектурные особенности:
- gRPC метод, сгенерированный через protobuf по proto-контракту
- HTTP Gateway с Swagger UI (OpenAPI, сгенерировано из proto)
- Генерация UUID v4 для каждой транзакции
- Логирование успешных оплат

#### Основная ручка:
Обрабатывает оплату заказов

1. `PayOrder(order_uuid, user_uuid, payment_method) transaction_uuid` — обработка команды на оплату заказа

   **Поведение:**
    - Валидирует входящие поля.
    - Генерирует `transaction_uuid` (UUID v4).
    - Логирует сообщение формата:
      ```
      Оплата прошла успешно, transaction_uuid: <uuid>
      ```
    - Возвращает `transaction_uuid` вызывающей стороне.
    - Состояние не сохраняется.

---

## AssemblyService
**Сервис асинхронной сборки кораблей**

Фоновый Kafka-консьюмер, который реагирует на событие оплаты заказа.  
После получения `OrderPaid` имитирует сборку, ждёт 10 секунд и пушит новое событие `ShipAssembled`.

#### Архитектурные особенности:
- Kafka consumer для входящих событий `order.paid`
- Kafka producer для исходящих событий `ship.assembled`
- Асинхронная обработка без HTTP/gRPC API

#### Основной сценарий работы:

1. Обработка события `OrderPaid`** — запуск процесса сборки

   **Поведение:**
    - Получает из Kafka событие `OrderPaid`.
    - Логирует начало сборки по заказу.
    - Засыпает на рандомное время от 1 до 10 секунд, имитируя процесс сборки корабля.
    - Генерирует новый `event_uuid` (UUID v4).
    - Формирует событие `ShipAssembled`.
    - Публикует его в Kafka в топик `ship.assembled`.

2. Публикация события `ShipAssembled`** — уведомление о завершении сборки

   **Поведение:**
    - Содержит `order_uuid`, `user_uuid` и `build_time_sec`.
    - Используется OrderService для обновления статуса заказа.
    - Гарантирует идемпотентность за счёт уникального `event_uuid`.

#### Входящее событие: `OrderPaid`
Содержит:
- `event_uuid` — уникальный ID события
- `order_uuid`
- `user_uuid`
- `payment_method`
- `transaction_uuid`

#### Исходящее событие: `ShipAssembled`
Содержит:
- `event_uuid` — уникальный ID события
- `order_uuid`
- `user_uuid`
- `build_time_sec`

## NotificationService
**Сервис отправки уведомлений пользователям**

Фоновый сервис, реагирующий на бизнес-события из Kafka и отправляющий уведомления в Telegram.
Используется для информирования пользователя о ключевых этапах жизненного цикла заказа.

#### Архитектурные особенности:
- Kafka consumer для входящих событий 
  - `order.paid`
  - `ship.assembled`
- Интеграция с Telegram Bot API через библиотеку go-telegram/bot
- Асинхронная обработка без HTTP/gRPC API
- Реализована политика ретраев при инициализации Telegram-бота:
  - количество попыток и задержка настраиваются через .env 
  - корректно обрабатывает временную недоступность Telegram API

#### Основной сценарий работы:

1. Обработка события `OrderPaid` — уведомление об оплате

   **Поведение:**
   - Получает из Kafka событие `OrderPaid`.
   - Формирует уведомление об оплате корабля.
   - Отправляет сообщение в Telegram-чат с заранее заданным `chat_id`..
   - Логирует успешную или неуспешную отправку.

2. Обработка события `ShipAssembled` — уведомление о сборке

   **Поведение:**
   - Получает из Kafka событие `ShipAssembled`.
   - Формирует уведомление о завершении сборки корабля.
   - Отправляет сообщение в Telegram-чат с заранее заданным `chat_id`..
   - Логирует успешную или неуспешную отправку.

3. Инициализация Telegram-бота

   **Поведение:**
   - При старте сервиса создаёт Telegram-бота.
   - Использует **retry-политику** при ошибках сети или Telegram API.
   - Регистрирует команду `/start`.

#### Входящие событие: 

`OrderPaid`

Содержит:
- `event_uuid` — уникальный ID события
- `order_uuid`
- `user_uuid`
- `payment_method`
- `transaction_uuid`

`ShipAssembled`

Содержит:
- `event_uuid` — уникальный ID события
- `order_uuid`
- `user_uuid`
- `build_time_sec`

## AuthService
**Сервис аутентификации и управления пользователями**

`AuthService` отвечает за регистрацию пользователей, аутентификацию, управление сессиями и получение информации о пользователе.
Используется как центральный сервис авторизации для всей микросервисной системы.

#### Архитектурные особенности:
- gRPC API для работы с пользователями и сессиями
- Хранение пользователей в PostgreSQL
- Хранение активных сессий в Redis
- Унифицированная модель авторизации для HTTP и gRPC

#### Основные ручки:

1. Регистрация пользователя (`Register`)

   **Поведение:**
   - Принимает данные нового пользователя (логин, пароль, email).
   - Хеширует пароль перед сохранением.
   - Сохраняет пользователя в PostgreSQL.
   - Сохраняет предпочтительные каналы уведомлений.
   - Возвращает `user_uuid` зарегистрированного пользователя.

2. Аутентификация пользователя (`Login`)

   **Поведение:**
   - Проверяет корректность логина и пароля.
   - Генерирует новую сессию (`session_uuid`).
   - Сохраняет сессию в Redis с `TTL` 24 часа. (настраивается в `env`)
   - Возвращает `session_uuid` клиенту.

3. Получение информации о текущем пользователе (`Whoami`)

   **Поведение:**
   - Извлекает `session_uuid` из `gRPC metadata`.
   - Валидирует сессию через Redis.
   - Возвращает информацию о текущем пользователе.
   - Используется другими сервисами для проверки авторизации.

4. Получение данных пользователя (`GetUser`)

   **Поведение:**
   - Принимает `user_uuid`.
   - Возвращает информацию о пользователе и его каналах уведомлений.
   - Может использоваться для проверки прав доступа.

#### Авторизация

- Реализована централизованная система авторизации для gRPC и HTTP микросервисов. 
- В gRPC используется `interceptor`, который извлекает `session-uuid` из `metadata` и валидирует сессию через `AuthService.Whoami`.
- В HTTP реализован `middleware`, который читает `X-Session-Uuid`, запрашивает данные пользователя у AuthService и добавляет их в контекст.
- Авторизация полностью интегрирована с `context.Context` и не требует ручной проверки в бизнес-логике.
- Сессия и информация о пользователе доступны через контекст и безопасно передаются между сервисами.
- Для межсервисных gRPC вызовов предусмотрена функция, которая передает session-uuid из контекста в metadata.

- Такая архитектура обеспечивает слабую связанность компонентов и высокий уровень переиспользования кода.
- Решение облегчает разработку защищённых микросервисов и упрощает поддержку кода.
- Поддерживает расширение (например, role-based access) без изменения бизнес-обработчиков.