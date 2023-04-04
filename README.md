# simple-ecommerce

This is a simple program for simple e-commerce website. It consists of Products, Customers, and Orders. For the database, there are 4 tables: products, customers, orders, and order_products. Table products is used to store all product data including name, price, description, image, quantity. Table customers is used to store all customer data including name, email, password, and role to differentiate customer and admin. Table orders is used to store order data that is created by customer, order data consists of customer_id, total_price, status, and date (order created and order paid). Table order_products is used to store all products that customer order, it consists of order_id, product_id, and quantity.

There are 5 endpoints in this program:

1. Register endpoint (/v1/register) to create customer
2. Login endpoint (/v1/login) to get access token for order endpoint
3. Create order endpoint (/v1/orders) to create order for a list of products by customer
4. Get order by customer endpoint (/v1/orders/customer/:customerId) for customer to get their own orders data
5. Get all orders endpoint (/v1/orders) for admin to check all orders data

<!-- GETTING STARTED -->

## Getting Started

There are 3 programs that can be run, the main server for simple e-commerce RESTful API, script to generat csv for reporting, and cron for order reminder scheduler.
To get a local copy up and running follow these simple example steps.

1. Create the database simple-ecommerce
2. Create table products, customers, orders, and order_products in simple-ecommerce database
3. Clone the repo
   ```sh
   git clone https://github.com/VincentTC/simple-ecommerce.git
   ```
4. Move to the repo
   ```sh
   cd simple-ecommerce
   ```
5. Install all libraries used
   ```sh
   go mod tidy
   ```
6. Run the main server for RESTful API
   ```sh
   go run cmd/*.go
   ```
7. Open new terminal
8. Run the cron for order reminder scheduler in the new terminal
   ```sh
   go run cmd/*.go run-cron
   ```
9. Open new terminal again
10. Run the script for reporting in the new terminal
    ```sh
    go run cmd/*.go run-script-report
    ```

Command run-cron covers:

- A scheduler to send email for pending orders
- It will run every day at midnight
- For this program, it will only log the message and not really sending the email

Command run-script-report covers:

- A script to generate a csv report for all orders
- It will create a csv file with filename "order-report-{date}.csv"
