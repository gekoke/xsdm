{ self, pkgs, ... }:
let
  tests = [
    {
      name = "smoke_test_starts_without_error";

      testScript = ''
        vm.wait_for_unit("display-manager.service")
      '';
    }
  ];
in
builtins.listToAttrs (
  map (test: {
    name = test.name;
    value = pkgs.testers.runNixOSTest (test // { nodes.vm = self.nixosConfigurations.testVm; });
  }) tests
)
