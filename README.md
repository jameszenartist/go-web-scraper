# Web Scraper Using Go & Rod

Greetings Everyone,

Welcome to my modest web scraper using Go and the infamous Go web scraping package: [Rod](https://github.com/go-rod/rod).

## Table of Contents

- [About](#about)
- [Features](#features)
- [Usage](#usage)
- [Support](#support)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## About

In this project, the goals I wanted to accomplish were to:

- Deepen my knowledge of Go and learn more about conquering combative challenges like CAPTCHAS and emulating user behavior while extracting dynamic content.
- Build flexible logic within the program to handle multiple sites.

Currently this scraper takes in various coins as command line arguments (using the Go [flag](https://pkg.go.dev/flag) package), then moves through the [coindesk.com](https://www.coindesk.com/) prices page to find data on the coins requested.
When finished, the scraper logs the coin data formatted into objects.

While I understand that there is an API for the site, I wanted to take it as a personal web scraping challenge by utilizing various Go packages to handle the problems of asynchronous Javascript loading on the front end.

## Features

While learning about the various paths I could take when building a web scraper in Go, I also had to think about the type of websites I wanted to scrape as well as the type of content I was scraping.
<br/>

Initially I thought about using the [Colly](https://go-colly.org/) framework, but then I quickly realized that it wasn't meant for my use cases, as I needed such things as the creation of a headless browser and user interactivity capabilities.
I soon then discovered the Rod library would work in relation to my goals.
<br/>

## Usage

Please feel free to clone the project yourself!

Of course if you want to scrape another source, you'll have to make the appropriate modifications.

I'd also like to add that if no keywords are added in the command line, default coin arguments (btc, eth, and xrp) will be used as the placeholders to search for.

To get started these basic commands should suffice:

```
// after cloning project:

cd web-scraper

// To run without creating the executable:

go run main.go

// To create the executable:

go build main.go

// run file with default entries (no keywords):

./main

// to run file with added coins to search for (keywords):

./main -keywords <comma separated coin names or codes as string>

```

## Support

Please [open an issue](https://github.com/jameszenartist/go-web-scraper/issues) for support.

## Contributing

Create a branch, add commits, and [open a pull request](https://github.com/jameszenartist/go-web-scraper/pulls).

## License

This project is licensed under the [MIT License](https://github.com/git/git-scm.com/blob/main/MIT-LICENSE.txt)

## Contact

Please feel free to contact me at
jameszenartist@gmail.com, <a href="https://syntaxsamurai.com/" target="blank"><img align="center" src="https://img.shields.io/badge/website-000000?style=for-the-badge&logo=About.me&logoColor=white" alt="syntaxsamurai" /></a> , or <a href="https://www.linkedin.com/in/jameshansen1981/" target="blank"><img align="center" src="https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white" alt="jameshansen1981" /></a>
