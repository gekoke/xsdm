{ inputs, self, ... }:
let
  inherit (inputs.nixpkgs.lib) nixosSystem;
in
{
  flake = {
    nixosConfigurations = {
      testVm = nixosSystem {
        system = "x86_64-linux";
        modules = [
          self.nixosModules.xsdm
          {
            services.displayManager.xsdm.enable = true;

            security.sudo.wheelNeedsPassword = false;

            users.users.alice = {
              isNormalUser = true;
              extraGroups = [ "wheel" ];
              initialPassword = "pass";
            };

            system.stateVersion = "24.05";
          }
        ];
      };
    };
  };
}
