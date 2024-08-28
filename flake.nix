{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-parts.url = "github:hercules-ci/flake-parts";
    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
  };

  outputs =
    inputs:
    inputs.flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      imports = [ inputs.pre-commit-hooks.flakeModule ];
      perSystem =
        {
          self',
          config,
          pkgs,
          lib,
          ...
        }:
        {
          packages.default = pkgs.buildGoModule {
            pname = "xsdm";
            version = "0.0.1";

            src = ./.;

            vendorHash = "sha256-mG6jwfWVCroZab6jrQk6DnhNabzbWG9XeN+NzemCZeQ=";

            buildInputs = [ pkgs.linux-pam ];

            meta = {
              description = "Extra Simple Display Manager";
              homepage = "https://github.com/gekoke/xsdm";
              mainProgram = "xsdm";
            };
          };

          devShells.default = pkgs.mkShell {
            inputsFrom = [ self'.packages.default ];
            packages = [ pkgs.golangci-lint ];
            shellHook = ''
              ${config.pre-commit.installationScript}
            '';
          };

          pre-commit = {
            check.enable = false;
            settings = {
              hooks = {
                gofmt = {
                  enable = true;
                };
                golangci-lint = {
                  enable = true;
                };
              };
            };
          };
        };
    };
}
