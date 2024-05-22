# Choppy!

Chop 3D models in half with a user-defined slice plane.

![Screenshot](http://i.imgur.com/DzixHKO.png)

### Prerequisites

First, [install Go](https://golang.org/dl/), set your `GOPATH`, and make sure `$GOPATH/bin` is on your `PATH`.

```bash
brew install go # if using homebrew

# put these in .bash_profile or .zshrc
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

You may need to [install prerequisites](https://github.com/go-gl/glfw#installation) for the `glfw` library.

### Installation

```
go install github.com/fogleman/choppy/cmd/choppy@latest
```

### Usage

```bash
choppy model.stl
```

### Controls

- Mouse: Arcball controls for the entire scene.
- Cmd + Mouse: Orient the model.
- Alt + Mouse: Orient the plane.
- Cmd + Shift + Mouse: Pan the model.
- Alt + Shift + Mouse: Pan the plane.
- Space: Chop! Writes two STL files to disk.

