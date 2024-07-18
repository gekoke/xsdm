{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells."${system}".default = pkgs.mkShell 
{
        packages = [
          pkgs.go
        ];
      };

      packages."${system}".default = pkgs.buildGoModule {
        pname = "xsdm";
        version = "0.0.1";

        src = ./.;

        vendorHash = "sha256-mG6jwfWVCroZab6jrQk6DnhNabzbWG9XeN+NzemCZeQ=";

        buildInputs = [ pkgs.linux-pam ];

        meta = {
          description = "Extra Simple Display Manager";
        };
      };
    };
}
