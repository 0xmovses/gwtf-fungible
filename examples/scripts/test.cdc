// test script to ensure code is running
import FungibleToken from "../contracts/FungibleToken.cdc"

pub fun main(account: Address): String {
    return getAccount(account).address.toString()
}
