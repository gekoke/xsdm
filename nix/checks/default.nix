{ self, ... }:
{
  perSystem =
    { pkgs, ... }:
    {
      checks = import ./vm-tests { inherit self pkgs; };

      pre-commit = {
        check.enable = false;
        settings = {
          hooks = {
            gofmt.enable = true;
            golangci-lint.enable = true;
          };
        };
      };
    };
}
