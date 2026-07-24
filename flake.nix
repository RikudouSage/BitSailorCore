{
  description = "Dev shell for CGO cross-compilation (386, armv7, arm64)";

  inputs = {
    nixpkgs.url = "nixpkgs";
    nixpkgs-legacy.url = "github:NixOS/nixpkgs/5d61a47f31092319ab82b44213420974b91c9bc1";
  };

  outputs = { self, nixpkgs, nixpkgs-legacy }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      forAllSystems = f:
        builtins.listToAttrs (map (system: {
          name = system;
          value = f system;
        }) systems);
    in {
      devShells = forAllSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };
          legacyPkgs = import nixpkgs-legacy { inherit system; };

          cc386   = legacyPkgs.pkgsCross.gnu32.stdenv.cc;
          ccArmv7 = legacyPkgs.pkgsCross.armv7l-hf-multiplatform.stdenv.cc;
          ccArm64 = legacyPkgs.pkgsCross.aarch64-multiplatform.stdenv.cc;
        in {
          default = pkgs.mkShell {
            buildInputs = [
              pkgs.go
              pkgs.golangci-lint
              pkgs.patchelf
              pkgs.revive
              pkgs.gnumake

              cc386
              ccArmv7
              ccArm64
            ];

            shellHook = ''
              unset GOROOT

              export CC_386=${cc386.targetPrefix}cc
              export CC_ARMV7=${ccArmv7.targetPrefix}cc
              export CC_ARM64=${ccArm64.targetPrefix}cc

              echo "Using cross compilers:"
              echo "  CC_386   = $CC_386"
              echo "  CC_ARMV7 = $CC_ARMV7"
              echo "  CC_ARM64 = $CC_ARM64"

              echo ""

              echo "Using go config:"
              echo "  GOROOT   = $(go env GOROOT)"
              echo "  GOCACHE  = $(go env GOCACHE)"
              echo "  GOPATH   = $(go env GOPATH)"
            '';
          };
        }
      );
    };
}
