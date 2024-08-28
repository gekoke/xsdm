_: {
  perSystem = _: {
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
