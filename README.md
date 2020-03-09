curl -i -X POST "localhost:7171/create-new-product?name=tomato&code=1&category=vegetables"

curl -i -X POST "localhost:7171/create-new-product?name=cucumber&code=2&category=vegetables"

curl -i -X POST "localhost:7171/create-new-product?name=pepper&code=3&category=vegetables"

curl -i -X GET "localhost:7171/get-product-list"

curl -i -X GET "localhost:7171/get-product-info?name=pepper"

curl -i -X DELETE "localhost:7171/delete-product?name=pepper"

curl -i -X POST "localhost:7171/change-product-by-name?old-name=cucumber&name=apple&code=7&category=fruits"