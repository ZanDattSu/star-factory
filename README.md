![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/ZanDattSu/f617c6f48434b663c898480190e13075/raw/coverage.json)

## ⚙️ Требования

- **Go** ≥ `1.24`
- **protoc** ≥ `3.21`
- **Buf CLI** (`task install-buf`) - для генерации protobuf
- **Node.js + npm** ≥ `18` - нужны только для сборки OpenAPI через **Redocly**
- **Taskfile CLI** → [инструкция по установке](https://taskfile.dev/#/installation)

Проверить версию Taskfile:
```bash
task --version
```

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Выполняется автоматическое извлечение версий из Taskfile.yml
  - Выполняется тестирование и подсчёт тестового покрытия
