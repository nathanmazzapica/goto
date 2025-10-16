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

### 1. Clone the Repo
```bash
git clone https://github.com/nathanmazzapica/goto.git
```

### 2. Install the Go Binary
```bash
go install goto/cmd/goto.go
```
>*This is only temporary while I figure out how GitHub releases work*

### 3. Add the Shell Function
##### 1. Add the following to your `.bashrc` or `.zshrc`:
```bash
function tp() {
    local dir
    case $1 in 
        -d|--delete|-a|--add|-l|--list)
            goto $@
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
| `tp --add <name>` or `-a`     | Adds a marker for the current dir   |
| `tp <name>`                   | Teleports to the saved directory    |
| `tp --list` or `-l`           | Lists all saved markers             |
| `tp --delete <name>` or `-d`  | Deletes the saved marker            |
| `tp --recall` or `-r`         | Returns to the previous directory            |


## Example

```bash
cd ~/github.com/mysuperlongname/my-project-with-a-long-name
tp --add proj

cd ~
tp proj     # instantly jumps back to ~/github.com/mysuperlongname/my-project-with-a-long-name

tp someOtherMarker
tp -r       # instantly returns to the previous directory (in this case home)
```

## Future Plans
- Ability to teleport files

## Contributions Welcome
If you have any suggestions or would like to help a young man build something cool, feel free to open an issue or PR.
