# Odoo one click

Odoo-one-click, help you to setup your odoo instances with just 1 command.

## Roadmap
- [x] Command init to prepare dependencies and requirement to run odoo
- [x] Setup and configure postgresql to run odoo without using root
- [x] Add some flags to customize your odoo installation
- [ ] Add new command `run` to help run existing odoo instances
- [ ] Make it available to other OSes (Mac, and other linux distributions)

## Installation
1. Download compressed binary file on [Releases](https://github.com/rockavoldy/odoo-one-click/releases/latest)
2. Extract using unzip
    ```sh
    unzip odoo-one-click_amd64.zip
    ```
3. Move the binary file to your PATH
    ```sh
    sudo mv odoo-one-click /usr/local/bin
    ```
4. Make it executable
    ```sh
    sudo chmod +x /usr/local/bin/odoo-one-click
    ```
5. Done!, you can now use `odoo-one-click` command

## Usage
- If it is your first time using `odoo-one-click`, you need to run `init` command to setup your system
    ```sh
    odoo-one-click init
    ```
- After it's done, you can now run command `install` to install your desired odoo version
    ```sh
    odoo-one-click install
    ```
    > NOTE: By default, this command will install odoo 15 community with python3.8.13
-  You can also customize your installation by using flags
    ```sh
    odoo-one-click install --odoo-version 13.0 --enterprise --python 3.7.4
    ```
    Command above will install odoo 13 enterprise with python3.7.4
- For other available flags, you can run `odoo-one-click install --help`
    ```sh
    odoo-one-click install --help
    ```
    ```sh
    Install and configure odoo with demo data

    Usage:
    odoo-one-click install [flags]

    Flags:
    -d, --db-name string          Database name to create or use
    -e, --enterprise              Install odoo enterprise
    -h, --help                    help for install
    -o, --odoo-version string     Odoo version to install
    -p, --python-version string   Python version to use

    Global Flags:
    -V, --verbose   Print logs to stdout
    ```