{ self, ... }:
{
  flake = {
    nixosModules = {
      xsdm = import ./xsdm.nix self;
      default = self.nixosModules.xsdm;
    };
  };
}
