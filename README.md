# Odoo one click

Odoo-one-click, help you to setup your odoo instances with just 1 command.

## Quick Start
1. Make sure you have already installed the dependencies
    ```sh
    sudo apt update && sudo apt install -y curl jq unzip
    ```
2. Run this installer
    ```sh
    curl https://raw.githubusercontent.com/rockavoldy/odoo-one-click/main/bin/ooc-installer | bash
    ```
3. Wait for the installer to finish, and the command will be available on your terminal. Try to run with
    ```sh
    odoo-one-click --help
    ```
4. Now you can continue initialize by using command
    ```sh
    odoo-one-click init
    ```
5. And when it's finished, you can continue install Odoo 15 by using
    ```sh
    odoo-one-click install
    ```
6. Done!, follow the first initialization after your Odoo 15 successfully installed

## Manual Installation
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
5. Restart your terminal, and Done!. You can now use `odoo-one-click` command

## Usage
- First thing first, if you're not currently using ubuntu, but you're sure your distro is derivatives of Ubuntu, please add this env to your system
    ```sh
    export SKIP_UBUNTU_CHECK=YES
    ```
    > Need a better way to determine if the OS or distro is supported or not
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
    odoo-one-click install --odoo 13.0 --enterprise --python 3.7.4 odoo13
    ```
    Command above will install odoo 13 enterprise with python3.7.4, and the instance name will be `odoo13`
- For other available flags, you can run `odoo-one-click install --help`
    ```sh
    odoo-one-click install --help
    
    Install and configure odoo

    Usage:
    odoo-one-click install [flags] directory_name

    Flags:
    -d, --db-name string   Database name to create or use
    -e, --enterprise       Install odoo enterprise
    -h, --help             help for install
    -o, --odoo string      Odoo version to install
    -p, --python string    Python version to use

    Global Flags:
    -v, --verbose   Print logs to stdout
    ```
    > NOTE: directory_name here is optional, if you have more than 1 instances with the same odoo version, you need to specify the directory name so it won't conflict with the other instance

## Roadmap
- [x] Command init to prepare dependencies and requirement to run odoo
- [x] Setup and configure postgresql to run odoo without using root
- [x] Add some flags to customize your odoo installation
- [x] Add validation on which OS can run the app
- [x] Add auto-install script
- [ ] Add a way to check and auto update the app
- [ ] Add new command `run` to help run existing odoo instances
- [ ] Make it available to other OSes (Mac, and other linux distributions)