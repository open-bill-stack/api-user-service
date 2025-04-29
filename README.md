## Міграції

### Для створення нової міграції

```bash
migrate create -ext sql -dir ./migrations -seq name_migrations
```

Де -seq `name_migrations` - назва вашої міграції.
-dir `./migrations` - директорія, в якій зберігаються ваші міграції.

Після створення міграції, ви можете редагувати SQL-файл, щоб додати необхідні зміни до бази даних.

### Для виконання міграцій

В консолі виконайте команду:

```bash
set -a; source .env; set +a;
```

Ця дія завантажить змінні середовища з файлу `.env` у вашу консоль, для наступного етапу

Після чого, виконуємо самі міграцію

```bash
migrate -path ./migrations -database "$DATABASE_URL?sslmode=disable" up
```

### Для rollback міграцій

```bash
migrate -path ./migrations -database "$DATABASE_URL?sslmode=disable" down 1
```

Це опустить міграцію на один крок назад. Ви можете вказати кількість кроків, на які потрібно повернутися, замість `1`
вказавши число.
