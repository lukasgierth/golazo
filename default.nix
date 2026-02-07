{pkgs ? import <nixpkgs> {}, ...}:
pkgs.buildGoModule {
  pname = "golazo";
  version = "0.21.0";
  vendorHash = "sha256-hPrWqDmsCjAnstKIV8W5tqCR4i6uRpnFIZWMr4OKEUo=";

  subPackages = ["."];

  src = builtins.path {
    path = ./.;
    name = "source";
  };
}
