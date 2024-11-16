function ExecuterEtZip2
{
    param (
        [string]$GOARCH,
        [string]$GOOS,
        [string]$GOARM="",
        [string]$fichierResultat,
        [string]$NomFichierZip
    )

    try {
    if (Test-Path $fichierResultat) {
        Write-Output "Suppression de $fichierResultat"
        Remove-Item $fichierResultat -verbose
    }
    if (Test-Path $NomFichierZip) {
        Write-Output "Suppression de $NomFichierZip"
        Remove-Item $NomFichierZip -verbose
    }

    $command = "set GOARCH=amd64"
    # Lance la commande et vérifie s'il y a une erreur
    Invoke-Expression $command -ErrorAction Stop
    Write-Output "La commande 1 a été exécutée avec succès."

    $command = "set GOOS=windows"
    # Lance la commande et vérifie s'il y a une erreur
    Invoke-Expression $command -ErrorAction Stop
    Write-Output "La commande 2 a été exécutée avec succès."
        
    if (![string]::IsNullOrEmpty($GOARM)) {
        $command = "set GOARM=5"
        # Lance la commande et vérifie s'il y a une erreur
        Invoke-Expression $command -ErrorAction Stop
        Write-Output "La commande 3 a été exécutée avec succès."
    }

    $command = "go build -o $fichierResultat ./cmd/gtools.go"
    # Lance la commande et vérifie s'il y a une erreur
    Invoke-Expression $command -ErrorAction Stop
    Write-Output "La commande 4 a été exécutée avec succès."

    Compress-Archive $fichierResultat $NomFichierZip
    Write-Output "Le fichier a ete compresse avec succes dans $NomFichierZip."

} catch {
    Write-Error "La commande a échoué. Arrêt du script."
    exit 1
}

}

Write-Output "Le repertoire du script est : $PSScriptRoot"

$repertoireDest = "$PSScriptRoot/dist"

Set-Location -Path $PSScriptRoot

$version="v1.3.0"

$GOARCH="amd64"
$GOOS="windows"
$fichierResultat="$PSScriptRoot/dist/gtools.exe"
$NomFichierZip="$PSScriptRoot/dist/gtools_$version_"+$GOOS+"_$GOARCH.zip"

ExecuterEtZip2 -GOARCH $GOARCH -GOOS $GOOS -fichierResultat $fichierResultat -NomFichierZip $NomFichierZip


$GOARCH="amd64"
$GOOS="linux"
$fichierResultat="$PSScriptRoot/dist/gtools"
$NomFichierZip="$PSScriptRoot/dist/gtools_$version_"+$GOOS+"_$GOARCH.zip"

ExecuterEtZip2 -GOARCH $GOARCH -GOOS $GOOS -fichierResultat $fichierResultat -NomFichierZip $NomFichierZip

$GOARCH="arm"
$GOOS="linux"
$GOARM="5"
$fichierResultat="$PSScriptRoot/dist/gtools"
$NomFichierZip="$PSScriptRoot/dist/gtools_$version_"+$GOOS+"_$GOARCH.zip"

ExecuterEtZip2 -GOARCH $GOARCH -GOOS $GOOS -GOARM $GOARM -fichierResultat $fichierResultat -NomFichierZip $NomFichierZip



