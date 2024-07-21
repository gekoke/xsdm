{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs =
    inputs:
    inputs.flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      perSystem =
        { self', pkgs, ... }:
        {
          packages.default = pkgs.buildGoModule {
            pname = "xsdm";
            version = "0.0.1";

            src = ./.;

            vendorHash = "sha256-mG6jwfWVCroZab6jrQk6DnhNabzbWG9XeN+NzemCZeQ=";

            buildInputs = [ pkgs.linux-pam ];

            meta = {
              description = "Extra Simple Display Manager";
            };
          };

          devShells.default = pkgs.mkShell {
            inputsFrom = [ self'.packages.default ];
          };
        };
    };
}
