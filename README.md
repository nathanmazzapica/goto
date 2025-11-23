# goto: A Terminal Teleporter
A lightweight command-line tool for setting quick navigation markers in your filesystem: rewritten in Go.

Jump instantly between your favorite directories using a simple command:

```bash
goto <marker_name>
# or if you're as lazy as me...
tp <marker_name>
```

## Features
- Save and name directory markers
- Instantly jump to saved locations

## Demo
WIP

## Installation

### MacOS (Apple Silicon)

```bash
curl -L https://github.com/nathanmazzapica/goto/releases/latest/download/goto-darwin-arm64 \
  -o /usr/local/bin/goto
chmod +x /usr/local/bin/goto
```

### Linux

```bash
curl -L https://github.com/nathanmazzapica/goto/releases/latest/download/goto-linux-amd64 \
  -o /usr/local/bin/goto
chmod +x /usr/local/bin/goto
```


### Adding the Shell Function

**Why is this required?**

Shells do not allow external programs to change the working directory of your terminal.

So while `goto` can retrieve the directory you want, it cannot run `cd` in your shell.

To solve this a small shell function is required that:

- Calls `goto`
- Captures the output
- and then performs `cd` in your running shell

```bash
function tp() {
    local dir
    case $1 in 
        -d|--delete|-a|--add|-l|--list)
            goto $@
            ;;
        -p|--print)
            dir=$(goto $2)
            echo $dir
            ;;
        -r|--recall|*)
            dir=$(goto $1)
            if [ -d $dir ]; then
                cd $dir
            else
                echo $dir
            fi
            ;;
    esac
}
```

>*You can name the function anything you'd like to customize your experience*

##### 2. Source your rc file
```bash
source ~/.zsrch
# or
source ~/.bashrc
```

## Usage

| Command                        | Description                         |
|-------------------------------|-------------------------------------|
| `tp --add <name>` or `-a`     | Adds a marker <name> for the current dir   |
| `tp <name>`                   | Teleports to the saved directory pointed to by marker <name>    |
| `tp --list` or `-l`           | Lists all saved markers             |
| `tp --delete <name>` or `-d`  | Deletes the saved marker            |
| `tp --recall` or `-r`         | Returns to the previous directory            |
| `tp --print <name>` or `-p`          | Outputs the directory pointed to by <name>|


## Examples

Below are some example usages:

### Teleporting & Recalling
```bash
cd ~/github.com/mysuperlongname/my-project-with-a-long-name
tp --add proj

cd ~
tp proj     # instantly jumps back to ~/github.com/mysuperlongname/my-project-with-a-long-name

tp someOtherMarker
tp -r       # instantly returns to the previous directory (in this case home)
```

### Moving a file
```bash
mv file_name $(tp -p proj)
```


## Future Plans
As of right now this tool does everything I need it to do, but as I run in to more problems I think fit in this problem domain, I'll continue to update it.

## Contributions Welcome
If you have any suggestions or would like to help a young man build something cool, feel free to open an issue or PR.

This repository runs a simple `go fmt ./...` and `staticcheck ./...` for every PR as part of a CI workflow. Please run these commands prior to your pull-request to ensure no tests are failed.
