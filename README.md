# DeepL(X) Translator

A translator library compatible with both the [DeepL](https://developers.deepl.com/docs) and [DeepLX](https://deeplx.owo.network/) APIs.

## Installation

Using the Go command, from inside your project:

```shell
go get -u github.com/xjasonlyu/deeplx-translator
```

## Usage

Import the package and create a `Translator`.

```go
package main

import (
    "fmt"
    "log"
    
    deepl "github.com/xjasonlyu/deeplx-translator"
)

func main() {
    authKey := "f63c02c5-f056-..."  // Replace with your key

    translator, err := deepl.NewTranslator(authKey)
    if err != nil {
        log.Fatal(err)
    }

    translations, err := translator.TranslateText([]string{"Hello, world!"}, "FR")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(translations[0].Text)  // "Bonjour, le monde !"
}
```

## Credits

- [cluttrdev/deepl-go](https://github.com/cluttrdev/deepl-go)
- [OwO-Network/DeepLX](https://github.com/OwO-Network/DeepLX)

## License

This project is open-sourced under the MIT license. See the [LICENSE](LICENSE) file for more details.
