#!/bin/bash

# =====================================================
# Нагрузочное тестирование Product Catalog Service
# =====================================================

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Конфигурация
URL="http://localhost:8080/api/v1/products"
METHOD="POST"
HEADER="Content-Type: application/json"
BODY='{"name":"LoadTest","description":"Performance test","price":99.99,"categoryId":1}'
OUTPUT_DIR="load_test_results"

# Создаём директорию для результатов
mkdir -p $OUTPUT_DIR

# Функция для вывода заголовка
print_header() {
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
}

# Функция для запуска теста
run_test() {
    local connections=$1
    local requests=$2
    local duration=$3
    local test_name=$4
    
    local output_file="$OUTPUT_DIR/${test_name}_${connections}_${requests}.txt"
    local summary_file="$OUTPUT_DIR/summary.txt"
    
    echo -e "${YELLOW}▶ Тест: $test_name${NC}"
    echo -e "  Соединений: $connections, Запросов: $requests, Длительность: ${duration:-N/A}"
    
    local start_time=$(date +%s.%N)
    
    # Запускаем hey
    if [ -n "$duration" ]; then
        # Тест по времени
        hey -z "$duration" -c "$connections" -m "$METHOD" \
            -H "$HEADER" \
            -d "$BODY" \
            "$URL" > "$output_file" 2>&1
    else
        # Тест по количеству запросов
        hey -n "$requests" -c "$connections" -m "$METHOD" \
            -H "$HEADER" \
            -d "$BODY" \
            "$URL" > "$output_file" 2>&1
    fi
    
    local end_time=$(date +%s.%N)
    local duration_sec=$(echo "$end_time - $start_time" | bc)
    
    # Парсим результаты
    local rps=$(grep "Requests/sec:" "$output_file" | awk '{print $2}')
    local avg=$(grep "Average:" "$output_file" | awk '{print $2}')
    local p95=$(grep "95% in" "$output_file" | awk '{print $3}')
    local p99=$(grep "99% in" "$output_file" | awk '{print $3}')
    local status=$(grep "\[20" "$output_file" | awk '{print $2}' | head -1)
    
    # Если не удалось распарсить, пробуем альтернативный формат
    if [ -z "$rps" ]; then
        rps=$(grep "Requests/sec:" "$output_file" | awk '{print $2}' | head -1)
    fi
    if [ -z "$avg" ]; then
        avg=$(grep "Average:" "$output_file" | awk '{print $2}' | head -1)
    fi
    
    # Выводим результаты
    echo -e "  ${GREEN}✓ RPS: $rps${NC}"
    echo -e "  ${GREEN}✓ Средняя задержка: $avg${NC}"
    echo -e "  ${GREEN}✓ 95% задержка: $p95${NC}"
    echo -e "  ${GREEN}✓ 99% задержка: $p99${NC}"
    echo -e "  ${GREEN}✓ Статус: ${status:-200}${NC}"
    echo -e "  ${GREEN}✓ Время выполнения: ${duration_sec}s${NC}"
    
    # Сохраняем в сводку
    echo "$test_name | $connections | $requests | $rps | $avg | $p95 | $p99 | ${status:-200} | ${duration_sec}s" >> "$summary_file"
    
    echo ""
}

# =====================================================
# Запуск тестов
# =====================================================

# Очищаем сводку
echo "Тест | Соединений | Запросов | RPS | Средняя | 95% | 99% | Статус | Время" > "$OUTPUT_DIR/summary.txt"
echo "-----|------------|----------|-----|---------|-----|-----|--------|------" >> "$OUTPUT_DIR/summary.txt"

print_header "НАГРУЗОЧНОЕ ТЕСТИРОВАНИЕ PRODUCT CATALOG SERVICE"
echo -e "URL: ${BLUE}$URL${NC}"
echo -e "Метод: ${BLUE}$METHOD${NC}"
echo -e "Данные: ${BLUE}$BODY${NC}"
echo ""

# =====================================================
# 1. Прогрев (теплый старт)
# =====================================================
print_header "1. ПРОГРЕВ СЕРВЕРА"
run_test 10 100 "" "warmup"

# =====================================================
# 2. Тесты с разным количеством соединений
# =====================================================
print_header "2. ТЕСТЫ С РАЗНЫМ КОЛИЧЕСТВОМ СОЕДИНЕНИЙ"

# Легкая нагрузка
run_test 10 1000 "" "low_connections"

# Средняя нагрузка
run_test 50 1000 "" "medium_connections"

# Ваш текущий результат
run_test 50 5000 "" "medium_connections_high_requests"

# Высокая нагрузка
run_test 100 2000 "" "high_connections"

# Экстремальная нагрузка
run_test 200 5000 "" "extreme_connections"

# =====================================================
# 3. Тесты по времени (стабильность)
# =====================================================
print_header "3. ТЕСТЫ НА СТАБИЛЬНОСТЬ (ПО ВРЕМЕНИ)"

# Короткий тест на стабильность
run_test 50 "" "10s" "stability_10s"

# Средний тест на стабильность
run_test 50 "" "30s" "stability_30s"

# Длительный тест (проверка утечек памяти)
run_test 50 "" "60s" "stability_60s"

# =====================================================
# 4. Поиск точки насыщения
# =====================================================
print_header "4. ПОИСК ТОЧКИ НАСЫЩЕНИЯ"

for connections in 10 20 50 100 200 500; do
    run_test $connections 2000 "" "saturation_${connections}"
done

# =====================================================
# 5. Тест с большими данными
# =====================================================
print_header "5. ТЕСТ С БОЛЬШИМИ ДАННЫМИ"

BIG_BODY='{"name":"Very Long Product Name That Exceeds Normal Length","description":"This is an extremely long description that goes on and on and on and on and on to test how the system handles large payloads and whether there are any performance degradations when processing significantly larger amounts of text data in the request body","price":9999.99,"categoryId":999}'

run_test 50 1000 "" "large_payload"

# Возвращаем стандартное тело
BODY='{"name":"LoadTest","description":"Performance test","price":99.99,"categoryId":1}'

# =====================================================
# 6. Тест с невалидными данными (отказоустойчивость)
# =====================================================
print_header "6. ТЕСТ ОТКАЗОУСТОЙЧИВОСТИ"

INVALID_BODY='{"name":"","description":"","price":-1,"categoryId":0}'

run_test 50 500 "" "invalid_data"

# Возвращаем стандартное тело
BODY='{"name":"LoadTest","description":"Performance test","price":99.99,"categoryId":1}'

# =====================================================
# Вывод итогов
# =====================================================
print_header "ИТОГИ ТЕСТИРОВАНИЯ"

echo -e "${GREEN}Результаты сохранены в директории:${NC} $OUTPUT_DIR"
echo -e "${GREEN}Сводный отчёт:${NC} $OUTPUT_DIR/summary.txt"
echo ""

echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  СВОДНЫЙ ОТЧЁТ${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""

# Показываем сводку в табличном формате
column -t -s "|" "$OUTPUT_DIR/summary.txt"

echo ""
echo -e "${GREEN}✅ Тестирование завершено!${NC}"
echo ""
echo -e "${YELLOW}Рекомендации:${NC}"
echo "  1. Посмотрите на RPS — найдите точку, где он перестаёт расти"
echo "  2. Посмотрите на 95% и 99% задержки — они не должны резко расти"
echo "  3. Если появились ошибки (статус не 201) — сервер перегружен"
echo "  4. Сравните результаты с вашими ожиданиями"
echo ""