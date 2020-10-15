// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity ^0.6.7;

// IRegistry is a simple contract addresses registry with ACL.
interface IRegistry {
    // METHODS
    function getContractAddress(string calldata contractName) external view returns (address);

    // EVENTS
    event ContractAddressLocked(bytes32 lockId, string contractName, address newValue);
    event ContractAddressConfirmed(bytes32 lockId, string contractName, address newValue);
}
