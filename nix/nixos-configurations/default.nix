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
          "${inputs.nixpkgs}/nixos/modules/virtualisation/qemu-vm.nix"
          self.nixosModules.xsdm
          (import ../vm-config.nix)
        ];
      };
    };
  };
}
