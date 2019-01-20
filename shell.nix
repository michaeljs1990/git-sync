let
  pkgs = import <nixpkgs> { };
in

pkgs.stdenv.mkDerivation {

  name = "go-env";

  src = ./.;

  buildInputs = [ 
    pkgs.go
  ];

  shellHook = ''
  export TEST=thing
  '';

}
