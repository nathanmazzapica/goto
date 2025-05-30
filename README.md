# goto: A Terminal Teleporter
# goto: A Terminal Teleporter
A lightweight command-line tool for setting quick navigation markers in your filesystem: rewritten in Go.

Jump instantly between your favorite directories using a simple command:

```bash
goto <marker_name>
# or if you're as lazy as me...
tp <marker_name>
```

## Features
Right now, nothing! This tool is a WIP rewrite of [Terminal Teleporter](github.com/nathanmazzapica/terminal-teleporter), originally written in Python.

I am choosing to rewrite this project in Go because I prefer having an executable binary over a Python script, and because I like writing Go.

## Demo

WIP

## Installation

WIP

## Usage

WIP

## Example

```bash
cd ~/github.com/mysuperlongname/my-project-with-a-long-name
tp --add proj

cd ~
tp proj     # instantly jumps back to ~/github.com/mysuperlongname/my-project-with-a-long-name
```

## Future Plans
- Fuzzy Matching

## Contributions Welcome
If you have any suggestions or would like to help a young man build something cool, feel free to open an issue or PR.
