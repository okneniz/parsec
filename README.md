<p align="center">
  <h3 align="center">Parsec</h3>

  <p align="center">
    Golang parser combinator library inspired by <a href="https://hackage.haskell.org/package/parsec">haskell parsec</a>.
    <br/>
    <a href="https://pkg.go.dev/github.com/okneniz/parsec"><strong>Explore the docs »</strong></a>
    <br/>
    <a href="https://github.com/okneniz/parsec/issues">Report Bug</a>.
    <a href="https://github.com/okneniz/parsec/issues">Request Feature</a>
  </p>
</p>

![Downloads](https://img.shields.io/github/downloads/okneniz/parsec/total) ![Contributors](https://img.shields.io/github/contributors/okneniz/parsec?color=dark-green) ![Forks](https://img.shields.io/github/forks/okneniz/parsec?style=social) ![Stargazers](https://img.shields.io/github/stars/okneniz/parsec?style=social) ![Issues](https://img.shields.io/github/issues/okneniz/parsec) ![License](https://img.shields.io/github/license/okneniz/parsec) 

## Table Of Contents

* [About the Project](#about-the-project)
* [Getting Started](#getting-started)
  * [Installation](#installation)
* [Usage](#usage)
* [Roadmap](#roadmap)
* [Contributing](#contributing)
* [License](#license)
* [Authors](#authors)
* [Acknowledgements](#acknowledgements)

## About The Project

Golang parser combinator library (inspired by haskell parsec).

> But what is parser combinator?

> In the parse combinatorial framework, a "parser" is a function that takes some semistructured input and produces some structured output, and "combinator" is a function that allows combining / composing things. So "parser combinators" is a way of expressing a system where you write a lot of small parsing functions and compose then together.


## Getting Started


### Installation

```bash
go get github.com/okneniz/parsec
```

## Examples


- text
  - [json](https://github.com/okneniz/parsec/tree/master/examples/strings/json)
  - [timestamps](https://github.com/okneniz/parsec/tree/master/examples/strings/timestamps)
  - [credit cards](https://github.com/okneniz/parsec/tree/master/examples/strings/cards)
- binary
  - [message pack](https://github.com/okneniz/parsec/tree/master/examples/bytes/message_pack)
  - [png](https://github.com/okneniz/parsec/tree/master/examples/bytes/png)

## Roadmap

See the [open issues](https://github.com/okneniz/parsec/issues) for a list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.
* If you have suggestions for adding or removing projects, feel free to [open an issue](https://github.com/okneniz/parsec/issues/new) to discuss it, or directly create a pull request after you edit the *README.md* file with necessary changes.
* Please make sure you check your spelling and grammar.
* Create individual PR for each suggestion.
* Please also read through the [Code Of Conduct](https://github.com/okneniz/parsec/blob/main/CODE_OF_CONDUCT.md) before posting your first idea as well.

### Creating A Pull Request

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Copyright (C) 2023 Andrey Zinenko

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
