# wgetgo

`wgetgo` is a simplified site downloader utility written in Go. It allows you to download HTML pages along with all embedded resources (CSS, JS, images, etc.) and recursively follow links within the same domain, producing a local copy of a website that can be viewed offline.

It is inspired by `wget -m` (mirror mode) but implemented as a lightweight Go application with support for concurrency, timeout handling, and optional `robots.txt` compliance.

---

## Features

- Download HTML pages and all linked resources (CSS, JS, images, videos, audio, etc.)
- Recursive download of pages within the same domain
- Proper handling of relative and absolute links
- Avoid duplicate downloads
- Correct local path generation and link rewriting for offline use
- Parallel downloading with configurable concurrency
- Optional respect for `robots.txt`
- Timeout handling for slow servers

---

## Project Structure

```
cmd/
└── wget/
│   └── main.go           # entry point for the CLI

internal/
├── config/
│   └── config.go        # CLI flags and configuration
├── crawler/
│   └── crawler.go       # site crawling logic
├── fetcher/
│   └── fetch.go         # HTTP fetching with timeout handling
├── files/
│   └── files.go         # local file saving utilities
├── flags/
│   └── flags.go         # command-line flags parser
├── parser/
│   ├── parse.go         # HTML parsing and link rewriting
│   └── utils.go         # helper utilities (e.g., IsHTML)
└── worker/
│   └── worker.go        # concurrency worker logic
```

---

## Installation

Build the binary using `make`:

```bash
make build
```

This will produce the executable in `bin/wgetgo`.

Clean up build artifacts and downloaded sites:

```bash
make clean
```

---

## Usage

Basic usage:

```bash
bin/wgetgo https://go.dev
```

Download with recursion depth 1 (download the page and links within the same domain):

```bash
bin/wgetgo -d 1 https://go.dev
```

Download with concurrency limit of 5 and a 15-second timeout:

```bash
bin/wgetgo -c 5 -t 15s www.scrapethissite.com
```

Respect `robots.txt` rules:

```bash
bin/wgetgo --robots www.scrapethissite.com
```

---

## Flags

| Flag            | Shorthand | Description                                       |
| --------------- | --------- | ------------------------------------------------- |
| `--depth`       | `-d`      | Maximum recursion depth (0 = only the given page) |
| `--concurrency` | `-c`      | Maximum number of parallel downloads              |
| `--timeout`     | `-t`      | Timeout for HTTP requests (e.g., 10s, 2m)         |
| `--robots`      | -         | Respect `robots.txt` rules                        |
| `--help`        | `-h`      | Show help message                                 |

---

## Example Output

After running the utility, all downloaded sites are stored in the `sites/` directory, with a structure like:

```
sites/
├── go.dev/
│   ├── index.html
│   ├── css/
│   │   └── styles.css
│   ├── js/
│   │   └── main.js
│   └── images/
│       └── logo.png
└── www.scrapethissite.com/
    ├── index.html
    ├── scripts/
    └── images/
```

Each HTML page has its links rewritten to point to the local copies of resources.

---
