# DeepL(X) Translator

A Golang translation library compatible with both the [DeepL](https://developers.deepl.com/docs) and [DeepLX](https://deeplx.owo.network/) APIs.

## Installation

Using the Go command, from inside your project:

```shell
go get -u github.com/xjasonlyu/deeplx-translator
```

## Usage

Import the package and create a `Translator`.

### Basic

```go
package main

import (
	"fmt"
	"log"
	"os"

	deeplx "github.com/xjasonlyu/deeplx-translator"
)

func main() {
	deeplAPIKey := os.Getenv("DEEPL_API_KEY")

	translator := deeplx.NewTranslator(deeplAPIKey)

	text, err := translator.TranslateText("Hello, world!", "ZH")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(text) // "你好，世界"
}
```

### Advanced

```go
package main

import (
	"fmt"
	"log"
	"os"

	deeplx "github.com/xjasonlyu/deeplx-translator"
)

func main() {
	deeplxAPIKey := os.Getenv("DEEPLX_API_KEY")
	deeplxAPIURL := os.Getenv("DEEPLX_API_URL")

	{ // Use Free/Pro DeepLX API (v1)
		translator := deeplx.NewTranslator(
			deeplxAPIKey,
			deeplx.WithBaseURL(deeplxAPIURL),
			deeplx.WithVersion(deeplx.VersionV1), // <-- Optional
		)

		text, err := translator.TranslateText("Hello, world!", "ZH")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(text) // "你好，世界"
	}

	{ // Use Official DeepLX API (v2)
		translator := deeplx.NewTranslator(
			deeplxAPIKey,
			deeplx.WithBaseURL(deeplxAPIURL+"/v2"), // <-- full URL is required
		)

		text, err := translator.TranslateText(
			[]string{"Hello, world!"}, // <- text can be either string or []string
			"FR",
			deeplx.WithSourceLang("EN-GB"),
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(text) // "Bonjour à tous !"
	}
}
```


## Credits

- [cluttrdev/deepl-go](https://github.com/cluttrdev/deepl-go)
- [OwO-Network/DeepLX](https://github.com/OwO-Network/DeepLX)

## License

This project is open-sourced under the MIT license. See the [LICENSE](LICENSE) file for more details.
