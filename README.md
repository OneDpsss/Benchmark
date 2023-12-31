# Benchmark
## 1.О Бенчмарке:
### Вступление:
#### В рамках данного исследования проводится бенчмарк (сравнительное тестирование производительности) нескольких популярных баз данных, реализованных на языке программирования ### Golang. В особом фокусе анализа оказываются DuckDB, PostgreSQL, SQLite, а также библиотеки Gota и Gorm.
### Цели бенчмарка
#### Сравнительная производительность время выполнения запросов
#### Гибкость и удобство использования: Анализ удобства работы с базой данных, возможности поддержки сложных запросов и структуры данных, предоставляемой каждой технологией.
#### Масштабируемость: Исследование способности каждой базы данных масштабироваться при увеличении объема данных и интенсивности запросов.
### Выбор Golang в качестве языка:
#### ![image](https://github.com/OneDpsss/Benchmark/assets/108849165/d2dba0a3-bfce-413c-a1ce-dec2ebff7cfc)
#### Использование Golang в этом исследовании обусловлено его выдающейся производительностью, простотой синтаксиса, обширным сообществом разработчиков, а также богатым набором библиотек и фреймворков для работы с базами данных.
## 2.Запросы
### 1)SELECT "VendorID", COUNT(*) FROM "trips" GROUP BY 1;
#### В этом запросе вы выбираете уникальные значения "VendorID" из таблицы "trips" и подсчитываете количество записей для каждого уникального "VendorID". Затем результат группируется по первому столбцу, что в данном случае соответствует "VendorID".
### 2)SELECT "passenger_count", AVG("total_amount") FROM "trips" GROUP BY 1;
#### Этот запрос выбирает уникальные значения "passenger_count" из таблицы "trips" и вычисляет среднее значение "total_amount" для каждого уникального "passenger_count". Результат группируется по первому столбцу ("passenger_count").
### 3)SELECT "passenger_count", EXTRACT(year FROM "tpep_pickup_datetime"), COUNT(*) FROM "trips" GROUP BY 1, 2;
#### Здесь вы выбираете "passenger_count" и извлекаете год из "tpep_pickup_datetime". Затем подсчитываете количество записей для каждой комбинации "passenger_count" и извлеченного года. Результат группируется по двум столбцам.
### 4)SELECT "passenger_count", EXTRACT(year FROM "tpep_pickup_datetime"), ROUND("trip_distance"), COUNT(*) FROM "trips" GROUP BY 1, 2, 3 ORDER BY 2, 4 DESC;
#### Этот запрос аналогичен третьему, но с добавлением столбца "trip_distance", который округляется. Затем результат снова группируется по трем столбцам и упорядочивается сначала по году (второй столбец), а затем по убыванию количества записей (четвертый столбец).
### Каждый из этих запросов предоставляет агрегированную информацию о данных в таблице "trips" в зависимости от выбранных столбцов и заданных условий группировки.
## 3.Модули,используемые в бенчмарке
### 	SQLite:
### ![19d437c355228dc4113ae](https://github.com/OneDpsss/Benchmark/assets/108849165/9480500a-306e-4825-bd7b-c0fd39368eb1)
#### •	Тип базы данных: Легковесная встроенная база данных.
#### •	Особенности: Хранится в одном файле, прекрасно подходит для малых и средних проектов. Поддерживает транзакции, но имеет ограниченные возможности масштабирования для больших проектов.
###   PostgreSQL:
### ![postgresql-banner](https://github.com/OneDpsss/Benchmark/assets/108849165/61952d03-6799-486d-bea8-bb9731347d6b)
#### •	Тип базы данных: Реляционная база данных с открытым исходным кодом.
#### •	Особенности: Мощная и надежная, поддерживает сложные SQL-запросы, транзакции и расширенные возможности. Используется в различных типах проектов, включая крупные и высоконагруженные системы.
###   DuckDB:
### ![image](https://github.com/OneDpsss/Benchmark/assets/108849165/ef212400-e7ae-45a7-ad98-55ca25d2891c)
#### •	Тип базы данных: Колоночная аналитическая база данных.
#### •	Особенности: Известна своей высокой производительностью и низкими накладными расходами при обработке аналитических запросов. Эффективна при работе с большими объемами данных.
###   Gota:
### ![image](https://github.com/OneDpsss/Benchmark/assets/108849165/d08d093a-abb2-442d-b030-4bbd6af4c327)
#### •	Тип библиотеки: Библиотека для обработки и анализа данных в стиле дата-фреймов.
#### •	Особенности: Предоставляет удобные средства для манипуляции и агрегации данных, вдохновленные функциональными возможностями языка R. Оптимизирована для работы с небольшими наборами данных.
### Gorm:
### ![image](https://github.com/OneDpsss/Benchmark/assets/108849165/4319497f-5538-4a54-b025-cf06b55063d8)
#### •	Тип библиотеки: ORM-библиотека (Object-Relational Mapping) для работы с базами данных в Golang.
#### •	Особенности: Упрощает взаимодействие с реляционными базами данных, предоставляя объектно-ориентированный интерфейс. Поддерживает функции автоматического создания, миграции и запросов к базе данных
## 4.Запуск:
### При запуске программы с использованием конфигурационного файла (cfg), выбор модуля осуществляется путем указания соответствующего параметра в конфигурационном файле. Каждый модуль, такой как SQLite, PostgreSQL, DuckDB, Gota или Gorm, All может быть выбран для использования в зависимости от требований проекта. 
### <img width="840" alt="Снимок экрана 2023-12-17 в 19 57 01" src="https://github.com/OneDpsss/Benchmark/assets/108849165/e2aa393e-c6a5-4522-8995-36a4cf5c3ac3">
## 5. О результате
### О процессе:
#### В процессе исследования были проведены бенчмарки для каждого из вышеупомянутых модулей баз данных – SQLite, PostgreSQL, DuckDB, Gota, Gorm. Медианные значения, полученные после 20 запусков каждого запроса, служат надежным показателем производительности и отражают эффективность каждого модуля в обработке данных.
### Графики  и результаты замеров
#### Время выполнения запросов на маленьком объеме данных:
##### <img width="554" alt="Снимок экрана 2023-12-17 в 20 03 39" src="https://github.com/OneDpsss/Benchmark/assets/108849165/19b9ec97-ae7a-4c91-a6bb-be9c984d52a8">
#### График:
##### <img width="795" alt="Снимок экрана 2023-12-17 в 20 05 34" src="https://github.com/OneDpsss/Benchmark/assets/108849165/863e2196-f4ba-4c32-bb9b-b22b039b0f3b">
#### Время выполнения запросов на большом объеме данных:
##### <img width="551" alt="Снимок экрана 2023-12-17 в 20 06 16" src="https://github.com/OneDpsss/Benchmark/assets/108849165/2cc96501-706a-4c5c-8b85-8b8f2afd283d">
#### График:
##### <img width="837" alt="Снимок экрана 2023-12-17 в 20 07 18" src="https://github.com/OneDpsss/Benchmark/assets/108849165/a2f2c4eb-74de-4fa1-9e60-01d24416fbcd">
## 6. Вывод и личное мнение:
### После проведения повторяющихся тестирований баз данных, можно выделить некоторые наблюдения:
#### Общие выводы:
##### Исходя из предоставленных данных (время выполнения запросов в миллисекундах), можно сделать следующие общие выводы:
###### DuckDB: Проявил выдающуюся производительность, демонстрируя минимальное время выполнения для каждого из представленных запросов. Это подтверждает его репутацию быстрой и эффективной аналитической базы данных.
###### SQLite: Показал умеренные значения времени выполнения, что свидетельствует о его пригодности для проектов средних размеров. Особенно хорошо сбалансированное время выполнения запросов.
##### Gota: Демонстрировал высокое время выполнения, особенно по сравнению с DuckDB и SQLite. Это может подразумевать, что Gota не всегда является оптимальным выбором для аналитических запросов с большим объемом данных.
##### PostgreSQL и Gorm: Показали средние значения времени выполнения. PostgreSQL, как мощная реляционная база данных, показала консистентность в производительности. Gorm, как ORM для работы с PostgreSQL, также предоставил приемлемые результаты.
#### Личное мнение:
##### DuckDB произвела впечатление своей выдающейся производительностью. 
##### Работа с Gorm была интересной, и он предоставил удобные инструменты.
##### Работать с Gota оказалось несколько труднее из-за сложности с написанием запросов.
### В ходе тестирования каждая база данных продемонстрировала свои характерные особенности, и выбор между ними будет зависеть от конкретных требований и контекста проекта.
