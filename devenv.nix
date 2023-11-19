{ pkgs, ... }:

{
  # https://devenv.sh/basics/
  env = {
    GREET = "devenv";
  };

  # https://devenv.sh/packages/
  packages = [ pkgs.git pkgs.zlib pkgs.glibc pkgs.stdenv.cc.cc.lib ];

  # https://devenv.sh/scripts/
  scripts.hello.exec = "echo hello from $GREET";

  enterShell = ''
    hello
    git --version
    go version
    python --version
  '';

  # https://devenv.sh/languages/
  languages.go.enable = true;
  languages.python = {
    enable = true;
    venv.enable = true;
    venv.requirements = ''
      scikit-learn
      facenet-pytorch
    '';
  };

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # https://devenv.sh/processes/
  # processes.ping.exec = "ping example.com";

  # See full reference at https://devenv.sh/reference/options/
}
