## Online shopping server backend

Проект представляет собой небольшой сервер для онлайн магазина, и на данный момент умеет после запуска с помощью предоставленного API сохранять, изменять и удалять информацию о товаре, также предоставлять json с информацией об одном товаре по имени, либо всех имеющихся товарах.

### Запуск сервера

`docker-compose up --build`

### Api

#### Post запросы

* Создать новый продукт с указанными `name`, `code` и `category`:
`curl -i -X POST "localhost:7171/create-new-product?name=tomato&code=1000&category=vegetables"`

* Изменить продукт по его `name`:
`curl -i -X POST "localhost:7171/change-product-by-name?old-name=tomato&name=apple&code=7&category=fruits"`

#### Get запросы

* Получить список товаров:
`curl -i -X GET "localhost:7171/get-product-list"`:

* Получить информацию о товаре по его `name`:
`curl -i -X GET "localhost:7171/get-product-info?name=apple"`

#### Delete запрос

* Удалить продукт по его `name`:
`curl -i -X DELETE "localhost:7171/delete-product?name=apple"`

