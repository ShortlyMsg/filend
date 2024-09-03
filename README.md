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
- [http://localhost:9090/ui/](http://localhost:9090/ui/) 
- ### **Api**
- http://localhost:9090/upload

- http://localhost:9090/download/:otp?fileName=name.jpeg&userSecurityCode=${userSecurityCode}
    - Example: http://localhost:9090/download/123456?fileName=Bike.jpeg&userSecurityCode=c0d3

- http://localhost:9090/getAllFiles/:otp
    - Example: http://localhost:9090/getAllFiles/123456