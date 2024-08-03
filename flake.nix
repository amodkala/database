{
    description = "Raft demo";

    inputs = {
        nixpkgs.url = "github:nixos/nixpkgs";
        flake-utils.url = "github:numtide/flake-utils";
    };

    outputs = { self, nixpkgs, flake-utils }:
        flake-utils.lib.eachDefaultSystem (system:
            with nixpkgs.legacyPackages.${system}; {

                defaultPackage = buildGoModule rec {
                    pname = "raft";
                    src = ./.;
                    version = "0.1";
                    vendorHash = null;
                    doCheck = false;
                };

                devShells.default = mkShell {
                    buildInputs = with pkgs; [
                        go
                        protobuf
                        protoc-gen-go
                        protoc-gen-go-grpc
                    ];
                };
            }
        );
}
