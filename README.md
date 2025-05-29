# Filend Project
### !! Known Issues
- `removeFile` function needs to be implemented on both frontend and backend.

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
## Tools
- Minio (Storage)
- PostgreSQL (Database)
- Firebase Cloud Messaging (FCM) - Real-time progress tracking -

## Notes
- **.env_sample** is a sample of the **.env** file and should be properly filled out.  
- **firebase-config-be sample.json** is a sample of the **firebase-config-be.json** file and should be properly filled out.  
- **filend-ui/src/config/firebaseConfig sample.json** is also a sample file and should be filled out accordingly. **filend-ui\public\firebase-messaging-sw.js**

- **Uploaded chunk sizes can be changed.**  
  **Default:** 2MB – `filend-ui/src/components/FileUpload.vue` – `chunkSize`

- **Files that have expired are automatically deleted from MinIO.**  
  **Default:** 24-hour retention period – `services/file_delete_service.go`

- **A scheduler checks whether any files are due for deletion at your specified interval.**  
  **Default:** 1 hour – `services/scheduler.go`

## Endpoints
- ### **Ui**  
- [http://localhost:9071/](http://localhost:9071/)
    - Upload & Download
- ### **Api / {Post}** 
- http://localhost:9091/generateOtp
- http://localhost:9091/checkFileHash
- http://localhost:9091/upload
- http://localhost:9091/SendUploadProgress
    - Front-End to FCM
- http://localhost:9091/subscribeToken
    -  for FCM

- http://localhost:9091/download/:{otp}?fileHash=${fileHash}
    - Example: http://localhost:9091/download/abc123?fileHash=hashcode

- http://localhost:9091/getAllFiles/:otp
    - Example: http://localhost:9091/getAllFiles/123456