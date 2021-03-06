# xpg

[![Go Report Card](https://goreportcard.com/badge/github.com/PavelVershinin/xpg)](https://goreportcard.com/report/github.com/PavelVershinin/xpg)

Обёртка для PostgreSQL в стиле ActiveRecords

Это не ORM, этот пакет, просто помогает упорядочить структуру и упрощает типовые задачи. 

# Установка

    go get -u github.com/PavelVershinin/xpg/...

# Тестирование

    $ go test -v

# Использование
```
	ctx := context.TODO()

	// Создание конфигурации подключения к БД
	config, err := pgxpool.ParseConfig("postgres://postgres:123456@localhost:5432/test")
	if err != nil {
		log.Fatal(err)
	}

	// Создание нового пула подключений, с именем test, и директорией с миграциями migrations
	if err := xpg.NewConnectionPool(ctx, "test", config, "migrations"); err != nil {
		log.Fatal(err)
	}

	// Отложенное закрытие всех подключений
	defer xpg.Close()

	// Проверка новых файлов миграций, для подключения main
	// Все миграции складываются в одну директорию, для main это будет migrations, с именами #SERIAL#_up.sql для поднятия
	// и #SERIAL#_down.sql для отката
	// #SERIAL# порядковый номер миграции
	// Для работы системы миграций будет создана таблица xpg_migrations её нельзя ни удалять ни изменять
	if err := migrations.Up(ctx, "test", -1); err != nil {
		log.Fatal(err)
	}

	// Создание таблицы в БД из модели
	// В структуре модели, у каждого поля должен быть прописан тег xpg, с описанием этого поля в SQL
	// Например:
	// type Role struct {
	//	 xpg.Model
	//	 Name string `xpg:"name VARCHAR(50) NOT NULL DEFAULT ''"`
	// }
	if err := migrations.Restore(ctx, &test.Role{}); err != nil {
		log.Fatal(err)
	}
	if err := migrations.Restore(ctx, &test.User{}); err != nil {
		log.Fatal(err)
	}

	// Запись в БД (INSERT)
	role := &test.Role{}
	role.Name = "user"
	if err := role.Save(ctx); err != nil {
		log.Fatal(err)
	}
	log.Printf("Запись сохранена в таблице %s, ей присвоен ID %d", role.Table(), role.ID)

	user := &test.User{}
	user.FirstName = "Ivan"
	user.SecondName = "Ivanovich"
	user.LastName = "Ivanov"
	user.Email = "ivan@mail.ru"
	user.Role = role
	if err := user.Save(ctx); err != nil {
		log.Fatal(err)
	}
	log.Printf("Запись сохранена в таблице %s, ей присвоен ID %d", user.Table(), user.ID)

	// Обновление (UPDATE)
	user.Phone = "+7 999 999-99-99"
	if err := user.Save(ctx); err != nil {
		log.Fatal(err)
	}
	log.Printf("Запись c ID %d, обновлена", user.ID)

	// Выборка (SELECT)
	users := []*test.User{}
	query := xpg.New(&test.User{}).
		WhereBetween("id", 1, 15).
		OrderBy("id", "DESC")
	rows, err := query.Select(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for row := range rows.Fetch() {
		users = append(users, row.(*test.User))
	}
	rows.Close()
```

# Автогенерация кода модели
* В файле где будет жить модель, создать структуру модели, с встроеным в неё xpg.Model 
* У каждого поля структуры, желательно, прописать тег xpg, с описанием этого поля в SQL

Пример структуры:

```
type User struct {
	xpg.Model
	FirstName  string `xpg:"first_name VARCHAR(50) NOT NULL DEFAULT ''"`
	SecondName string `xpg:"second_name VARCHAR(50) NOT NULL DEFAULT ''"`
	LastName   string `xpg:"last_name VARCHAR(50) NOT NULL DEFAULT ''"`
	Email      string `xpg:"email VARCHAR(254) NOT NULL DEFAULT ''"`
	Phone      string `xpg:"phone VARCHAR(18) NOT NULL DEFAULT ''"`
	Role       int64  `xpg:"role BIGINT NOT NULL DEFAULT 0"`
	Balance    int64  `xpg:"balance BIGINT NOT NULL DEFAULT 0"`
}
```

Вызвать команду:

    $ xpgen -path=/absolute/path/model.go -connect=main
 
Полученым кодом, дополнить файл с моделью.