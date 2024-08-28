{ inputs, self, ... }:
{
  flake = {
    overlays.default =
      final: _prev:
      let
        system = final.stdenv.hostPlatform.system;
      in
      {
        xsdm = self.packages.${system}.xsdm;
      };
  };

  perSystem =
    { pkgs, system, ... }:
    {
      _module.args.pkgs = import inputs.nixpkgs {
        inherit system;
        overlays = [ self.overlays.default ];
      };
    };
}
