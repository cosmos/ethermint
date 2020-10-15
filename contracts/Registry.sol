// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity ^0.6.7;

import "./util/Upgradeable.sol";
import "./util/CustodianUpgradeable.sol";

import "./IRegistry.sol";

contract Registry is Upgradeable, CustodianUpgradeable, IRegistry {

    // a substitution to constructor to support ugradeability
    function init(address _custodian) public initializer {
        CustodianUpgradeable(this).initialize(_custodian);
    }

    mapping (bytes32 => PendingAddressValue) private pendingAddressValues;
    
    // contract name => contract address
    mapping (string => address) private contractAddresses;

    function getContractAddress(string memory _contractName) override public view returns (address) {
        return contractAddresses[_contractName];
    }

    function requestContractAddressChange(string memory _contractName, address _newContractAddress) public returns (bytes32 lockId) {
       require(_newContractAddress != contractAddresses[_contractName]);
       
       lockId = generateLockId();
       pendingAddressValues[lockId] = PendingAddressValue({
       	   contractName: _contractName,
           addressValue: _newContractAddress,
           set: true
       });

       emit ContractAddressLocked(lockId, _contractName, _newContractAddress);
    }

    function confirmContractAddressChange(bytes32 _lockId) public onlyCustodian {
        PendingAddressValue storage value = pendingAddressValues[_lockId];
        require(value.set == true);
        contractAddresses[value.contractName] = value.addressValue;
        emit ContractAddressConfirmed(_lockId, value.contractName, value.addressValue);
        delete pendingAddressValues[_lockId];
    }

    struct PendingAddressValue {
    	string  contractName; 
        address addressValue;
        bool    set;
    }
}
