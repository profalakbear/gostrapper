# Gostrapper CLI tool

A CLI tool to generate a predefined folder structure and files inside of it based on a provided structure file for a go project(containing go.mod file).

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Example Structure File](#example-structure-file)

## Installation

To install the CLI tool, you'll need to have Go installed. Then, follow these steps:

1. Clone this repository or download the source code.

2. Navigate to the directory containing the source code.

3. Build the tool using the `go build` command:

   ```bash
   go build -o gostrapper main.go

4. Move the generated binary to a directory included in your system's PATH:

    ```bash
    sudo mv gostrapper /usr/local/bin

## Usage

1. To use the CLI tool, you need to provide both the rootpath as a -r and structure as -s flags. Here's the basic usage:

    ```bash 
    gostrapper -p /path/to/your/project -s /path/to/structure.txt

## Example Structure File

1. Create a text file named structure.txt (or any name you prefer) with the following content(you can describe your project structure as you wish) as an example:

    ```bash
    cmd/app/main.go
    database/migrations/
    database/database.go
    profiles/default.env
    internal/config/builder.go
    internal/server/server.go
    internal/server/handler/handler.go
    internal/server/middleware/
    internal/server/router/router.go
    internal/server/dto/
    internal/service/service.go
    internal/utils/
    .gitignore
    .gitlab-ci.yml
    Dockerfile
    Makefile
    README.md
