package main

import (
	"log"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
	
)

func main() {

	//This method starts an in memory flow emulator
	// - it then looks at all the contracts in the deployment block for emulator and deploys them
	// - then it looks at all the accounts that does not have contracts in them and create those accounts. These can be used as stakeholders in your "storyline" below.
	// - when referencing accounts in the "storyline" below note that the default option is to prepened the network to the account name, This is done so that it is easy to run a storyline against emulator, tesnet and mainnet. This can be disabled with the `DoNotPrependNetworkToAccountNames` method on the g object below.

	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	g.TransactionFromFile("setup_account").SignProposeAndPayAs("account").RunPrintEventsFull()

	g.TransactionFromFile("mint_tokens").SignProposeAndPayAsService().AccountArgument("account").UFix64Argument("33.0").RunPrintEventsFull()

	// Run script that returns
	result := g.ScriptFromFile("get_balance").AccountArgument("account").RunFailOnError()
	log.Printf("Script returned %s", result)
}
