# basic_aws-s3_operation_Golang

Step 1 : Go to project file and run the command
```bash
go get
```

step 2 : Replace your aws credentials in .env file (your region,access key, secret key, and bucket name).

```go
AWS_REGION = your_region
AWS_ACCESS_KEY_ID = your_access_key_id
AWS_SECRET_ACCESS_KEY = your_secret_access_key
BUCKET_NAME = your_bucket_name
```

Now run the command
```bash
go run main.go
```
App is now running on http://localhost:4000

Now open Postman

  1. To uplaod object give **POST** request to http://localhost:4000/upload and body input **"document"** key and select file that you want to upload and press send.

  2. To get list of object give **GET** request to http://localhost:4000/list and press send.

  3. To download an object give **POST** request to http://localhost:4000/download and body input **"document"** key and enter correct file name that you want to download and press send.

  4. To delete an object give **POST** request to http://localhost:4000/delete and body input **"document"** key and enter correct file name that you want to delete and press send.