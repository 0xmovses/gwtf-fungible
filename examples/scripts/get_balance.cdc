// This script reads the balance field of an account's FlowToken Balance

import FungibleToken from 0xf8d6e0586b0a20c7
import ExampleToken from 0xf8d6e0586b0a20c7

pub fun main(account: Address): UFix64 {
    let acct = getAccount(account)
    let vaultRef = acct.getCapability(/public/exampleTokenBalance)
        .borrow<&ExampleToken.Vault{FungibleToken.Balance}>()
        ?? panic("Could not borrow Balance reference to the Vault")

    return vaultRef.balance
}