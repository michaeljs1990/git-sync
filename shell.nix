let
  pkgs = import <nixpkgs> { };
  godeps = pkgs.buildGoPackage rec {
    name = "dep2nix";
    goPackagePath = "src/test";
    src = ./src/test;
    goDeps = ./src/test/deps.nix;
  };
in

pkgs.stdenv.mkDerivation {

  name = "go-env";

  src = ./.;

  buildInputs = [ 
    pkgs.go
    pkgs.dep
    pkgs.dep2nix
  ];

  shellHook = ''
  export TEST=thing
  '';

}
