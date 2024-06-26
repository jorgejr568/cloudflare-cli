# Cloudflare CLI

This project provides a command-line interface (CLI) for interacting with Cloudflare's API. It is written in Go and uses the Cobra library for CLI commands.

## Features

- Configure settings
- Add DNS records
- Delete DNS records
- List DNS records

## Commands

### Configure Settings

The `config` command allows you to configure settings for the CLI. This includes setting up your Cloudflare API key.

The `config list` command allows you to print a table with all your configurations.

The `config set` command allows you to set configurations by key in pairs.
```sh
cloudflare-cli config set cloudflare_api_key YOUR_API_KEY
```

The `config get` command allows you to get configurations by key.
```sh
cloudflare-cli config get cloudflare_api_key
```

### Add DNS Record

The `add` command allows you to add a new DNS record. The following flags are required:

- `--domain`: The zone to list records for. Eg. example.com
- `--name`: Name of the record
- `--type`: Type of the record
- `--content`: Content of the record

Optional flags include:

- `--ttl`: TTL of the record (1-86400)
- `--proxied`: Proxied status of the record
- `--tags`: Tags of the record
- `--comment`: Comment of the record

### Delete DNS Record

The `delete` command allows you to delete an existing DNS record. You can specify the record by its ID or by its name and type. The following flags are required:

- `--domain`: The zone to list records for. Eg. example.com
- `--id`: The ID of the record to delete
- `--name`: The name of the record to delete
- `--type`: The type of the record to delete


### List DNS Records

The `list` command allows you to list all DNS records for a specific zone. The following flag is required:

- `--domain`: The zone to list records for. Eg. example.com

## Setup

To use this CLI, you need to have Go installed on your machine. After cloning the repository, you can build the project using `go build`.

## Ignored Files

The `.gitignore` file is set to ignore the following:

- `.idea/`: This directory is created by JetBrains IDEs and contains user-specific settings.
- `vendor/`: This directory contains Go dependencies.
- `cloudflare-cli`: This is the compiled binary of the CLI.

### Vendors
- [Cobra](https://github.com/spf13/cobra)
- [Cloudflare Docs](https://developers.cloudflare.com/api/)
- [Go](https://golang.org/)
- [Dig](https://pkg.go.dev/go.uber.org/dig)