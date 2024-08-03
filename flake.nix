{
    description = "Raft demo";

    inputs = {
        nixpkgs.url = "github:nixos/nixpkgs";
        flake-utils.url = "github:numtide/flake-utils";
    };

    outputs = { self, nixpkgs, flake-utils }:
        flake-utils.lib.eachDefaultSystem (system:
            with nixpkgs.legacyPackages.${system}; rec {

                version = builtins.substring 0 8 self.lastModifiedDate;

                packages = rec {
                    bin = buildGoModule {
                        pname = "raft";
                        src = ./.;
                        inherit version;
                        vendorHash = null;
                        doCheck = false;
                    };

                    docker = dockerTools.buildLayeredImage {
                        name = "raft";
                        tag = "latest";
                        config.Cmd = "${bin}/bin/cmd";
                    };
                };

                devShells = {
                    default = mkShell {
                        buildInputs = with pkgs; [
                            go
                            protobuf
                            protoc-gen-go
                            protoc-gen-go-grpc
                        ];
                    };
                };
            }
        );
}
