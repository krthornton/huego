#!/usr/bin/env bash

echo "Building and installing huego binary..."
go install .

echo "Installing huego desktop file..."
cp ./huego.desktop ~/.local/share/applications
update-desktop-database ~/.local/share/applications -v

echo "Installation complete."
