# Go Storage Client

This is a simple cloud storage client made with Go.

## Prerequisites

You need to have Go installed on your machine. You can download it
[here](https://golang.org/dl/).

## Installation

```bash
go get github.com/dptsi/go-storage
```

## Usage

### Using Base Go by DPTSI

```go
// Instantiate storage client
config := storageapi.Config{
    ClientID:        os.Getenv("OIDC_CLIENT_ID"),
    ClientSecret:    os.Getenv("OIDC_CLIENT_SECRET"),
    OidcProviderURL: os.Getenv("OIDC_PROVIDER"),
    StorageApiURL:   os.Getenv("STORAGE_API_URL"),
}
storageApi, err := storageapi.NewStorageApi(ctx, config)

// Example upload file
r.POST("/ping", func(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        // Handle error
    }
    uploadResponse, err := storageApi.Upload(c, file)
    if err != nil {
        // Handle error
    }
    // Handle success
})

// Example get file
r.GET("/ping/:id", func(c *gin.Context) {
    id := c.Param("id")
    getFileByIdResponse, err := storageApi.Get(c, id)
    if err != nil {
        // Handle error
    }
    // Handle success
})

// Example delete file
r.DELETE("/ping/:id", func(c *gin.Context) {
    id := c.Param("id")
    deleteFileByIdResponse, err := storageApi.Delete(c, id)
    if err != nil {
        // Handle error
    }
    // Handle success
})
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

## License

[GNU GPLv3](https://choosealicense.com/licenses/gpl-3.0/)
