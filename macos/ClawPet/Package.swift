// swift-tools-version: 6.0

import PackageDescription

let package = Package(
    name: "ClawPet",
    platforms: [
        .macOS(.v14)
    ],
    products: [
        .executable(name: "ClawPet", targets: ["ClawPet"])
    ],
    targets: [
        .executableTarget(
            name: "ClawPet",
            path: "Sources"
        ),
        .testTarget(
            name: "ClawPetTests",
            dependencies: ["ClawPet"],
            path: "Tests"
        )
    ]
)
