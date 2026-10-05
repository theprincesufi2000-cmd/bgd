package main

import (
    "fmt"
    "os"
    "path/filepath"

    msi "go.digitalxero.dev/go-msi"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Fprintln(os.Stderr, "usage: msibuilder <stage-dir> <output.msi>")
        os.Exit(2)
    }
    stage, err := filepath.Abs(os.Args[1])
    if err != nil {
        panic(err)
    }
    outPath, err := filepath.Abs(os.Args[2])
    if err != nil {
        panic(err)
    }

    launcherPath := filepath.Join(stage, "Launch-Highway-Racing.bat")
    gamePath := filepath.Join(stage, "Game")
    if _, err := os.Stat(launcherPath); err != nil {
        fmt.Fprintln(os.Stderr, "missing launcher:", err)
        os.Exit(1)
    }
    if st, err := os.Stat(gamePath); err != nil || !st.IsDir() {
        fmt.Fprintln(os.Stderr, "missing Game directory")
        os.Exit(1)
    }

    launcher, err := msi.FileSourceFromPath(launcherPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, "read launcher:", err)
        os.Exit(1)
    }

    b := msi.NewPackage().
        WithProductName("Mustafa Highway Racing").
        WithManufacturer("Mustafa Ahmed Ali").
        WithVersion("1.0.0").
        WithProductCode("{D19D2C66-DA88-4D24-969D-0640D1732E61}").
        WithUpgradeCode("{A96E2EEB-2704-4CA4-9CC0-8EB2BA20E477}").
        InstallToProgramFiles()

    b.Feature("MainFeature").WithTitle("Mustafa Highway Racing").WithLevel(1)

    root := b.RootDirectory("INSTALLFOLDER", "Mustafa Highway Racing")
    root.Subdirectory("GAMEDIR", "Game")

    launcherComp := root.Component("Launcher").AssociateToFeature("MainFeature")
    launcherComp.WithFile("Launch-Highway-Racing.bat", launcher)

    launcherComp.Shortcut("Mustafa Highway Racing.lnk", "[INSTALLFOLDER]Launch-Highway-Racing.bat").
        InDirectory("ProgramMenuFolder").
        Description("Launch Mustafa Highway Racing")

    launcherComp.Shortcut("Mustafa Highway Racing.lnk", "[INSTALLFOLDER]Launch-Highway-Racing.bat").
        InDirectory("DesktopFolder").
        Description("Launch Mustafa Highway Racing")

    if err := b.AddTree(os.DirFS(gamePath), "GAMEDIR", "MainFeature"); err != nil {
        fmt.Fprintln(os.Stderr, "add game tree:", err)
        os.Exit(1)
    }

    pkg, err := b.Build()
    if err != nil {
        fmt.Fprintln(os.Stderr, "build MSI:", err)
        os.Exit(1)
    }

    f, err := os.Create(outPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, "create output:", err)
        os.Exit(1)
    }
    if err := pkg.WriteMSI(f); err != nil {
        _ = f.Close()
        fmt.Fprintln(os.Stderr, "write MSI:", err)
        os.Exit(1)
    }
    if err := f.Close(); err != nil {
        fmt.Fprintln(os.Stderr, "close output:", err)
        os.Exit(1)
    }
    fmt.Println(outPath)
}
