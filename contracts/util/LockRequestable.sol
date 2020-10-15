// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity ^0.6.7;

/** @title  A contract for generating unique identifiers
  *
  * @notice  A contract that provides a identifier generation scheme,
  * guaranteeing uniqueness across all contracts that inherit from it,
  * as well as unpredictability of future identifiers.
  *
  * @dev  This contract is intended to be inherited by any contract that
  * implements the callback software pattern for cooperative custodianship.
  */
contract LockRequestable {

    // MEMBERS
    /// @notice  the count of all invocations of `generateLockId`.
    uint256 public lockRequestCount;

    // FUNCTIONS
    /** @notice  Returns a fresh unique identifier.
      *
      * @dev the generation scheme uses three components.
      * First, the blockhash of the previous block.
      * Second, the deployed address.
      * Third, the next value of the counter.
      * This ensure that identifiers are unique across all contracts
      * following this scheme, and that future identifiers are
      * unpredictable.
      *
      * @return lockId LockID a 32-byte unique identifier.
      */
    function generateLockId() internal returns (bytes32 lockId) {
        return keccak256(
        abi.encodePacked(block.number - 1, address(this), ++lockRequestCount)
        );
    }
}
