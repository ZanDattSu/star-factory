```
task: Available tasks for this project:
* down:                              Остановить и удалить все сервисы по очереди вместе с зависимостями
* down-auth:                         Остановить и удалить AUTH сервис и все его зависимости
* down-core:                         Остановить и удалить core контейнеры
* down-inventory:                    Остановить и удалить Inventory сервис и все его зависимости
* down-order:                        Остановить и удалить Order сервис и все его зависимости
* format:                            Форматирует весь проект gofumpt + gci, исключая mocks
* gen:                               Генерация всех proto и OpenAPI деклараций
* install-buf:                       Устанавливает Buf в каталог bin
* install-formatters:                Устанавливает форматтеры gci и gofumpt в ./bin
* install-golangci-lint:             Устанавливает golangci-lint в каталог bin
* lint:                              Запускает golangci-lint для всех модулей
* run:                               Запустить проект
* stop:                              Остановить все сервисы
* test:                              Запускает юнит-тесты для всех модулей
* test-api:                          Запуск тестов для проверки API микросервисов
* test-api-auth-rejection:           Тестирование корректности отклонения неавторизованных запросов
* test-api-with-auth:                Тестирование API микросервисов с корректной аутентификацией
* test-auth:                         Тестирование AUTH сервиса (регистрация, логин, whoami)
* test-coverage:                     Тесты с покрытием бизнес-логики (service/repository), отчёт по каждому модулю + общий
* tree:                              Устанавливает tree и выводит структуру проекта без node_modules
* up:                                Поднять все сервисы по очереди вместе с зависимостями
* up-auth:                           Поднять AUTH сервис и все его зависимости
* up-core:                           Поднять core контейнеры
* up-inventory:                      Поднять Inventory сервис и все его зависимости
* up-order:                          Поднять Order сервис и все его зависимости
* coverage:html:                     Генерирует HTML-отчёт покрытия и открывает его в браузере
* deps:update:                       Обновление зависимостей в go.mod во всех модулях
* env:gen:                           Генерирует .env файлы для всех сервисов из шаблонов и единого файла конфигурации
* env:install-envsubst:              Устанавливает envsubst в bin/
* grpcurl:install:                   Устанавливает grpcurl в каталог bin
* mockery:gen:                       Генерирует моки интерфейсов с помощью mockery
* mockery:install:                   Устанавливает mockery в ./bin
* ogen:gen:                          Генерация Go-кода из всех OpenAPI-деклараций с x-ogen
* ogen:install:                      Скачивает ogen в папку bin
* proto:gen:                         Генерация Go-кода из .proto
* proto:install-plugins:             Устанавливает protoc плагины в каталог bin
* proto:lint:                        Проверка .proto-файлов на соответствие стилю
* proto:update-deps:                 Обновляет зависимости protobuf из удаленных репозиториев (googleapis и т.д.)
* redocly-cli:bundle:                Собрать все схемы OpenAPI в общие файлы через локальный redocly
* redocly-cli:install:               Установить локально Redocly CLI
* redocly-cli:order-v1-bundle:       Собрать OpenAPI в один файл через локальный redocly
* yq:install:                        Устанавливает yq в bin/ при необходимости
```