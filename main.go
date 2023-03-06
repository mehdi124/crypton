package main

import (
	"log"

	"github.com/mehdi124/crypton/cli"
	"github.com/mehdi124/crypton/storage"
)

func main() {
	/*p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}*/

	wallets, err := storage.GetWalletsList()
	if err != nil {
		log.Println(err, "error")
	}

	if len(wallets) > 0 {

		cli.WalletList(wallets)

	} else {

		cli.StartList()

	}

}
