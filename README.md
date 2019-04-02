## Building Footprints

### Requirements
1. Install and run MongoDB
1. Download Building Footprints data (600+Mb):
    ```shell
    wget https://data.cityofnewyork.us/api/views/mtik-6c5q/rows.json
    ```
1. Compile the app:
    ```shell
    make build
    ```
1. Load the data into Mongo:
    ```shell
    ./building-footprints load-data -mongo-url 127.0.0.1:27017 -mongo-db-name building-footprints -data-file rows.json
    ```
1. Start the API server:
    ```shell
    ./building-footprints
    ```
1. Visit http://localhost:3000/api/v0/avg_height?borough_code=3