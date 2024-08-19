[![build](https://github.com/magmaheat/avito_tech/actions/workflows/build.yaml/badge.svg?branch=master)](https://github.com/magmaheat/avito_tech/actions/workflows/build.yaml)
[![coverage](https://img.shields.io/codecov/c/github/magmaheat/avito_tech/master?label=coverage)](https://codecov.io/gh/magmaheat/avito_tech)





Coverage: ...

## Проблемы, вопросы, рассуждения.
#### 1. Шире функционал, больше полей в сущности. 
- В начале задания, казалось бы, четко определена структура двух сущностей, но читая предоставленную API конфигурацию мы видим, что у этих сущностей от места к месту меняются обязательные поля. Стараясь грамотно разделить на слои, я пытался изолировать данные сущности от остальной бизнес логики и чтобы не перегружать методы аргументами пришел к минимально необходимому количеству полей в структуре сущностей(address, year, developer), время же генерировалось по факту отправки и в явном виде из структуры в нем не было необходимости. Но чем шире становился функционал, тем больше требовалось данных. 

#### 2. Так же были вопросы к номерам(ID) домов. 
- Очевидно, что они уникальны, но встает вопрос, является ли id дома автоинкрементрируемым значением в базе или же данное значение задает модератор при создании дома? Тут опять протеворечивая информация, в строке примера поля address присутствет номер дома относительно улицы и хочется сделать id автоинкрементируемым, но если опираться на API конфигурацию, то у сущности house - id фигурирует как обязательный параметр. Так как в первом случае мы никак не сможем проконтролировать наличие номера дома в адресе, оставляю id как заполняемое модератором поле.

#### 3. Добавить поле id_user в таблицу flats.
- Была идея добавить id пользователя со связью многие к одному с таблицей users. Данное поле помогало бы отследить кто именно создал квартиру, тип пользователя, а так же настроить каскадное удаление при удалении самого пользователя из таблицы users. Так же данный столбец дает более удобное применение санкций или определенных льгот к объявлениям пользователя. Столбец все же был удален из таблицы flats из-за наличия упрощенного получения токена и потребутся достаточно много времени, чтобы продумать корректное взаимодействие данного столбца с упрощенной авторизацией. По итогу поле было добавлено, необходимое взаимодействие было настроено. При упрощенной авторизации создается шаблонный пользователь. От каскадного удаления решил отказаться, при ошибочном или намеренном удалении дома/пользователя, хочу оставить больше информации. При необходимости модератор явно удалит записи.

#### 4. Зависимость тестов и создание мусора.
- Функциональные тесты очень сильно загрязняют таблицу и для уменьшения мусора пришлось создать некоторые зависимости(создание пользователя с определенной ролью происходит единожды, сохранение ID созданного дома для последующего к нему обращения и т.д.). Ни наличие мусора в таблице, ни зависимость тестов друг от друга мне не нравится, но на данный момент первоочередно - протестровать функционал, остальное оставляю на update после выполнения всех поставленных задач. Пока среди идей только создать функцию к Storage для удаления данных из необходимых таблиц, достаточно просто и в лоб.

#### 5. Много нюансов.
 - На самом деле при написании данного проекта остается множество нюансов, которые приходится часто пропускать, руализуя основную логику и если их устранять, то допиливать этот проект можно еще очень долго. Во всяком случае тестовое оказалось достаточно интересным и потраченного на него времени не жаль.
