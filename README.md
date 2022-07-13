
# MDS Challenge
One of the main challenges of building an ecommerce platform is to keep an accurate
list of products and their stocks up to date.
Based on that, we want to build a system that allows us to manage products for an
hypothetical ecommerce platform.

For this system a product should have an unique SKU and could be commercialized in
multiple countries. Each product can then have different stock per country.

Design and build a system that manages products and their stock with the following
requirements:

1. Provide a products API
    - Get a product by SKU
    - b. Consume stock from a product.
    - Should validate if the stock requested is available first, and then decrease it.

2. Provide an API that allows a bulk update of products from a CSV.
    - For each CSV line, the stock update could be positive or negative
    - If a product doesn’t exist, it should be created.

## High Level System Design
![MD5 Challenge System Design](https://res.cloudinary.com/layitheinfotechguru/image/upload/v1657747076/MDS_system_design_xmsq4u.jpg "System Design Image")

The implementation of the task uses hexagonal archictectural style whereby the business logic is centred in the service layers and every other external dependencies are completely isolated from the business layer.

### Endpoints
Baseurl: localhost:3100

**GET products by sku**
- baseurl/api/v1/products/sku/:sku

**Sample response**

```json
{
  "message": "Successfully found product",
  "success": true,
  "payload": [
    {
      "country": "ke",
      "sku": "9befa247cd11",
      "name": "Chung PLC Table",
      "stock_change": -3059,
      "createdAt": "2022-07-12T15:04:02.062Z",
      "updatedAt": "2022-07-13T19:26:19.67Z"
    },
    {
      "country": "ng",
      "sku": "9befa247cd11",
      "name": "Chung PLC Table",
      "stock_change": -3022,
      "createdAt": "2022-07-12T15:04:02.062Z",
      "updatedAt": "2022-07-13T19:26:19.67Z"
    },
]
}
```
**Consume product stock**
- baseurl/api/v1/products/stocks?amount=-10&sku=9befa247cd11

**Sample Response**
```json
{
  "message": "Successfully consumed stock",
  "success": true
}
```
**Bulk update**
- baseurl/api/v1/products/csv/uploads?csv=sample

where the query params csv is assumed that the file exists on the server or uploaded

**Sample Response**

```json
{
  "message": "Acknowledged"
}
```
### HTTP Status Codes Used
- 200 OK
- 200 Accepted
- 302 Found
- 400 Bad Request
- 500 Internal Server Error

## To run the application 
1. cd cmd/api
2. go build
3. ./api 

## To run all test with coverage 
- go test -cover ./...

The service layer has a test coverage of 91.3%




