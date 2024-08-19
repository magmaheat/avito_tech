[![build](https://github.com/magmaheat/avito_tech/actions/workflows/build.yaml/badge.svg?branch=master)](https://github.com/magmaheat/avito_tech/actions/workflows/build.yaml)
[![coverage](https://img.shields.io/codecov/c/github/magmaheat/avito_tech/master?label=coverage)](https://codecov.io/gh/magmaheat/avito_tech)

### Выполнено.
- Сделаны все ручки, включая дополнительные задания вроде подписки и пользовательской регистрации
- 

#### 1. Номера(ID) домов. 
- Очевидно, что они уникальны, но встает вопрос, является ли id дома автоинкрементрируемым значением в базе или же данное значение задает модератор при создании дома? В строке примера поля address присутствет номер дома относительно улицы и в самом запросе на создание дома среди обязательных полей нет id. При данных вводных хочется сделать id автоинкрементируемым, но если опираться именно на сущность house в API конфигурации - id фигурирует как обязательный параметр. Так как в первом случае мы никак не сможем проконтролировать наличие номера дома в адресе, оставляю id как заполняемое модератором поле.

#### 2. Дополнительные поля.
- Была идея добавить id пользователя со связью многие к одному с таблицей users. Данное поле помогало бы отследить кто именно создал квартиру, тип пользователя, а так же настроить каскадное удаление при удалении самого пользователя из таблицы users. Так же данный столбец дает более удобное применение санкций или определенных льгот к объявлениям пользователя. Поле было добавлено, необходимое взаимодействие было настроено. При упрощенной авторизации создается шаблонный пользователь. От каскадного удаления решил отказаться, при ошибочном или намеренном удалении дома/пользователя, хочу оставить больше информации. При необходимости модератор явно удалит записи.
- Так же для реализации логики добавлены две таблиы users и subscriptions.

#### 3. Валидация данных.
 - Про валидацию данных не было никаких условий и в данном контексте я решил ограничиться определенной настройкой полей в бд. Основным минусом данного способа является излишняя нагрузка на бд, ведь по факту сама бд будет выступать в роли валидатора данных.

#### 4. Зависимость тестов и создание мусора.
- Функциональные тесты очень сильно загрязняют таблицу и для уменьшения мусора пришлось создать некоторые зависимости(создание пользователя с определенной ролью происходит единожды, сохранение ID созданного дома для последующего к нему обращения и т.д.). Ни наличие мусора в таблице, ни зависимость тестов друг от друга мне не нравится, но на данный момент первоочередно - протестровать функционал, остальное оставляю на update после выполнения всех поставленных задач. Пока среди идей только создать функцию к Storage для удаления данных из необходимых таблиц, достаточно просто и в лоб.
