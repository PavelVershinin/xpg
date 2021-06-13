package migrations

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PavelVershinin/xpg"
)

// Down выполнит SQL запросы из файлов отката миграций
func Down(ctx context.Context, poolName string, to int) error {
	var reTest = regexp.MustCompile(`^[0-9]+_down\.sql$`)
	var migrationPath = xpg.MigrationsPath(poolName)
	var migrations []int
	var objMigration = &migration{}
	objMigration.SetPool(poolName)
	if err := Restore(ctx, objMigration); err != nil {
		return err
	}

	var from int
	if res, err := xpg.New(objMigration).OrderBy("created_at", "DESC").First(ctx); err != nil && err.Error() != "xpg: No records found" {
		return err
	} else if err == nil {
		file := res.(*migration).File
		if regexp.MustCompile(`^[0-9]+_up\.sql$`).MatchString(file) {
			from, _ = strconv.Atoi(strings.Split(file, "_")[0])
		}
	}

	if from <= to {
		return nil
	}
	files, err := ioutil.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		if reTest.MatchString(f.Name()) {
			num, _ := strconv.Atoi(strings.Split(f.Name(), "_")[0])
			if num > from {
				break
			}
			if num <= to {
				continue
			}
			migrations = append(migrations, num)
		}
	}

	sort.Ints(migrations)

	for i := len(migrations) - 1; i >= 0; i-- {
		fileNum := migrations[i]
		b, err := ioutil.ReadFile(migrationPath + string(os.PathSeparator) + strconv.Itoa(fileNum) + "_down.sql")
		if err != nil {
			return err
		}
		sql := strings.TrimSpace(string(b))
		if sql != "" {
			if _, err := xpg.DB(poolName).Exec(ctx, sql); err != nil {
				return err
			}
			log.Printf("%s executed!\n", strconv.Itoa(fileNum)+"_down.sql")
		}
		if err := xpg.New(objMigration).Where("file", "=", strconv.Itoa(fileNum)+"_up.sql").Delete(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Up Выполнит SQL запросы из файлов миграции
func Up(ctx context.Context, poolName string, to int) error {
	var reTest = regexp.MustCompile(`^[0-9]+_up\.sql$`)
	var migrationPath = xpg.MigrationsPath(poolName)
	var migrations []int
	var objMigration = &migration{}
	objMigration.SetPool(poolName)
	if err := Restore(ctx, objMigration); err != nil {
		return err
	}

	if to > -1 {
		if ok, err := xpg.New(objMigration).Where("file", "=", strconv.Itoa(to)+"_up.sql").Exists(ctx); err != nil {
			return err
		} else if ok {
			return nil
		}
	}

	var from int
	if res, err := xpg.New(objMigration).OrderBy("created_at", "DESC").First(ctx); err != nil && err.Error() != "xpg: No records found" {
		return err
	} else if err == nil {
		file := res.(*migration).File
		if reTest.MatchString(file) {
			from, _ = strconv.Atoi(strings.Split(file, "_")[0])
		}
	}

	files, err := ioutil.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		if reTest.MatchString(f.Name()) {
			num, _ := strconv.Atoi(strings.Split(f.Name(), "_")[0])
			if to > -1 && num > to {
				break
			}
			if num <= from {
				continue
			}
			migrations = append(migrations, num)
		}
	}

	sort.Ints(migrations)

	for _, fileNum := range migrations {
		b, err := ioutil.ReadFile(migrationPath + string(os.PathSeparator) + strconv.Itoa(fileNum) + "_up.sql")
		if err != nil {
			return err
		}
		sql := strings.TrimSpace(string(b))
		if sql != "" {
			if _, err := xpg.DB(poolName).Exec(ctx, sql); err != nil {
				return err
			}
			log.Printf("%s executed!\n", strconv.Itoa(fileNum)+"_up.sql")
		}
		row := &migration{}
		row.SetPool(poolName)
		row.File = strconv.Itoa(fileNum) + "_up.sql"
		if err := row.Save(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Restore Сверит структуру с базой данных, создаст таблицу, если её нет и добавит недостающие колонки
func Restore(ctx context.Context, model xpg.Modeler) error {
	var valueOf = reflect.ValueOf(model)
	var tableName = model.Table()
	var columns = []string{
		`"id" BIGSERIAL NOT NULL`,
	}

	if valueOf.Elem().IsValid() {
		for i := 0; i < valueOf.Elem().NumField(); i++ {
			field := valueOf.Elem().Field(i)
			if field.IsValid() && field.CanInterface() {
				if column := valueOf.Elem().Type().Field(i).Tag.Get("xpg"); column != "" {
					columns = append(columns, column)
				}
			}
		}
	} else {
		return errors.New("migrations: Modeler is not valid")
	}

	columns = append(columns, []string{
		`"created_at" TIMESTAMP NOT NULL DEFAULT NOW()`,
		`"updated_at" TIMESTAMP DEFAULT NULL`,
	}...)

	pool := xpg.New(model)
	tables, err := pool.Tables(ctx)
	if err != nil {
		return err
	}
	exists := false
	for _, table := range tables {
		if table == tableName {
			exists = true
			break
		}
	}

	if !exists {
		_, err := xpg.DB(model.PoolName()).Exec(ctx, `CREATE TABLE "`+tableName+`" (`+strings.Join(columns, ", ")+`)`)
		return err
	}

	existsColumns, err := pool.Columns(ctx)
	if err != nil {
		return err
	}
	for _, column := range columns {
		name := strings.Trim(strings.ToLower(strings.Fields(strings.TrimSpace(column))[0]), `"`)
		exists := false
		for _, col := range existsColumns {
			if strings.Trim(strings.ToLower(col.Name), `"`) == name {
				exists = true
				break
			}
		}
		if !exists {
			if _, err := xpg.DB(model.PoolName()).Exec(ctx, `ALTER TABLE "`+tableName+`" ADD COLUMN `+column); err != nil {
				return err
			}
		}
	}

	return nil
}
