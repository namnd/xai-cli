{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gopls
    cobra-cli
    sqlite
  ];

  DB = "/home/namnguyen/.xai/chat_history.db";
}

