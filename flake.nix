{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
          devShells.default = pkgs.mkShell {
              buildInputs = with pkgs; [
                go
                protobuf
                protoc-gen-go
                protoc-gen-go-grpc
              ];
          };
      });
}
