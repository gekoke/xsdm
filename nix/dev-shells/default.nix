_: {
  perSystem =
    {
      self',
      config,
      pkgs,
      ...
    }:
    {
      devShells.default = pkgs.mkShell {
        inputsFrom = [ self'.packages.default ];
        packages = [ pkgs.golangci-lint ];
        shellHook = ''
          ${config.pre-commit.installationScript}
        '';
      };
    };
}

