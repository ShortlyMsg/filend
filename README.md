# Filend Project

## Installation
1. Clone the Repository:
    ```sh
    git clone https://github.com/ShortlyMsg/filend.git
    ```
2. Navigate into the project directory:
    ```sh
    cd filend
    ```
3. Start the server:
    ```sh
    go run main.go
    ```
## For Ui - Prerequisites
1. Navigate into the Ui:
    ```sh
   cd filend-ui
    ```
2. Download Node Modules:
    ```sh
    npm install
    ```
3. Start the server:
    ```sh
    npm run serve
    ```

## Endpoints
- ### **Ui**  
- [http://localhost:9091/ui/](http://localhost:9091/ui/) 
- ### **Api / {Post}** 
- http://localhost:9091/checkFileHash
- http://localhost:9091/upload

- http://localhost:9091/download/:{otp}?fileHash=${fileHash}
    - Example: http://localhost:9091/download/abc123?fileHash=hashcode

- http://localhost:9091/getAllFiles/:otp
    - Example: http://localhost:9091/getAllFiles/123456