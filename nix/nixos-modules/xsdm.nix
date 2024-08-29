self:
{
  config,
  lib,
  pkgs,
  ...
}:
let
  package = self.packages.${pkgs.stdenv.hostPlatform.system}.xsdm;
  cfg = config.services.displayManager.xsdm;
  inherit (lib) mkEnableOption mkOption types;
in
{
  options.services.displayManager.xsdm = {
    enable = mkEnableOption "xsdm display manager";
    package = mkOption {
      type = types.nullOr types.package;
      default = package;
    };
    tty = mkOption {
      type = types.int;
      default = 1;
      description = ''
        The tty which contains xsdm.
      '';
    };
  };

  config =
    let
      tty = "tty${toString (cfg.tty)}";
    in
    lib.mkIf cfg.enable {
      security.pam.services.xsdm = {
        startSession = true;
        allowNullPassword = true;
      };

      services.displayManager.enable = true;

      systemd = {
        defaultUnit = "graphical.target";

        services = {
          xsdm = {
            aliases = [ "display-manager.service" ];

            unitConfig = {
              Wants = [ "systemd-user-sessions.service" ];
              After = [
                "systemd-user-sessions.service"
                "plymouth-quit-wait.service"
                "getty@${tty}.service"
              ];
              Conflicts = [ "getty@${tty}.service" ];
            };

            serviceConfig = {
              ExecStart = "${cfg.package}/bin/xsdm";
              StandardInput = "tty";
              TTYPath = "/dev/${tty}";
              TTYReset = "yes";
              TTYVHangup = "yes";
              Type = "idle";
            };

            restartIfChanged = false;
            wantedBy = [ "graphical.target" ];
          };
        };
      };
    };
}
