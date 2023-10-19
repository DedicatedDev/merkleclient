# How to use client app.

1. Upload file.

   - create local files, for example, 1.txt and 2.txt
   - confirm server is running.
   - change server url `config/config.go` file.
   - `go run main.go upload -f 1.txt -f 2.txt`
   - Remember file ID for download test.

2. How to download file.
   `go run main.go download -i <fileID>`
   example:
   `go run main.go download -i 4bac27393bdd9777ce02453256c5577cd02275510b2227f473d03f533924f877`
   `
