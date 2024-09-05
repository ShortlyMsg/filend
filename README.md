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

## Endpoints
- ### **Ui**  
- [http://localhost:9091/ui/](http://localhost:9091/ui/) 
- ### **Api**
- http://localhost:9091/upload

- http://localhost:9091/download/:otp?fileName=name.jpeg&userSecurityCode=${userSecurityCode}
    - Example: http://localhost:9091/download/123456?fileName=Bike.jpeg&userSecurityCode=c0d3

- http://localhost:9091/getAllFiles/:otp
    - Example: http://localhost:9091/getAllFiles/123456