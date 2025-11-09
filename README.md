![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/olezhek28/bf33a2bda0693f1162c4323702033d27/raw/coverage.json)

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:
https://taskfile.dev/docs/installation


## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Тестирование
  - В
  - Выполняется автоматическое извлечение версий из Taskfile.yml
