### Tech blog scrapper (Go)

### Features

---

- Scrapes Stripe, Uber & Netflix tech blogs
- Shows all the articles from these blogs at one place

### Usage

---

- dev (air):

```shell
go mod tidy
make templ-generate
make tailwind-build
make dev
```

- dev (with docker):

```shell
docker compose up --build
```

### Dependencies

---

- [Rod](https://github.com/go-rod/go-rod.github.io.git) - for chrome browser based scrapping
- [templ](https://github.com/a-h/templ) - for html templates
- [htmx](https://github.com/bigskysoftware/htmx) - for dynamic html

### Folder Structure

---

```
└── 📁tech-blog-scrapper
    └── Makefile         #contains all the scripts to run the app
    └── 📁cmd
        └── 📁scrapper
            └── main.go  #main file
    └── 📁internal
        └── 📁handlers  # api handlers
        └── 📁scrapper  # scrapper logic
        └── 📁utils     # utility
        └── 📁view      # html templates (templ)
    └── 📁static
        └── 📁css       # contains custom & tailwindcss built css
        └── 📁script    # htmx & custom js scripts
```

<hr style="height:3px; border:none; background-color:#2e2e2e;" />

> [manishmandal.com](https://manishmandal.com) • <span style="opacity:0.6;">GitHub </span> [@manishMandal02](https://github.com/manishMandal02) • <span style="opacity:0.6;">X</span> [@manishMandalJ](https://twitter.com/manishMandalJ)
