curl -i -X POST "localhost:8000/create-new-product?name=tomato&code=1&category=vegetables"

curl -i -X POST "localhost:8000/create-new-product?name=cucumber&code=2&category=vegetables"

curl -i -X POST "localhost:8000/create-new-product?name=pepper&code=3&category=vegetables"

curl -i -X GET "localhost:8000/get-product-list"

curl -i -X GET "localhost:8000/get-product-info?name=pepper"

curl -i -X DELETE "localhost:8000/delete-product?name=pepper"

curl -i -X POST "localhost:8000/change-product-by-name?old-name=cucumber&name=apple&code=7&category=fruits"