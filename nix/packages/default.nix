_: {
  perSystem =
    { pkgs, ... }:
    {
      packages.default = pkgs.buildGoModule {
        pname = "xsdm";
        version = "0.0.1";

        src = ../../.;

        vendorHash = "sha256-mG6jwfWVCroZab6jrQk6DnhNabzbWG9XeN+NzemCZeQ=";

        buildInputs = [ pkgs.linux-pam ];

        meta = {
          description = "Extra Simple Display Manager";
          homepage = "https://github.com/gekoke/xsdm";
          mainProgram = "xsdm";
        };
      };
    };
}
