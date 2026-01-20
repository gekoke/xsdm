{

  boot.loader = {
    systemd-boot.enable = true;
    efi.canTouchEfiVariables = true;
  };

  services.displayManager.xsdm.enable = true;

  security.sudo.wheelNeedsPassword = false;

  users.users.alice = {
    isNormalUser = true;
    extraGroups = [ "wheel" ];
    initialPassword = "pass";
  };

  system.stateVersion = "24.05";
}
