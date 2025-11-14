# HP MSA Exporter

Prometheus экспортер для HP MSA Storage SAN.

**Примечание**: Проект переписан на Go для лучшей производительности и меньшего потребления ресурсов.

## Установка

### Из исходников

```bash
# Сборка
go build -o msa_exporter

# Запуск тестов
go test -v

# Проверка покрытия тестами
go test -cover
```

### Использование Docker

```bash
docker build -t msa_exporter .
docker run -e HOST=msa_hostname -e LOGIN=username -e PASSWORD=password -p 8000:8000 msa_exporter
```

**Размер Docker образа**: ~26MB (vs ~192MB Python версия - **уменьшение на 86%**)

**Безопасность**: Контейнер запускается от непривилегированного пользователя (`exporter:exporter` с UID/GID 10001)

## Использование

### Командная строка

```bash
# Использование флагов
./msa_exporter -hostname msa_san_hostname -login msa_san_username -password msa_san_password -port 8000 -interval 60 -timeout 60

# Использование позиционных аргументов (обратная совместимость)
./msa_exporter msa_san_hostname msa_san_username msa_san_password
```

### Параметры

- `-hostname string` - Имя хоста MSA storage (обязательно)
- `-login string` - Логин для MSA storage (обязательно)
- `-password string` - Пароль для MSA storage (обязательно)
- `-port int` - Порт экспортера (по умолчанию: 8000)
- `-interval int` - Интервал сбора метрик в секундах (по умолчанию: 60)
- `-timeout int` - Таймаут сбора в секундах (по умолчанию: 60)

## Метрики

Экспортер предоставляет следующие метрики:

| Название                              | Описание                        | Метки                        |
|---------------------------------------|---------------------------------|------------------------------|
| msa_hostport_data_read                | Прочитано данных                | port                         |
| msa_hostport_data_written             | Записано данных                 | port                         |
| msa_hostport_avg_resp_time_read       | Время отклика чтения            | port                         |
| msa_hostport_avg_resp_time_write      | Время отклика записи            | port                         |
| msa_hostport_avg_resp_time            | Время отклика I/O               | port                         |
| msa_hostport_queue_depth              | Глубина очереди                 | port                         |
| msa_hostport_reads                    | Операции чтения                 | port                         |
| msa_hostport_writes                   | Операции записи                 | port                         |
| msa_disk_temperature                  | Температура                     | location, serial             |
| msa_disk_iops                         | IOPS                            | location, serial             |
| msa_disk_bps                          | Байт в секунду                  | location, serial             |
| msa_disk_avg_resp_time                | Среднее время отклика I/O       | location, serial             |
| msa_disk_ssd_life_left                | Остаток ресурса SSD             | location, serial             |
| msa_disk_health                       | Состояние здоровья              | location, serial             |
| msa_disk_power_on_hours               | Часов работы                    | location, serial             |
| msa_disk_errors                       | Ошибки                          | location, port, serial, type |
| msa_volume_health                     | Состояние здоровья              | volume                       |
| msa_volume_iops                       | IOPS                            | volume                       |
| msa_volume_bps                        | Байт в секунду                  | volume                       |
| msa_volume_reads                      | Операции чтения                 | volume                       |
| msa_volume_writes                     | Операции записи                 | volume                       |
| msa_volume_data_read                  | Прочитано данных                | volume                       |
| msa_volume_data_written               | Записано данных                 | volume                       |
| msa_volume_shared_pages               | Общие страницы                  | volume                       |
| msa_volume_read_hits                  | Попадания в кеш чтения          | volume                       |
| msa_volume_read_misses                | Промахи кеша чтения             | volume                       |
| msa_volume_write_hits                 | Попадания в кеш записи          | volume                       |
| msa_volume_write_misses               | Промахи кеша записи             | volume                       |
| msa_volume_small_destage              | Малые сбросы                    | volume                       |
| msa_volume_full_stripe_write_destages | Полные сбросы stripe            | volume                       |
| msa_volume_read_ahead_ops             | Операции опережающего чтения    | volume                       |
| msa_volume_write_cache_space          | Пространство кеша записи        | volume                       |
| msa_volume_write_cache_percent        | Процент кеша записи             | volume                       |
| msa_volume_size                       | Размер                          | volume                       |
| msa_volume_total_size                 | Полный размер                   | volume                       |
| msa_volume_allocated_size             | Выделенный размер               | volume                       |
| msa_volume_blocks                     | Блоки                           | volume                       |
| msa_volume_tier_distribution          | Распределение по тирам          | tier, volume                 |
| msa_pool_data_read                    | Прочитано данных                | serial, pool                 |
| msa_pool_data_written                 | Записано данных                 | serial, pool                 |
| msa_pool_avg_resp_time                | Время отклика I/O               | serial, pool                 |
| msa_pool_avg_resp_time_read           | Время отклика чтения            | serial, pool                 |
| msa_pool_total_size                   | Полный размер                   | serial, pool                 |
| msa_pool_available_size               | Доступный размер                | serial, pool                 |
| msa_pool_snapshot_size                | Размер снапшотов                | serial, pool                 |
| msa_pool_allocated_pages              | Выделенные страницы             | serial, pool                 |
| msa_pool_available_pages              | Доступные страницы              | serial, pool                 |
| msa_pool_metadata_volume_size         | Размер метаданных               | serial, pool                 |
| msa_pool_total_rfc_size               | Полный размер RFC               | serial, pool                 |
| msa_pool_available_rfc_size           | Доступный размер RFC            | serial, pool                 |
| msa_pool_reserved_size                | Зарезервированный размер        | serial, pool                 |
| msa_pool_unallocated_reserved_size    | Невыделенный резерв             | serial, pool                 |
| msa_tier_reads                        | Операции чтения                 | serial, pool, tier           |
| msa_tier_writes                       | Операции записи                 | serial, pool, tier           |
| msa_tier_data_read                    | Прочитано данных                | serial, pool, tier           |
| msa_tier_data_written                 | Записано данных                 | serial, pool, tier           |
| msa_tier_avg_resp_time                | Время отклика I/O               | serial, pool, tier           |
| msa_tier_avg_resp_time_read           | Время отклика чтения            | serial, pool, tier           |
| msa_tier_avg_resp_time_write          | Время отклика записи            | serial, pool, tier           |
| msa_enclosure_power                   | Потребление энергии в ваттах    | wwn, id                      |
| msa_controller_cpu                    | Загрузка CPU                    | controller                   |
| msa_controller_iops                   | IOPS                            | controller                   |
| msa_controller_bps                    | Байт в секунду                  | controller                   |
| msa_controller_read_hits              | Попадания в кеш чтения          | controller                   |
| msa_controller_read_misses            | Промахи кеша чтения             | controller                   |
| msa_controller_write_hits             | Попадания в кеш записи          | controller                   |
| msa_controller_write_misses           | Промахи кеша записи             | controller                   |
| msa_psu_health                        | Состояние блока питания         | psu, serial                  |
| msa_psu_status                        | Статус блока питания            | psu, serial                  |
| msa_system_health                     | Состояние системы               |                              |

## Совместимое оборудование

Экспортер протестирован на следующем оборудовании:

 - HPE MSA 2060

Вероятно, будет работать на:

 - HP MSA 2050/2052 серии с 3.5" и 2.5" корзинами с внешними JBOD или без них
 - DELL ME4024 с 2.5" корзинами
 - Dothill/Seagate AssuredSan
 - Lenovo DS S2200 / S3200

## Почему Go?

Проект переписан с Python на Go со следующими преимуществами:

- **Лучшая производительность**: Написан на Go с конкурентной обработкой
- **Меньший размер**:
  - Бинарный файл: ~12MB (автономный)
  - Docker образ: ~26MB vs ~192MB (Python версия - **уменьшение на 86%**)
- **Быстрый старт**: Без инициализации интерпретатора
- **Те же метрики**: Все метрики из Python версии сохранены
- **Обратная совместимость**: Поддержка флагов и позиционных аргументов
- **Безопасность**: Docker контейнер запускается от непривилегированного пользователя

### Ключевые улучшения:

1. **Нет зависимостей**: Один статический бинарный файл, не требуется Python runtime или pip пакеты
2. **Гибкие аргументы**: Поддержка флагов `-hostname` и позиционных аргументов для обратной совместимости
3. **Автоматическое управление зависимостями**: Go модули автоматически управляют всеми зависимостями при сборке

## Тестирование

Проект включает полное покрытие тестами:

```bash
# Запустить все тесты
go test -v

# Запустить тесты с покрытием
go test -cover

# Сгенерировать детальный отчет о покрытии
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Покрытие тестами

- **Общее покрытие**: ~70%
- **Helper функции**: 100%
- **MetricStore**: 100%
- **Определения метрик**: 100%
- **MSA API клиент**: 81.5%
- **Функция сбора**: 81.2%

Тесты включают:
- Unit-тесты для всех вспомогательных функций
- Валидацию парсинга XML
- Mock MSA сервер для тестирования API
- Валидацию определений метрик
- Проверку маппинга меток
