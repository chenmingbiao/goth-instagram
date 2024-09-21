# goth-instagram

This package implements the [Instagram Graph API](https://developers.facebook.com/docs/instagram-basic-display-api/) for [Goth](https://github.com/markbates/goth) OAuth library.

## Installation

go get github.com/chenmingbiao/goth-instagram

## Usage

To use this provider with Goth, you need to set it up with your Instagram application credentials. Here's an example of how to configure and use the Instagram provider:

```go
package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/chenmingbiao/goth-instagram"
)

func main() {
    goth.UseProviders(
        instagram.New(os.Getenv("INSTAGRAM_CLIENT_ID"), os.Getenv("INSTAGRAM_CLIENT_SECRET"), "http://localhost:3000/auth/instagram/callback"),
	)
	http.HandleFunc("/auth/instagram", func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	})
    http.HandleFunc("/auth/instagram/callback", func(w http.ResponseWriter, r *http.Request) {
        user, err := gothic.CompleteUserAuth(w, r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "%+v", user)
	})
    http.ListenAndServe(":3000", nil)
}
```

Make sure to set the `INSTAGRAM_CLIENT_ID` and `INSTAGRAM_CLIENT_SECRET` environment variables with your Instagram application credentials.

## Configuration

The provider is initialized with `New()` and takes three arguments:
* clientKey: Your Instagram application's App ID
* secret: Your Instagram application's App Secret
* callbackURL: The URL to redirect to after authentication
You can also pass additional scopes as variadic arguments to `New()`. If no scopes are provided, it defaults to `user_profile`.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.