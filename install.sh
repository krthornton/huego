#!/usr/bin/env bash

# validate we can install first
if [[ ! -n "$GOPATH" ]]; then
    echo "GOPATH not set. Please ensure GOPATH is set before installing."
    exit 1
fi

echo "Building and installing huego binary..."
go install .

echo "Installing huego desktop file..."
content="$(sed -e "s#{{GOPATH}}#${GOPATH}#" ./huego.desktop)"
echo "$content" > ~/.local/share/applications/huego.desktop
update-desktop-database ~/.local/share/applications

echo "Installation complete."
